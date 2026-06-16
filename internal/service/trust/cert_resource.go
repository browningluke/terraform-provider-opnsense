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

var _ resource.Resource = &certResource{}
var _ resource.ResourceWithConfigure = &certResource{}
var _ resource.ResourceWithImportState = &certResource{}

func newCertResource() resource.Resource {
	return &certResource{}
}

type certResource struct {
	client opnsense.Client
}

func (r *certResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trust_cert"
}

func (r *certResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = certResourceSchema()
}

func (r *certResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *certResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *certResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cert, err := convertCertSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse certificate, got error: %s", err))
		return
	}

	id, err := r.client.Trust().AddCert(ctx, cert)
	if err != nil {
		if id != "" {
			data.Id = types.StringValue(id)
			if readStruct, readErr := r.client.Trust().GetCert(ctx, id); readErr == nil {
				if readModel, convErr := convertCertStructToSchema(readStruct); convErr == nil {
					readModel.Id = data.Id
					preserveCertStateFields(readModel, data)
					data = readModel
				}
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create certificate, got error: %s", err))
		return
	}

	data.Id = types.StringValue(id)

	certStruct, err := r.client.Trust().GetCert(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read certificate after creation, got error: %s", err))
		return
	}

	certModel, err := convertCertStructToSchema(certStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse certificate after creation, got error: %s", err))
		return
	}

	certModel.Id = data.Id
	preserveCertStateFields(certModel, data)

	tflog.Trace(ctx, "created trust_cert resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &certModel)...)
}

func (r *certResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *certResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cert, err := r.client.Trust().GetCert(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, fmt.Sprintf("trust_cert not present in remote, removing from state"))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read certificate, got error: %s", err))
		return
	}

	certModel, err := convertCertStructToSchema(cert)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse certificate, got error: %s", err))
		return
	}

	certModel.Id = data.Id
	preserveCertStateFields(certModel, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &certModel)...)
}

func (r *certResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *certResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cert, err := convertCertSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse certificate, got error: %s", err))
		return
	}

	err = r.client.Trust().UpdateCert(ctx, data.Id.ValueString(), cert)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update certificate, got error: %s", err))
		return
	}

	certStruct, err := r.client.Trust().GetCert(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read certificate after update, got error: %s", err))
		return
	}

	certModel, err := convertCertStructToSchema(certStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse certificate after update, got error: %s", err))
		return
	}

	certModel.Id = data.Id
	preserveCertStateFields(certModel, data)

	tflog.Trace(ctx, "updated trust_cert resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &certModel)...)
}

func (r *certResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *certResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Trust().DeleteCert(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete certificate, got error: %s", err))
		return
	}
}

func (r *certResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// preserveCertStateFields copies fields from prior state that either don't
// roundtrip correctly or can vary between API calls into the freshly-read model.
func preserveCertStateFields(model *certResourceModel, state *certResourceModel) {
	if !state.Action.IsNull() && !state.Action.IsUnknown() {
		model.Action = state.Action
	}
	// Preserve key material and validity timestamps from prior state once set.
	// OPNsense may return slightly different timestamps between reads.
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
