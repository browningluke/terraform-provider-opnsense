package service

import (
	"context"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/firewall"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"regexp"
	"terraform-provider-opnsense/internal/tools"
)

type firewallLocation struct {
	Net    types.String `tfsdk:"net"`
	Port   types.String `tfsdk:"port"`
	Invert types.Bool   `tfsdk:"invert"`
}

// FirewallFilterResourceModel describes the resource data model.
type FirewallFilterResourceModel struct {
	Enabled  types.Bool   `tfsdk:"enabled"`
	Sequence types.Int64  `tfsdk:"sequence"`
	Action   types.String `tfsdk:"action"`
	Quick    types.Bool   `tfsdk:"quick"`

	Interface types.Set    `tfsdk:"interface"`
	Direction types.String `tfsdk:"direction"`

	IPProtocol types.String `tfsdk:"ip_protocol"`
	Protocol   types.String `tfsdk:"protocol"`

	Source      *firewallLocation `tfsdk:"source"`
	Destination *firewallLocation `tfsdk:"destination"`

	Gateway types.String `tfsdk:"gateway"`
	Log     types.Bool   `tfsdk:"log"`

	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func FirewallFilterResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Firewall filter rules can be used to restrict or allow traffic from and/or to specific networks as well as influence how traffic should be forwarded",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this firewall filter rule. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"sequence": schema.Int64Attribute{
				MarkdownDescription: "Specify the order of this filter rule. Defaults to `1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Choose what to do with packets that match the criteria specified below. Hint: the difference between block and reject is that with reject, a packet (TCP RST or ICMP port unreachable for UDP) is returned to the sender, whereas with block the packet is dropped silently. In either case, the original packet is discarded. Available values: `pass`, `block`, `reject`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("pass", "block", "reject"),
				},
			},
			"quick": schema.BoolAttribute{
				MarkdownDescription: "If a packet matches a rule specifying quick, then that rule is considered the last matching rule and the specified action is taken. When a rule does not have quick enabled, the last matching rule wins. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"interface": schema.SetAttribute{
				MarkdownDescription: "Choose on which interface(s) packets must come in to match this rule. Must specify at least 1.",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"direction": schema.StringAttribute{
				MarkdownDescription: "Direction of the traffic. The default policy is to filter inbound traffic, which sets the policy to the interface originally receiving the traffic. Available values: `in`, `out`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("in", "out"),
				},
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
			"gateway": schema.StringAttribute{
				MarkdownDescription: "Leave as `\"\"` to use the system routing table. Or choose a gateway to utilize policy based routing. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"log": schema.BoolAttribute{
				MarkdownDescription: "Log packets that are handled by this rule. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Optional:            true,
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

func FirewallFilterDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Firewall filter rules can be used to restrict or allow traffic from and/or to specific networks as well as influence how traffic should be forwarded",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this firewall filter rule.",
				Computed:            true,
			},
			"sequence": dschema.Int64Attribute{
				MarkdownDescription: "Specify the order of this filter rule.",
				Computed:            true,
			},
			"action": dschema.StringAttribute{
				MarkdownDescription: "Choose what to do with packets that match the criteria specified below. Hint: the difference between block and reject is that with reject, a packet (TCP RST or ICMP port unreachable for UDP) is returned to the sender, whereas with block the packet is dropped silently. In either case, the original packet is discarded. Available values: `pass`, `block`, `reject`.",
				Computed:            true,
			},
			"quick": dschema.BoolAttribute{
				MarkdownDescription: "If a packet matches a rule specifying quick, then that rule is considered the last matching rule and the specified action is taken. When a rule does not have quick enabled, the last matching rule wins.",
				Computed:            true,
			},
			"interface": dschema.SetAttribute{
				MarkdownDescription: "The interface(s) on which the packets must come in to match this rule.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"direction": dschema.StringAttribute{
				MarkdownDescription: "Direction of the traffic. The default policy is to filter inbound traffic, which sets the policy to the interface originally receiving the traffic. Available values: `in`, `out`.",
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
			"gateway": dschema.StringAttribute{
				MarkdownDescription: "Leave as `\"\"` to use the system routing table. Or choose a gateway to utilize policy based routing.",
				Computed:            true,
			},
			"log": dschema.BoolAttribute{
				MarkdownDescription: "Log packets that are handled by this rule.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed). Must be between 1 and 255 characters. Must be a character in set `[a-zA-Z0-9 .]`.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9 .]*$`),
						"must only contain only alphanumeric characters, spaces or `.`",
					),
					stringvalidator.LengthBetween(1, 255),
				},
			},
		},
	}
}

func convertFirewallFilterSchemaToStruct(d *FirewallFilterResourceModel) (*firewall.Filter, error) {
	// Parse 'Interface'
	var interfaceList []string
	d.Interface.ElementsAs(context.Background(), &interfaceList, false)

	return &firewall.Filter{
		Enabled:           tools.BoolToString(d.Enabled.ValueBool()),
		Sequence:          tools.Int64ToString(d.Sequence.ValueInt64()),
		Action:            api.SelectedMap(d.Action.ValueString()),
		Quick:             tools.BoolToString(d.Quick.ValueBool()),
		Interface:         interfaceList,
		Direction:         api.SelectedMap(d.Direction.ValueString()),
		IPProtocol:        api.SelectedMap(d.IPProtocol.ValueString()),
		Protocol:          api.SelectedMap(d.Protocol.ValueString()),
		SourceNet:         d.Source.Net.ValueString(),
		SourcePort:        d.Source.Port.ValueString(),
		SourceInvert:      tools.BoolToString(d.Source.Invert.ValueBool()),
		DestinationNet:    d.Destination.Net.ValueString(),
		DestinationPort:   d.Destination.Port.ValueString(),
		DestinationInvert: tools.BoolToString(d.Destination.Invert.ValueBool()),
		Gateway:           api.SelectedMap(d.Gateway.ValueString()),
		Log:               tools.BoolToString(d.Log.ValueBool()),
		Description:       d.Description.ValueString(),
	}, nil
}

func convertFirewallFilterStructToSchema(d *firewall.Filter) (*FirewallFilterResourceModel, error) {
	model := &FirewallFilterResourceModel{
		Enabled:    types.BoolValue(tools.StringToBool(d.Enabled)),
		Sequence:   tools.StringToInt64Null(d.Sequence),
		Action:     types.StringValue(d.Action.String()),
		Quick:      types.BoolValue(tools.StringToBool(d.Quick)),
		Interface:  types.SetNull(types.StringType),
		Direction:  types.StringValue(d.Direction.String()),
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
		Gateway:     types.StringValue(d.Gateway.String()),
		Log:         types.BoolValue(tools.StringToBool(d.Log)),
		Description: tools.StringOrNull(d.Description),
	}

	// Parse 'Interface'
	var interfaceList []attr.Value
	for _, i := range d.Interface {
		interfaceList = append(interfaceList, basetypes.NewStringValue(i))
	}
	interfaceTypeList, _ := types.SetValue(types.StringType, interfaceList)
	model.Interface = interfaceTypeList

	return model, nil
}
