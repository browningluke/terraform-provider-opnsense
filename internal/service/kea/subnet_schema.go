package kea

import (
	"context"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"strings"

	"github.com/browningluke/opnsense-go/pkg/kea"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// subnetResourceModel describes the resource data model.
type subnetResourceModel struct {
	Subnet types.String `tfsdk:"subnet"`
	Pools  types.Set    `tfsdk:"pools"`

	MatchClientId types.Bool `tfsdk:"match_client_id"`

	AutoCollect types.Bool `tfsdk:"auto_collect"`

	Routers      types.Set `tfsdk:"routers"`
	StaticRoutes types.Set `tfsdk:"static_routes"`

	DomainNameServers types.Set    `tfsdk:"dns_servers"`
	DomainName        types.String `tfsdk:"domain_name"`
	DomainSearch      types.Set    `tfsdk:"domain_search"`

	NTPServers  types.Set `tfsdk:"ntp_servers"`
	TimeServers types.Set `tfsdk:"time_servers"`

	NextServer   types.String `tfsdk:"next_server"`
	TFPTServer   types.String `tfsdk:"tfpt_server"`
	TFTPBootfile types.String `tfsdk:"tftp_bootfile"`

	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

type staticRouteModel struct {
	DestinationIp types.String `tfsdk:"destination_ip"`
	RouterIp      types.String `tfsdk:"router_ip"`
}

var keaStaticRouteTypes = map[string]attr.Type{
	"destination_ip": types.StringType,
	"router_ip":      types.StringType,
}

func subnetResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure DHCP subnets for Kea.",

		Attributes: map[string]schema.Attribute{
			"subnet": schema.StringAttribute{
				MarkdownDescription: "Subnet to use (e.g. `\"192.0.2.64/26\"`), should be large enough to hold the specified pools and reservations.",
				Required:            true,
			},
			"pools": schema.SetAttribute{
				MarkdownDescription: "Set of pools in range or subnet format (e.g. `\"192.168.0.100 - 192.168.0.200\"` , `\"192.0.2.64/26\"`). Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"match_client_id": schema.BoolAttribute{
				MarkdownDescription: "By default, KEA uses client-identifiers instead of MAC addresses to locate clients, disabling this option changes back to matching on MAC address which is used by most dhcp implementations. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"auto_collect": schema.BoolAttribute{
				MarkdownDescription: "Automatically update option data from the GUI for relevant attributes. When set, values for `routers`, `dns_servers` and `ntp_servers` will be ignored. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"routers": schema.SetAttribute{
				MarkdownDescription: "Default gateways to offer to the clients. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"static_routes": schema.SetNestedAttribute{
				MarkdownDescription: "Static routes that the client should install in its routing cache. Defaults to `[]`.",
				// Required:            true,
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"destination_ip": schema.StringAttribute{
							MarkdownDescription: "Destination IP address for static route.",
							Required:            true,
						},
						"router_ip": schema.StringAttribute{
							MarkdownDescription: "Gateway IP for static route.",
							Required:            true,
						},
					},
				},
				Default: setdefault.StaticValue(
					tools.EmptySetValue(
						types.ObjectType{
							AttrTypes: keaStaticRouteTypes,
						},
					),
				),
			},
			"dns_servers": schema.SetAttribute{
				MarkdownDescription: "DNS servers to offer to the clients. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"domain_name": schema.StringAttribute{
				MarkdownDescription: "Domain name to offer to the client, set to this firewall's domain name when left empty. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"domain_search": schema.SetAttribute{
				MarkdownDescription: "Set of Domain Names to be used by the client to locate not-fully-qualified domain names. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"ntp_servers": schema.SetAttribute{
				MarkdownDescription: "Set of IP addresses indicating NTP (RFC 5905) servers available to the client. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},

			"time_servers": schema.SetAttribute{
				MarkdownDescription: "Set of RFC 868 time servers available to the client. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"next_server": schema.StringAttribute{
				MarkdownDescription: "Next server IP address. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},

			"tfpt_server": schema.StringAttribute{
				MarkdownDescription: "TFTP server address or fqdn. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tftp_bootfile": schema.StringAttribute{
				MarkdownDescription: "Boot filename to request. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the subnet.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func subnetDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure DHCP subnets for Kea.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"subnet": dschema.StringAttribute{
				MarkdownDescription: "Subnet in use (e.g. `\"192.0.2.64/26\"`).",
				Computed:            true,
			},
			"pools": dschema.SetAttribute{
				MarkdownDescription: "Set of pools in range or subnet format (e.g. `\"192.168.0.100 - 192.168.0.200\"` , `\"192.0.2.64/26\"`).",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"match_client_id": dschema.BoolAttribute{
				MarkdownDescription: "By default, KEA uses client-identifiers instead of MAC addresses to locate clients, disabling this option changes back to matching on MAC address which is used by most dhcp implementations. Defaults to `true`.",
				Computed:            true,
			},
			"auto_collect": dschema.BoolAttribute{
				MarkdownDescription: "Automatically update option data from the GUI for relevant attributes. When set, values for `routers`, `dns_servers` and `ntp_servers` will be ignored.",
				Computed:            true,
			},
			"routers": dschema.SetAttribute{
				MarkdownDescription: "Default gateways to offer to the clients.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"static_routes": dschema.SetNestedAttribute{
				MarkdownDescription: "Static routes that the client should install in its routing cache.",
				Computed:            true,
				NestedObject: dschema.NestedAttributeObject{
					Attributes: map[string]dschema.Attribute{
						"destination_ip": dschema.StringAttribute{
							MarkdownDescription: "Destination IP address for static route.",
							Computed:            true,
						},
						"router_ip": dschema.StringAttribute{
							MarkdownDescription: "Gateway IP for static route.",
							Computed:            true,
						},
					},
				},
			},
			"dns_servers": dschema.SetAttribute{
				MarkdownDescription: "DNS servers to offer to the clients.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"domain_name": dschema.StringAttribute{
				MarkdownDescription: "Domain name to offer to the client, set to this firewall's domain name when left empty.",
				Optional:            true,
			},
			"domain_search": dschema.SetAttribute{
				MarkdownDescription: "Set of Domain Names to be used by the client to locate not-fully-qualified domain names.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ntp_servers": dschema.SetAttribute{
				MarkdownDescription: "Set of IP addresses indicating NTP (RFC 5905) servers available to the client.",
				Computed:            true,
				ElementType:         types.StringType,
			},

			"time_servers": dschema.SetAttribute{
				MarkdownDescription: "Set of RFC 868 time servers available to the client.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"next_server": dschema.StringAttribute{
				MarkdownDescription: "Next server IP address.",
				Computed:            true,
			},

			"tfpt_server": dschema.StringAttribute{
				MarkdownDescription: "TFTP server address or fqdn.",
				Computed:            true,
			},
			"tftp_bootfile": dschema.StringAttribute{
				MarkdownDescription: "Boot filename to request.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Computed:            true,
			},
		},
	}
}

func convertSubnetSchemaToStruct(d *subnetResourceModel) (*kea.Subnet, error) {
	// Parse static routes
	var routesList []staticRouteModel
	d.StaticRoutes.ElementsAs(context.Background(), &routesList, false)

	// Convert routes to string
	var routePieceList []string
	for _, route := range routesList {
		routePieceList = append(routePieceList, route.DestinationIp.ValueString()+","+
			route.RouterIp.ValueString())
	}

	return &kea.Subnet{
		Subnet:                d.Subnet.ValueString(),
		NextServer:            d.NextServer.ValueString(),
		Pools:                 tools.SetToString(d.Pools, "\n"),
		MatchClientId:         tools.BoolToString(d.MatchClientId.ValueBool()),
		OptionDataAutoCollect: tools.BoolToString(d.AutoCollect.ValueBool()),
		OptionData: kea.OptionData{
			DomainNameServers: tools.SetToStringSlice(d.DomainNameServers),
			DomainSearch:      tools.SetToStringSlice(d.DomainSearch),
			Routers:           tools.SetToStringSlice(d.Routers),
			StaticRoutes:      strings.Join(routePieceList, ";"),
			DomainName:        d.DomainName.ValueString(),
			NtpServers:        tools.SetToStringSlice(d.NTPServers),
			TimeServers:       tools.SetToStringSlice(d.TimeServers),
			TftpServerName:    d.TFPTServer.ValueString(),
			BootFileName:      d.TFTPBootfile.ValueString(),
		},
		Description: d.Description.ValueString(),
	}, nil
}

func convertSubnetStructToSchema(d *kea.Subnet) (*subnetResourceModel, error) {
	model := &subnetResourceModel{
		Subnet:            types.StringValue(d.Subnet),
		Pools:             tools.StringSliceToSet(strings.Split(d.Pools, "\n")),
		MatchClientId:     types.BoolValue(tools.StringToBool(d.MatchClientId)),
		AutoCollect:       types.BoolValue(tools.StringToBool(d.OptionDataAutoCollect)),
		Routers:           tools.StringSliceToSet(d.OptionData.Routers),
		DomainNameServers: tools.StringSliceToSet(d.OptionData.DomainNameServers),
		DomainName:        types.StringValue(d.OptionData.DomainName),
		DomainSearch:      tools.StringSliceToSet(d.OptionData.DomainSearch),
		NTPServers:        tools.StringSliceToSet(d.OptionData.NtpServers),
		TimeServers:       tools.StringSliceToSet(d.OptionData.TimeServers),
		NextServer:        types.StringValue(d.NextServer),
		TFPTServer:        types.StringValue(d.OptionData.TftpServerName),
		TFTPBootfile:      types.StringValue(d.OptionData.BootFileName),
		Description:       types.StringValue(d.Description),
	}

	// Create empty set first
	v, _ := types.SetValue(
		types.ObjectType{
			AttrTypes: keaStaticRouteTypes,
		},
		[]attr.Value{},
	)

	// Parse static routes
	routeSlice := strings.Split(d.OptionData.StaticRoutes, ";")
	routes := []staticRouteModel{}

	// Try to fill list
	if len(routeSlice) > 0 {
		for _, i := range routeSlice {
			routePiece := strings.Split(i, ",")
			// Skip if not properly formed route
			if len(routePiece) != 2 {
				continue
			}

			routes = append(routes, staticRouteModel{
				DestinationIp: types.StringValue(routePiece[0]),
				RouterIp:      types.StringValue(routePiece[1]),
			})
		}

		v, _ = types.SetValueFrom(
			context.Background(),
			types.ObjectType{
				AttrTypes: keaStaticRouteTypes,
			},
			routes,
		)
	}
	model.StaticRoutes = v

	return model, nil
}
