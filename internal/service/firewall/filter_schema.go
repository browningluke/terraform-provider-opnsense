package firewall

import (
	"context"
	"regexp"
	"sort"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/firewall"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/browningluke/terraform-provider-opnsense/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type firewallLocation struct {
	Net    types.String `tfsdk:"net"`
	Port   types.String `tfsdk:"port"`
	Invert types.Bool   `tfsdk:"invert"`
}

type filterInterfaceBlock struct {
	Invert    types.Bool `tfsdk:"invert"`
	Interface types.Set  `tfsdk:"interface"`
}

type filterFilterBlock struct {
	Quick         types.Bool        `tfsdk:"quick"`
	Action        types.String      `tfsdk:"action"`
	AllowOptions  types.Bool        `tfsdk:"allow_options"`
	Direction     types.String      `tfsdk:"direction"`
	IPProtocol    types.String      `tfsdk:"ip_protocol"`
	Protocol      types.String      `tfsdk:"protocol"`
	ICMPType      types.Set         `tfsdk:"icmp_type"`
	Source        *firewallLocation `tfsdk:"source"`
	Destination   *firewallLocation `tfsdk:"destination"`
	Log           types.Bool        `tfsdk:"log"`
	TCPFlags      types.Set         `tfsdk:"tcp_flags"`
	TCPFlagsOutOf types.Set         `tfsdk:"tcp_flags_out_of"`
	Schedule      types.String      `tfsdk:"schedule"`
}

type filterAdaptiveTimeouts struct {
	Start types.Int64 `tfsdk:"start"`
	End   types.Int64 `tfsdk:"end"`
}

type filterNewConnections struct {
	Count   types.Int64 `tfsdk:"count"`
	Seconds types.Int64 `tfsdk:"seconds"`
}

type filterMax struct {
	States            types.Int64           `tfsdk:"states"`
	SourceNodes       types.Int64           `tfsdk:"source_nodes"`
	SourceStates      types.Int64           `tfsdk:"source_states"`
	SourceConnections types.Int64           `tfsdk:"source_connections"`
	NewConnections    *filterNewConnections `tfsdk:"new_connections"`
}

type filterStatefulFirewallBlock struct {
	Type             types.String            `tfsdk:"type"`
	Policy           types.String            `tfsdk:"policy"`
	Timeout          types.Int64             `tfsdk:"timeout"`
	AdaptiveTimeouts *filterAdaptiveTimeouts `tfsdk:"adaptive_timeouts"`
	Max              *filterMax              `tfsdk:"max"`
	OverloadTable    types.String            `tfsdk:"overload_table"`
	NoPfsync         types.Bool              `tfsdk:"no_pfsync"`
}

type filterTrafficShapingBlock struct {
	Shaper        types.String `tfsdk:"shaper"`
	ReverseShaper types.String `tfsdk:"reverse_shaper"`
}

type filterSourceRoutingBlock struct {
	Gateway        types.String `tfsdk:"gateway"`
	DisableReplyTo types.Bool   `tfsdk:"disable_reply_to"`
	ReplyTo        types.String `tfsdk:"reply_to"`
}

type filterPriorityBlock struct {
	Match       types.Int64  `tfsdk:"match"`
	Set         types.Int64  `tfsdk:"set"`
	LowDelaySet types.Int64  `tfsdk:"low_delay_set"`
	MatchTOS    types.String `tfsdk:"match_tos"`
}

type filterInternalTaggingBlock struct {
	SetLocal   types.String `tfsdk:"set_local"`
	MatchLocal types.String `tfsdk:"match_local"`
}

// filterResourceModel describes the resource data model.
type filterResourceModel struct {
	Enabled      types.Bool   `tfsdk:"enabled"`
	Sequence     types.Int64  `tfsdk:"sequence"`
	NoXMLRPCSync types.Bool   `tfsdk:"no_xmlrpc_sync"`
	Description  types.String `tfsdk:"description"`
	Categories   types.Set    `tfsdk:"categories"`

	Interface        *filterInterfaceBlock        `tfsdk:"interface"`
	Filter           *filterFilterBlock           `tfsdk:"filter"`
	StatefulFirewall *filterStatefulFirewallBlock `tfsdk:"stateful_firewall"`
	TrafficShaping   *filterTrafficShapingBlock   `tfsdk:"traffic_shaping"`
	SourceRouting    *filterSourceRoutingBlock    `tfsdk:"source_routing"`
	Priority         *filterPriorityBlock         `tfsdk:"priority"`
	InternalTagging  *filterInternalTaggingBlock  `tfsdk:"internal_tagging"`

	Id types.String `tfsdk:"id"`
}

// filterResourceModelV0 describes the OLD v0 schema with flat structure (pre-nested blocks).
// This is used for state migration from v0 to v1.
type filterResourceModelV0 struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Sequence    types.Int64  `tfsdk:"sequence"`
	Action      types.String `tfsdk:"action"`
	Quick       types.Bool   `tfsdk:"quick"`
	Interface   types.Set    `tfsdk:"interface"`
	Direction   types.String `tfsdk:"direction"`
	IPProtocol  types.String `tfsdk:"ip_protocol"`
	Protocol    types.String `tfsdk:"protocol"`
	Source      types.Object `tfsdk:"source"`
	Destination types.Object `tfsdk:"destination"`
	Gateway     types.String `tfsdk:"gateway"`
	Log         types.Bool   `tfsdk:"log"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
}

func filterResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Firewall filter rules can be used to restrict or allow traffic from and/or to specific networks as well as influence how traffic should be forwarded",
		Version:             1,

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
			"no_xmlrpc_sync": schema.BoolAttribute{
				MarkdownDescription: "Whether to exclude this item from the HA synchronization process. An already existing item with the same UUID on the synchronization target will not be altered or deleted as long as this is active. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Optional:            true,
			},
			"categories": schema.SetAttribute{
				MarkdownDescription: "For grouping purposes, provide the IDs of multiple groups here to organize items. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(validators.IsUUIDv4()),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"interface": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"invert": schema.BoolAttribute{
						MarkdownDescription: "Whether to use all but selected interfaces. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"interface": schema.SetAttribute{
						MarkdownDescription: "The interfaces to apply the filter rule on.",
						Required:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
				},
			},
			"filter": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"quick": schema.BoolAttribute{
						MarkdownDescription: "If a packet matches a rule specifying quick, then that rule is considered the last matching rule and the specified action is taken. When a rule does not have quick enabled, the last matching rule wins. Defaults to `true`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
					},
					"action": schema.StringAttribute{
						MarkdownDescription: "Choose what to do with packets that match the criteria specified below. Hint: the difference between block and reject is that with reject, a packet (TCP RST or ICMP port unreachable for UDP) is returned to the sender, whereas with block the packet is dropped silently. In either case, the original packet is discarded.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("pass", "block", "reject"),
						},
					},
					"allow_options": schema.BoolAttribute{
						MarkdownDescription: "Whether to allow packets with IP options to pass. Otherwise they are blocked. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"direction": schema.StringAttribute{
						MarkdownDescription: "Direction of the traffic. The default policy is to filter inbound traffic, which sets the policy to the interface originally receiving the traffic.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("in", "out"),
						},
					},
					"ip_protocol": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Validators: []validator.String{
							stringvalidator.OneOf("inet", "inet6", "inet46"),
						},
						Default: stringdefault.StaticString("inet"),
					},
					"protocol": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"3PC",
								"A/N",
								"AH",
								"ARGUS",
								"ARIS",
								"AX.25",
								"BBN-RCC",
								"BNA",
								"BR-SAT-MON",
								"CARP",
								"CBT",
								"CFTP",
								"CHAOS",
								"COMPAQ-PEER",
								"CPHB",
								"CPNX",
								"CRTP",
								"CRUDP",
								"DCCP",
								"DCN",
								"DDP",
								"DDX",
								"DGP",
								"DIVERT",
								"DSR",
								"EGP",
								"EIGRP",
								"EMCON",
								"ENCAP",
								"ESP",
								"ETHERIP",
								"FC",
								"GGP",
								"GMTP",
								"GRE",
								"HIP",
								"HMP",
								"I-NLSP",
								"IATP",
								"ICMP",
								"IDPR",
								"IDPR-CMTP",
								"IDRP",
								"IFMP",
								"IGMP",
								"IGP",
								"IL",
								"IPCOMP",
								"IPCV",
								"IPENCAP",
								"IPIP",
								"IPPC",
								"IPV6",
								"IPV6-ICMP",
								"IPX-IN-IP",
								"IRTP",
								"ISIS",
								"ISO-IP",
								"ISO-TP4",
								"KRYPTOLAN",
								"L2TP",
								"LARP",
								"LEAF-1",
								"LEAF-2",
								"MANET",
								"MERIT-INP",
								"MFE-NSP",
								"MICP",
								"MOBILE",
								"MPLS-IN-IP",
								"MTP",
								"MUX",
								"NARP",
								"NETBLT",
								"NSFNET-IGP",
								"NVP",
								"OSPF",
								"PFSYNC",
								"PGM",
								"PIM",
								"PIPE",
								"PNNI",
								"PRM",
								"PTP",
								"PUP",
								"PVP",
								"QNX",
								"RDP",
								"ROHC",
								"RSVP",
								"RSVP-E2E-IGNORE",
								"RVD",
								"SAT-EXPAK",
								"SAT-MON",
								"SCC-SP",
								"SCPS",
								"SCTP",
								"SDRP",
								"SECURE-VMTP",
								"SHIM6",
								"SKIP",
								"SM",
								"SMP",
								"SNP",
								"SPRITE-RPC",
								"SPS",
								"SRP",
								"ST2",
								"STP",
								"SUN-ND",
								"SWIPE",
								"TCF",
								"TCP",
								"TCP/UDP",
								"TLSP",
								"TP++",
								"TRUNK-1",
								"TRUNK-2",
								"TTP",
								"UDP",
								"UDPLITE",
								"UTI",
								"VINES",
								"VISA",
								"VMTP",
								"WB-EXPAK",
								"WB-MON",
								"WESP",
								"WSN",
								"XNET",
								"XNS-IDP",
								"XTP",
								"any",
							),
						},
					},
					"icmp_type": schema.SetAttribute{
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(stringvalidator.OneOf(
								"althost",
								"echorep",
								"echoreq",
								"inforep",
								"inforeq",
								"maskrep",
								"maskreq",
								"paramprob",
								"redir",
								"routeradv",
								"routersol",
								"squench",
								"timerep",
								"timereq",
								"timex",
								"unreach",
							)),
						},
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
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("any"),
							},
							"port": schema.StringAttribute{
								MarkdownDescription: "Source port number or well known name (imap, imaps, http, https, ...), for ranges use a dash. Defaults to `\"\"`.",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString(""),
								Validators: []validator.String{
									stringvalidator.RegexMatches(regexp.MustCompile("^(\\d|-)+$|^(\\w){0,32}$"),
										"must be number (80), range (80-443), well known name (http) or alias name"),
								},
							},
							"invert": schema.BoolAttribute{
								MarkdownDescription: "Whether to invert the sense of the match. Defaults to `false`.",
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
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString("any"),
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
								MarkdownDescription: "Whether to invert the sense of the match. Defaults to `false`.",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
						},
					},
					"log": schema.BoolAttribute{
						MarkdownDescription: "Whether to log packets that are handled by this rule. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"tcp_flags": schema.SetAttribute{
						MarkdownDescription: "The TCP flags that must be set this rule to match. Defaults to `[]`.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
					"tcp_flags_out_of": schema.SetAttribute{
						MarkdownDescription: "The TCP flags that must be cleared for this rule to match. Defaults to `[]`.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
					"schedule": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
				},
			},
			"stateful_firewall": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						map[string]attr.Type{
							"type":    types.StringType,
							"policy":  types.StringType,
							"timeout": types.Int64Type,
							"adaptive_timeouts": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"start": types.Int64Type,
									"end":   types.Int64Type,
								},
							},
							"max": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"states":             types.Int64Type,
									"source_nodes":       types.Int64Type,
									"source_states":      types.Int64Type,
									"source_connections": types.Int64Type,
									"new_connections": types.ObjectType{
										AttrTypes: map[string]attr.Type{
											"count":   types.Int64Type,
											"seconds": types.Int64Type,
										},
									},
								},
							},
							"overload_table": types.StringType,
							"no_pfsync":      types.BoolType,
						},
						map[string]attr.Value{
							"type":    types.StringValue("keep"),
							"policy":  types.StringValue(""),
							"timeout": types.Int64Value(-1),
							"adaptive_timeouts": types.ObjectValueMust(
								map[string]attr.Type{
									"start": types.Int64Type,
									"end":   types.Int64Type,
								},
								map[string]attr.Value{
									"start": types.Int64Value(-1),
									"end":   types.Int64Value(-1),
								},
							),
							"max": types.ObjectValueMust(
								map[string]attr.Type{
									"states":             types.Int64Type,
									"source_nodes":       types.Int64Type,
									"source_states":      types.Int64Type,
									"source_connections": types.Int64Type,
									"new_connections": types.ObjectType{
										AttrTypes: map[string]attr.Type{
											"count":   types.Int64Type,
											"seconds": types.Int64Type,
										},
									},
								},
								map[string]attr.Value{
									"states":             types.Int64Value(-1),
									"source_nodes":       types.Int64Value(-1),
									"source_states":      types.Int64Value(-1),
									"source_connections": types.Int64Value(-1),
									"new_connections": types.ObjectValueMust(
										map[string]attr.Type{
											"count":   types.Int64Type,
											"seconds": types.Int64Type,
										},
										map[string]attr.Value{
											"count":   types.Int64Value(-1),
											"seconds": types.Int64Value(-1),
										},
									),
								},
							),
							"overload_table": types.StringValue(""),
							"no_pfsync":      types.BoolValue(false),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "State tracking mechanism to use, default is full stateful tracking, sloppy ignores sequence numbers, use none for stateless rules. Defaults to `\"keep\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("keep"),
						Validators: []validator.String{
							stringvalidator.OneOf("keep", "modulate", "none", "sloppy", "synproxy"),
						},
					},
					"policy": schema.StringAttribute{
						MarkdownDescription: "How states created by this rule are treated, default (as defined in advanced), floating in which case states are valid on all interfaces or interface bound. Interface bound states are more secure, floating more flexible. Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.OneOf("", "floating", "if-bound"),
						},
					},
					"timeout": schema.Int64Attribute{
						MarkdownDescription: "State Timeout in seconds (TCP only). Defaults to `-1`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(-1),
					},
					"adaptive_timeouts": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Default: objectdefault.StaticValue(
							types.ObjectValueMust(
								map[string]attr.Type{
									"start": types.Int64Type,
									"end":   types.Int64Type,
								},
								map[string]attr.Value{
									"start": types.Int64Value(-1),
									"end":   types.Int64Value(-1),
								},
							),
						),
						Attributes: map[string]schema.Attribute{
							"start": schema.Int64Attribute{
								MarkdownDescription: "When the number of state entries exceeds this value, adaptive scaling begins. All timeout values are scaled linearly with factor `(adaptive.end - number of states) / (adaptive.end - adaptive.start)`. Defaults to `-1`.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(-1),
							},
							"end": schema.Int64Attribute{
								MarkdownDescription: "When reaching this number of state entries, all timeout values become zero, effectively purging all state entries immediately. This value is used to define the scale factor, it should not actually be reached (set a lower state limit). Defaults to `-1`.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(-1),
							},
						},
					},
					"max": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Default: objectdefault.StaticValue(
							types.ObjectValueMust(
								map[string]attr.Type{
									"states":             types.Int64Type,
									"source_nodes":       types.Int64Type,
									"source_states":      types.Int64Type,
									"source_connections": types.Int64Type,
									"new_connections": types.ObjectType{
										AttrTypes: map[string]attr.Type{
											"count":   types.Int64Type,
											"seconds": types.Int64Type,
										},
									},
								},
								map[string]attr.Value{
									"states":             types.Int64Value(-1),
									"source_nodes":       types.Int64Value(-1),
									"source_states":      types.Int64Value(-1),
									"source_connections": types.Int64Value(-1),
									"new_connections": types.ObjectValueMust(
										map[string]attr.Type{
											"count":   types.Int64Type,
											"seconds": types.Int64Type,
										},
										map[string]attr.Value{
											"count":   types.Int64Value(-1),
											"seconds": types.Int64Value(-1),
										},
									),
								},
							),
						),
						Attributes: map[string]schema.Attribute{
							"states": schema.Int64Attribute{
								MarkdownDescription: "Limits the number of concurrent states the rule may create. When this limit is reached, further packets that would create state are dropped until existing states time out. Defaults to `-1`.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(-1),
							},
							"source_nodes": schema.Int64Attribute{
								MarkdownDescription: "Limits the maximum number of source addresses which can simultaneously have state table entries. Defaults to `-1`.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(-1),
							},
							"source_states": schema.Int64Attribute{
								MarkdownDescription: "Limits the maximum number of simultaneous state entries that a single source address can create with this rule. Defaults to `-1`.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(-1),
							},
							"source_connections": schema.Int64Attribute{
								MarkdownDescription: "Limit the maximum number of simultaneous TCP connections which have completed the 3-way handshake that a single host can make. Defaults to `-1`.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(-1),
							},
							"new_connections": schema.SingleNestedAttribute{
								Optional: true,
								Computed: true,
								Default: objectdefault.StaticValue(
									types.ObjectValueMust(
										map[string]attr.Type{
											"count":   types.Int64Type,
											"seconds": types.Int64Type,
										},
										map[string]attr.Value{
											"count":   types.Int64Value(-1),
											"seconds": types.Int64Value(-1),
										},
									),
								),
								Attributes: map[string]schema.Attribute{
									"count": schema.Int64Attribute{
										MarkdownDescription: "Maximum new connections per host, measured over time. Defaults to `-1`.",
										Optional:            true,
										Computed:            true,
										Default:             int64default.StaticInt64(-1),
									},
									"seconds": schema.Int64Attribute{
										MarkdownDescription: "Time interval (seconds) to measure the number of connections. Defaults to `-1`.",
										Optional:            true,
										Computed:            true,
										Default:             int64default.StaticInt64(-1),
									},
								},
							},
						},
					},
					"overload_table": schema.StringAttribute{
						MarkdownDescription: "Overload table used when max new connections per time interval has been reached. The default virusprot table comes with a default block rule in floating rules, alternatively specify your own table here. Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"no_pfsync": schema.BoolAttribute{
						MarkdownDescription: "Whether to prevent states created by this rule to be synced with pfsync. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
			"traffic_shaping": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						map[string]attr.Type{
							"shaper":         types.StringType,
							"reverse_shaper": types.StringType,
						},
						map[string]attr.Value{
							"shaper":         types.StringValue(""),
							"reverse_shaper": types.StringValue(""),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"shaper": schema.StringAttribute{
						MarkdownDescription: "Shape packets using the selected pipe or queue in the rule direction. Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						Validators: []validator.String{
							validators.IsUUIDv4(),
						},
					},
					"reverse_shaper": schema.StringAttribute{
						MarkdownDescription: "Shape packets using the selected pipe or queue in the reverse rule direction. Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						Validators: []validator.String{
							validators.IsUUIDv4(),
						},
					},
				},
			},
			"source_routing": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						map[string]attr.Type{
							"gateway":          types.StringType,
							"disable_reply_to": types.BoolType,
							"reply_to":         types.StringType,
						},
						map[string]attr.Value{
							"gateway":          types.StringValue(""),
							"disable_reply_to": types.BoolValue(false),
							"reply_to":         types.StringValue(""),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"gateway": schema.StringAttribute{
						MarkdownDescription: "Leave as 'default' to use the system routing table. Or choose a gateway to utilize policy based routing. Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"disable_reply_to": schema.BoolAttribute{
						MarkdownDescription: "Whether to explicitly disable reply-to for this rule. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"reply_to": schema.StringAttribute{
						MarkdownDescription: "Determines how packets route back in the opposite direction (replies), when set to default, packets on WAN type interfaces reply to their connected gateway on the interface (unless globally disabled). A specific gateway may be chosen as well here. This setting is only relevant in the context of a state, for stateless rules there is no defined opposite direction. Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
				},
			},
			"priority": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						map[string]attr.Type{
							"match":         types.Int64Type,
							"set":           types.Int64Type,
							"low_delay_set": types.Int64Type,
							"match_tos":     types.StringType,
						},
						map[string]attr.Value{
							"match":         types.Int64Value(-1),
							"set":           types.Int64Value(-1),
							"low_delay_set": types.Int64Value(-1),
							"match_tos":     types.StringValue(""),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"match": schema.Int64Attribute{
						MarkdownDescription: "Only match packets which have the given queueing priority assigned. Defaults to `-1`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(-1),
						Validators: []validator.Int64{
							int64validator.Between(-1, 7),
						},
					},
					"set": schema.Int64Attribute{
						MarkdownDescription: "Packets matching this rule will be assigned a specific queueing priority. If the packet is transmitted on a vlan(4) interface, the queueing priority will be written as the priority code point in the 802.1Q VLAN header. Defaults to `-1`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(-1),
						Validators: []validator.Int64{
							int64validator.Between(-1, 7),
						},
					},
					"low_delay_set": schema.Int64Attribute{
						MarkdownDescription: "Used in combination with set priority, packets which have a TOS of lowdelay and TCP ACKs with no data payload will be assigned this priority when offered. Defaults to `-1`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(-1),
						Validators: []validator.Int64{
							int64validator.Between(-1, 7),
						},
					},
					"match_tos": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.OneOf(
								"",
								"0x0", "0x1", "0x2", "0x3", "0x4", "0x5", "0x6", "0x7", "0x8", "0x9", "0xa", "0xb", "0xc", "0xd", "0xe", "0xf",
								"0x10", "0x11", "0x12", "0x13", "0x14", "0x15", "0x16", "0x17", "0x18", "0x19", "0x1a", "0x1b", "0x1c", "0x1d", "0x1e", "0x1f",
								"0x20", "0x21", "0x22", "0x23", "0x24", "0x25", "0x26", "0x27", "0x28", "0x29", "0x2a", "0x2b", "0x2c", "0x2d", "0x2e", "0x2f",
								"0x30", "0x31", "0x32", "0x33", "0x34", "0x35", "0x36", "0x37", "0x38", "0x39", "0x3a", "0x3b", "0x3c", "0x3d", "0x3e", "0x3f",
								"0x40", "0x41", "0x42", "0x43", "0x44", "0x45", "0x46", "0x47", "0x48", "0x49", "0x4a", "0x4b", "0x4c", "0x4d", "0x4e", "0x4f",
								"0x50", "0x51", "0x52", "0x53", "0x54", "0x55", "0x56", "0x57", "0x58", "0x59", "0x5a", "0x5b", "0x5c", "0x5d", "0x5e", "0x5f",
								"0x60", "0x61", "0x62", "0x63", "0x64", "0x65", "0x66", "0x67", "0x68", "0x69", "0x6a", "0x6b", "0x6c", "0x6d", "0x6e", "0x6f",
								"0x70", "0x71", "0x72", "0x73", "0x74", "0x75", "0x76", "0x77", "0x78", "0x79", "0x7a", "0x7b", "0x7c", "0x7d", "0x7e", "0x7f",
								"0x80", "0x81", "0x82", "0x83", "0x84", "0x85", "0x86", "0x87", "0x88", "0x89", "0x8a", "0x8b", "0x8c", "0x8d", "0x8e", "0x8f",
								"0x90", "0x91", "0x92", "0x93", "0x94", "0x95", "0x96", "0x97", "0x98", "0x99", "0x9a", "0x9b", "0x9c", "0x9d", "0x9e", "0x9f",
								"0xa0", "0xa1", "0xa2", "0xa3", "0xa4", "0xa5", "0xa6", "0xa7", "0xa8", "0xa9", "0xaa", "0xab", "0xac", "0xad", "0xae", "0xaf",
								"0xb0", "0xb1", "0xb2", "0xb3", "0xb4", "0xb5", "0xb6", "0xb7", "0xb8", "0xb9", "0xba", "0xbb", "0xbc", "0xbd", "0xbe", "0xbf",
								"0xc0", "0xc1", "0xc2", "0xc3", "0xc4", "0xc5", "0xc6", "0xc7", "0xc8", "0xc9", "0xca", "0xcb", "0xcc", "0xcd", "0xce", "0xcf",
								"0xd0", "0xd1", "0xd2", "0xd3", "0xd4", "0xd5", "0xd6", "0xd7", "0xd8", "0xd9", "0xda", "0xdb", "0xdc", "0xdd", "0xde", "0xdf",
								"0xe0", "0xe1", "0xe2", "0xe3", "0xe4", "0xe5", "0xe6", "0xe7", "0xe8", "0xe9", "0xea", "0xeb", "0xec", "0xed", "0xee", "0xef",
								"0xf0", "0xf1", "0xf2", "0xf3", "0xf4", "0xf5", "0xf6", "0xf7", "0xf8", "0xf9", "0xfa", "0xfb", "0xfc", "0xfd", "0xfe", "0xff",
								"af11", "af12", "af13", "af21", "af22", "af23", "af31", "af32", "af33", "af41", "af42", "af43",
								"critical", "cs0", "cs1", "cs2", "cs3", "cs4", "cs5", "cs6", "cs7",
								"ef", "inetcontrol", "lowdelay", "netcontrol", "reliability", "throughput",
							),
						},
					},
				},
			},
			"internal_tagging": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						map[string]attr.Type{
							"set_local":   types.StringType,
							"match_local": types.StringType,
						},
						map[string]attr.Value{
							"set_local":   types.StringValue(""),
							"match_local": types.StringValue(""),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"set_local": schema.StringAttribute{
						MarkdownDescription: "Packets matching this rule will be tagged with the specified string. The tag acts as an internal marker that can be used to identify these packets later on. This can be used, for example, to provide trust between interfaces and to determine if packets have been processed by translation rules. Tags are \"sticky\", meaning that the packet will be tagged even if the rule is not the last matching rule. Further matching rules can replace the tag with a new one but will not remove a previously applied tag. A packet is only ever assigned one tag at a time. Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"match_local": schema.StringAttribute{
						MarkdownDescription: "Used to specify that packets must already be tagged with the given tag in order to match the rule. Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
				},
			},
		},
	}
}

func filterDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Firewall filter rules can be used to restrict or allow traffic from and/or to specific networks as well as influence how traffic should be forwarded",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether this firewall filter rule is enabled.",
				Computed:            true,
			},
			"sequence": dschema.Int64Attribute{
				MarkdownDescription: "The order of this filter rule.",
				Computed:            true,
			},
			"no_xmlrpc_sync": dschema.BoolAttribute{
				MarkdownDescription: "Whether this item is excluded from the HA synchronization process. An already existing item with the same UUID on the synchronization target will not be altered or deleted as long as this is active.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description for reference (not parsed).",
				Computed:            true,
			},
			"categories": dschema.SetAttribute{
				MarkdownDescription: "The IDs of multiple groups for organizing items.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"interface": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"invert": dschema.BoolAttribute{
						MarkdownDescription: "Whether all but selected interfaces are used.",
						Computed:            true,
					},
					"interface": dschema.SetAttribute{
						MarkdownDescription: "The interfaces the filter rule is applied on.",
						Computed:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"filter": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"quick": dschema.BoolAttribute{
						MarkdownDescription: "Whether a packet matching this rule is the last matching rule. If quick is enabled, the specified action is taken immediately.",
						Computed:            true,
					},
					"action": dschema.StringAttribute{
						MarkdownDescription: "What to do with packets that match the criteria. Hint: the difference between block and reject is that with reject, a packet (TCP RST or ICMP port unreachable for UDP) is returned to the sender, whereas with block the packet is dropped silently. In either case, the original packet is discarded.",
						Computed:            true,
					},
					"allow_options": dschema.BoolAttribute{
						MarkdownDescription: "Whether packets with IP options are allowed to pass. Otherwise they are blocked.",
						Computed:            true,
					},
					"direction": dschema.StringAttribute{
						MarkdownDescription: "The direction of the traffic. The default policy is to filter inbound traffic, which sets the policy to the interface originally receiving the traffic.",
						Computed:            true,
					},
					"ip_protocol": dschema.StringAttribute{
						Computed: true,
					},
					"protocol": dschema.StringAttribute{
						Computed: true,
					},
					"icmp_type": dschema.SetAttribute{
						Computed:    true,
						ElementType: types.StringType,
					},
					"source": dschema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]dschema.Attribute{
							"net": dschema.StringAttribute{
								Computed: true,
							},
							"port": dschema.StringAttribute{
								MarkdownDescription: "Source port number or well known name (imap, imaps, http, https, ...), for ranges use a dash.",
								Computed:            true,
							},
							"invert": dschema.BoolAttribute{
								MarkdownDescription: "Whether the sense of the match is inverted.",
								Computed:            true,
							},
						},
					},
					"destination": dschema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]dschema.Attribute{
							"net": dschema.StringAttribute{
								Computed: true,
							},
							"port": dschema.StringAttribute{
								MarkdownDescription: "Destination port number or well known name (imap, imaps, http, https, ...), for ranges use a dash.",
								Computed:            true,
							},
							"invert": dschema.BoolAttribute{
								MarkdownDescription: "Whether the sense of the match is inverted.",
								Computed:            true,
							},
						},
					},
					"log": dschema.BoolAttribute{
						MarkdownDescription: "Whether packets handled by this rule are logged.",
						Computed:            true,
					},
					"tcp_flags": dschema.SetAttribute{
						MarkdownDescription: "The TCP flags that must be set for this rule to match.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"tcp_flags_out_of": dschema.SetAttribute{
						MarkdownDescription: "The TCP flags that must be cleared for this rule to match.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"schedule": dschema.StringAttribute{
						Computed: true,
					},
				},
			},
			"stateful_firewall": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"type": dschema.StringAttribute{
						MarkdownDescription: "The state tracking mechanism used, default is full stateful tracking, sloppy ignores sequence numbers, use none for stateless rules.",
						Computed:            true,
					},
					"policy": dschema.StringAttribute{
						MarkdownDescription: "How states created by this rule are treated, default (as defined in advanced), floating in which case states are valid on all interfaces or interface bound. Interface bound states are more secure, floating more flexible.",
						Computed:            true,
					},
					"timeout": dschema.Int64Attribute{
						MarkdownDescription: "State Timeout in seconds (TCP only).",
						Computed:            true,
					},
					"adaptive_timeouts": dschema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]dschema.Attribute{
							"start": dschema.Int64Attribute{
								MarkdownDescription: "When the number of state entries exceeds this value, adaptive scaling begins. All timeout values are scaled linearly with factor `(adaptive.end - number of states) / (adaptive.end - adaptive.start)`.",
								Computed:            true,
							},
							"end": dschema.Int64Attribute{
								MarkdownDescription: "When reaching this number of state entries, all timeout values become zero, effectively purging all state entries immediately. This value is used to define the scale factor, it should not actually be reached (set a lower state limit).",
								Computed:            true,
							},
						},
					},
					"max": dschema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]dschema.Attribute{
							"states": dschema.Int64Attribute{
								MarkdownDescription: "The limit on the number of concurrent states the rule may create. When this limit is reached, further packets that would create state are dropped until existing states time out.",
								Computed:            true,
							},
							"source_nodes": dschema.Int64Attribute{
								MarkdownDescription: "The maximum number of source addresses which can simultaneously have state table entries.",
								Computed:            true,
							},
							"source_states": dschema.Int64Attribute{
								MarkdownDescription: "The maximum number of simultaneous state entries that a single source address can create with this rule.",
								Computed:            true,
							},
							"source_connections": dschema.Int64Attribute{
								MarkdownDescription: "The maximum number of simultaneous TCP connections which have completed the 3-way handshake that a single host can make.",
								Computed:            true,
							},
							"new_connections": dschema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]dschema.Attribute{
									"count": dschema.Int64Attribute{
										MarkdownDescription: "Maximum new connections per host, measured over time.",
										Computed:            true,
									},
									"seconds": dschema.Int64Attribute{
										MarkdownDescription: "Time interval (seconds) to measure the number of connections.",
										Computed:            true,
									},
								},
							},
						},
					},
					"overload_table": dschema.StringAttribute{
						MarkdownDescription: "The overload table used when max new connections per time interval has been reached. The default virusprot table comes with a default block rule in floating rules, alternatively specify your own table here.",
						Computed:            true,
					},
					"no_pfsync": dschema.BoolAttribute{
						MarkdownDescription: "Whether states created by this rule are prevented from being synced with pfsync.",
						Computed:            true,
					},
				},
			},
			"traffic_shaping": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"shaper": dschema.StringAttribute{
						MarkdownDescription: "The selected pipe or queue used to shape packets in the rule direction.",
						Computed:            true,
					},
					"reverse_shaper": dschema.StringAttribute{
						MarkdownDescription: "The selected pipe or queue used to shape packets in the reverse rule direction.",
						Computed:            true,
					},
				},
			},
			"source_routing": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"gateway": dschema.StringAttribute{
						MarkdownDescription: "The gateway used for routing. 'default' uses the system routing table. A specific gateway can be chosen to utilize policy based routing.",
						Computed:            true,
					},
					"disable_reply_to": dschema.BoolAttribute{
						MarkdownDescription: "Whether reply-to is explicitly disabled for this rule.",
						Computed:            true,
					},
					"reply_to": dschema.StringAttribute{
						MarkdownDescription: "How packets route back in the opposite direction (replies), when set to default, packets on WAN type interfaces reply to their connected gateway on the interface (unless globally disabled). A specific gateway may be chosen as well here. This setting is only relevant in the context of a state, for stateless rules there is no defined opposite direction.",
						Computed:            true,
					},
				},
			},
			"priority": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"match": dschema.Int64Attribute{
						MarkdownDescription: "Only matches packets which have the given queueing priority assigned.",
						Computed:            true,
					},
					"set": dschema.Int64Attribute{
						MarkdownDescription: "The specific queueing priority assigned to packets matching this rule. If the packet is transmitted on a vlan(4) interface, the queueing priority will be written as the priority code point in the 802.1Q VLAN header.",
						Computed:            true,
					},
					"low_delay_set": dschema.Int64Attribute{
						MarkdownDescription: "Used in combination with set priority, packets which have a TOS of lowdelay and TCP ACKs with no data payload will be assigned this priority when offered.",
						Computed:            true,
					},
					"match_tos": dschema.StringAttribute{
						Computed: true,
					},
				},
			},
			"internal_tagging": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"set_local": dschema.StringAttribute{
						MarkdownDescription: "Packets matching this rule are tagged with the specified string. The tag acts as an internal marker that can be used to identify these packets later on. This can be used, for example, to provide trust between interfaces and to determine if packets have been processed by translation rules. Tags are \"sticky\", meaning that the packet will be tagged even if the rule is not the last matching rule. Further matching rules can replace the tag with a new one but will not remove a previously applied tag. A packet is only ever assigned one tag at a time.",
						Computed:            true,
					},
					"match_local": dschema.StringAttribute{
						MarkdownDescription: "Specifies that packets must already be tagged with the given tag in order to match the rule.",
						Computed:            true,
					},
				},
			},
		},
	}
}

func convertFilterSchemaToStruct(d *filterResourceModel) (*firewall.Filter, error) {
	result := &firewall.Filter{
		Enabled:      tools.BoolToString(d.Enabled.ValueBool()),
		Sequence:     tools.Int64ToString(d.Sequence.ValueInt64()),
		NoXMLRPCSync: tools.BoolToString(d.NoXMLRPCSync.ValueBool()),
		Description:  d.Description.ValueString(),
	}

	// Parse 'Categories'
	if !d.Categories.IsNull() && !d.Categories.IsUnknown() {
		var categoryList []string
		d.Categories.ElementsAs(context.Background(), &categoryList, false)
		sort.Strings(categoryList)
		result.Categories = api.SelectedMapList(categoryList)
	}

	// Parse 'Interface' block
	if d.Interface != nil {
		result.InvertInterface = tools.BoolToString(d.Interface.Invert.ValueBool())
		var interfaceList []string
		d.Interface.Interface.ElementsAs(context.Background(), &interfaceList, false)
		sort.Strings(interfaceList)
		result.Interface = api.SelectedMapList(interfaceList)
	}

	// Parse 'Filter' block
	if d.Filter != nil {
		result.Quick = tools.BoolToString(d.Filter.Quick.ValueBool())
		result.Action = api.SelectedMap(d.Filter.Action.ValueString())
		result.AllowOptions = tools.BoolToString(d.Filter.AllowOptions.ValueBool())
		result.Direction = api.SelectedMap(d.Filter.Direction.ValueString())
		result.IPProtocol = api.SelectedMap(d.Filter.IPProtocol.ValueString())
		result.Protocol = api.SelectedMap(d.Filter.Protocol.ValueString())
		result.Log = tools.BoolToString(d.Filter.Log.ValueBool())
		result.Schedule = api.SelectedMap(d.Filter.Schedule.ValueString())

		// ICMP Type
		if !d.Filter.ICMPType.IsNull() && !d.Filter.ICMPType.IsUnknown() {
			var icmpTypeList []string
			d.Filter.ICMPType.ElementsAs(context.Background(), &icmpTypeList, false)
			sort.Strings(icmpTypeList)
			result.ICMPType = api.SelectedMapList(icmpTypeList)
		}

		// TCP Flags
		if !d.Filter.TCPFlags.IsNull() && !d.Filter.TCPFlags.IsUnknown() {
			var tcpFlagsList []string
			d.Filter.TCPFlags.ElementsAs(context.Background(), &tcpFlagsList, false)
			sort.Strings(tcpFlagsList)
			result.TCPFlags = api.SelectedMapList(tcpFlagsList)
		}

		// TCP Flags Out Of
		if !d.Filter.TCPFlagsOutOf.IsNull() && !d.Filter.TCPFlagsOutOf.IsUnknown() {
			var tcpFlagsOutOfList []string
			d.Filter.TCPFlagsOutOf.ElementsAs(context.Background(), &tcpFlagsOutOfList, false)
			sort.Strings(tcpFlagsOutOfList)
			result.TCPFlagsOutOf = api.SelectedMapList(tcpFlagsOutOfList)
		}

		// Source
		if d.Filter.Source != nil {
			result.SourceNet = d.Filter.Source.Net.ValueString()
			result.SourcePort = d.Filter.Source.Port.ValueString()
			result.SourceInvert = tools.BoolToString(d.Filter.Source.Invert.ValueBool())
		}

		// Destination
		if d.Filter.Destination != nil {
			result.DestinationNet = d.Filter.Destination.Net.ValueString()
			result.DestinationPort = d.Filter.Destination.Port.ValueString()
			result.DestinationInvert = tools.BoolToString(d.Filter.Destination.Invert.ValueBool())
		}
	}

	// Parse 'StatefulFirewall' block
	if d.StatefulFirewall != nil {
		result.StateType = api.SelectedMap(d.StatefulFirewall.Type.ValueString())
		result.StatePolicy = api.SelectedMap(d.StatefulFirewall.Policy.ValueString())
		result.StateTimeout = tools.Int64ToStringNegative(d.StatefulFirewall.Timeout.ValueInt64())
		result.OverloadTable = api.SelectedMap(d.StatefulFirewall.OverloadTable.ValueString())
		result.NoPfsync = tools.BoolToString(d.StatefulFirewall.NoPfsync.ValueBool())

		if d.StatefulFirewall.AdaptiveTimeouts != nil {
			result.AdaptiveTimeoutsStart = tools.Int64ToStringNegative(d.StatefulFirewall.AdaptiveTimeouts.Start.ValueInt64())
			result.AdaptiveTimeoutsEnd = tools.Int64ToStringNegative(d.StatefulFirewall.AdaptiveTimeouts.End.ValueInt64())
		}

		if d.StatefulFirewall.Max != nil {
			result.MaxStates = tools.Int64ToStringNegative(d.StatefulFirewall.Max.States.ValueInt64())
			result.MaxSourceNodes = tools.Int64ToStringNegative(d.StatefulFirewall.Max.SourceNodes.ValueInt64())
			result.MaxSourceStates = tools.Int64ToStringNegative(d.StatefulFirewall.Max.SourceStates.ValueInt64())
			result.MaxSourceConnections = tools.Int64ToStringNegative(d.StatefulFirewall.Max.SourceConnections.ValueInt64())

			if d.StatefulFirewall.Max.NewConnections != nil {
				result.MaxNewConnectionsCount = tools.Int64ToStringNegative(d.StatefulFirewall.Max.NewConnections.Count.ValueInt64())
				result.MaxNewConnectionsSeconds = tools.Int64ToStringNegative(d.StatefulFirewall.Max.NewConnections.Seconds.ValueInt64())
			}
		}
	}

	// Parse 'TrafficShaping' block
	if d.TrafficShaping != nil {
		result.TrafficShaper = api.SelectedMap(d.TrafficShaping.Shaper.ValueString())
		result.TrafficShaperReverse = api.SelectedMap(d.TrafficShaping.ReverseShaper.ValueString())
	}

	// Parse 'SourceRouting' block
	if d.SourceRouting != nil {
		result.Gateway = api.SelectedMap(d.SourceRouting.Gateway.ValueString())
		result.DisableReplyTo = tools.BoolToString(d.SourceRouting.DisableReplyTo.ValueBool())
		result.ReplyTo = api.SelectedMap(d.SourceRouting.ReplyTo.ValueString())
	}

	// Parse 'Priority' block
	if d.Priority != nil {
		result.MatchPriority = api.SelectedMap(tools.Int64ToStringNegative(d.Priority.Match.ValueInt64()))
		result.SetPriority = api.SelectedMap(tools.Int64ToStringNegative(d.Priority.Set.ValueInt64()))
		result.SetPriorityLowDelay = api.SelectedMap(tools.Int64ToStringNegative(d.Priority.LowDelaySet.ValueInt64()))
		result.MatchTOS = api.SelectedMap(d.Priority.MatchTOS.ValueString())
	}

	// Parse 'InternalTagging' block
	if d.InternalTagging != nil {
		result.SetLocalTag = d.InternalTagging.SetLocal.ValueString()
		result.MatchLocalTag = d.InternalTagging.MatchLocal.ValueString()
	}

	return result, nil
}

func convertFilterStructToSchema(d *firewall.Filter) (*filterResourceModel, error) {
	model := &filterResourceModel{
		Enabled:      types.BoolValue(tools.StringToBool(d.Enabled)),
		Sequence:     tools.StringToInt64Null(d.Sequence),
		NoXMLRPCSync: types.BoolValue(tools.StringToBool(d.NoXMLRPCSync)),
		Description:  tools.StringOrNull(d.Description),
	}

	// Parse 'Categories'
	categories := tools.StringSliceToSet(d.Categories)
	model.Categories = categories

	// Parse 'Interface' block
	interfaceSet := tools.StringSliceToSet(d.Interface)
	model.Interface = &filterInterfaceBlock{
		Invert:    types.BoolValue(tools.StringToBool(d.InvertInterface)),
		Interface: interfaceSet,
	}

	// Parse 'Filter' block
	model.Filter = &filterFilterBlock{
		Quick:        types.BoolValue(tools.StringToBool(d.Quick)),
		Action:       types.StringValue(d.Action.String()),
		AllowOptions: types.BoolValue(tools.StringToBool(d.AllowOptions)),
		Direction:    types.StringValue(d.Direction.String()),
		IPProtocol:   types.StringValue(d.IPProtocol.String()),
		Protocol:     types.StringValue(d.Protocol.String()),
		Log:          types.BoolValue(tools.StringToBool(d.Log)),
		Schedule:     types.StringValue(d.Schedule.String()),
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
	}

	// ICMP Type
	icmpTypeSet := tools.StringSliceToSet(d.ICMPType)
	model.Filter.ICMPType = icmpTypeSet

	// TCP Flags
	tcpFlagsSet := tools.StringSliceToSet(d.TCPFlags)
	model.Filter.TCPFlags = tcpFlagsSet

	// TCP Flags Out Of
	tcpFlagsOutOfSet := tools.StringSliceToSet(d.TCPFlagsOutOf)
	model.Filter.TCPFlagsOutOf = tcpFlagsOutOfSet

	// Parse 'StatefulFirewall' block
	model.StatefulFirewall = &filterStatefulFirewallBlock{
		Type:          types.StringValue(d.StateType.String()),
		Policy:        types.StringValue(d.StatePolicy.String()),
		Timeout:       types.Int64Value(tools.StringToInt64(d.StateTimeout)),
		OverloadTable: types.StringValue(d.OverloadTable.String()),
		NoPfsync:      types.BoolValue(tools.StringToBool(d.NoPfsync)),
		AdaptiveTimeouts: &filterAdaptiveTimeouts{
			Start: types.Int64Value(tools.StringToInt64(d.AdaptiveTimeoutsStart)),
			End:   types.Int64Value(tools.StringToInt64(d.AdaptiveTimeoutsEnd)),
		},
		Max: &filterMax{
			States:            types.Int64Value(tools.StringToInt64(d.MaxStates)),
			SourceNodes:       types.Int64Value(tools.StringToInt64(d.MaxSourceNodes)),
			SourceStates:      types.Int64Value(tools.StringToInt64(d.MaxSourceStates)),
			SourceConnections: types.Int64Value(tools.StringToInt64(d.MaxSourceConnections)),
			NewConnections: &filterNewConnections{
				Count:   types.Int64Value(tools.StringToInt64(d.MaxNewConnectionsCount)),
				Seconds: types.Int64Value(tools.StringToInt64(d.MaxNewConnectionsSeconds)),
			},
		},
	}

	// Parse 'TrafficShaping' block
	model.TrafficShaping = &filterTrafficShapingBlock{
		Shaper:        types.StringValue(d.TrafficShaper.String()),
		ReverseShaper: types.StringValue(d.TrafficShaperReverse.String()),
	}

	// Parse 'SourceRouting' block
	model.SourceRouting = &filterSourceRoutingBlock{
		Gateway:        types.StringValue(d.Gateway.String()),
		DisableReplyTo: types.BoolValue(tools.StringToBool(d.DisableReplyTo)),
		ReplyTo:        types.StringValue(d.ReplyTo.String()),
	}

	// Parse 'Priority' block
	model.Priority = &filterPriorityBlock{
		Match:       types.Int64Value(tools.StringToInt64(d.MatchPriority.String())),
		Set:         types.Int64Value(tools.StringToInt64(d.SetPriority.String())),
		LowDelaySet: types.Int64Value(tools.StringToInt64(d.SetPriorityLowDelay.String())),
		MatchTOS:    types.StringValue(d.MatchTOS.String()),
	}

	// Parse 'InternalTagging' block
	model.InternalTagging = &filterInternalTaggingBlock{
		SetLocal:   types.StringValue(d.SetLocalTag),
		MatchLocal: types.StringValue(d.MatchLocalTag),
	}

	return model, nil
}

// filterResourceSchemaV0 returns the v0 (flat) schema for state migration.
// This is the schema before nested blocks were introduced.
func filterResourceSchemaV0() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Firewall filter rules can be used to restrict or allow traffic from and/or to specific networks as well as influence how traffic should be forwarded",
		Version:             0,

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
				MarkdownDescription: "Select the Internet Protocol version this rule applies to. Available values: `inet`, `inet6`, `inet46`. Defaults to `inet`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("inet", "inet6", "inet46"),
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
						MarkdownDescription: "Destination port number, well known name (imap, imaps, http, https, ...) or alias name, for ranges use a dash. Defaults to `\"\"`.",
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
