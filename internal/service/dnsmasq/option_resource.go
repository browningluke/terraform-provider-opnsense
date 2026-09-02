package dnsmasq

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

var _ resource.Resource = &optionResource{}
var _ resource.ResourceWithConfigure = &optionResource{}
var _ resource.ResourceWithImportState = &optionResource{}
var _ resource.ResourceWithValidateConfig = &optionResource{}

func newOptionResource() resource.Resource {
	return &optionResource{}
}

type optionResource struct {
	client opnsense.Client
}

func (r *optionResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data optionResourceModel

	// Config in Model dekodieren
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	v4OptionEmpty := (data.OptionV4.IsNull() || data.OptionV4.IsUnknown())
	v6OptionEmpty := (data.OptionV6.IsNull() || data.OptionV6.IsUnknown())
	if v4OptionEmpty && v6OptionEmpty {
		resp.Diagnostics.AddError(
			"DHCP option v4 or v6 need to be set",
			"",
		)
	}
	if !v4OptionEmpty && !v6OptionEmpty {
		resp.Diagnostics.AddError(
			"Only one DHCP option v4 or v6 can be set",
			"",
		)
	}

	hasTags := !data.TypeSetTags.IsNull() &&
		!data.TypeSetTags.IsUnknown() &&
		len(data.TypeSetTags.Elements()) > 0
	hasSetTag := !data.TypeMatchTag.IsNull() &&
		!data.TypeMatchTag.IsUnknown() &&
		data.TypeMatchTag.ValueString() != ""

	switch data.Type.ValueString() {
	case "set":
		if hasSetTag {
			resp.Diagnostics.AddAttributeError(
				path.Root("set_tag"),
				"Invalid DHCP option configuration",
				"set_tag is only valid when type is \"match\".",
			)
		}
	case "match":
		if hasTags {
			resp.Diagnostics.AddAttributeError(
				path.Root("tag"),
				"Invalid DHCP option configuration",
				"tag is only valid when type is \"set\".",
			)
		}
		if !hasSetTag {
			resp.Diagnostics.AddAttributeError(
				path.Root("set_tag"),
				"Missing DHCP match tag",
				"set_tag must be set when type is \"match\".",
			)
		}
	}
}

func (r *optionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dnsmasq_option"
}

func (r *optionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = optionResourceSchema()
}

func (r *optionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *optionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *optionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	option, err := convertOptionSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse option, got error: %s", err))
		return
	}

	// Add option to dnsmasq
	id, err := r.client.Dnsmasq().AddOption(ctx, option)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create option, got error: %s", err))
		return
	}

	// Option new resource with ID from OPNsense
	data.Id = types.StringValue(id)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *optionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *optionResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get resource from OPNsense dnsmasq API
	option, err := r.client.Dnsmasq().GetOption(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, "option not present in remote, removing from state")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read option, got error: %s", err))
		return
	}

	// Convert OPNsense struct to TF schema
	optionModel, err := convertOptionStructToSchema(option)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read option, got error: %s", err))
		return
	}
	optionModel.Id = data.Id

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &optionModel)...)
}

func (r *optionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *optionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	option, err := convertOptionSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse option, got error: %s", err))
		return
	}

	// Update dhcp option in dnsmasq
	err = r.client.Dnsmasq().UpdateOption(ctx, data.Id.ValueString(), option)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update option, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *optionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *optionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Dnsmasq().DeleteOption(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete option, got error: %s", err))
		return
	}
}

func (r *optionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
