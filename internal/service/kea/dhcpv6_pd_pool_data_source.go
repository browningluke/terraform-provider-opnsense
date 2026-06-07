package kea

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &dhcpv6PdPoolDataSource{}
var _ datasource.DataSourceWithConfigure = &dhcpv6PdPoolDataSource{}

func newDhcpv6PdPoolDataSource() datasource.DataSource {
	return &dhcpv6PdPoolDataSource{}
}

// dhcpv6PdPoolDataSource defines the data source implementation.
type dhcpv6PdPoolDataSource struct {
	client opnsense.Client
}

func (d *dhcpv6PdPoolDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kea_dhcpv6_pd_pool"
}

func (d *dhcpv6PdPoolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dhcpv6PdPoolDataSourceSchema()
}

func (d *dhcpv6PdPoolDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dhcpv6PdPoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *dhcpv6PdPoolResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceStruct, err := d.client.Kea().GetPDPool(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read kea pd_pool, got error: %s", err))
		return
	}

	resourceModel, err := convertDhcpv6PdPoolStructToSchema(resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read kea pd_pool, got error: %s", err))
		return
	}

	resourceModel.Id = data.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}
