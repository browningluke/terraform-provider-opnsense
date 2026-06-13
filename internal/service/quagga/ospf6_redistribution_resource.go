package quagga

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

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ospf6RedistributionResource{}
var _ resource.ResourceWithConfigure = &ospf6RedistributionResource{}
var _ resource.ResourceWithImportState = &ospf6RedistributionResource{}

func newOSPF6RedistributionResource() resource.Resource {
	return &ospf6RedistributionResource{}
}

// ospf6RedistributionResource defines the resource implementation.
type ospf6RedistributionResource struct {
	client opnsense.Client
}

func (r *ospf6RedistributionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quagga_ospf6_redistribution"
}

func (r *ospf6RedistributionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ospf6RedistributionResourceSchema()
}

func (r *ospf6RedistributionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ospf6RedistributionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ospf6RedistributionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ospf6Redistribution, err := convertOSPF6RedistributionSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ospf6 redistribution, got error: %s", err))
		return
	}

	id, err := r.client.Quagga().AddOSPF6Redistribution(ctx, ospf6Redistribution)
	if err != nil {
		if id != "" {
			data.Id = types.StringValue(id)

			if readStruct, readErr := r.client.Quagga().GetOSPF6Redistribution(ctx, id); readErr == nil {
				if readModel, convErr := convertOSPF6RedistributionStructToSchema(readStruct); convErr == nil {
					readModel.Id = data.Id
					data = readModel
				}
			}

			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create ospf6 redistribution, got error: %s", err))
		return
	}

	data.Id = types.StringValue(id)

	tflog.Trace(ctx, "created a resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ospf6RedistributionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ospf6RedistributionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ospf6Redistribution, err := r.client.Quagga().GetOSPF6Redistribution(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, fmt.Sprintf("ospf6 redistribution not present in remote, removing from state"))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ospf6 redistribution, got error: %s", err))
		return
	}

	ospf6RedistributionModel, err := convertOSPF6RedistributionStructToSchema(ospf6Redistribution)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ospf6 redistribution, got error: %s", err))
		return
	}

	ospf6RedistributionModel.Id = data.Id

	resp.Diagnostics.Append(resp.State.Set(ctx, &ospf6RedistributionModel)...)
}

func (r *ospf6RedistributionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ospf6RedistributionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ospf6Redistribution, err := convertOSPF6RedistributionSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ospf6 redistribution, got error: %s", err))
		return
	}

	err = r.client.Quagga().UpdateOSPF6Redistribution(ctx, data.Id.ValueString(), ospf6Redistribution)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create ospf6 redistribution, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ospf6RedistributionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ospf6RedistributionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Quagga().DeleteOSPF6Redistribution(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete ospf6 redistribution, got error: %s", err))
		return
	}
}

func (r *ospf6RedistributionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
