package firewall

import (
	"github.com/browningluke/opnsense-go/pkg/firewall"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type geoipSettingsResourceModel struct {
	Id  types.String `tfsdk:"id"`
	Url types.String `tfsdk:"url"`
}

func geoipSettingsResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages the GeoIP database source used by `geoip` type firewall aliases (Firewall > Aliases > GeoIP). " +
			"This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_firewall_geoip_settings.settings firewall_geoip_settings\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `firewall_geoip_settings`. Use this value when importing: `terraform import opnsense_firewall_geoip_settings.settings firewall_geoip_settings`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Download URL for the GeoIP country database (CSV, zipped). For MaxMind GeoLite2, include the license key in the URL query string, " +
					"e.g. `https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country-CSV&license_key=<key>&suffix=zip`. " +
					"Applying this resource triggers a database download on the OPNsense host.",
				Required: true,
			},
		},
	}
}

func geoipSettingsDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Reads the GeoIP database source configuration from the upstream system.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `firewall_geoip_settings`.",
			},
			"url": dschema.StringAttribute{
				MarkdownDescription: "Download URL configured for the GeoIP country database.",
				Computed:            true,
			},
		},
	}
}

func convertGeoIPSettingsSchemaToStruct(d *geoipSettingsResourceModel) (*firewall.GeoIPSettingsSetParams, error) {
	return &firewall.GeoIPSettingsSetParams{
		GeoIP: firewall.GeoIPSettingsUrl{
			Url: d.Url.ValueString(),
		},
	}, nil
}

func convertGeoIPSettingsStructToSchema(d *firewall.GeoIPSettings) (*geoipSettingsResourceModel, error) {
	return &geoipSettingsResourceModel{
		Url: types.StringValue(d.Url),
	}, nil
}
