package trust

import (
	"context"
	"errors"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/errs"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &caResource{}
var _ resource.ResourceWithConfigure = &caResource{}
var _ resource.ResourceWithImportState = &caResource{}


type caResource struct {
	client opnsense.Client
}

func (r *caResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trust_ca"
}

func (r *caResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = caResourceSchema()
}

func (r *caResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = opnsense.NewClient(apiClient)
}

func (r *caResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *caResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ca, err := convertCaSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse CA, got error: %s", err))
		return
	}

	id, err := r.client.Trust().AddCa(ctx, ca)
	if err != nil {
		if id != "" {
			data.Id = types.StringValue(id)
			if readStruct, readErr := r.client.Trust().GetCa(ctx, id); readErr == nil {
				if readModel, convErr := convertCaStructToSchema(readStruct); convErr == nil {
					readModel.Id = data.Id
					preserveCaStateFields(readModel, data)
					data = readModel
				}
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create CA, got error: %s", err))
		return
	}

	data.Id = types.StringValue(id)

	caStruct, err := r.client.Trust().GetCa(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read CA after creation, got error: %s", err))
		return
	}

	caModel, err := convertCaStructToSchema(caStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse CA after creation, got error: %s", err))
		return
	}

	caModel.Id = data.Id
	// Preserve user-specified fields that don't roundtrip
	preserveCaStateFields(caModel, data)

	tflog.Trace(ctx, "created trust_ca resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &caModel)...)
}

func (r *caResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *caResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ca, err := r.client.Trust().GetCa(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, fmt.Sprintf("trust_ca not present in remote, removing from state"))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read CA, got error: %s", err))
		return
	}

	caModel, err := convertCaStructToSchema(ca)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse CA, got error: %s", err))
		return
	}

	caModel.Id = data.Id
	// Preserve user-specified fields that don't roundtrip from the API
	preserveCaStateFields(caModel, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &caModel)...)
}

func (r *caResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *caResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ca, err := convertCaSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse CA, got error: %s", err))
		return
	}

	err = r.client.Trust().UpdateCa(ctx, data.Id.ValueString(), ca)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update CA, got error: %s", err))
		return
	}

	caStruct, err := r.client.Trust().GetCa(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read CA after update, got error: %s", err))
		return
	}

	caModel, err := convertCaStructToSchema(caStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse CA after update, got error: %s", err))
		return
	}

	caModel.Id = data.Id
	preserveCaStateFields(caModel, data)

	tflog.Trace(ctx, "updated trust_ca resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &caModel)...)
}

func (r *caResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *caResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Trust().DeleteCa(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete CA, got error: %s", err))
		return
	}
}

func (r *caResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// preserveCaStateFields copies fields from prior state that either don't
// roundtrip correctly or can vary between API calls into the freshly-read model.
func preserveCaStateFields(model *caResourceModel, state *caResourceModel) {
	if !state.Action.IsNull() && !state.Action.IsUnknown() {
		model.Action = state.Action
	}
	if !state.Lifetime.IsNull() && !state.Lifetime.IsUnknown() {
		model.Lifetime = state.Lifetime
	}
	// Preserve key material and validity timestamps from prior state once set.
	if !state.Crt.IsNull() && !state.Crt.IsUnknown() && state.Crt.ValueString() != "" {
		model.Crt = state.Crt
	}
	if !state.Prv.IsNull() && !state.Prv.IsUnknown() && state.Prv.ValueString() != "" {
		model.Prv = state.Prv
	}
	if !state.CrtPayload.IsNull() && !state.CrtPayload.IsUnknown() && state.CrtPayload.ValueString() != "" {
		model.CrtPayload = state.CrtPayload
	}
	if !state.PrvPayload.IsNull() && !state.PrvPayload.IsUnknown() && state.PrvPayload.ValueString() != "" {
		model.PrvPayload = state.PrvPayload
	}
	if !state.ValidFrom.IsNull() && !state.ValidFrom.IsUnknown() && state.ValidFrom.ValueString() != "" {
		model.ValidFrom = state.ValidFrom
	}
	if !state.ValidTo.IsNull() && !state.ValidTo.IsUnknown() && state.ValidTo.ValueString() != "" {
		model.ValidTo = state.ValidTo
	}
}
