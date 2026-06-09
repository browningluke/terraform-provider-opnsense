package openvpn

import (
	"context"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/openvpn"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clientOverwriteResourceModel describes the resource data model.
type clientOverwriteResourceModel struct {
	Enabled         types.Bool   `tfsdk:"enabled"`
	CommonName      types.String `tfsdk:"common_name"`
	Description     types.String `tfsdk:"description"`
	Servers         types.Set    `tfsdk:"servers"`
	Block           types.Bool   `tfsdk:"block"`
	PushReset       types.Bool   `tfsdk:"push_reset"`
	TunnelNetwork   types.String `tfsdk:"tunnel_network"`
	TunnelNetworkV6 types.String `tfsdk:"tunnel_network_v6"`
	LocalNetworks   types.Set    `tfsdk:"local_networks"`
	RemoteNetworks  types.Set    `tfsdk:"remote_networks"`
	RouteGateway    types.String `tfsdk:"route_gateway"`
	RedirectGateway types.Set    `tfsdk:"redirect_gateway"`
	RegisterDNS     types.Bool   `tfsdk:"register_dns"`
	DNSDomain       types.Set    `tfsdk:"dns_domain"`
	DNSDomainSearch types.Set    `tfsdk:"dns_domain_search"`
	DNSServers      types.Set    `tfsdk:"dns_servers"`
	NTPServers      types.Set    `tfsdk:"ntp_servers"`
	WINSServers     types.Set    `tfsdk:"wins_servers"`

	Id types.String `tfsdk:"id"`
}

func clientOverwriteResourceSchema() schema.Schema {
	emptySet := setdefault.StaticValue(tools.EmptySetValue(types.StringType))

	return schema.Schema{
		MarkdownDescription: "Client-specific overrides applied by an OpenVPN server when a client with a matching common name connects (e.g. fixed tunnel address, pushed routes, DNS).",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this override. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"common_name": schema.StringAttribute{
				MarkdownDescription: "Client common name to match. Wildcards are not allowed; use one resource per common name.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of this override. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"servers": schema.SetAttribute{
				MarkdownDescription: "UUIDs of the OpenVPN server instances this override applies to. When empty, the override applies to all servers. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"block": schema.BoolAttribute{
				MarkdownDescription: "When `true`, blocks this client from connecting. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"push_reset": schema.BoolAttribute{
				MarkdownDescription: "When `true`, sends `push-reset` so the client ignores any server-pushed options not explicitly redefined here. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"tunnel_network": schema.StringAttribute{
				MarkdownDescription: "IPv4 tunnel network to assign this client (CIDR notation). Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tunnel_network_v6": schema.StringAttribute{
				MarkdownDescription: "IPv6 tunnel network to assign this client (CIDR notation). Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"local_networks": schema.SetAttribute{
				MarkdownDescription: "Local IPv4/IPv6 networks accessible from this client (CIDR). Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"remote_networks": schema.SetAttribute{
				MarkdownDescription: "Remote networks reachable behind this client (CIDR), pushed as iroute. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"route_gateway": schema.StringAttribute{
				MarkdownDescription: "Override the default route gateway for this client. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"redirect_gateway": schema.SetAttribute{
				MarkdownDescription: "OpenVPN `redirect-gateway` flags to push. Any of `local`, `autolocal`, `def1`, `bypass-dhcp`, `bypass-dns`, `block-local`, `ipv6`, `!ipv4`. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"register_dns": schema.BoolAttribute{
				MarkdownDescription: "Push `register-dns` so Windows clients refresh DNS settings. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"dns_domain": schema.SetAttribute{
				MarkdownDescription: "Push DNS domains to the client. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"dns_domain_search": schema.SetAttribute{
				MarkdownDescription: "Push DNS search domains to the client. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"dns_servers": schema.SetAttribute{
				MarkdownDescription: "Push DNS servers to the client. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"ntp_servers": schema.SetAttribute{
				MarkdownDescription: "Push NTP servers to the client. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"wins_servers": schema.SetAttribute{
				MarkdownDescription: "Push WINS servers to the client. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the override.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func clientOverwriteDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Read an existing OpenVPN client-specific override by UUID.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled":           dschema.BoolAttribute{MarkdownDescription: "Whether this override is enabled.", Computed: true},
			"common_name":       dschema.StringAttribute{MarkdownDescription: "Client common name to match.", Computed: true},
			"description":       dschema.StringAttribute{MarkdownDescription: "Description of this override.", Computed: true},
			"servers":           dschema.SetAttribute{MarkdownDescription: "UUIDs of the OpenVPN server instances this override applies to.", Computed: true, ElementType: types.StringType},
			"block":             dschema.BoolAttribute{MarkdownDescription: "Whether this client is blocked.", Computed: true},
			"push_reset":        dschema.BoolAttribute{MarkdownDescription: "Whether push-reset is sent.", Computed: true},
			"tunnel_network":    dschema.StringAttribute{MarkdownDescription: "IPv4 tunnel network.", Computed: true},
			"tunnel_network_v6": dschema.StringAttribute{MarkdownDescription: "IPv6 tunnel network.", Computed: true},
			"local_networks":    dschema.SetAttribute{MarkdownDescription: "Local networks pushed to this client.", Computed: true, ElementType: types.StringType},
			"remote_networks":   dschema.SetAttribute{MarkdownDescription: "Remote networks reachable behind this client.", Computed: true, ElementType: types.StringType},
			"route_gateway":     dschema.StringAttribute{MarkdownDescription: "Route gateway override.", Computed: true},
			"redirect_gateway":  dschema.SetAttribute{MarkdownDescription: "redirect-gateway flags.", Computed: true, ElementType: types.StringType},
			"register_dns":      dschema.BoolAttribute{MarkdownDescription: "Push register-dns to Windows clients.", Computed: true},
			"dns_domain":        dschema.SetAttribute{MarkdownDescription: "DNS domains pushed to the client.", Computed: true, ElementType: types.StringType},
			"dns_domain_search": dschema.SetAttribute{MarkdownDescription: "DNS search domains pushed to the client.", Computed: true, ElementType: types.StringType},
			"dns_servers":       dschema.SetAttribute{MarkdownDescription: "DNS servers pushed to the client.", Computed: true, ElementType: types.StringType},
			"ntp_servers":       dschema.SetAttribute{MarkdownDescription: "NTP servers pushed to the client.", Computed: true, ElementType: types.StringType},
			"wins_servers":      dschema.SetAttribute{MarkdownDescription: "WINS servers pushed to the client.", Computed: true, ElementType: types.StringType},
		},
	}
}

func convertClientOverwriteSchemaToStruct(d *clientOverwriteResourceModel) (*openvpn.ClientOverwrite, error) {
	ctx := context.Background()
	toList := func(s types.Set) []string {
		var out []string
		s.ElementsAs(ctx, &out, false)
		return out
	}

	return &openvpn.ClientOverwrite{
		Enabled:         tools.BoolToString(d.Enabled.ValueBool()),
		CommonName:      d.CommonName.ValueString(),
		Description:     d.Description.ValueString(),
		Servers:         api.SelectedMapList(toList(d.Servers)),
		Block:           tools.BoolToString(d.Block.ValueBool()),
		PushReset:       tools.BoolToString(d.PushReset.ValueBool()),
		TunnelNetwork:   d.TunnelNetwork.ValueString(),
		TunnelNetworkV6: d.TunnelNetworkV6.ValueString(),
		LocalNetworks:   api.SelectedMapList(toList(d.LocalNetworks)),
		RemoteNetworks:  api.SelectedMapList(toList(d.RemoteNetworks)),
		RouteGateway:    d.RouteGateway.ValueString(),
		RedirectGateway: api.SelectedMapList(toList(d.RedirectGateway)),
		RegisterDNS:     tools.BoolToString(d.RegisterDNS.ValueBool()),
		DNSDomain:       api.SelectedMapList(toList(d.DNSDomain)),
		DNSDomainSearch: api.SelectedMapList(toList(d.DNSDomainSearch)),
		DNSServers:      api.SelectedMapList(toList(d.DNSServers)),
		NTPServers:      api.SelectedMapList(toList(d.NTPServers)),
		WINSServers:     api.SelectedMapList(toList(d.WINSServers)),
	}, nil
}

func convertClientOverwriteStructToSchema(d *openvpn.ClientOverwrite) (*clientOverwriteResourceModel, error) {
	return &clientOverwriteResourceModel{
		Enabled:         types.BoolValue(tools.StringToBool(d.Enabled)),
		CommonName:      types.StringValue(d.CommonName),
		Description:     types.StringValue(d.Description),
		Servers:         tools.StringSliceToSet(d.Servers),
		Block:           types.BoolValue(tools.StringToBool(d.Block)),
		PushReset:       types.BoolValue(tools.StringToBool(d.PushReset)),
		TunnelNetwork:   types.StringValue(d.TunnelNetwork),
		TunnelNetworkV6: types.StringValue(d.TunnelNetworkV6),
		LocalNetworks:   tools.StringSliceToSet(d.LocalNetworks),
		RemoteNetworks:  tools.StringSliceToSet(d.RemoteNetworks),
		RouteGateway:    types.StringValue(d.RouteGateway),
		RedirectGateway: tools.StringSliceToSet(d.RedirectGateway),
		RegisterDNS:     types.BoolValue(tools.StringToBool(d.RegisterDNS)),
		DNSDomain:       tools.StringSliceToSet(d.DNSDomain),
		DNSDomainSearch: tools.StringSliceToSet(d.DNSDomainSearch),
		DNSServers:      tools.StringSliceToSet(d.DNSServers),
		NTPServers:      tools.StringSliceToSet(d.NTPServers),
		WINSServers:     tools.StringSliceToSet(d.WINSServers),
	}, nil
}
