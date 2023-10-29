package service

import (
	"context"
	"github.com/browningluke/opnsense-go/pkg/diagnostics"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/tools"
)

type InterfaceDataSourceModel struct {
	Device     types.String `tfsdk:"device"`
	Media      types.String `tfsdk:"media"`
	MediaRaw   types.String `tfsdk:"media_raw"`
	MacAddr    types.String `tfsdk:"macaddr"`
	IsPhysical types.Bool   `tfsdk:"is_physical"`
	MTU        types.Int64  `tfsdk:"mtu"`
	Status     types.String `tfsdk:"status"`

	Flags          types.Set `tfsdk:"flags"`
	Capabilities   types.Set `tfsdk:"capabilities"`
	Options        types.Set `tfsdk:"options"`
	SupportedMedia types.Set `tfsdk:"supported_media"`
	Groups         types.Set `tfsdk:"groups"`

	Ipv4 types.List `tfsdk:"ipv4"`
	Ipv6 types.List `tfsdk:"ipv6"`
}

type Ipv4Model struct {
	Ipaddr     types.String `tfsdk:"ipaddr"`
	SubnetBits types.Int64  `tfsdk:"subnetbits"`
	Tunnel     types.Bool   `tfsdk:"tunnel"`
}

type Ipv6Model struct {
	Ipaddr     types.String `tfsdk:"ipaddr"`
	SubnetBits types.Int64  `tfsdk:"subnetbits"`
	Tunnel     types.Bool   `tfsdk:"tunnel"`
	Autoconf   types.Bool   `tfsdk:"autoconf"`
	Deprecated types.Bool   `tfsdk:"deprecated"`
	LinkLocal  types.Bool   `tfsdk:"link_local"`
	Tentative  types.Bool   `tfsdk:"tentative"`
}

var ipv4AttrTypes = map[string]attr.Type{
	"ipaddr":     types.StringType,
	"subnetbits": types.Int64Type,
	"tunnel":     types.BoolType,
}

var ipv6AttrTypes = map[string]attr.Type{
	"ipaddr":     types.StringType,
	"subnetbits": types.Int64Type,
	"tunnel":     types.BoolType,
	"autoconf":   types.BoolType,
	"deprecated": types.BoolType,
	"link_local": types.BoolType,
	"tentative":  types.BoolType,
}

func InterfaceDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Interfaces can be used to get configurations of OPNsense interfaces.",

		Attributes: map[string]schema.Attribute{
			"device": schema.StringAttribute{
				MarkdownDescription: "Name of the interface device.",
				Required:            true,
			},
			"media": schema.StringAttribute{
				MarkdownDescription: "Interface media type settings (see https://man.openbsd.org/ifmedia.4).",
				Computed:            true,
			},
			"media_raw": schema.StringAttribute{
				MarkdownDescription: "User-friendly interface media type.",
				Computed:            true,
			},
			"macaddr": schema.StringAttribute{
				MarkdownDescription: "MAC address assigned to the interface.",
				Computed:            true,
			},
			"is_physical": schema.BoolAttribute{
				MarkdownDescription: "Whether the interface is physical or virtual.",
				Computed:            true,
			},
			"mtu": schema.Int64Attribute{
				MarkdownDescription: "Maximum Transmission Unit for the interface. This is typically 1500 bytes but can vary in some circumstances.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Status of the interface (e.g. `\"active\"`).",
				Computed:            true,
			},
			"flags": schema.SetAttribute{
				MarkdownDescription: "List of flags configured on the interface (equiv. to flags=xxxx in output of ifconfig).",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"capabilities": schema.SetAttribute{
				MarkdownDescription: "List of capabilities the interface supports.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"options": schema.SetAttribute{
				MarkdownDescription: "List of options configured on the interface (equiv. to options=xx in output of ifconfig).",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"supported_media": schema.SetAttribute{
				MarkdownDescription: "List of supported media type settings (see https://man.openbsd.org/ifmedia.4).",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"groups": schema.SetAttribute{
				MarkdownDescription: "List of groups the interface is a member of.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ipv4": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ipaddr": schema.StringAttribute{
							MarkdownDescription: "IPv4 address assigned to the interface.",
							Computed:            true,
						},
						"subnetbits": schema.Int64Attribute{
							MarkdownDescription: "Number of subnet bits (i.e. CIDR).",
							Computed:            true,
						},
						"tunnel": schema.BoolAttribute{
							MarkdownDescription: "Whether IPv4 tunnelling is enabled.",
							Computed:            true,
						},
					},
				}},
			"ipv6": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ipaddr": schema.StringAttribute{
							MarkdownDescription: "IPv6 address assigned to the interface.",
							Computed:            true,
						},
						"subnetbits": schema.Int64Attribute{
							MarkdownDescription: "Number of subnet bits (i.e. CIDR).",
							Computed:            true,
						},
						"tunnel": schema.BoolAttribute{
							MarkdownDescription: "Whether IPv6 tunnelling is enabled.",
							Computed:            true,
						},
						"autoconf": schema.BoolAttribute{
							MarkdownDescription: "Whether auto-configuration is enabled for the address.",
							Computed:            true,
						},
						"deprecated": schema.BoolAttribute{
							MarkdownDescription: "Whether the address is deprecated.",
							Computed:            true,
						},
						"link_local": schema.BoolAttribute{
							MarkdownDescription: "Whether the address is link-local.",
							Computed:            true,
						},
						"tentative": schema.BoolAttribute{
							MarkdownDescription: "Whether the address is tentative.",
							Computed:            true,
						},
					},
				}},
		},
	}
}

func convertInterfaceConfigStructToSchema(d *diagnostics.Interface) (*InterfaceDataSourceModel, error) {
	model := &InterfaceDataSourceModel{
		Device:         types.StringValue(d.Device),
		Media:          types.StringValue(d.Media),
		MediaRaw:       types.StringValue(d.MediaRaw),
		MacAddr:        types.StringValue(d.MacAddr),
		IsPhysical:     types.BoolValue(d.IsPhysical),
		MTU:            tools.StringToInt64Null(d.MTU),
		Status:         types.StringValue(d.Status),
		Flags:          tools.StringSliceToSet(d.Flags),
		Capabilities:   tools.StringSliceToSet(d.Capabilities),
		Options:        tools.StringSliceToSet(d.Options),
		SupportedMedia: tools.StringSliceToSet(d.SupportedMedia),
		Groups:         tools.StringSliceToSet(d.Groups),
	}

	// Creating an empty slice results in `[]` rather than `null` if OPNsense API returned an empty list.
	ipv4s := []Ipv4Model{}
	for _, elem := range d.Ipv4 {
		ipv4s = append(ipv4s, Ipv4Model{
			Ipaddr:     types.StringValue(elem.IpAddr),
			SubnetBits: types.Int64Value(elem.SubnetBits),
			Tunnel:     types.BoolValue(elem.Tunnel),
		})
	}

	ipv6s := []Ipv6Model{}
	for _, elem := range d.Ipv6 {
		ipv6s = append(ipv6s, Ipv6Model{
			Ipaddr:     types.StringValue(elem.IpAddr),
			SubnetBits: types.Int64Value(elem.SubnetBits),
			Tunnel:     types.BoolValue(elem.Tunnel),
			Autoconf:   types.BoolValue(elem.Autoconf),
			Deprecated: types.BoolValue(elem.Deprecated),
			LinkLocal:  types.BoolValue(elem.LinkLocal),
			Tentative:  types.BoolValue(elem.Tentative),
		})
	}

	model.Ipv4, _ = types.ListValueFrom(
		context.Background(),
		types.ObjectType{}.WithAttributeTypes(ipv4AttrTypes),
		ipv4s,
	)

	model.Ipv6, _ = types.ListValueFrom(
		context.Background(),
		types.ObjectType{}.WithAttributeTypes(ipv6AttrTypes),
		ipv6s,
	)

	return model, nil
}
