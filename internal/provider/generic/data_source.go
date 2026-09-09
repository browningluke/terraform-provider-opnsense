package generic

import (
	"context"
	"errors"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/errs"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type GenericDataSource[OpnSenseType any, TerraformSchema any] struct {
	Definition Definition[OpnSenseType, TerraformSchema]
	Client     opnsense.Client
}

func (d *GenericDataSource[OpnSenseType, TerraformSchema]) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s_%s", req.ProviderTypeName, d.Definition.ServiceName, d.Definition.TypeName)
}

func (d *GenericDataSource[OpnSenseType, TerraformSchema]) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = d.Definition.DataSourceSchema
}

func (d *GenericDataSource[OpnSenseType, TerraformSchema]) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.Client = opnsense.NewClient(apiClient)
}

func (d *GenericDataSource[OpnSenseType, TerraformSchema]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.Definition.Read == nil {
		resp.Diagnostics.AddError("Invalid Resource Definition",
			fmt.Sprintf("Read operation is not defined for '%s_%s'.", d.Definition.ServiceName, d.Definition.TypeName),
		)
		return
	}

	var data *TerraformSchema

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get firewall alias from OPNsense unbound API
	resourceStruct, err := d.Definition.Read(ctx, d.Client, d.Definition.GetID(data).ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, fmt.Sprintf("%s %s not present in remote, removing from state", d.Definition.ServiceName, d.Definition.TypeName))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read '%s' from '%s', got error: %s", d.Definition.TypeName, d.Definition.ServiceName, err))
		return
	}

	// Convert OPNsense struct to TF schema
	resourceModel, err := d.Definition.ToModel(resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read '%s' from '%s', got error: %s", d.Definition.TypeName, d.Definition.ServiceName, err))
		return
	}

	// ID cannot be added by convert... func, have to add here
	d.Definition.SetID(resourceModel, d.Definition.GetID(data))

	if d.Definition.ReadHook != nil {
		if err := d.Definition.ReadHook(ctx, d.Client, resourceModel); err != nil {
			resp.Diagnostics.AddError("Read Hook Error", err.Error())
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}
