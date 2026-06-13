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
var _ resource.Resource = &bgpPeerGroupResource{}
var _ resource.ResourceWithConfigure = &bgpPeerGroupResource{}
var _ resource.ResourceWithImportState = &bgpPeerGroupResource{}

func newBGPPeerGroupResource() resource.Resource {
	return &bgpPeerGroupResource{}
}

// bgpPeerGroupResource defines the resource implementation.
type bgpPeerGroupResource struct {
	client opnsense.Client
}

func (r *bgpPeerGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quagga_bgp_peer_group"
}

func (r *bgpPeerGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = bgpPeerGroupResourceSchema()
}

func (r *bgpPeerGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (r *bgpPeerGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *bgpPeerGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	bgpPeerGroup, err := convertBGPPeerGroupSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse bgp peer group, got error: %s", err))
		return
	}

	// Add bgp peer group to OPNsense
	id, err := r.client.Quagga().AddBGPPeerGroup(ctx, bgpPeerGroup)
	if err != nil {
		if id != "" {
			data.Id = types.StringValue(id)

			// Read back so state captures API-normalised values (defaults,
			// sorting, trimming); fall back to plan-only state if the
			// read-back fails so the upstream resource isn't orphaned.
			if readStruct, readErr := r.client.Quagga().GetBGPPeerGroup(ctx, id); readErr == nil {
				if readModel, convErr := convertBGPPeerGroupStructToSchema(readStruct); convErr == nil {
					readModel.Id = data.Id
					data = readModel
				}
			}

			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create bgp peer group, got error: %s", err))
		return
	}

	// Tag new resource with ID from OPNsense
	data.Id = types.StringValue(id)

	// Read back to populate computed fields (e.g. listen_ranges) with concrete values.
	if readStruct, readErr := r.client.Quagga().GetBGPPeerGroup(ctx, id); readErr == nil {
		if readModel, convErr := convertBGPPeerGroupStructToSchema(readStruct); convErr == nil {
			readModel.Id = data.Id
			data = readModel
		}
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *bgpPeerGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *bgpPeerGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get bgp peer group from OPNsense API
	bgpPeerGroup, err := r.client.Quagga().GetBGPPeerGroup(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, fmt.Sprintf("bgp peer group not present in remote, removing from state"))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read bgp peer group, got error: %s", err))
		return
	}

	// Convert OPNsense struct to TF schema
	bgpPeerGroupModel, err := convertBGPPeerGroupStructToSchema(bgpPeerGroup)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read bgp peer group, got error: %s", err))
		return
	}

	// ID cannot be added by convert... func, have to add here
	bgpPeerGroupModel.Id = data.Id

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &bgpPeerGroupModel)...)
}

func (r *bgpPeerGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *bgpPeerGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	bgpPeerGroup, err := convertBGPPeerGroupSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse bgp peer group, got error: %s", err))
		return
	}

	// Update bgp peer group in OPNsense
	err = r.client.Quagga().UpdateBGPPeerGroup(ctx, data.Id.ValueString(), bgpPeerGroup)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update bgp peer group, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *bgpPeerGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *bgpPeerGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Quagga().DeleteBGPPeerGroup(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete bgp peer group, got error: %s", err))
		return
	}
}

func (r *bgpPeerGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
