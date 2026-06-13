package trust

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &settingsResource{}
var _ resource.ResourceWithConfigure = &settingsResource{}
var _ resource.ResourceWithImportState = &settingsResource{}

func newSettingsResource() resource.Resource {
	return &settingsResource{}
}

type settingsResource struct {
	client opnsense.Client
}

func (r *settingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trust_settings"
}

func (r *settingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = settingsResourceSchema()
}

func (r *settingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *settingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Cannot Create Singleton Resource",
		"This resource manages existing upstream Trust configuration that cannot be created or destroyed.\n\n"+
			"To manage this resource, you must import it first:\n"+
			"  terraform import opnsense_trust_settings.<name> trust_settings\n\n"+
			"After importing, you can manage the configuration with terraform apply.",
	)
}

func (r *settingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *settingsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.Trust().SettingsGet(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read trust settings, got error: %s", err))
		return
	}

	resourceModel, err := convertSettingsStructToSchema(&result.Trust)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse trust settings, got error: %s", err))
		return
	}

	resourceModel.Id = data.Id

	tflog.Trace(ctx, "read trust settings resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}

func (r *settingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *settingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceStruct, err := convertSettingsSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse trust settings, got error: %s", err))
		return
	}

	_, err = r.client.Trust().SettingsSet(ctx, resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update trust settings, got error: %s", err))
		return
	}

	_, err = r.client.Trust().SettingsReconfigure(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to reconfigure trust settings, got error: %s", err))
		return
	}

	result, err := r.client.Trust().SettingsGet(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read updated trust settings, got error: %s", err))
		return
	}

	resourceModel, err := convertSettingsStructToSchema(&result.Trust)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse updated trust settings, got error: %s", err))
		return
	}

	resourceModel.Id = data.Id

	tflog.Trace(ctx, "updated trust settings resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}

func (r *settingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *settingsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx,
		"Singleton resource removed from Terraform state. "+
			"The upstream Trust configuration remains unchanged and will not be deleted.")

	resp.Diagnostics.AddWarning(
		"Singleton Resource Removed From State Only",
		"This resource has been removed from Terraform state, but the upstream "+
			"Trust configuration has NOT been deleted or modified. The settings "+
			"remain active in the upstream system.\n\n"+
			"To manage this resource again in the future, re-import it:\n"+
			"  terraform import opnsense_trust_settings.<name> trust_settings",
	)
}

func (r *settingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "trust_settings" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"This is a singleton resource and must be imported using the ID 'trust_settings'.\n\n"+
				"Usage:\n"+
				"  terraform import opnsense_trust_settings.<name> trust_settings\n\n"+
				fmt.Sprintf("You provided: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)

	tflog.Info(ctx, "imported trust settings resource", map[string]any{"id": req.ID})
}
