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
var _ resource.Resource = &bfdResource{}
var _ resource.ResourceWithConfigure = &bfdResource{}
var _ resource.ResourceWithImportState = &bfdResource{}

func newBFDResource() resource.Resource {
	return &bfdResource{}
}

// bfdResource defines the resource implementation.
// This is a SINGLETON resource - it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type bfdResource struct {
	client opnsense.Client
}

func (r *bfdResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quagga_bfd"
}

func (r *bfdResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = bfdResourceSchema()
}

func (r *bfdResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *bfdResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Cannot Create Singleton Resource",
		"This resource manages existing upstream Quagga BFD configuration that cannot be created or destroyed.\n\n"+
			"To manage this resource, you must import it first:\n"+
			"  terraform import opnsense_quagga_bfd.<name> quagga_bfd\n\n"+
			"After importing, you can manage the configuration with terraform apply.",
	)
}

func (r *bfdResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *bfdResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.Quagga().BFDGet(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read quagga BFD settings, got error: %s", err))
		return
	}

	resourceModel, err := convertBFDStructToSchema(&result.BFD)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse quagga BFD settings, got error: %s", err))
		return
	}

	resourceModel.Id = data.Id

	tflog.Trace(ctx, "read quagga BFD settings resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}

func (r *bfdResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *bfdResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceStruct, err := convertBFDSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse quagga BFD settings, got error: %s", err))
		return
	}

	_, err = r.client.Quagga().BFDSet(ctx, resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update quagga BFD settings, got error: %s", err))
		return
	}

	_, err = r.client.Quagga().ServiceReconfigure(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to reconfigure quagga after updating BFD settings, got error: %s", err))
		return
	}

	result, err := r.client.Quagga().BFDGet(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read updated quagga BFD settings, got error: %s", err))
		return
	}

	resourceModel, err := convertBFDStructToSchema(&result.BFD)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse updated quagga BFD settings, got error: %s", err))
		return
	}

	resourceModel.Id = data.Id

	tflog.Trace(ctx, "updated quagga BFD settings resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}

// Delete removes the resource from Terraform state but does NOT modify upstream.
func (r *bfdResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *bfdResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx,
		"Singleton resource removed from Terraform state. "+
			"The upstream Quagga BFD configuration remains unchanged and will not be deleted. "+
			"To manage this resource again, re-import it with: "+
			"terraform import opnsense_quagga_bfd.<name> quagga_bfd")

	resp.Diagnostics.AddWarning(
		"Singleton Resource Removed From State Only",
		"This resource has been removed from Terraform state, but the upstream "+
			"Quagga BFD configuration has NOT been deleted or modified. The settings "+
			"remain active in the upstream system.\n\n"+
			"To manage this resource again in the future, re-import it:\n"+
			"  terraform import opnsense_quagga_bfd.<name> quagga_bfd",
	)
}

// ImportState imports the singleton resource using the fixed ID "quagga_bfd".
func (r *bfdResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "quagga_bfd" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"This is a singleton resource and must be imported using the ID 'quagga_bfd'.\n\n"+
				"Usage:\n"+
				"  terraform import opnsense_quagga_bfd.<name> quagga_bfd\n\n"+
				fmt.Sprintf("You provided: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)

	tflog.Info(ctx, "imported quagga BFD settings resource", map[string]any{"id": req.ID})
}
