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
var _ resource.Resource = &ospfPrefixListResource{}
var _ resource.ResourceWithConfigure = &ospfPrefixListResource{}
var _ resource.ResourceWithImportState = &ospfPrefixListResource{}

func newOSPFPrefixListResource() resource.Resource {
	return &ospfPrefixListResource{}
}

// ospfPrefixListResource defines the resource implementation.
type ospfPrefixListResource struct {
	client opnsense.Client
}

func (r *ospfPrefixListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quagga_ospf_prefix_list"
}

func (r *ospfPrefixListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ospfPrefixListResourceSchema()
}

func (r *ospfPrefixListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ospfPrefixListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ospfPrefixListResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ospfPrefixList, err := convertOSPFPrefixListSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ospf prefix list, got error: %s", err))
		return
	}

	id, err := r.client.Quagga().AddOSPFPrefixList(ctx, ospfPrefixList)
	if err != nil {
		if id != "" {
			data.Id = types.StringValue(id)

			if readStruct, readErr := r.client.Quagga().GetOSPFPrefixList(ctx, id); readErr == nil {
				if readModel, convErr := convertOSPFPrefixListStructToSchema(readStruct); convErr == nil {
					readModel.Id = data.Id
					data = readModel
				}
			}

			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create ospf prefix list, got error: %s", err))
		return
	}

	data.Id = types.StringValue(id)

	tflog.Trace(ctx, "created a resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ospfPrefixListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ospfPrefixListResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ospfPrefixList, err := r.client.Quagga().GetOSPFPrefixList(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, fmt.Sprintf("ospf prefix list not present in remote, removing from state"))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ospf prefix list, got error: %s", err))
		return
	}

	ospfPrefixListModel, err := convertOSPFPrefixListStructToSchema(ospfPrefixList)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ospf prefix list, got error: %s", err))
		return
	}

	ospfPrefixListModel.Id = data.Id

	resp.Diagnostics.Append(resp.State.Set(ctx, &ospfPrefixListModel)...)
}

func (r *ospfPrefixListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ospfPrefixListResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ospfPrefixList, err := convertOSPFPrefixListSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ospf prefix list, got error: %s", err))
		return
	}

	err = r.client.Quagga().UpdateOSPFPrefixList(ctx, data.Id.ValueString(), ospfPrefixList)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create ospf prefix list, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ospfPrefixListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ospfPrefixListResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Quagga().DeleteOSPFPrefixList(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete ospf prefix list, got error: %s", err))
		return
	}
}

func (r *ospfPrefixListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
