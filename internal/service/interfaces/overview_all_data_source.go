package interfaces

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &overviewAllDataSource{}
var _ datasource.DataSourceWithConfigure = &overviewAllDataSource{}

func newOverviewAllDataSource() datasource.DataSource {
	return &overviewAllDataSource{}
}

// overviewAllDataSource defines the data source implementation.
type overviewAllDataSource struct {
	client opnsense.Client
}

type overviewAllDataSourceModel struct {
	Interfaces types.List `tfsdk:"interfaces"`
}

func (d *overviewAllDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interfaces_overview_all"
}

func (d *overviewAllDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Interfaces Overview All provides live state for all OPNsense interfaces, including IP addresses, CARP status, VLAN, and LAGG information.",

		Attributes: map[string]schema.Attribute{
			"interfaces": schema.ListNestedAttribute{
				MarkdownDescription: "A list of all interfaces present in OPNsense.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: overviewInterfaceDataSourceSchema().Attributes,
				},
			},
		},
	}
}

func (d *overviewAllDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = opnsense.NewClient(apiClient)
}

func (d *overviewAllDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *overviewAllDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get all interfaces from OPNsense API
	result, err := d.client.Interfaces().OverviewGet(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read interfaces overview, got error: %s", err))
		return
	}

	var ifaces []overviewInterfaceModel
	for i := range result.Rows {
		model, err := convertOverviewInterfaceStructToSchema(&result.Rows[i])
		if err != nil {
			resp.Diagnostics.AddError("Client Error",
				fmt.Sprintf("Unable to convert interface overview, got error: %s", err))
			return
		}
		ifaces = append(ifaces, *model)
	}

	// Create empty list first
	v, _ := types.ListValue(
		types.ObjectType{AttrTypes: overviewInterfaceAttrTypes},
		[]attr.Value{},
	)
	// Try to fill list
	if len(ifaces) > 0 {
		v, _ = types.ListValueFrom(
			context.Background(),
			types.ObjectType{AttrTypes: overviewInterfaceAttrTypes},
			ifaces,
		)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &overviewAllDataSourceModel{Interfaces: v})...)
}
