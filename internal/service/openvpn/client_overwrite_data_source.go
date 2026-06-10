package openvpn

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &clientOverwriteDataSource{}
var _ datasource.DataSourceWithConfigure = &clientOverwriteDataSource{}

func newClientOverwriteDataSource() datasource.DataSource {
	return &clientOverwriteDataSource{}
}

type clientOverwriteDataSource struct {
	client opnsense.Client
}

func (d *clientOverwriteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_openvpn_client_overwrite"
}

func (d *clientOverwriteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = clientOverwriteDataSourceSchema()
}

func (d *clientOverwriteDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = opnsense.NewClient(apiClient)
}

func (d *clientOverwriteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *clientOverwriteResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.Openvpn().GetClientOverwrite(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read openvpn client overwrite, got error: %s", err))
		return
	}

	model, err := convertClientOverwriteStructToSchema(res)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read openvpn client overwrite, got error: %s", err))
		return
	}
	model.Id = data.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
