package service

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/tools"
)

// QuaggaBGPNeighborResourceModel describes the resource data model.
type QuaggaBGPNeighborResourceModel struct {
	Enabled               types.Bool   `tfsdk:"enabled"`
	Description           types.String `tfsdk:"description"`
	PeerIP                types.String `tfsdk:"peer_ip"`
	RemoteAS              types.Int64  `tfsdk:"remote_as"`
	Password              types.String `tfsdk:"md5_password"`
	Weight                types.Int64  `tfsdk:"weight"`
	LocalIP               types.String `tfsdk:"local_ip"`
	UpdateSource          types.String `tfsdk:"update_source"`
	LinkLocalInterface    types.String `tfsdk:"link_local_interface"`
	NextHopSelf           types.Bool   `tfsdk:"next_hop_self"`
	NextHopSelfAll        types.Bool   `tfsdk:"next_hop_self_all"`
	MultiHop              types.Bool   `tfsdk:"multi_hop"`
	MultiProtocol         types.Bool   `tfsdk:"multi_protocol"`
	RRClient              types.Bool   `tfsdk:"rr_client"`
	BFD                   types.Bool   `tfsdk:"bfd"`
	KeepAlive             types.Int64  `tfsdk:"keep_alive"`
	HoldDown              types.Int64  `tfsdk:"hold_down"`
	ConnectTimer          types.Int64  `tfsdk:"connect_timer"`
	DefaultRoute          types.Bool   `tfsdk:"default_route"`
	ASOverride            types.Bool   `tfsdk:"as_override"`
	DisableConnectedCheck types.Bool   `tfsdk:"disable_connected_check"`
	AttributeUnchanged    types.String `tfsdk:"attribute_unchanged"`
	PrefixListIn          types.String `tfsdk:"prefix_list_in"`
	PrefixListOut         types.String `tfsdk:"prefix_list_out"`
	RouteMapIn            types.String `tfsdk:"route_map_in"`
	RouteMapOut           types.String `tfsdk:"route_map_out"`

	Id types.String `tfsdk:"id"`
}

func quaggaBGPNeighborResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure neighbors for BGP.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this neighbor. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "An optional description for this neighbor. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"peer_ip": schema.StringAttribute{
				MarkdownDescription: "The IP of your neighbor.",
				Required:            true,
			},
			"remote_as": schema.Int64Attribute{
				MarkdownDescription: "The neighbor AS.",
				Required:            true,
			},
			"md5_password": schema.StringAttribute{
				MarkdownDescription: "The password for BGP authentication. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"weight": schema.Int64Attribute{
				MarkdownDescription: "Specify a default weight value for the neighbor’s routes. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"local_ip": schema.StringAttribute{
				MarkdownDescription: "The local IP connecting to the neighbor. This is only required for BGP authentication. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"update_source": schema.StringAttribute{
				MarkdownDescription: "Physical name of the IPv4 interface facing the peer. Must be a valid OPNsense interface in lowercase (e.g. `wan`). Please refer to the FRR documentation for more information. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"link_local_interface": schema.StringAttribute{
				MarkdownDescription: "Interface to use for IPv6 link-local neighbours. Must be a valid OPNsense interface in lowercase (e.g. `wan`). Please refer to the FRR documentation for more information. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"next_hop_self": schema.BoolAttribute{
				MarkdownDescription: "Enable the next-hop-self command. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"next_hop_self_all": schema.BoolAttribute{
				MarkdownDescription: "Add the parameter \"all\" after next-hop-self command. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"multi_hop": schema.BoolAttribute{
				MarkdownDescription: "Enable multi-hop. Specifying ebgp-multihop allows sessions with eBGP neighbors to establish when they are multiple hops away. When the neighbor is not directly connected and this knob is not enabled, the session will not establish. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"multi_protocol": schema.BoolAttribute{
				MarkdownDescription: "Mark this neighbor as multiprotocol capable per RFC 2283. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"rr_client": schema.BoolAttribute{
				MarkdownDescription: "Enable route reflector client. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"bfd": schema.BoolAttribute{
				MarkdownDescription: "Enable BFD support for this neighbor. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"keep_alive": schema.Int64Attribute{
				MarkdownDescription: "Enable Keepalive timer to check if the neighbor is still up. Defaults to `60`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(60),
			},
			"hold_down": schema.Int64Attribute{
				MarkdownDescription: "The time in seconds when a neighbor is considered dead. This is usually 3 times the keepalive timer. Defaults to `180`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(180),
			},
			"connect_timer": schema.Int64Attribute{
				MarkdownDescription: "The time in seconds how fast a neighbor tries to reconnect. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"default_route": schema.BoolAttribute{
				MarkdownDescription: "Enable to send Defaultroute. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"as_override": schema.BoolAttribute{
				MarkdownDescription: "Override AS number of the originating router with the local AS number. This command is only allowed for eBGP peers. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"disable_connected_check": schema.BoolAttribute{
				MarkdownDescription: "Enable to allow peerings between directly connected eBGP peers using loopback addresses. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"attribute_unchanged": schema.StringAttribute{
				MarkdownDescription: "Specify attribute to be left unchanged when sending advertisements to a peer. Read more at FRR documentation. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf("", "as-path", "next-hop", "med"),
				},
			},
			"prefix_list_in": schema.StringAttribute{
				MarkdownDescription: "The prefix list ID for inbound direction. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"prefix_list_out": schema.StringAttribute{
				MarkdownDescription: "The prefix list ID for outbound direction. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"route_map_in": schema.StringAttribute{
				MarkdownDescription: "The route map ID for inbound direction. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"route_map_out": schema.StringAttribute{
				MarkdownDescription: "The route map ID for outbound direction. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the neighbor.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func QuaggaBGPNeighborDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure neighbors for BGP.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this neighbor.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "An optional description for this neighbor.",
				Computed:            true,
			},
			"peer_ip": dschema.StringAttribute{
				MarkdownDescription: "The IP of your neighbor.",
				Computed:            true,
			},
			"remote_as": dschema.Int64Attribute{
				MarkdownDescription: "The neighbor AS.",
				Computed:            true,
			},
			"md5_password": dschema.StringAttribute{
				MarkdownDescription: "The password for BGP authentication.",
				Computed:            true,
			},
			"weight": dschema.Int64Attribute{
				MarkdownDescription: "Specify a default weight value for the neighbor’s routes.",
				Computed:            true,
			},
			"local_ip": dschema.StringAttribute{
				MarkdownDescription: "The local IP connecting to the neighbor. This is only required for BGP authentication.",
				Computed:            true,
			},
			"update_source": dschema.StringAttribute{
				MarkdownDescription: "Physical name of the IPv4 interface facing the peer. Must be a valid OPNsense interface in lowercase (e.g. `wan`). Please refer to the FRR documentation for more information.",
				Computed:            true,
			},
			"link_local_interface": dschema.StringAttribute{
				MarkdownDescription: "Interface to use for IPv6 link-local neighbours. Must be a valid OPNsense interface in lowercase (e.g. `wan`). Please refer to the FRR documentation for more information.",
				Computed:            true,
			},
			"next_hop_self": dschema.BoolAttribute{
				MarkdownDescription: "Enable the next-hop-self command.",
				Computed:            true,
			},
			"next_hop_self_all": dschema.BoolAttribute{
				MarkdownDescription: "Add the parameter \"all\" after next-hop-self command.",
				Computed:            true,
			},
			"multi_hop": dschema.BoolAttribute{
				MarkdownDescription: "Enable multi-hop. Specifying ebgp-multihop allows sessions with eBGP neighbors to establish when they are multiple hops away. When the neighbor is not directly connected and this knob is not enabled, the session will not establish.",
				Computed:            true,
			},
			"multi_protocol": dschema.BoolAttribute{
				MarkdownDescription: "Mark this neighbor as multiprotocol capable per RFC 2283.",
				Computed:            true,
			},
			"rr_client": dschema.BoolAttribute{
				MarkdownDescription: "Enable route reflector client.",
				Computed:            true,
			},
			"bfd": dschema.BoolAttribute{
				MarkdownDescription: "Enable BFD support for this neighbor.",
				Computed:            true,
			},
			"keep_alive": dschema.Int64Attribute{
				MarkdownDescription: "Enable Keepalive timer to check if the neighbor is still up.",
				Computed:            true,
			},
			"hold_down": dschema.Int64Attribute{
				MarkdownDescription: "The time in seconds when a neighbor is considered dead. This is usually 3 times the keepalive timer.",
				Computed:            true,
			},
			"connect_timer": dschema.Int64Attribute{
				MarkdownDescription: "The time in seconds how fast a neighbor tries to reconnect.",
				Computed:            true,
			},
			"default_route": dschema.BoolAttribute{
				MarkdownDescription: "Enable to send Defaultroute.",
				Computed:            true,
			},
			"as_override": dschema.BoolAttribute{
				MarkdownDescription: "Override AS number of the originating router with the local AS number. This command is only allowed for eBGP peers.",
				Computed:            true,
			},
			"disable_connected_check": dschema.BoolAttribute{
				MarkdownDescription: "Enable to allow peerings between directly connected eBGP peers using loopback addresses.",
				Computed:            true,
			},
			"attribute_unchanged": dschema.StringAttribute{
				MarkdownDescription: "Specify attribute to be left unchanged when sending advertisements to a peer. Read more at FRR documentation.",
				Computed:            true,
			},
			"prefix_list_in": dschema.StringAttribute{
				MarkdownDescription: "The prefix list ID for inbound direction.",
				Computed:            true,
			},
			"prefix_list_out": dschema.StringAttribute{
				MarkdownDescription: "The prefix list ID for outbound direction.",
				Computed:            true,
			},
			"route_map_in": dschema.StringAttribute{
				MarkdownDescription: "The route map ID for inbound direction.",
				Computed:            true,
			},
			"route_map_out": dschema.StringAttribute{
				MarkdownDescription: "The route map ID for outbound direction.",
				Computed:            true,
			},
		},
	}
}

func convertQuaggaBGPNeighborSchemaToStruct(d *QuaggaBGPNeighborResourceModel) (*quagga.BGPNeighbor, error) {
	return &quagga.BGPNeighbor{
		Enabled:               tools.BoolToString(d.Enabled.ValueBool()),
		Description:           d.Description.ValueString(),
		PeerIP:                d.PeerIP.ValueString(),
		RemoteAS:              tools.Int64ToString(d.RemoteAS.ValueInt64()),
		Password:              d.Password.ValueString(),
		Weight:                tools.Int64ToStringNegative(d.Weight.ValueInt64()),
		LocalIP:               d.LocalIP.ValueString(),
		UpdateSource:          api.SelectedMap(d.UpdateSource.ValueString()),
		LinkLocalInterface:    api.SelectedMap(d.LinkLocalInterface.ValueString()),
		NextHopSelf:           tools.BoolToString(d.NextHopSelf.ValueBool()),
		NextHopSelfAll:        tools.BoolToString(d.NextHopSelfAll.ValueBool()),
		MultiHop:              tools.BoolToString(d.MultiHop.ValueBool()),
		MultiProtocol:         tools.BoolToString(d.MultiProtocol.ValueBool()),
		RRClient:              tools.BoolToString(d.RRClient.ValueBool()),
		BFD:                   tools.BoolToString(d.BFD.ValueBool()),
		KeepAlive:             tools.Int64ToString(d.KeepAlive.ValueInt64()),
		HoldDown:              tools.Int64ToString(d.HoldDown.ValueInt64()),
		ConnectTimer:          tools.Int64ToStringNegative(d.ConnectTimer.ValueInt64()),
		DefaultRoute:          tools.BoolToString(d.DefaultRoute.ValueBool()),
		ASOverride:            tools.BoolToString(d.ASOverride.ValueBool()),
		DisableConnectedCheck: tools.BoolToString(d.DisableConnectedCheck.ValueBool()),
		AttributeUnchanged:    api.SelectedMap(d.AttributeUnchanged.ValueString()),
		PrefixListIn:          api.SelectedMap(d.PrefixListIn.ValueString()),
		PrefixListOut:         api.SelectedMap(d.PrefixListOut.ValueString()),
		RouteMapIn:            api.SelectedMap(d.RouteMapIn.ValueString()),
		RouteMapOut:           api.SelectedMap(d.RouteMapOut.ValueString()),
	}, nil
}

func convertQuaggaBGPNeighborStructToSchema(d *quagga.BGPNeighbor) (*QuaggaBGPNeighborResourceModel, error) {
	return &QuaggaBGPNeighborResourceModel{
		Enabled:               types.BoolValue(tools.StringToBool(d.Enabled)),
		Description:           types.StringValue(d.Description),
		PeerIP:                types.StringValue(d.PeerIP),
		RemoteAS:              types.Int64Value(tools.StringToInt64(d.RemoteAS)),
		Password:              types.StringValue(d.Password),
		Weight:                types.Int64Value(tools.StringToInt64(d.Weight)),
		LocalIP:               types.StringValue(d.LocalIP),
		UpdateSource:          types.StringValue(d.UpdateSource.String()),
		LinkLocalInterface:    types.StringValue(d.LinkLocalInterface.String()),
		NextHopSelf:           types.BoolValue(tools.StringToBool(d.NextHopSelf)),
		NextHopSelfAll:        types.BoolValue(tools.StringToBool(d.NextHopSelfAll)),
		MultiHop:              types.BoolValue(tools.StringToBool(d.MultiHop)),
		MultiProtocol:         types.BoolValue(tools.StringToBool(d.MultiProtocol)),
		RRClient:              types.BoolValue(tools.StringToBool(d.RRClient)),
		BFD:                   types.BoolValue(tools.StringToBool(d.BFD)),
		KeepAlive:             types.Int64Value(tools.StringToInt64(d.KeepAlive)),
		HoldDown:              types.Int64Value(tools.StringToInt64(d.HoldDown)),
		ConnectTimer:          types.Int64Value(tools.StringToInt64(d.ConnectTimer)),
		DefaultRoute:          types.BoolValue(tools.StringToBool(d.DefaultRoute)),
		ASOverride:            types.BoolValue(tools.StringToBool(d.ASOverride)),
		DisableConnectedCheck: types.BoolValue(tools.StringToBool(d.DisableConnectedCheck)),
		AttributeUnchanged:    types.StringValue(d.AttributeUnchanged.String()),
		PrefixListIn:          types.StringValue(d.PrefixListIn.String()),
		PrefixListOut:         types.StringValue(d.PrefixListOut.String()),
		RouteMapIn:            types.StringValue(d.RouteMapIn.String()),
		RouteMapOut:           types.StringValue(d.RouteMapOut.String()),
	}, nil
}
