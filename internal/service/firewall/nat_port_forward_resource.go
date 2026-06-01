package firewall

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
var _ resource.Resource = &natPortForwardResource{}
var _ resource.ResourceWithConfigure = &natPortForwardResource{}
var _ resource.ResourceWithImportState = &natPortForwardResource{}
var _ resource.ResourceWithUpgradeState = &natPortForwardResource{}

func newNATPortForwardResource() resource.Resource {
	return &natPortForwardResource{}
}

// natPortForwardResource defines the resource implementation.
type natPortForwardResource struct {
	client opnsense.Client
}

func (r *natPortForwardResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_nat_port_forward"
}

func (r *natPortForwardResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = natPortForwardResourceSchema()
}

func (r *natPortForwardResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *natPortForwardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *natPortForwardResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	portForward, err := convertNATPortForwardSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse firewall nat port forward, got error: %s", err))
		return
	}

	// Add firewall nat port forward
	id, err := r.client.Firewall().AddNatPortForward(ctx, portForward)
	if err != nil {
		if id != "" {
			// Tag new resource with ID from OPNsense
			data.Id = types.StringValue(id)

			// Save data into Terraform state
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create firewall nat port forward, got error: %s", err))
		return
	}

	// Tag new resource with ID from OPNsense
	data.Id = types.StringValue(id)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *natPortForwardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *natPortForwardResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get firewall nat port forward from OPNsense API
	resourceStruct, err := r.client.Firewall().GetNatPortForward(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, fmt.Sprintf("firewall nat port forward not present in remote, removing from state"))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read firewall nat port forward, got error: %s", err))
		return
	}

	// Convert OPNsense struct to TF schema
	resourceModel, err := convertNATPortForwardStructToSchema(resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read firewall nat port forward, got error: %s", err))
		return
	}

	// ID cannot be added by convert... func, have to add here
	resourceModel.Id = data.Id

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}

func (r *natPortForwardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *natPortForwardResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	portForward, err := convertNATPortForwardSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse firewall nat port forward, got error: %s", err))
		return
	}

	// Update firewall nat port forward
	err = r.client.Firewall().UpdateNatPortForward(ctx, data.Id.ValueString(), portForward)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update firewall nat port forward, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *natPortForwardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *natPortForwardResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Firewall().DeleteNatPortForward(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete firewall nat port forward, got error: %s", err))
		return
	}
}

func (r *natPortForwardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *natPortForwardResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	schemaV1 := natPortForwardResourceSchemaV1()
	return map[int64]resource.StateUpgrader{
		1: {
			PriorSchema:   &schemaV1,
			StateUpgrader: upgradeNATPortForwardStateV1toV2,
		},
	}
}

func upgradeNATPortForwardStateV1toV2(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	tflog.Info(ctx, "Upgrading NAT port forward resource state from v1 to v2")

	var oldState natPortForwardResourceModelV1
	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to read old NAT port forward state during upgrade")
		return
	}

	newState := &natPortForwardResourceModel{
		Enabled:       oldState.Enabled,
		Sequence:      oldState.Sequence,
		Interface:     natPortForwardInterfaceStringToSet(oldState.Interface.ValueString()),
		IPProtocol:    oldState.IPProtocol,
		Protocol:      oldState.Protocol,
		Source:        oldState.Source,
		Destination:   oldState.Destination,
		Target:        oldState.Target,
		Log:           oldState.Log,
		NatReflection: oldState.NatReflection,
		Description:   oldState.Description,
		Id:            oldState.Id,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)

	if !resp.Diagnostics.HasError() {
		tflog.Info(ctx, "Successfully upgraded NAT port forward resource state from v1 to v2", map[string]any{
			"id": oldState.Id.ValueString(),
		})
	}
}
