package unbound

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &settingsResource{}
var _ resource.ResourceWithConfigure = &settingsResource{}
var _ resource.ResourceWithImportState = &settingsResource{}

func newSettingsResource() resource.Resource {
	return &settingsResource{}
}

// settingsResource defines the resource implementation.
// This is a SINGLETON resource - it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type settingsResource struct {
	client opnsense.Client
}

func (r *settingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_unbound_settings"
}

func (r *settingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = settingsResourceSchema()
}

func (r *settingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *settingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *settingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// IMPORTANT: Block resource creation with a clear error message
	resp.Diagnostics.AddError(
		"Cannot Create Singleton Resource",
		"This resource manages existing upstream Unbound DNS configuration that cannot be created or destroyed.\n\n"+
			"To manage this resource, you must import it first:\n"+
			"  terraform import opnsense_unbound_settings.<name> unbound_settings\n\n"+
			"After importing, you can manage the configuration with terraform apply.",
	)
}

// Read fetches the current state from the upstream system.
func (r *settingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *settingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read from upstream using the settings endpoint
	settings, err := r.client.Unbound().SettingsGet(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read unbound settings, got error: %s", err))
		return
	}

	// Convert upstream struct to TF schema
	resourceModel, err := convertSettingsStructToSchema(settings)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse unbound settings, got error: %s", err))
		return
	}

	// Preserve the ID from state (always "unbound_settings")
	resourceModel.Id = data.Id

	tflog.Trace(ctx, "read unbound settings resource")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}

// Update modifies the upstream singleton configuration.
func (r *settingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *settingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema to upstream struct
	resourceStruct, err := convertSettingsSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse unbound settings, got error: %s", err))
		return
	}

	// Update upstream configuration
	_, err = r.client.Unbound().SettingsUpdate(ctx, resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update unbound settings, got error: %s", err))
		return
	}

	// Reconfigure the service to apply changes
	_, err = r.client.Unbound().SettingsReconfigure(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to reconfigure unbound service, got error: %s", err))
		return
	}

	// Check if general.* settings changed - if so, call ReconfigureGeneral
	var generalPath = path.Root("general")
	var planGeneral, stateGeneral types.Object

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, generalPath, &planGeneral)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, generalPath, &stateGeneral)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !planGeneral.Equal(stateGeneral) {
		tflog.Info(ctx, "General settings changed, calling ReconfigureGeneral")
		_, err = r.client.Unbound().SettingsReconfigureGeneral(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Client Error",
				fmt.Sprintf("Unable to reconfigure unbound general settings, got error: %s", err))
			return
		}
	}

	// Read back the updated settings to ensure state consistency
	settings, err := r.client.Unbound().SettingsGet(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read updated unbound settings, got error: %s", err))
		return
	}

	// Convert back to schema
	resourceModel, err := convertSettingsStructToSchema(settings)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse updated unbound settings, got error: %s", err))
		return
	}

	// Preserve the ID
	resourceModel.Id = data.Id

	tflog.Trace(ctx, "updated unbound settings resource")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}

// Delete removes the resource from Terraform state but does NOT modify upstream.
// This is the key behavior for singleton resources - they can't be destroyed.
func (r *settingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *settingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log a warning that the upstream configuration is not being deleted
	tflog.Warn(ctx,
		"Singleton resource removed from Terraform state. "+
			"The upstream Unbound DNS configuration remains unchanged and will not be deleted. "+
			"To manage this resource again, re-import it with: "+
			"terraform import opnsense_unbound_settings.<name> unbound_settings")

	// Add a warning to the user output
	resp.Diagnostics.AddWarning(
		"Singleton Resource Removed From State Only",
		"This resource has been removed from Terraform state, but the upstream "+
			"Unbound DNS configuration has NOT been deleted or modified. The settings "+
			"remain active in the upstream system.\n\n"+
			"To manage this resource again in the future, re-import it:\n"+
			"  terraform import opnsense_unbound_settings.<name> unbound_settings",
	)

	// Terraform automatically removes the resource from state
	// We do NOT make any API calls to delete/reset upstream configuration
}

// ImportState imports the singleton resource using the fixed ID "unbound_settings".
func (r *settingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Validate that the import ID is "unbound_settings"
	if req.ID != "unbound_settings" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"This is a singleton resource and must be imported using the ID 'unbound_settings'.\n\n"+
				"Usage:\n"+
				"  terraform import opnsense_unbound_settings.<name> unbound_settings\n\n"+
				fmt.Sprintf("You provided: %q", req.ID),
		)
		return
	}

	// Set the ID attribute in state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)

	tflog.Info(ctx, "imported unbound settings resource", map[string]any{"id": req.ID})
}
