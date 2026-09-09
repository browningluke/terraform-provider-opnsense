package wireguard

import (
	"context"

	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/browningluke/opnsense-go/pkg/wireguard"
	"github.com/browningluke/terraform-provider-opnsense/internal/provider/generic"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var serviceName = "wireguard"

var settingsDefinition = generic.Definition[
	wireguard.WireguardGeneral,
	settingsResourceModel,
]{
	ServiceName: serviceName,
	TypeName:    "settings",

	ResourceSchema:   settingsResourceSchema(),
	DataSourceSchema: settingsDataSourceSchema(),

	FromModel: convertSettingsSchemaToStruct,
	ToModel:   convertSettingsStructToSchema,

	GetID: func(m *settingsResourceModel) types.String {
		return types.StringValue(serviceName + "_settings")
	},
	SetID: func(m *settingsResourceModel, id types.String) {
		m.Id = id
	},

	Read: func(ctx context.Context, client opnsense.Client, _ string) (*wireguard.WireguardGeneral, error) {
		result, err := client.Wireguard().GeneralGet(ctx)
		if err != nil {
			return nil, err
		}

		return &result.General, nil
	},

	Update: func(ctx context.Context, client opnsense.Client, _ string, data *wireguard.WireguardGeneral) error {
		_, err := client.Wireguard().GeneralSet(ctx, data)
		return err
	},
}

// settingsResourceModel describes the resource data model.
// This is a SINGLETON resource — it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type settingsResourceModel struct {
	Id      types.String `tfsdk:"id"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func settingsResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages WireGuard general settings. This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_wireguard_settings.settings wireguard_settings\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `wireguard_settings`. Use this value when importing: `terraform import opnsense_wireguard_settings.settings wireguard_settings`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "When enabled, the WireGuard daemon is active. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func settingsDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Reads WireGuard general settings from the upstream system.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `wireguard_settings`.",
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether the WireGuard daemon is enabled.",
				Computed:            true,
			},
		},
	}
}

func convertSettingsSchemaToStruct(d *settingsResourceModel) (*wireguard.WireguardGeneral, error) {
	return &wireguard.WireguardGeneral{
		Enabled: tools.BoolToString(d.Enabled.ValueBool()),
	}, nil
}

func convertSettingsStructToSchema(d *wireguard.WireguardGeneral) (*settingsResourceModel, error) {
	return &settingsResourceModel{
		Enabled: types.BoolValue(tools.StringToBool(d.Enabled)),
	}, nil
}
