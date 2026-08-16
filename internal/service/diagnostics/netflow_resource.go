package diagnostics

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &netflowResource{}
var _ resource.ResourceWithConfigure = &netflowResource{}
var _ resource.ResourceWithImportState = &netflowResource{}

func newNetflowResource() resource.Resource {
	return &netflowResource{}
}

// netflowResource defines the resource implementation.
// This is a SINGLETON resource - it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type netflowResource struct {
	client opnsense.Client
}

func (r *netflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_diagnostics_netflow"
}

func (r *netflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = netflowResourceSchema()
}

func (r *netflowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create is blocked for singleton resources. Users must import the resource first.
func (r *netflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *netflowResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// IMPORTANT: Block resource creation with a clear error message
	resp.Diagnostics.AddError(
		"Cannot Create Singleton Resource",
		"This resource manages existing upstream Netflow configuration that cannot be created or destroyed.\n\n"+
			"To manage this resource, you must import it first:\n"+
			"  terraform import opnsense_diagnostics_netflow.<name> netflow_settings\n\n"+
			"After importing, you can manage the configuration with terraform apply.",
	)
}

// Read fetches the current state from the upstream system.
func (r *netflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *netflowResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read from upstream
	netflow, err := r.client.Diagnostics().NetflowGetConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read netflow config, got error: %s", err))
		return
	}

	// Convert upstream struct to TF schema
	resourceModel, err := convertNetflowStructToSchema(netflow)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse netflow config, got error: %s", err))
		return
	}

	// Preserve the ID from state (always "netflow_settings")
	resourceModel.Id = data.Id

	tflog.Trace(ctx, "read diagnostics netflow resource")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}

// Update modifies the upstream singleton configuration.
func (r *netflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *netflowResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema to upstream struct
	resourceStruct, err := convertNetflowSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse netflow config, got error: %s", err))
		return
	}

	// Update upstream configuration
	_, err = r.client.Diagnostics().NetflowSetConfig(ctx, resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update netflow config, got error: %s", err))
		return
	}

	// Reconfigure the service to apply changes
	_, err = r.client.Diagnostics().NetflowReconfigure(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to reconfigure netflow, got error: %s", err))
		return
	}

	// Read back the updated config to ensure state consistency
	netflow, err := r.client.Diagnostics().NetflowGetConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read updated netflow config, got error: %s", err))
		return
	}

	// Convert back to schema
	resourceModel, err := convertNetflowStructToSchema(netflow)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse updated netflow config, got error: %s", err))
		return
	}

	// Preserve the ID
	resourceModel.Id = data.Id

	tflog.Trace(ctx, "updated diagnostics netflow resource")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}

// Delete removes the resource from Terraform state but does NOT modify upstream.
// This is the key behavior for singleton resources - they can't be destroyed.
func (r *netflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *netflowResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log a warning that the upstream configuration is not being deleted
	tflog.Warn(ctx,
		"Singleton resource removed from Terraform state. "+
			"The upstream Netflow configuration remains unchanged and will not be deleted. "+
			"To manage this resource again, re-import it with: "+
			"terraform import opnsense_diagnostics_netflow.<name> netflow_settings")

	// Add a warning to the user output
	resp.Diagnostics.AddWarning(
		"Singleton Resource Removed From State Only",
		"This resource has been removed from Terraform state, but the upstream "+
			"Netflow configuration has NOT been deleted or modified. The settings "+
			"remain active in the upstream system.\n\n"+
			"To manage this resource again in the future, re-import it:\n"+
			"  terraform import opnsense_diagnostics_netflow.<name> netflow_settings",
	)

	// Terraform automatically removes the resource from state
	// We do NOT make any API calls to delete/reset upstream configuration
}

// ImportState imports the singleton resource using the fixed ID "netflow_settings".
func (r *netflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Validate that the import ID is "netflow_settings"
	if req.ID != "netflow_settings" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"This is a singleton resource and must be imported using the ID 'netflow_settings'.\n\n"+
				"Usage:\n"+
				"  terraform import opnsense_diagnostics_netflow.<name> netflow_settings\n\n"+
				fmt.Sprintf("You provided: %q", req.ID),
		)
		return
	}

	// Set the ID attribute in state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)

	tflog.Info(ctx, "imported diagnostics netflow resource", map[string]any{"id": req.ID})
}
