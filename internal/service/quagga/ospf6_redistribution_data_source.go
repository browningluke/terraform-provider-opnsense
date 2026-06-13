package quagga

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ospf6RedistributionDataSource{}
var _ datasource.DataSourceWithConfigure = &ospf6RedistributionDataSource{}

func newOSPF6RedistributionDataSource() datasource.DataSource {
	return &ospf6RedistributionDataSource{}
}

// ospf6RedistributionDataSource defines the data source implementation.
type ospf6RedistributionDataSource struct {
	client opnsense.Client
}

func (d *ospf6RedistributionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quagga_ospf6_redistribution"
}

func (d *ospf6RedistributionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ospf6RedistributionDataSourceSchema()
}

func (d *ospf6RedistributionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ospf6RedistributionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ospf6RedistributionResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := d.client.Quagga().GetOSPF6Redistribution(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ospf6 redistribution, got error: %s", err))
		return
	}

	resourceModel, err := convertOSPF6RedistributionStructToSchema(resource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ospf6 redistribution, got error: %s", err))
		return
	}

	resourceModel.Id = data.Id

	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}
