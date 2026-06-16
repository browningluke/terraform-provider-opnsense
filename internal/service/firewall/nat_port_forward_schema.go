package firewall

import (
	"context"
	"regexp"
	"sort"
	"strings"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/firewall"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
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
)

// natPortForwardResourceModel describes the resource data model.
type natPortForwardResourceModel struct {
	Enabled types.Bool `tfsdk:"enabled"`

	Sequence  types.Int64 `tfsdk:"sequence"`
	Interface types.Set   `tfsdk:"interface"`

	IPProtocol types.String `tfsdk:"ip_protocol"`
	Protocol   types.String `tfsdk:"protocol"`

	Source      *firewallLocation `tfsdk:"source"`
	Destination *firewallLocation `tfsdk:"destination"`
	Target      *firewallTarget   `tfsdk:"target"`

	Log           types.Bool   `tfsdk:"log"`
	NatReflection types.String `tfsdk:"nat_reflection"`
	Description   types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

type natPortForwardResourceModelV1 struct {
	Enabled types.Bool `tfsdk:"enabled"`

	Sequence  types.Int64  `tfsdk:"sequence"`
	Interface types.String `tfsdk:"interface"`

	IPProtocol types.String `tfsdk:"ip_protocol"`
	Protocol   types.String `tfsdk:"protocol"`

	Source      *firewallLocation `tfsdk:"source"`
	Destination *firewallLocation `tfsdk:"destination"`
	Target      *firewallTarget   `tfsdk:"target"`

	Log           types.Bool   `tfsdk:"log"`
	NatReflection types.String `tfsdk:"nat_reflection"`
	Description   types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func natPortForwardResourceSchema() schema.Schema {
	return schema.Schema{
		Version:             2,
		MarkdownDescription: "Destination NAT (port forwarding) redirects traffic arriving on an external interface to an internal host. Use this to expose internal services (e.g. web servers, SSH) to the outside network.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this port forwarding rule. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"sequence": schema.Int64Attribute{
				MarkdownDescription: "Specify the order of this port forwarding rule. Defaults to `1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"interface": schema.SetAttribute{
				MarkdownDescription: "Choose on which interface packets must come in to match this rule. Must specify at least 1.",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
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
						MarkdownDescription: "Specify the IP address, CIDR or alias for the source of the packet. For `<INT> net`, enter `<int>` (e.g. `lan`). For `<INT> address`, enter `<int>ip` (e.g. `lanip`). Defaults to `any`.",
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
						MarkdownDescription: "Specify the IP address, CIDR or alias for the destination of the packet. For `<INT> net`, enter `<int>` (e.g. `lan`). For `<INT> address`, enter `<int>ip` (e.g. `lanip`). Defaults to `any`.",
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
						MarkdownDescription: "Specify the internal IP address or alias to forward packets to.",
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
			"nat_reflection": schema.StringAttribute{
				MarkdownDescription: "NAT reflection mode. One of `default`, `enable`, or `disable`. `default` means OPNsense uses the global firewall NAT reflection setting. Defaults to `default`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("default"),
				Validators: []validator.String{
					stringvalidator.OneOf("default", "enable", "disable"),
				},
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

func natPortForwardResourceSchemaV1() schema.Schema {
	s := natPortForwardResourceSchema()
	s.Version = 1
	s.Attributes["interface"] = schema.StringAttribute{
		MarkdownDescription: "Choose on which interface packets must come in to match this rule.",
		Required:            true,
	}
	return s
}

func natPortForwardDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Destination NAT (port forwarding) redirects traffic arriving on an external interface to an internal host. Use this to expose internal services (e.g. web servers, SSH) to the outside network.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether this port forwarding rule is enabled.",
				Computed:            true,
			},
			"sequence": dschema.Int64Attribute{
				MarkdownDescription: "The order of this port forwarding rule.",
				Computed:            true,
			},
			"interface": dschema.SetAttribute{
				MarkdownDescription: "The interfaces on which packets must come in to match this rule.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ip_protocol": dschema.StringAttribute{
				MarkdownDescription: "The Internet Protocol version this rule applies to. Available values: `inet`, `inet6`.",
				Computed:            true,
			},
			"protocol": dschema.StringAttribute{
				MarkdownDescription: "The IP protocol this rule matches.",
				Computed:            true,
			},
			"source": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"net": dschema.StringAttribute{
						MarkdownDescription: "The IP address, CIDR or alias for the source of the packet.",
						Computed:            true,
					},
					"port": dschema.StringAttribute{
						MarkdownDescription: "The source port for this rule.",
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
						MarkdownDescription: "The IP address, CIDR or alias for the destination of the packet.",
						Computed:            true,
					},
					"port": dschema.StringAttribute{
						MarkdownDescription: "The port for the destination of the packet.",
						Computed:            true,
					},
					"invert": dschema.BoolAttribute{
						MarkdownDescription: "Whether the sense of the match is inverted.",
						Computed:            true,
					},
				},
			},
			"target": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"ip": dschema.StringAttribute{
						MarkdownDescription: "The internal IP address or alias packets are forwarded to.",
						Computed:            true,
					},
					"port": dschema.StringAttribute{
						MarkdownDescription: "The internal port number or well known name packets are forwarded to.",
						Computed:            true,
					},
				},
			},
			"log": dschema.BoolAttribute{
				MarkdownDescription: "Whether packets handled by this rule are logged.",
				Computed:            true,
			},
			"nat_reflection": dschema.StringAttribute{
				MarkdownDescription: "NAT reflection mode. One of `default`, `enable`, or `disable`.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description for reference.",
				Computed:            true,
			},
		},
	}
}

// natReflectionSchemaToAPI converts the user-friendly Terraform value
// ("default"/"enable"/"disable") to the OPNsense API value ("", "purenat", "disable").
func natReflectionSchemaToAPI(s string) string {
	switch s {
	case "enable":
		return "purenat"
	case "disable":
		return "disable"
	default:
		return ""
	}
}

// natReflectionAPIToSchema converts the OPNsense API value back to the
// user-friendly Terraform value.
func natReflectionAPIToSchema(s string) string {
	switch s {
	case "purenat":
		return "enable"
	case "disable":
		return "disable"
	default:
		return "default"
	}
}

func natPortForwardInterfaceSchemaToAPI(s types.Set) api.SelectedMapList {
	var interfaces []string
	s.ElementsAs(context.Background(), &interfaces, false)
	sort.Strings(interfaces)
	return api.SelectedMapList(interfaces)
}

func natPortForwardInterfaceStringToSet(s string) types.Set {
	var interfaces []attr.Value
	for _, iface := range strings.Split(s, ",") {
		iface = strings.TrimSpace(iface)
		if iface != "" {
			interfaces = append(interfaces, types.StringValue(iface))
		}
	}
	return types.SetValueMust(types.StringType, interfaces)
}

func natPortForwardInterfaceSliceToSet(s []string) types.Set {
	seen := make(map[string]struct{})
	var list []attr.Value
	for _, iface := range s {
		if _, ok := seen[iface]; ok {
			continue
		}
		seen[iface] = struct{}{}
		list = append(list, types.StringValue(iface))
	}
	sv, _ := types.SetValue(types.StringType, list)
	return sv
}

func convertNATPortForwardSchemaToStruct(d *natPortForwardResourceModel) (*firewall.NatPortForward, error) {
	return &firewall.NatPortForward{
		// Schema uses "enabled" (user-friendly), API uses "disabled" (inverted).
		Disabled:   tools.BoolToString(!d.Enabled.ValueBool()),
		Sequence:   tools.Int64ToString(d.Sequence.ValueInt64()),
		Interface:  natPortForwardInterfaceSchemaToAPI(d.Interface),
		IPProtocol: api.SelectedMap(d.IPProtocol.ValueString()),
		Protocol:   api.SelectedMap(d.Protocol.ValueString()),
		Source: firewall.NatPortForwardLocation{
			Network: d.Source.Net.ValueString(),
			Port:    d.Source.Port.ValueString(),
			Invert:  tools.BoolToString(d.Source.Invert.ValueBool()),
		},
		Destination: firewall.NatPortForwardLocation{
			Network: d.Destination.Net.ValueString(),
			Port:    d.Destination.Port.ValueString(),
			Invert:  tools.BoolToString(d.Destination.Invert.ValueBool()),
		},
		Target:        d.Target.IP.ValueString(),
		TargetPort:    d.Target.Port.ValueString(),
		Log:           tools.BoolToString(d.Log.ValueBool()),
		NatReflection: api.SelectedMap(natReflectionSchemaToAPI(d.NatReflection.ValueString())),
		Description:   d.Description.ValueString(),
	}, nil
}

func convertNATPortForwardStructToSchema(d *firewall.NatPortForward) (*natPortForwardResourceModel, error) {
	sourceNet := d.Source.Network
	if sourceNet == "" {
		sourceNet = "any"
	}
	destinationNet := d.Destination.Network
	if destinationNet == "" {
		destinationNet = "any"
	}

	return &natPortForwardResourceModel{
		// API uses "disabled" (inverted), schema uses "enabled" (user-friendly).
		Enabled:    types.BoolValue(!tools.StringToBool(d.Disabled)),
		Sequence:   tools.StringToInt64Null(d.Sequence),
		Interface:  natPortForwardInterfaceSliceToSet([]string(d.Interface)),
		IPProtocol: types.StringValue(d.IPProtocol.String()),
		Protocol:   types.StringValue(d.Protocol.String()),
		Source: &firewallLocation{
			Net:    types.StringValue(sourceNet),
			Port:   types.StringValue(d.Source.Port),
			Invert: types.BoolValue(tools.StringToBool(d.Source.Invert)),
		},
		Destination: &firewallLocation{
			Net:    types.StringValue(destinationNet),
			Port:   types.StringValue(d.Destination.Port),
			Invert: types.BoolValue(tools.StringToBool(d.Destination.Invert)),
		},
		Target: &firewallTarget{
			IP:   types.StringValue(d.Target),
			Port: types.StringValue(d.TargetPort),
		},
		Log:           types.BoolValue(tools.StringToBool(d.Log)),
		NatReflection: types.StringValue(natReflectionAPIToSchema(d.NatReflection.String())),
		Description:   tools.StringOrNull(d.Description),
	}, nil
}
