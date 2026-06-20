package interfaces

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &overviewInterfaceDataSource{}
var _ datasource.DataSourceWithConfigure = &overviewInterfaceDataSource{}


// overviewInterfaceDataSource defines the data source implementation.
type overviewInterfaceDataSource struct {
	client opnsense.Client
}

func (d *overviewInterfaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interfaces_overview"
}

func (d *overviewInterfaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = overviewInterfaceDataSourceSchema()
}

func (d *overviewInterfaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *overviewInterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *overviewInterfaceModel

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

	// Find the interface matching the requested device name
	device := data.Device.ValueString()
	for i := range result.Rows {
		if result.Rows[i].Device == device {
			model, err := convertOverviewInterfaceStructToSchema(&result.Rows[i])
			if err != nil {
				resp.Diagnostics.AddError("Client Error",
					fmt.Sprintf("Unable to convert interface overview, got error: %s", err))
				return
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
			return
		}
	}

	resp.Diagnostics.AddError("Not Found",
		fmt.Sprintf("Interface with device %q not found in interfaces overview", device))
}
