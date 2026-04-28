package kea

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &dhcpv6SubnetDataSource{}
var _ datasource.DataSourceWithConfigure = &dhcpv6SubnetDataSource{}

func newDhcpv6SubnetDataSource() datasource.DataSource {
	return &dhcpv6SubnetDataSource{}
}

// dhcpv6SubnetDataSource defines the data source implementation.
type dhcpv6SubnetDataSource struct {
	client opnsense.Client
}

func (d *dhcpv6SubnetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kea_dhcpv6_subnet"
}

func (d *dhcpv6SubnetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dhcpv6SubnetDataSourceSchema()
}

func (d *dhcpv6SubnetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dhcpv6SubnetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *dhcpv6SubnetResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceStruct, err := d.client.Kea().GetSubnetV6(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read kea subnet, got error: %s", err))
		return
	}

	resourceModel, err := convertDhcpv6SubnetStructToSchema(resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read kea subnet, got error: %s", err))
		return
	}

	resourceModel.Id = data.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}
