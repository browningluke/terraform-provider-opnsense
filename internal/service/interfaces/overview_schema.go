package interfaces

import (
	"context"

	"github.com/browningluke/opnsense-go/pkg/interfaces"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type overviewInterfaceModel struct {
	Identifier  types.String `tfsdk:"identifier"`
	Description types.String `tfsdk:"description"`
	Device      types.String `tfsdk:"device"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	LinkType    types.String `tfsdk:"link_type"`

	Addr4     types.String `tfsdk:"addr4"`
	Addr6     types.String `tfsdk:"addr6"`
	MacAddr   types.String `tfsdk:"macaddr"`
	MacAddrHw types.String `tfsdk:"macaddr_hw"`
	MTU       types.String `tfsdk:"mtu"`
	Status    types.String `tfsdk:"status"`
	Media     types.String `tfsdk:"media"`
	MediaRaw  types.String `tfsdk:"media_raw"`

	IsPhysical  types.Bool   `tfsdk:"is_physical"`
	VLANTag     types.String `tfsdk:"vlan_tag"`
	LaggProto   types.String `tfsdk:"lagg_proto"`
	LaggHash    types.String `tfsdk:"lagg_hash"`
	LaggOptions types.String `tfsdk:"lagg_options"`

	Flags          types.Set `tfsdk:"flags"`
	Capabilities   types.Set `tfsdk:"capabilities"`
	Options        types.Set `tfsdk:"options"`
	SupportedMedia types.Set `tfsdk:"supported_media"`
	Groups         types.Set `tfsdk:"groups"`
	Gateways       types.Set `tfsdk:"gateways"`
	Routes         types.Set `tfsdk:"routes"`
	ND6Flags       types.Set `tfsdk:"nd6_flags"`

	IfctlNameserver   types.Set `tfsdk:"ifctl_nameserver"`
	IfctlRouter       types.Set `tfsdk:"ifctl_router"`
	IfctlPrefix       types.Set `tfsdk:"ifctl_prefix"`
	IfctlSearchdomain types.Set `tfsdk:"ifctl_searchdomain"`

	IPv4 types.List   `tfsdk:"ipv4"`
	IPv6 types.List   `tfsdk:"ipv6"`
	VLAN types.Object `tfsdk:"vlan"`
}

type overviewIPv4Model struct {
	IPAddr     types.String `tfsdk:"ipaddr"`
	Vhid       types.String `tfsdk:"vhid"`
	CarpStatus types.String `tfsdk:"carp_status"`
	AdvBase    types.String `tfsdk:"adv_base"`
	AdvSkew    types.String `tfsdk:"adv_skew"`
	Peer       types.String `tfsdk:"peer"`
	Peer6      types.String `tfsdk:"peer6"`
}

type overviewIPv6Model struct {
	IPAddr     types.String `tfsdk:"ipaddr"`
	Vhid       types.String `tfsdk:"vhid"`
	CarpStatus types.String `tfsdk:"carp_status"`
	AdvBase    types.String `tfsdk:"adv_base"`
	AdvSkew    types.String `tfsdk:"adv_skew"`
	Peer       types.String `tfsdk:"peer"`
	Peer6      types.String `tfsdk:"peer6"`
}

type overviewVLANModel struct {
	Tag    types.String `tfsdk:"tag"`
	Proto  types.String `tfsdk:"proto"`
	PCP    types.String `tfsdk:"pcp"`
	Parent types.String `tfsdk:"parent"`
}

var overviewIPv4AttrTypes = map[string]attr.Type{
	"ipaddr":      types.StringType,
	"vhid":        types.StringType,
	"carp_status": types.StringType,
	"adv_base":    types.StringType,
	"adv_skew":    types.StringType,
	"peer":        types.StringType,
	"peer6":       types.StringType,
}

var overviewIPv6AttrTypes = map[string]attr.Type{
	"ipaddr":      types.StringType,
	"vhid":        types.StringType,
	"carp_status": types.StringType,
	"adv_base":    types.StringType,
	"adv_skew":    types.StringType,
	"peer":        types.StringType,
	"peer6":       types.StringType,
}

var overviewVLANAttrTypes = map[string]attr.Type{
	"tag":    types.StringType,
	"proto":  types.StringType,
	"pcp":    types.StringType,
	"parent": types.StringType,
}

var overviewInterfaceAttrTypes = map[string]attr.Type{
	"identifier":         types.StringType,
	"description":        types.StringType,
	"device":             types.StringType,
	"enabled":            types.BoolType,
	"link_type":          types.StringType,
	"addr4":              types.StringType,
	"addr6":              types.StringType,
	"macaddr":            types.StringType,
	"macaddr_hw":         types.StringType,
	"mtu":                types.StringType,
	"status":             types.StringType,
	"media":              types.StringType,
	"media_raw":          types.StringType,
	"is_physical":        types.BoolType,
	"vlan_tag":           types.StringType,
	"lagg_proto":         types.StringType,
	"lagg_hash":          types.StringType,
	"lagg_options":       types.StringType,
	"flags":              types.SetType{ElemType: types.StringType},
	"capabilities":       types.SetType{ElemType: types.StringType},
	"options":            types.SetType{ElemType: types.StringType},
	"supported_media":    types.SetType{ElemType: types.StringType},
	"groups":             types.SetType{ElemType: types.StringType},
	"gateways":           types.SetType{ElemType: types.StringType},
	"routes":             types.SetType{ElemType: types.StringType},
	"nd6_flags":          types.SetType{ElemType: types.StringType},
	"ifctl_nameserver":   types.SetType{ElemType: types.StringType},
	"ifctl_router":       types.SetType{ElemType: types.StringType},
	"ifctl_prefix":       types.SetType{ElemType: types.StringType},
	"ifctl_searchdomain": types.SetType{ElemType: types.StringType},
	"ipv4":               types.ListType{ElemType: types.ObjectType{AttrTypes: overviewIPv4AttrTypes}},
	"ipv6":               types.ListType{ElemType: types.ObjectType{AttrTypes: overviewIPv6AttrTypes}},
	"vlan":               types.ObjectType{AttrTypes: overviewVLANAttrTypes},
}

var overviewIPv4Attrs = map[string]schema.Attribute{
	"ipaddr": schema.StringAttribute{
		MarkdownDescription: "IPv4 address.",
		Computed:            true,
	},
	"vhid": schema.StringAttribute{
		MarkdownDescription: "CARP virtual host ID.",
		Computed:            true,
	},
	"carp_status": schema.StringAttribute{
		MarkdownDescription: "CARP status (e.g. `\"MASTER\"`, `\"BACKUP\"`).",
		Computed:            true,
	},
	"adv_base": schema.StringAttribute{
		MarkdownDescription: "CARP advertisement base.",
		Computed:            true,
	},
	"adv_skew": schema.StringAttribute{
		MarkdownDescription: "CARP advertisement skew.",
		Computed:            true,
	},
	"peer": schema.StringAttribute{
		MarkdownDescription: "CARP peer address.",
		Computed:            true,
	},
	"peer6": schema.StringAttribute{
		MarkdownDescription: "CARP peer IPv6 address.",
		Computed:            true,
	},
}

var overviewIPv6Attrs = map[string]schema.Attribute{
	"ipaddr": schema.StringAttribute{
		MarkdownDescription: "IPv6 address.",
		Computed:            true,
	},
	"vhid": schema.StringAttribute{
		MarkdownDescription: "CARP virtual host ID.",
		Computed:            true,
	},
	"carp_status": schema.StringAttribute{
		MarkdownDescription: "CARP status (e.g. `\"MASTER\"`, `\"BACKUP\"`).",
		Computed:            true,
	},
	"adv_base": schema.StringAttribute{
		MarkdownDescription: "CARP advertisement base.",
		Computed:            true,
	},
	"adv_skew": schema.StringAttribute{
		MarkdownDescription: "CARP advertisement skew.",
		Computed:            true,
	},
	"peer": schema.StringAttribute{
		MarkdownDescription: "CARP peer address.",
		Computed:            true,
	},
	"peer6": schema.StringAttribute{
		MarkdownDescription: "CARP peer IPv6 address.",
		Computed:            true,
	},
}

func overviewInterfaceDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Interfaces Overview provides live state of an OPNsense interface, including IP addresses, CARP status, VLAN, and LAGG information.",

		Attributes: map[string]schema.Attribute{
			"device": schema.StringAttribute{
				MarkdownDescription: "Kernel device name of the interface (e.g. `\"vtnet0\"`).",
				Required:            true,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "OPNsense logical interface identifier (e.g. `\"wan\"`, `\"lan\"`).",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "User-facing interface description.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the interface is enabled.",
				Computed:            true,
			},
			"link_type": schema.StringAttribute{
				MarkdownDescription: "Link type of the interface (e.g. `\"ether\"`, `\"dhcp\"`).",
				Computed:            true,
			},
			"addr4": schema.StringAttribute{
				MarkdownDescription: "Primary IPv4 address of the interface.",
				Computed:            true,
			},
			"addr6": schema.StringAttribute{
				MarkdownDescription: "Primary IPv6 address of the interface.",
				Computed:            true,
			},
			"macaddr": schema.StringAttribute{
				MarkdownDescription: "Current MAC address of the interface.",
				Computed:            true,
			},
			"macaddr_hw": schema.StringAttribute{
				MarkdownDescription: "Hardware MAC address of the interface.",
				Computed:            true,
			},
			"mtu": schema.StringAttribute{
				MarkdownDescription: "Maximum Transmission Unit of the interface.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Link status of the interface (e.g. `\"active\"`).",
				Computed:            true,
			},
			"media": schema.StringAttribute{
				MarkdownDescription: "Interface media type settings.",
				Computed:            true,
			},
			"media_raw": schema.StringAttribute{
				MarkdownDescription: "User-friendly interface media type.",
				Computed:            true,
			},
			"is_physical": schema.BoolAttribute{
				MarkdownDescription: "Whether the interface is a physical interface.",
				Computed:            true,
			},
			"vlan_tag": schema.StringAttribute{
				MarkdownDescription: "VLAN tag of the interface.",
				Computed:            true,
			},
			"lagg_proto": schema.StringAttribute{
				MarkdownDescription: "LAGG aggregation protocol (e.g. `\"lacp\"`).",
				Computed:            true,
			},
			"lagg_hash": schema.StringAttribute{
				MarkdownDescription: "LAGG hash configuration.",
				Computed:            true,
			},
			"lagg_options": schema.StringAttribute{
				MarkdownDescription: "LAGG options.",
				Computed:            true,
			},
			"flags": schema.SetAttribute{
				MarkdownDescription: "Interface flags.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"capabilities": schema.SetAttribute{
				MarkdownDescription: "Interface capabilities.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"options": schema.SetAttribute{
				MarkdownDescription: "Interface options.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"supported_media": schema.SetAttribute{
				MarkdownDescription: "Supported media types.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"groups": schema.SetAttribute{
				MarkdownDescription: "Interface groups.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"gateways": schema.SetAttribute{
				MarkdownDescription: "Gateways associated with the interface.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"routes": schema.SetAttribute{
				MarkdownDescription: "Routes associated with the interface.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"nd6_flags": schema.SetAttribute{
				MarkdownDescription: "ND6 flags of the interface.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ifctl_nameserver": schema.SetAttribute{
				MarkdownDescription: "Name servers assigned to the interface via ifctl.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ifctl_router": schema.SetAttribute{
				MarkdownDescription: "Routers assigned to the interface via ifctl.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ifctl_prefix": schema.SetAttribute{
				MarkdownDescription: "Prefixes assigned to the interface via ifctl.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ifctl_searchdomain": schema.SetAttribute{
				MarkdownDescription: "Search domains assigned to the interface via ifctl.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ipv4": schema.ListNestedAttribute{
				MarkdownDescription: "IPv4 addresses assigned to the interface, including CARP virtual IPs.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: overviewIPv4Attrs,
				},
			},
			"ipv6": schema.ListNestedAttribute{
				MarkdownDescription: "IPv6 addresses assigned to the interface, including CARP virtual IPs.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: overviewIPv6Attrs,
				},
			},
			"vlan": schema.SingleNestedAttribute{
				MarkdownDescription: "VLAN configuration of the interface.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"tag": schema.StringAttribute{
						MarkdownDescription: "VLAN tag.",
						Computed:            true,
					},
					"proto": schema.StringAttribute{
						MarkdownDescription: "VLAN protocol.",
						Computed:            true,
					},
					"pcp": schema.StringAttribute{
						MarkdownDescription: "VLAN priority code point.",
						Computed:            true,
					},
					"parent": schema.StringAttribute{
						MarkdownDescription: "Parent interface of the VLAN.",
						Computed:            true,
					},
				},
			},
		},
	}
}

func convertOverviewInterfaceStructToSchema(d *interfaces.InterfaceInfo) (*overviewInterfaceModel, error) {
	model := &overviewInterfaceModel{
		Identifier:        types.StringValue(d.Identifier),
		Description:       types.StringValue(d.Description),
		Device:            types.StringValue(d.Device),
		Enabled:           types.BoolValue(d.Enabled),
		LinkType:          types.StringValue(d.LinkType),
		Addr4:             types.StringValue(d.Addr4),
		Addr6:             types.StringValue(d.Addr6),
		MacAddr:           types.StringValue(d.MacAddr),
		MacAddrHw:         types.StringValue(d.MacAddrHw),
		MTU:               types.StringValue(d.MTU),
		Status:            types.StringValue(d.Status),
		Media:             types.StringValue(d.Media),
		MediaRaw:          types.StringValue(d.MediaRaw),
		IsPhysical:        types.BoolValue(d.IsPhysical),
		VLANTag:           types.StringValue(d.VLANTag),
		LaggProto:         types.StringValue(d.LaggProto),
		LaggHash:          types.StringValue(d.LaggHash),
		LaggOptions:       types.StringValue(d.LaggOptions),
		Flags:             tools.StringSliceToSet(d.Flags),
		Capabilities:      tools.StringSliceToSet(d.Capabilities),
		Options:           tools.StringSliceToSet(d.Options),
		SupportedMedia:    tools.StringSliceToSet(d.SupportedMedia),
		Groups:            tools.StringSliceToSet(d.Groups),
		Gateways:          tools.StringSliceToSet(d.Gateways),
		Routes:            tools.StringSliceToSet(d.Routes),
		ND6Flags:          tools.StringSliceToSet(d.ND6.Flags),
		IfctlNameserver:   tools.StringSliceToSet(d.IfctlNameserver),
		IfctlRouter:       tools.StringSliceToSet(d.IfctlRouter),
		IfctlPrefix:       tools.StringSliceToSet(d.IfctlPrefix),
		IfctlSearchdomain: tools.StringSliceToSet(d.IfctlSearchdomain),
	}

	ipv4s := []overviewIPv4Model{}
	for _, ip := range d.IPv4 {
		ipv4s = append(ipv4s, overviewIPv4Model{
			IPAddr:     types.StringValue(ip.IPAddr),
			Vhid:       types.StringValue(ip.Vhid),
			CarpStatus: types.StringValue(ip.CarpStatus),
			AdvBase:    types.StringValue(ip.AdvBase),
			AdvSkew:    types.StringValue(ip.AdvSkew),
			Peer:       types.StringValue(ip.Peer),
			Peer6:      types.StringValue(ip.Peer6),
		})
	}

	ipv6s := []overviewIPv6Model{}
	for _, ip := range d.IPv6 {
		ipv6s = append(ipv6s, overviewIPv6Model{
			IPAddr:     types.StringValue(ip.IPAddr),
			Vhid:       types.StringValue(ip.Vhid),
			CarpStatus: types.StringValue(ip.CarpStatus),
			AdvBase:    types.StringValue(ip.AdvBase),
			AdvSkew:    types.StringValue(ip.AdvSkew),
			Peer:       types.StringValue(ip.Peer),
			Peer6:      types.StringValue(ip.Peer6),
		})
	}

	model.IPv4, _ = types.ListValueFrom(
		context.Background(),
		types.ObjectType{}.WithAttributeTypes(overviewIPv4AttrTypes),
		ipv4s,
	)

	model.IPv6, _ = types.ListValueFrom(
		context.Background(),
		types.ObjectType{}.WithAttributeTypes(overviewIPv6AttrTypes),
		ipv6s,
	)

	model.VLAN, _ = types.ObjectValueFrom(
		context.Background(),
		overviewVLANAttrTypes,
		overviewVLANModel{
			Tag:    types.StringValue(d.VLAN.Tag),
			Proto:  types.StringValue(d.VLAN.Proto),
			PCP:    types.StringValue(d.VLAN.PCP),
			Parent: types.StringValue(d.VLAN.Parent),
		},
	)

	return model, nil
}
