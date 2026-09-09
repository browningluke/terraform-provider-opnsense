package generic

import (
	"context"
	"errors"
	"fmt"

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
	resp.TypeName = buildTerraformTypeName(req.ProviderTypeName, d.Definition)
}

func (d *GenericDataSource[OpnSenseType, TerraformSchema]) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = d.Definition.DataSourceSchema
}

func (d *GenericDataSource[OpnSenseType, TerraformSchema]) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, err := configureClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", err.Error())
		return
	}

	d.Client = client
}

func (d *GenericDataSource[OpnSenseType, TerraformSchema]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *TerraformSchema

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	model, err := readResource(ctx, d.Client, d.Definition, data)
	if err != nil {
		var notFoundError *errs.NotFoundError

		if errors.As(err, &notFoundError) {
			tflog.Warn(
				ctx,
				fmt.Sprintf("%s %s not present in remote", d.Definition.ServiceName, d.Definition.TypeName),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
