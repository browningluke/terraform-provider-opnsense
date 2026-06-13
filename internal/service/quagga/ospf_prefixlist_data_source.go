package quagga

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ospfPrefixListDataSource{}
var _ datasource.DataSourceWithConfigure = &ospfPrefixListDataSource{}

func newOSPFPrefixListDataSource() datasource.DataSource {
	return &ospfPrefixListDataSource{}
}

// ospfPrefixListDataSource defines the data source implementation.
type ospfPrefixListDataSource struct {
	client opnsense.Client
}

func (d *ospfPrefixListDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quagga_ospf_prefix_list"
}

func (d *ospfPrefixListDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ospfPrefixListDataSourceSchema()
}

func (d *ospfPrefixListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ospfPrefixListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ospfPrefixListResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := d.client.Quagga().GetOSPFPrefixList(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ospf prefix list, got error: %s", err))
		return
	}

	resourceModel, err := convertOSPFPrefixListStructToSchema(resource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read ospf prefix list, got error: %s", err))
		return
	}

	resourceModel.Id = data.Id

	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}
