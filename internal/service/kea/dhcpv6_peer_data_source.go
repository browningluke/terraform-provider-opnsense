package kea

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &dhcpv6PeerDataSource{}
var _ datasource.DataSourceWithConfigure = &dhcpv6PeerDataSource{}

func newDhcpv6PeerDataSource() datasource.DataSource {
	return &dhcpv6PeerDataSource{}
}

// dhcpv6PeerDataSource defines the data source implementation.
type dhcpv6PeerDataSource struct {
	client opnsense.Client
}

func (d *dhcpv6PeerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kea_dhcpv6_peer"
}

func (d *dhcpv6PeerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dhcpv6PeerDataSourceSchema()
}

func (d *dhcpv6PeerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dhcpv6PeerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *dhcpv6PeerResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceStruct, err := d.client.Kea().GetPeerV6(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read kea peer, got error: %s", err))
		return
	}

	resourceModel, err := convertDhcpv6PeerStructToSchema(resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read kea peer, got error: %s", err))
		return
	}

	resourceModel.Id = data.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}
