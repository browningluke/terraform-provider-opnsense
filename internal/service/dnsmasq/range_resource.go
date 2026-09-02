package dnsmasq

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/errs"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &rangeResource{}
var _ resource.ResourceWithConfigure = &rangeResource{}
var _ resource.ResourceWithImportState = &rangeResource{}
var _ resource.ResourceWithValidateConfig = &rangeResource{}

func newRangeResource() resource.Resource {
	return &rangeResource{}
}

type rangeResource struct {
	client opnsense.Client
}

func (r *rangeResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data rangeResourceModel

	// Config in Model dekodieren
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	startAddress := data.StartAddress.ValueString()

	if startAddress != "" {
		ip := net.ParseIP(startAddress)
		if ip == nil {
			resp.Diagnostics.AddError(
				"Invalid IP address",
				fmt.Sprintf("start_address '%s' is not a valid IPv4 or IPv6 address.", startAddress),
			)
			return
		}

		if ip.To4() == nil {
			if data.Interface.IsNull() {
				resp.Diagnostics.AddError(
					"Invalid IPv6 configuration",
					"When using an IPv6 start_address, interface must be specified.",
				)
			}

			if data.Interface.IsNull() || data.Interface.ValueString() == "" {
				resp.Diagnostics.AddError(
					"Invalid IPv6 configuration",
					"When ipv6 is enabled, interface must be specified.",
				)
			}

			if data.RaMode.IsNull() {
				resp.Diagnostics.AddError(
					"Missing RA mode",
					"ra_mode must be set when ipv6 is enabled.",
				)
			}
		}
	}
}

func (r *rangeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dnsmasq_range"
}

func (r *rangeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rangeResourceSchema()
}

func (r *rangeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *rangeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *rangeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	rangeObj, err := convertRangeSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse range, got error: %s", err))
		return
	}

	// Add range to dnsmasq
	id, err := r.client.Dnsmasq().AddRange(ctx, rangeObj)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create range, got error: %s", err))
		return
	}

	// Range new resource with ID from OPNsense
	data.Id = types.StringValue(id)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rangeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *rangeResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get resource from OPNsense dnsmasq API
	rangeObj, err := r.client.Dnsmasq().GetRange(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, "range not present in remote, removing from state")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read range, got error: %s", err))
		return
	}

	// Convert OPNsense struct to TF schema
	rangeModel, err := convertRangeStructToSchema(rangeObj)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read range, got error: %s", err))
		return
	}
	rangeModel.Id = data.Id

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &rangeModel)...)
}

func (r *rangeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *rangeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	rangeObj, err := convertRangeSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse range, got error: %s", err))
		return
	}

	// Update range in dnsmasq
	err = r.client.Dnsmasq().UpdateRange(ctx, data.Id.ValueString(), rangeObj)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update range, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rangeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *rangeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Dnsmasq().DeleteRange(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete range, got error: %s", err))
		return
	}
}

func (r *rangeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
