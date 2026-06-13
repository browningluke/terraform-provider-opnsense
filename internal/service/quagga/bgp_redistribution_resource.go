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
var _ resource.Resource = &bgpRedistributionResource{}
var _ resource.ResourceWithConfigure = &bgpRedistributionResource{}
var _ resource.ResourceWithImportState = &bgpRedistributionResource{}

func newBGPRedistributionResource() resource.Resource {
	return &bgpRedistributionResource{}
}

// bgpRedistributionResource defines the resource implementation.
type bgpRedistributionResource struct {
	client opnsense.Client
}

func (r *bgpRedistributionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quagga_bgp_redistribution"
}

func (r *bgpRedistributionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = bgpRedistributionResourceSchema()
}

func (r *bgpRedistributionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *bgpRedistributionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *bgpRedistributionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	bgpRedistribution, err := convertBGPRedistributionSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse bgp redistribution, got error: %s", err))
		return
	}

	// Add bgp redistribution to OPNsense
	id, err := r.client.Quagga().AddBGPRedistribution(ctx, bgpRedistribution)
	if err != nil {
		if id != "" {
			data.Id = types.StringValue(id)

			// Read back so state captures API-normalised values (defaults,
			// sorting, trimming); fall back to plan-only state if the
			// read-back fails so the upstream resource isn't orphaned.
			if readStruct, readErr := r.client.Quagga().GetBGPRedistribution(ctx, id); readErr == nil {
				if readModel, convErr := convertBGPRedistributionStructToSchema(readStruct); convErr == nil {
					readModel.Id = data.Id
					data = readModel
				}
			}

			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create bgp redistribution, got error: %s", err))
		return
	}

	// Tag new resource with ID from OPNsense
	data.Id = types.StringValue(id)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *bgpRedistributionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *bgpRedistributionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get bgp redistribution from OPNsense API
	bgpRedistribution, err := r.client.Quagga().GetBGPRedistribution(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, fmt.Sprintf("bgp redistribution not present in remote, removing from state"))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read bgp redistribution, got error: %s", err))
		return
	}

	// Convert OPNsense struct to TF schema
	bgpRedistributionModel, err := convertBGPRedistributionStructToSchema(bgpRedistribution)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read bgp redistribution, got error: %s", err))
		return
	}

	// ID cannot be added by convert... func, have to add here
	bgpRedistributionModel.Id = data.Id

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &bgpRedistributionModel)...)
}

func (r *bgpRedistributionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *bgpRedistributionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	bgpRedistribution, err := convertBGPRedistributionSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse bgp redistribution, got error: %s", err))
		return
	}

	// Update bgp redistribution in OPNsense
	err = r.client.Quagga().UpdateBGPRedistribution(ctx, data.Id.ValueString(), bgpRedistribution)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update bgp redistribution, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *bgpRedistributionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *bgpRedistributionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Quagga().DeleteBGPRedistribution(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete bgp redistribution, got error: %s", err))
		return
	}
}

func (r *bgpRedistributionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
