package generic

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

type GenericResource[OpnSenseType any, TerraformSchema any] struct {
	Definition Definition[OpnSenseType, TerraformSchema]
	Client     opnsense.Client
}

func (r *GenericResource[OpnSenseType, TerraformSchema]) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s_%s", req.ProviderTypeName, r.Definition.ServiceName, r.Definition.TypeName)
}

func (r *GenericResource[OpnSenseType, TerraformSchema]) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.Definition.ResourceSchema
}

func (r *GenericResource[OpnSenseType, TerraformSchema]) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.Client = opnsense.NewClient(apiClient)
}

func (r *GenericResource[OpnSenseType, TerraformSchema]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.Definition.Create == nil {
		resp.Diagnostics.AddError("Invalid Resource Definition",
			fmt.Sprintf("Create operation is not defined for '%s_%s'.", r.Definition.ServiceName, r.Definition.TypeName),
		)
		return
	}

	var data *TerraformSchema

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	opnsenseData, err := r.Definition.FromModel(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse '%s' from '%s', got error: %s", r.Definition.TypeName, r.Definition.ServiceName, err))
		return
	}

	// Add resource to service
	id, err := r.Definition.Create(ctx, r.Client, opnsenseData)
	if err != nil {
		if id != "" {
			// Tag new resource with ID from OPNsense
			r.Definition.SetID(data, types.StringValue(id))

			// Read back so state captures API-normalised values (defaults,
			// sorting, trimming); fall back to plan-only state if the
			// read-back fails so the upstream resource isn't orphaned.
			if readStruct, readErr := r.Definition.Read(ctx, r.Client, id); readErr == nil {
				if readModel, convErr := r.Definition.ToModel(readStruct); convErr == nil {
					r.Definition.SetID(readModel, r.Definition.GetID(data))
					data = readModel
				}
			}

			// Save data into Terraform state
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create '%s' from '%s', got error: %s", r.Definition.TypeName, r.Definition.ServiceName, err))
		return
	}
	// Tag new resource with ID from OPNsense
	r.Definition.SetID(data, types.StringValue(id))

	if r.Definition.CreateHook != nil {
		if err := r.Definition.CreateHook(ctx, r.Client, data, id); err != nil {
			resp.Diagnostics.AddError("Create Hook Error", err.Error())
			return
		}
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GenericResource[OpnSenseType, TerraformSchema]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.Definition.Read == nil {
		resp.Diagnostics.AddError("Invalid Resource Definition",
			fmt.Sprintf("Read operation is not defined for '%s_%s'.", r.Definition.ServiceName, r.Definition.TypeName),
		)
		return
	}

	var data *TerraformSchema

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get resource from OPNsense API
	resourceStruct, err := r.Definition.Read(ctx, r.Client, r.Definition.GetID(data).ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, fmt.Sprintf("%s %s not present in remote, removing from state", r.Definition.ServiceName, r.Definition.TypeName))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read '%s' from '%s', got error: %s", r.Definition.TypeName, r.Definition.ServiceName, err))
		return
	}

	// Convert OPNsense struct to TF schema
	model, err := r.Definition.ToModel(resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read '%s' from '%s', got error: %s", r.Definition.TypeName, r.Definition.ServiceName, err))
		return
	}

	// ID cannot be added by convert... func, have to add here
	r.Definition.SetID(model, r.Definition.GetID(data))

	if r.Definition.ReadHook != nil {
		if err := r.Definition.ReadHook(ctx, r.Client, model); err != nil {
			resp.Diagnostics.AddError("Read Hook Error", err.Error())
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *GenericResource[OpnSenseType, TerraformSchema]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.Definition.Update == nil {
		resp.Diagnostics.AddError("Invalid Resource Definition",
			fmt.Sprintf("Update operation is not defined for '%s_%s'.", r.Definition.ServiceName, r.Definition.TypeName),
		)
		return
	}

	var data *TerraformSchema

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	resourceStruct, err := r.Definition.FromModel(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse '%s' from '%s', got error: %s", r.Definition.TypeName, r.Definition.ServiceName, err))
		return
	}

	// Update resource in service
	err = r.Definition.Update(ctx, r.Client, r.Definition.GetID(data).ValueString(), resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create '%s' from '%s', got error: %s", r.Definition.TypeName, r.Definition.ServiceName, err))
		return
	}

	if r.Definition.UpdateHook != nil {
		if err := r.Definition.UpdateHook(ctx, r.Client, data); err != nil {
			resp.Diagnostics.AddError("Update Hook Error", err.Error())
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GenericResource[OpnSenseType, TerraformSchema]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.Definition.Delete == nil {
		resp.Diagnostics.AddError("Invalid Resource Definition",
			fmt.Sprintf("Delete operation is not defined for '%s_%s'.", r.Definition.ServiceName, r.Definition.TypeName),
		)
		return
	}

	var data *TerraformSchema

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.Definition.Delete(ctx, r.Client, r.Definition.GetID(data).ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete '%s' from '%s', got error: %s", r.Definition.TypeName, r.Definition.ServiceName, err))
		return
	}

	if r.Definition.DeleteHook != nil {
		if err := r.Definition.DeleteHook(ctx, r.Client, data); err != nil {
			resp.Diagnostics.AddError("Delete Hook Error", err.Error())
			return
		}
	}
}

func (r *GenericResource[OpnSenseType, TerraformSchema]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.Definition.ImportState != nil {
		r.Definition.ImportState(ctx, req, resp)
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *GenericResource[OpnSenseType, TerraformSchema]) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if r.Definition.ConfigValidators == nil {
		return nil
	}
	return r.Definition.ConfigValidators(ctx)
}

func (r *GenericResource[OpnSenseType, TerraformSchema]) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	if r.Definition.UpgradeState == nil {
		return nil
	}
	return r.Definition.UpgradeState(ctx)
}
