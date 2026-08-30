package crowdsec

import (
	"github.com/browningluke/opnsense-go/pkg/crowdsec"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// settingsResourceModel describes the resource data model.
// This is a SINGLETON resource — it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type settingsResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	AgentEnabled            types.Bool   `tfsdk:"agent_enabled"`
	LapiEnabled             types.Bool   `tfsdk:"lapi_enabled"`
	FirewallBouncerEnabled  types.Bool   `tfsdk:"firewall_bouncer_enabled"`
	LapiManualConfiguration types.Bool   `tfsdk:"lapi_manual_configuration"`
	LapiListenAddress       types.String `tfsdk:"lapi_listen_address"`
	LapiListenPort          types.String `tfsdk:"lapi_listen_port"`
	RulesEnabled            types.Bool   `tfsdk:"rules_enabled"`
	RulesLog                types.Bool   `tfsdk:"rules_log"`
	RulesTag                types.String `tfsdk:"rules_tag"`
	EnrollKey               types.String `tfsdk:"enroll_key"`
	CrowdsecFirewallVerbose types.Bool   `tfsdk:"crowdsec_firewall_verbose"`
}

func settingsResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages CrowdSec plugin general settings. This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_crowdsec_settings.settings crowdsec_settings\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `crowdsec_settings`. Use this value when importing: `terraform import opnsense_crowdsec_settings.settings crowdsec_settings`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"agent_enabled": schema.BoolAttribute{
				MarkdownDescription: "When enabled, the CrowdSec agent (log processor / IDS) is active. Keep this enabled to detect attacks and receive alerts from the CrowdSec central service. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"lapi_enabled": schema.BoolAttribute{
				MarkdownDescription: "When enabled, the CrowdSec Local API (LAPI) is active. Keep this enabled unless you connect to a LAPI on another machine. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"firewall_bouncer_enabled": schema.BoolAttribute{
				MarkdownDescription: "When enabled, the remediation component (IPS / firewall bouncer) is active. Keep this enabled to block packets from attacking IP addresses. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"lapi_manual_configuration": schema.BoolAttribute{
				MarkdownDescription: "When enabled, avoids overwriting LAPI settings for `config.yaml`, `local_api_credentials.yaml`, and `crowdsec-firewall-bouncer.yaml`. The `lapi_listen_address` and `lapi_listen_port` attributes are ignored when this is enabled. Allows unsupported configurations such as linking together multiple OPNsense instances or connecting to an existing CrowdSec multi-server setup. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"lapi_listen_address": schema.StringAttribute{
				MarkdownDescription: "Where to listen for LAPI connections (IP address). Change it to a LAN address to connect from other agents/machines and bouncers. Ignored when `lapi_manual_configuration` is enabled. Defaults to `\"127.0.0.1\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("127.0.0.1"),
			},
			"lapi_listen_port": schema.StringAttribute{
				MarkdownDescription: "Where to listen for LAPI connections (port). Change it to avoid conflicts with existing services. Ignored when `lapi_manual_configuration` is enabled. Defaults to `\"8080\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("8080"),
			},
			"rules_enabled": schema.BoolAttribute{
				MarkdownDescription: "When enabled, block rules are generated from the CrowdSec blocklists and applied to all interfaces, IPv4/v6, ingress and egress. If disabled, you'll have to write your own rules to block anything. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"rules_log": schema.BoolAttribute{
				MarkdownDescription: "When enabled, log collection is enabled for CrowdSec's block rules. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"rules_tag": schema.StringAttribute{
				MarkdownDescription: "Tag added to packets dropped by CrowdSec rules, for diagnostic purposes. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"enroll_key": schema.StringAttribute{
				MarkdownDescription: "Enrollment key from https://app.crowdsec.net. Click \"Enroll command\" on the website and copy the key here. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Sensitive:           true,
			},
			"crowdsec_firewall_verbose": schema.BoolAttribute{
				MarkdownDescription: "When enabled, verbose logging is written to `/var/log/crowdsec/crowdsec-firewall-bouncer.log`. Enable this for debugging. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func settingsDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Reads CrowdSec plugin general settings from the upstream system.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `crowdsec_settings`.",
			},
			"agent_enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether the CrowdSec agent (log processor / IDS) is enabled.",
				Computed:            true,
			},
			"lapi_enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether the CrowdSec Local API (LAPI) is enabled.",
				Computed:            true,
			},
			"firewall_bouncer_enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether the remediation component (IPS / firewall bouncer) is enabled.",
				Computed:            true,
			},
			"lapi_manual_configuration": dschema.BoolAttribute{
				MarkdownDescription: "Whether manual LAPI configuration is enabled, avoiding overwrites of config.yaml, local_api_credentials.yaml, and crowdsec-firewall-bouncer.yaml.",
				Computed:            true,
			},
			"lapi_listen_address": dschema.StringAttribute{
				MarkdownDescription: "IP address the LAPI listens on.",
				Computed:            true,
			},
			"lapi_listen_port": dschema.StringAttribute{
				MarkdownDescription: "Port the LAPI listens on.",
				Computed:            true,
			},
			"rules_enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether block rules are generated from the CrowdSec blocklists.",
				Computed:            true,
			},
			"rules_log": dschema.BoolAttribute{
				MarkdownDescription: "Whether log collection is enabled for CrowdSec's block rules.",
				Computed:            true,
			},
			"rules_tag": dschema.StringAttribute{
				MarkdownDescription: "Tag added to packets dropped by CrowdSec rules.",
				Computed:            true,
			},
			"enroll_key": dschema.StringAttribute{
				MarkdownDescription: "Enrollment key from https://app.crowdsec.net.",
				Computed:            true,
				Sensitive:           true,
			},
			"crowdsec_firewall_verbose": dschema.BoolAttribute{
				MarkdownDescription: "Whether verbose logging is written to the firewall bouncer log.",
				Computed:            true,
			},
		},
	}
}

func convertSettingsSchemaToStruct(d *settingsResourceModel) (*crowdsec.CrowdsecGeneral, error) {
	return &crowdsec.CrowdsecGeneral{
		AgentEnabled:            tools.BoolToString(d.AgentEnabled.ValueBool()),
		LapiEnabled:             tools.BoolToString(d.LapiEnabled.ValueBool()),
		FirewallBouncerEnabled:  tools.BoolToString(d.FirewallBouncerEnabled.ValueBool()),
		LapiManualConfiguration: tools.BoolToString(d.LapiManualConfiguration.ValueBool()),
		LapiListenAddress:       d.LapiListenAddress.ValueString(),
		LapiListenPort:          d.LapiListenPort.ValueString(),
		RulesEnabled:            tools.BoolToString(d.RulesEnabled.ValueBool()),
		RulesLog:                tools.BoolToString(d.RulesLog.ValueBool()),
		RulesTag:                d.RulesTag.ValueString(),
		EnrollKey:               d.EnrollKey.ValueString(),
		CrowdsecFirewallVerbose: tools.BoolToString(d.CrowdsecFirewallVerbose.ValueBool()),
	}, nil
}

func convertSettingsStructToSchema(d *crowdsec.CrowdsecGeneral) (*settingsResourceModel, error) {
	return &settingsResourceModel{
		AgentEnabled:            types.BoolValue(tools.StringToBool(d.AgentEnabled)),
		LapiEnabled:             types.BoolValue(tools.StringToBool(d.LapiEnabled)),
		FirewallBouncerEnabled:  types.BoolValue(tools.StringToBool(d.FirewallBouncerEnabled)),
		LapiManualConfiguration: types.BoolValue(tools.StringToBool(d.LapiManualConfiguration)),
		LapiListenAddress:       types.StringValue(d.LapiListenAddress),
		LapiListenPort:          types.StringValue(d.LapiListenPort),
		RulesEnabled:            types.BoolValue(tools.StringToBool(d.RulesEnabled)),
		RulesLog:                types.BoolValue(tools.StringToBool(d.RulesLog)),
		RulesTag:                types.StringValue(d.RulesTag),
		EnrollKey:               types.StringValue(d.EnrollKey),
		CrowdsecFirewallVerbose: types.BoolValue(tools.StringToBool(d.CrowdsecFirewallVerbose)),
	}, nil
}
