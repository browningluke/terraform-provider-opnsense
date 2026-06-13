package quagga

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ospfRouteMapDataSource{}
var _ datasource.DataSourceWithConfigure = &ospfRouteMapDataSource{}

func newOSPFRouteMapDataSource() datasource.DataSource {
	return &ospfRouteMapDataSource{}
}

// ospfRouteMapDataSource defines the data source implementation.
type ospfRouteMapDataSource struct {
	client opnsense.Client
}

func (d *ospfRouteMapDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quagga_ospf_route_map"
}

func (d *ospfRouteMapDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ospfRouteMapDataSourceSchema()
}

func (d *ospfRouteMapDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ospfRouteMapDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ospfRouteMapResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := d.client.Quagga().GetOSPFRouteMap(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ospf route map, got error: %s", err))
		return
	}

	resourceModel, err := convertOSPFRouteMapStructToSchema(resource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ospf route map, got error: %s", err))
		return
	}

	resourceModel.Id = data.Id

	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}
