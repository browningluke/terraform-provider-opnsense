package quagga

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
var _ resource.Resource = &ospfResource{}
var _ resource.ResourceWithConfigure = &ospfResource{}
var _ resource.ResourceWithImportState = &ospfResource{}

func newOSPFResource() resource.Resource {
	return &ospfResource{}
}

// ospfResource defines the resource implementation.
// This is a SINGLETON resource - it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type ospfResource struct {
	client opnsense.Client
}

func (r *ospfResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quagga_ospf"
}

func (r *ospfResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ospfResourceSchema()
}

func (r *ospfResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ospfResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Cannot Create Singleton Resource",
		"This resource manages existing upstream Quagga OSPF configuration that cannot be created or destroyed.\n\n"+
			"To manage this resource, you must import it first:\n"+
			"  terraform import opnsense_quagga_ospf.<name> quagga_ospf\n\n"+
			"After importing, you can manage the configuration with terraform apply.",
	)
}

func (r *ospfResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ospfResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.Quagga().OSPFGet(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read quagga OSPF settings, got error: %s", err))
		return
	}

	resourceModel, err := convertOSPFStructToSchema(&result.OSPF)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse quagga OSPF settings, got error: %s", err))
		return
	}

	resourceModel.Id = data.Id

	tflog.Trace(ctx, "read quagga OSPF settings resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}

func (r *ospfResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ospfResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceStruct, err := convertOSPFSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse quagga OSPF settings, got error: %s", err))
		return
	}

	_, err = r.client.Quagga().OSPFSet(ctx, resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update quagga OSPF settings, got error: %s", err))
		return
	}

	_, err = r.client.Quagga().ServiceReconfigure(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to reconfigure quagga after updating OSPF settings, got error: %s", err))
		return
	}

	result, err := r.client.Quagga().OSPFGet(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read updated quagga OSPF settings, got error: %s", err))
		return
	}

	resourceModel, err := convertOSPFStructToSchema(&result.OSPF)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse updated quagga OSPF settings, got error: %s", err))
		return
	}

	resourceModel.Id = data.Id

	tflog.Trace(ctx, "updated quagga OSPF settings resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}

// Delete removes the resource from Terraform state but does NOT modify upstream.
func (r *ospfResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ospfResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx,
		"Singleton resource removed from Terraform state. "+
			"The upstream Quagga OSPF configuration remains unchanged and will not be deleted. "+
			"To manage this resource again, re-import it with: "+
			"terraform import opnsense_quagga_ospf.<name> quagga_ospf")

	resp.Diagnostics.AddWarning(
		"Singleton Resource Removed From State Only",
		"This resource has been removed from Terraform state, but the upstream "+
			"Quagga OSPF configuration has NOT been deleted or modified. The settings "+
			"remain active in the upstream system.\n\n"+
			"To manage this resource again in the future, re-import it:\n"+
			"  terraform import opnsense_quagga_ospf.<name> quagga_ospf",
	)
}

// ImportState imports the singleton resource using the fixed ID "quagga_ospf".
func (r *ospfResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "quagga_ospf" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"This is a singleton resource and must be imported using the ID 'quagga_ospf'.\n\n"+
				"Usage:\n"+
				"  terraform import opnsense_quagga_ospf.<name> quagga_ospf\n\n"+
				fmt.Sprintf("You provided: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)

	tflog.Info(ctx, "imported quagga OSPF settings resource", map[string]any{"id": req.ID})
}
