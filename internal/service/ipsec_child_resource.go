package service

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
var _ resource.Resource = &IpsecChildResource{}
var _ resource.ResourceWithImportState = &IpsecChildResource{}

func NewIpsecChildResource() resource.Resource {
	return &IpsecChildResource{}
}

// IpsecChildResource defines the resource implementation.
type IpsecChildResource struct {
	client opnsense.Client
}

func (r *IpsecChildResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipsec_child"
}

func (r *IpsecChildResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = IpsecChildResourceSchema()
}

func (r *IpsecChildResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (r *IpsecChildResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *IpsecChildResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema to OPNsense struct
	child, err := convertIpsecChildSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ipsec child, got error: %s", err))
		return
	}

	// Add IPsec Child to OPNsense
	id, err := r.client.Ipsec().AddIPsecChild(ctx, child)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create ipsec child, got error: %s", err))
		return
	}

	// Tag new resource with ID from OPNsense
	data.Id = types.StringValue(id)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IpsecChildResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *IpsecChildResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IPsec Child from OPNsense API
	child, err := r.client.Ipsec().GetIPsecChild(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, "ipsec child not present in remote, removing from state")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ipsec child, got error: %s", err))
		return
	}

	// Convert OPNsense struct to TF schema
	childModel, err := convertIpsecChildStructToSchema(child)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ipsec child, got error: %s", err))
		return
	}

	// ID cannot be added by convert... func, have to add here
	childModel.Id = data.Id

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &childModel)...)
}

func (r *IpsecChildResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *IpsecChildResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema to OPNsense struct
	child, err := convertIpsecChildSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse ipsec child, got error: %s", err))
		return
	}

	// Update IPsec Child in OPNsense
	err = r.client.Ipsec().UpdateIPsecChild(ctx, data.Id.ValueString(), child)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update ipsec child, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IpsecChildResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *IpsecChildResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Ipsec().DeleteIPsecChild(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete ipsec child, got error: %s", err))
		return
	}
}

func (r *IpsecChildResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}