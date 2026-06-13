package quagga

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// generalResourceModel describes the resource data model.
// This is a SINGLETON resource — it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type generalResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	Profile      types.String `tfsdk:"profile"`
	EnableCarp   types.Bool   `tfsdk:"enable_carp"`
	EnableSyslog types.Bool   `tfsdk:"enable_syslog"`
	EnableSNMP   types.Bool   `tfsdk:"enable_snmp"`
	SyslogLevel  types.String `tfsdk:"syslog_level"`
	FWRules      types.Bool   `tfsdk:"fw_rules"`
}

func generalResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages Quagga general settings. This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_quagga_general.general quagga_general\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `quagga_general`. Use this value when importing: `terraform import opnsense_quagga_general.general quagga_general`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable the Quagga routing daemon. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"profile": schema.StringAttribute{
				MarkdownDescription: "Select the Quagga profile. `traditional` follows the original FRRouting defaults; `datacenter` enables a more aggressive timers and additional features suitable for datacenter deployments. Defaults to `\"traditional\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("traditional"),
				Validators: []validator.String{
					stringvalidator.OneOf("traditional", "datacenter"),
				},
			},
			"enable_carp": schema.BoolAttribute{
				MarkdownDescription: "Enable CARP support. When enabled, Quagga adjusts its behaviour based on the CARP master/backup state. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"enable_syslog": schema.BoolAttribute{
				MarkdownDescription: "Send Quagga log output to syslog. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"enable_snmp": schema.BoolAttribute{
				MarkdownDescription: "Enable SNMP support via the agentx protocol. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"syslog_level": schema.StringAttribute{
				MarkdownDescription: "The syslog verbosity level for Quagga log messages. Defaults to `\"notifications\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("notifications"),
				Validators: []validator.String{
					stringvalidator.OneOf(
						"critical",
						"emergencies",
						"errors",
						"alerts",
						"warnings",
						"notifications",
						"informational",
						"debugging",
					),
				},
			},
			"fw_rules": schema.BoolAttribute{
				MarkdownDescription: "Automatically install firewall rules to allow routing protocol traffic (BGP, OSPF, etc.). Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
		},
	}
}

func convertGeneralSchemaToStruct(d *generalResourceModel) (*quagga.QuaggaGeneral, error) {
	return &quagga.QuaggaGeneral{
		Enabled:      tools.BoolToString(d.Enabled.ValueBool()),
		Profile:      api.SelectedMap(d.Profile.ValueString()),
		EnableCarp:   tools.BoolToString(d.EnableCarp.ValueBool()),
		EnableSyslog: tools.BoolToString(d.EnableSyslog.ValueBool()),
		EnableSNMP:   tools.BoolToString(d.EnableSNMP.ValueBool()),
		SyslogLevel:  api.SelectedMap(d.SyslogLevel.ValueString()),
		FWRules:      tools.BoolToString(d.FWRules.ValueBool()),
	}, nil
}

func convertGeneralStructToSchema(d *quagga.QuaggaGeneral) (*generalResourceModel, error) {
	return &generalResourceModel{
		Enabled:      types.BoolValue(tools.StringToBool(d.Enabled)),
		Profile:      types.StringValue(d.Profile.String()),
		EnableCarp:   types.BoolValue(tools.StringToBool(d.EnableCarp)),
		EnableSyslog: types.BoolValue(tools.StringToBool(d.EnableSyslog)),
		EnableSNMP:   types.BoolValue(tools.StringToBool(d.EnableSNMP)),
		SyslogLevel:  types.StringValue(d.SyslogLevel.String()),
		FWRules:      types.BoolValue(tools.StringToBool(d.FWRules)),
	}, nil
}
