package generic

import (
	"context"
	"errors"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/errs"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type GenericSingletonResource[OpnSenseType any, TerraformSchema any] struct {
	Definition       Definition[OpnSenseType, TerraformSchema]
	Client           opnsense.Client
	ProviderTypeName string
}

func (r *GenericSingletonResource[OpnSenseType, TerraformSchema]) buildId() string {
	return r.Definition.GetID(nil).ValueString()
}

func (r *GenericSingletonResource[OpnSenseType, TerraformSchema]) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	resp.TypeName = buildTerraformTypeName(req.ProviderTypeName, r.Definition)
}

func (r *GenericSingletonResource[OpnSenseType, TerraformSchema]) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.Definition.ResourceSchema
}

func (r *GenericSingletonResource[OpnSenseType, TerraformSchema]) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, err := configureClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", err.Error())
		return
	}

	r.Client = client
}

func (r *GenericSingletonResource[OpnSenseType, TerraformSchema]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Cannot Create Singleton Resource",
		fmt.Sprintf(
			"This resource manages existing upstream Trust configuration that cannot be created or destroyed.\n\n"+
				"To manage this resource, you must import it first:\n"+
				"  terraform import %s.<name> %s\n\n"+
				"After importing, you can manage the configuration with terraform apply.",
			buildTerraformTypeName(r.ProviderTypeName, r.Definition),
			r.buildId(),
		),
	)
}

func (r *GenericSingletonResource[OpnSenseType, TerraformSchema]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *TerraformSchema

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := readResource(ctx, r.Client, r.Definition, data)
	if err != nil {
		var notFoundError *errs.NotFoundError

		if errors.As(err, &notFoundError) {
			tflog.Warn(
				ctx,
				fmt.Sprintf("%s %s not present in remote, removing from state", r.Definition.ServiceName, r.Definition.TypeName),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *GenericSingletonResource[OpnSenseType, TerraformSchema]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.Definition.Update == nil {
		resp.Diagnostics.AddError("Invalid Resource Definition",
			fmt.Sprintf("Update operation is not defined for '%s %s'.", r.Definition.ServiceName, r.Definition.TypeName),
		)
		return
	}

	var data *TerraformSchema

	// 1. Diagnostic
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 2. Convert
	resourceStruct, err := r.Definition.FromModel(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse '%s' from '%s', got error: %s", r.Definition.TypeName, r.Definition.ServiceName, err))
		return
	}

	// 3. Update on Server
	err = r.Definition.Update(ctx, r.Client, r.buildId(), resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create '%s' from '%s', got error: %s", r.Definition.TypeName, r.Definition.ServiceName, err))
		return
	}

	// Optional: 3.a Reconfigure
	if r.Definition.Reconfigure != nil {
		_, err = r.Definition.Reconfigure(ctx, r.Client)
		if err != nil {
			resp.Diagnostics.AddError("Client Error",
				fmt.Sprintf("Unable to reconfigure %s settings, got error: %s", r.Definition.ServiceName, err))
			return
		}
	}

	// 4. Load Settings again
	resourceStruct, err = r.Definition.Read(ctx, r.Client, r.buildId())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create '%s' from '%s', got error: %s", r.Definition.TypeName, r.Definition.ServiceName, err))
		return
	}
	resourceModel, err := r.Definition.ToModel(resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse '%s' from '%s', got error: %s", r.Definition.TypeName, r.Definition.ServiceName, err))
		return
	}
	r.Definition.SetID(resourceModel, types.StringValue(r.buildId()))

	tflog.Trace(ctx, fmt.Sprintf("updated %s settings resource", r.Definition.ServiceName))

	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}

func (r *GenericSingletonResource[OpnSenseType, TerraformSchema]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TerraformSchema

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx,
		"Singleton resource removed from Terraform state. "+
			"The upstream Trust configuration remains unchanged and will not be deleted.")

	resp.Diagnostics.AddWarning(
		"Singleton Resource Removed From State Only",
		fmt.Sprintf(
			"This resource has been removed from Terraform state, but the upstream "+
				"Trust configuration has NOT been deleted or modified. The settings "+
				"remain active in the upstream system.\n\n"+
				"To manage this resource again in the future, re-import it:\n"+
				"  terraform import %s.<name> %s",
			buildTerraformTypeName(r.ProviderTypeName, r.Definition),
			r.buildId(),
		),
	)
}

func (r *GenericSingletonResource[OpnSenseType, TerraformSchema]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != r.buildId() {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf(
				"This is a singleton resource and must be imported using the ID '%s'.\n\n"+
					"Usage:\n"+
					"  terraform import %s.<name> %s\n\n"+
					"You provided: %q",
				r.buildId(),
				buildTerraformTypeName(r.ProviderTypeName, r.Definition),
				r.buildId(),
				req.ID,
			),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)

	tflog.Info(ctx, fmt.Sprintf("imported %s settings resource", r.Definition.ServiceName), map[string]any{"id": req.ID})
}
