package mdnsrepeater

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/mdnsrepeater"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/browningluke/terraform-provider-opnsense/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const settingsId = "mdnsrepeater_settings"

type settingsResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	EnableCarp types.Bool   `tfsdk:"enable_carp"`
	Interfaces types.Set    `tfsdk:"interfaces"`
	Blocklist  types.Set    `tfsdk:"blocklist"`
}

func settingsResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages mDNS Repeater settings. Requires the `os-mdns-repeater` plugin. This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_mdnsrepeater_settings.settings " + settingsId + "\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `" + settingsId + "`. Use this value when importing.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether to enable the repeater. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"enable_carp": schema.BoolAttribute{
				MarkdownDescription: "When enabled, the repeater service only runs on the CARP master. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"interfaces": schema.SetAttribute{
				MarkdownDescription: "Listen interfaces (e.g. `lan`, `opt1`). At least 2 interfaces must be set when `enabled` is `true`. The daemon supports at most 5 interfaces.",
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 5),
				},
			},
			"blocklist": schema.SetAttribute{
				MarkdownDescription: "Up to 16 client IPv4 addresses or networks, in CIDR notation, to exclude. Defaults to `[]`.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
				Validators: []validator.Set{
					setvalidator.SizeAtMost(16),
					setvalidator.ValueStringsAre(validators.CIDR()),
				},
			},
		},
	}
}

func settingsDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Reads mDNS Repeater settings. Requires the `os-mdns-repeater` plugin.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `" + settingsId + "`.",
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether the repeater is enabled.",
				Computed:            true,
			},
			"enable_carp": dschema.BoolAttribute{
				MarkdownDescription: "Whether the repeater service only runs on the CARP master.",
				Computed:            true,
			},
			"interfaces": dschema.SetAttribute{
				MarkdownDescription: "Listen interfaces.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"blocklist": dschema.SetAttribute{
				MarkdownDescription: "Client IPv4 addresses or networks excluded from repeating.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func convertSettingsSchemaToStruct(d *settingsResourceModel) (*mdnsrepeater.Settings, error) {
	return &mdnsrepeater.Settings{
		Enabled:    tools.BoolToString(d.Enabled.ValueBool()),
		EnableCarp: tools.BoolToString(d.EnableCarp.ValueBool()),
		Interfaces: api.SelectedMapList(tools.SetToStringSlice(d.Interfaces)),
		Blocklist:  api.SelectedMapList(tools.SetToStringSlice(d.Blocklist)),
	}, nil
}

func convertSettingsStructToSchema(d *mdnsrepeater.Settings) (*settingsResourceModel, error) {
	return &settingsResourceModel{
		Enabled:    types.BoolValue(tools.StringToBool(d.Enabled)),
		EnableCarp: types.BoolValue(tools.StringToBool(d.EnableCarp)),
		Interfaces: tools.StringSliceToSet(d.Interfaces),
		Blocklist:  tools.StringSliceToSet(d.Blocklist),
	}, nil
}
