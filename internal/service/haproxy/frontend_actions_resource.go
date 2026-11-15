package haproxy

import (
	"context"
	"errors"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/errs"
	"github.com/browningluke/opnsense-go/pkg/haproxy"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &frontendActionsResource{}
var _ resource.ResourceWithConfigure = &frontendActionsResource{}
var _ resource.ResourceWithImportState = &frontendActionsResource{}

func newFrontendActionsResource() resource.Resource {
	return &frontendActionsResource{}
}

type frontendActionsResource struct {
	client opnsense.Client
}

func (r *frontendActionsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_haproxy_frontend_actions"
}

func (r *frontendActionsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = frontendActionsResourceSchema()
}

func (r *frontendActionsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *opnsense.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = opnsense.NewClient(apiClient)
}

func (r *frontendActionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *frontendActionsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	frontendID := data.FrontendId.ValueString()
	// Fail with a clear message instead of a validation error when the frontend does not exist.
	if _, err := r.client.Haproxy().GetFrontend(ctx, frontendID); err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read haproxy frontend %s, got error: %s", frontendID, err))
		return
	}

	frontend, diags := convertFrontendActionsSchemaToStruct(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Haproxy().UpdateFrontend(ctx, frontendID, frontend); err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to set haproxy frontend actions, got error: %s", err))
		return
	}

	data.Id = data.FrontendId
	tflog.Trace(ctx, "created a resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *frontendActionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *frontendActionsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	frontend, err := r.client.Haproxy().GetFrontend(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, "haproxy frontend not present in remote, removing from state")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read haproxy frontend actions, got error: %s", err))
		return
	}

	model, diags := convertFrontendActionsStructToSchema(ctx, data.Id.ValueString(), frontend)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *frontendActionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *frontendActionsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	frontend, diags := convertFrontendActionsSchemaToStruct(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Haproxy().UpdateFrontend(ctx, data.Id.ValueString(), frontend); err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to set haproxy frontend actions, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *frontendActionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *frontendActionsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The frontend is not owned here, so only its rule list is cleared.
	err := r.client.Haproxy().UpdateFrontend(ctx, data.Id.ValueString(), &haproxy.Frontend{LinkedActions: api.SelectedMapListOrdered{}})
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			return
		}
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to clear haproxy frontend actions, got error: %s", err))
		return
	}
}

func (r *frontendActionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
