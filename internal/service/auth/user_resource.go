package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/errs"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &userResource{}
var _ resource.ResourceWithConfigure = &userResource{}
var _ resource.ResourceWithImportState = &userResource{}
var _ resource.ResourceWithConfigValidators = &userResource{}

func newUserResource() resource.Resource {
	return &userResource{}
}

func (r *userResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		passwordScrambleConflictValidator{},
	}
}

type userResource struct {
	client opnsense.Client
	api    *api.Client
}

func (r *userResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_user"
}

func (r *userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = userResourceSchema()
}

func (r *userResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.api = apiClient
	r.client = opnsense.NewClient(apiClient)
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *userResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupGIDs, err := resolveGroupMembershipsToGIDs(ctx, r.api, tools.SetToStringSlice(data.Groups))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to resolve groups: %s", err))
		return
	}

	user, err := convertUserSchemaToStruct(data, groupGIDs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to parse auth user, got error: %s", err))
		return
	}

	id, err := r.client.Auth().AddUser(ctx, user)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create auth user, got error: %s", err))
		return
	}

	data.Id = types.StringValue(id)

	readModel, err := r.readUserModel(ctx, id, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Created auth user but failed to refresh: %s", err))
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	tflog.Trace(ctx, "created auth user")
	resp.Diagnostics.Append(resp.State.Set(ctx, &readModel)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *userResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readModel, err := r.readUserModel(ctx, data.Id.ValueString(), data)
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, "auth user not present in remote, removing from state")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read auth user, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &readModel)...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *userResourceModel
	var state *userResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupGIDs, err := resolveGroupMembershipsToGIDs(ctx, r.api, tools.SetToStringSlice(data.Groups))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to resolve groups: %s", err))
		return
	}

	// Password is not returned by the API. Only send a password change when
	// explicitly set in config; otherwise keep the existing hash untouched.
	// Do not mutate planned scramble_password (would cause inconsistent state).
	apiData := *data
	if data.Password.IsNull() || data.Password.IsUnknown() || data.Password.ValueString() == "" {
		apiData.ScramblePassword = types.BoolValue(false)
	}

	user, err := convertUserSchemaToStruct(&apiData, groupGIDs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to parse auth user, got error: %s", err))
		return
	}

	err = r.client.Auth().UpdateUser(ctx, data.Id.ValueString(), user)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update auth user, got error: %s", err))
		return
	}

	readModel, err := r.readUserModel(ctx, data.Id.ValueString(), data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Updated auth user but failed to refresh: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &readModel)...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *userResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Auth().DeleteUser(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete auth user, got error: %s", err))
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *userResource) readUserModel(ctx context.Context, id string, prior *userResourceModel) (*userResourceModel, error) {
	user, err := r.client.Auth().GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	groupNames, err := groupGIDsToNames(ctx, r.api, []string(user.GroupMemberships))
	if err != nil {
		return nil, err
	}

	model, err := convertUserStructToSchema(user, groupNames)
	if err != nil {
		return nil, err
	}
	model.Id = types.StringValue(id)

	// Password / scramble flags are write-only from Terraform's point of view.
	if prior != nil {
		model.Password = prior.Password
		if !prior.ScramblePassword.IsNull() && !prior.ScramblePassword.IsUnknown() {
			model.ScramblePassword = prior.ScramblePassword
		} else {
			model.ScramblePassword = types.BoolValue(true)
		}
	}

	if model.Groups.IsNull() {
		model.Groups = tools.StringSliceToSet([]string{})
	}
	if model.Privileges.IsNull() {
		model.Privileges = tools.StringSliceToSet([]string{})
	}
	if model.AuthorizedKeys.IsNull() {
		model.AuthorizedKeys = tools.StringSliceToSet([]string{})
	}

	return model, nil
}
