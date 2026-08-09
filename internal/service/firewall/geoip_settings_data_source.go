package firewall

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &geoipSettingsDataSource{}
var _ datasource.DataSourceWithConfigure = &geoipSettingsDataSource{}

func newGeoIPSettingsDataSource() datasource.DataSource {
	return &geoipSettingsDataSource{}
}

type geoipSettingsDataSource struct {
	client opnsense.Client
}

func (d *geoipSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_geoip_settings"
}

func (d *geoipSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = geoipSettingsDataSourceSchema()
}

func (d *geoipSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = opnsense.NewClient(apiClient)
}

func (d *geoipSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *geoipSettingsResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.Firewall().GeoIPSettingsGet(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read firewall geoip settings, got error: %s", err))
		return
	}

	resourceModel, err := convertGeoIPSettingsStructToSchema(&result.Alias.GeoIP)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse firewall geoip settings, got error: %s", err))
		return
	}

	resourceModel.Id = types.StringValue("firewall_geoip_settings")

	tflog.Trace(ctx, "read firewall geoip settings data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}
