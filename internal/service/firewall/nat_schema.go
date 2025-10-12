package firewall

import (
	"regexp"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/firewall"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type firewallTarget struct {
	IP   types.String `tfsdk:"ip"`
	Port types.String `tfsdk:"port"`
}

// natResourceModel describes the resource data model.
type natResourceModel struct {
	Enabled    types.Bool `tfsdk:"enabled"`
	DisableNAT types.Bool `tfsdk:"disable_nat"`

	Sequence  types.Int64  `tfsdk:"sequence"`
	Interface types.String `tfsdk:"interface"`

	IPProtocol types.String `tfsdk:"ip_protocol"`
	Protocol   types.String `tfsdk:"protocol"`

	Source      *firewallLocation `tfsdk:"source"`
	Destination *firewallLocation `tfsdk:"destination"`
	Target      *firewallTarget   `tfsdk:"target"`

	Log         types.Bool   `tfsdk:"log"`
	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func natResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Network Address Translation (abbreviated to NAT) is a way to separate external and internal networks (WANs and LANs), and to share an external IP between clients on the internal network.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this firewall NAT rule. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"disable_nat": schema.BoolAttribute{
				MarkdownDescription: "Enabling this option will disable NAT for traffic matching this rule and stop processing Outbound NAT rules. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"sequence": schema.Int64Attribute{
				MarkdownDescription: "Specify the order of this NAT rule. Defaults to `1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"interface": schema.StringAttribute{
				MarkdownDescription: "Choose on which interface(s) packets must come in to match this rule.",
				Required:            true,
			},
			"ip_protocol": schema.StringAttribute{
				MarkdownDescription: "Select the Internet Protocol version this rule applies to. Available values: `inet`, `inet6`. Defaults to `inet`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("inet", "inet6"),
				},
				Default: stringdefault.StaticString("inet"),
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Choose which IP protocol this rule should match.",
				Required:            true,
			},
			"source": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						map[string]attr.Type{
							"net":    types.StringType,
							"port":   types.StringType,
							"invert": types.BoolType,
						},
						map[string]attr.Value{
							"net":    types.StringValue("any"),
							"port":   types.StringValue(""),
							"invert": types.BoolValue(false),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"net": schema.StringAttribute{
						MarkdownDescription: "Specify the IP address, CIDR or alias for the source of the packet for this mapping. For `<INT> net`, enter `<int>` (e.g. `lan`). For `<INT> address`, enter `<int>ip` (e.g. `lanip`). Defaults to `any`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("any"),
					},
					"port": schema.StringAttribute{
						MarkdownDescription: "Specify the source port for this rule. This is usually random and almost never equal to the destination port range (and should usually be `\"\"`). Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile("^(\\d|-)+$|^([a-z])+$"),
								"must be number (80), range (80-443) or well known name (http)"),
						},
					},
					"invert": schema.BoolAttribute{
						MarkdownDescription: "Use this option to invert the sense of the match. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
			"destination": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						map[string]attr.Type{
							"net":    types.StringType,
							"port":   types.StringType,
							"invert": types.BoolType,
						},
						map[string]attr.Value{
							"net":    types.StringValue("any"),
							"port":   types.StringValue(""),
							"invert": types.BoolValue(false),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"net": schema.StringAttribute{
						MarkdownDescription: "Specify the IP address, CIDR or alias for the destination of the packet for this mapping. For `<INT> net`, enter `<int>` (e.g. `lan`). For `<INT> address`, enter `<int>ip` (e.g. `lanip`). Defaults to `any`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("any"),
					},
					"port": schema.StringAttribute{
						MarkdownDescription: "Destination port number or well known name (imap, imaps, http, https, ...), for ranges use a dash. Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile("^(\\d|-)+$|^(\\w){0,32}$"),
								"must be number (80), range (80-443), well known name (http) or alias name"),
						},
					},
					"invert": schema.BoolAttribute{
						MarkdownDescription: "Use this option to invert the sense of the match. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
			"target": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"ip": schema.StringAttribute{
						MarkdownDescription: "Specify the IP address or alias for the packets to be mapped to. For `<INT> address`, enter `<int>ip` (e.g. `lanip`).",
						Required:            true,
					},
					"port": schema.StringAttribute{
						MarkdownDescription: "Destination port number or well known name (imap, imaps, http, https, ...), for ranges use a dash. Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile("^(\\d|-)+$|^([a-z])+$"),
								"must be number (80), range (80-443) or well known name (http)"),
						},
					},
				},
			},
			"log": schema.BoolAttribute{
				MarkdownDescription: "Log packets that are handled by this rule. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed). Must be between 1 and 255 characters. Must be a character in set `[a-zA-Z0-9 .]`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9 .]*$`),
						"must only contain only alphanumeric characters, spaces or `.`",
					),
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func natDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Network Address Translation (abbreviated to NAT) is a way to separate external and internal networks (WANs and LANs), and to share an external IP between clients on the internal network.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this firewall NAT rule.",
				Computed:            true,
			},
			"disable_nat": schema.BoolAttribute{
				MarkdownDescription: "Enabling this option will disable NAT for traffic matching this rule and stop processing Outbound NAT rules.",
				Computed:            true,
			},
			"sequence": dschema.Int64Attribute{
				MarkdownDescription: "Specify the order of this NAT rule.",
				Computed:            true,
			},
			"interface": dschema.StringAttribute{
				MarkdownDescription: "The interface on which packets must come in to match this rule.",
				Computed:            true,
			},
			"ip_protocol": dschema.StringAttribute{
				MarkdownDescription: "Select the Internet Protocol version this rule applies to. Available values: `inet`, `inet6`.",
				Computed:            true,
			},
			"protocol": dschema.StringAttribute{
				MarkdownDescription: "Choose which IP protocol this rule should match.",
				Computed:            true,
			},
			"source": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"net": dschema.StringAttribute{
						MarkdownDescription: "Specify the IP address, CIDR or alias for the source of the packet for this mapping.",
						Computed:            true,
					},
					"port": dschema.StringAttribute{
						MarkdownDescription: "Specify the source port for this rule. This is usually random and almost never equal to the destination port range (and should usually be `\"\"`).",
						Computed:            true,
					},
					"invert": dschema.BoolAttribute{
						MarkdownDescription: "Use this option to invert the sense of the match.",
						Computed:            true,
					},
				},
			},
			"destination": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"net": dschema.StringAttribute{
						MarkdownDescription: "Specify the IP address, CIDR or alias for the destination of the packet for this mapping.",
						Computed:            true,
					},
					"port": dschema.StringAttribute{
						MarkdownDescription: "Specify the port for the destination of the packet for this mapping.",
						Computed:            true,
					},
					"invert": dschema.BoolAttribute{
						MarkdownDescription: "Use this option to invert the sense of the match.",
						Computed:            true,
					},
				},
			},
			"target": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"ip": schema.StringAttribute{
						MarkdownDescription: "Specify the IP address or alias for the packets to be mapped to.",
						Computed:            true,
					},
					"port": schema.StringAttribute{
						MarkdownDescription: "Destination port number or well known name (imap, imaps, http, https, ...), for ranges use a dash.",
						Computed:            true,
					},
				},
			},
			"log": dschema.BoolAttribute{
				MarkdownDescription: "Log packets that are handled by this rule.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Computed:            true,
			},
		},
	}
}

func convertNATSchemaToStruct(d *natResourceModel) (*firewall.NAT, error) {
	return &firewall.NAT{
		Enabled:           tools.BoolToString(d.Enabled.ValueBool()),
		DisableNAT:        tools.BoolToString(d.DisableNAT.ValueBool()),
		Sequence:          tools.Int64ToString(d.Sequence.ValueInt64()),
		Interface:         api.SelectedMap(d.Interface.ValueString()),
		IPProtocol:        api.SelectedMap(d.IPProtocol.ValueString()),
		Protocol:          api.SelectedMap(d.Protocol.ValueString()),
		SourceNet:         d.Source.Net.ValueString(),
		SourcePort:        d.Source.Port.ValueString(),
		SourceInvert:      tools.BoolToString(d.Source.Invert.ValueBool()),
		DestinationNet:    d.Destination.Net.ValueString(),
		DestinationPort:   d.Destination.Port.ValueString(),
		DestinationInvert: tools.BoolToString(d.Destination.Invert.ValueBool()),
		Target:            d.Target.IP.ValueString(),
		TargetPort:        d.Target.Port.ValueString(),
		Log:               tools.BoolToString(d.Log.ValueBool()),
		Description:       d.Description.ValueString(),
	}, nil
}

func convertNATStructToSchema(d *firewall.NAT) (*natResourceModel, error) {
	return &natResourceModel{
		Enabled:    types.BoolValue(tools.StringToBool(d.Enabled)),
		DisableNAT: types.BoolValue(tools.StringToBool(d.DisableNAT)),
		Sequence:   tools.StringToInt64Null(d.Sequence),
		Interface:  types.StringValue(d.Interface.String()),
		IPProtocol: types.StringValue(d.IPProtocol.String()),
		Protocol:   types.StringValue(d.Protocol.String()),
		Source: &firewallLocation{
			Net:    types.StringValue(d.SourceNet),
			Port:   types.StringValue(d.SourcePort),
			Invert: types.BoolValue(tools.StringToBool(d.SourceInvert)),
		},
		Destination: &firewallLocation{
			Net:    types.StringValue(d.DestinationNet),
			Port:   types.StringValue(d.DestinationPort),
			Invert: types.BoolValue(tools.StringToBool(d.DestinationInvert)),
		},
		Target: &firewallTarget{
			IP:   types.StringValue(d.Target),
			Port: types.StringValue(d.TargetPort),
		},
		Log:         types.BoolValue(tools.StringToBool(d.Log)),
		Description: tools.StringOrNull(d.Description),
	}, nil
}
