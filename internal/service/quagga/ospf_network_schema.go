package quagga

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ospfNetworkResourceModel describes the resource data model.
type ospfNetworkResourceModel struct {
	Enabled        types.Bool   `tfsdk:"enabled"`
	IPAddr         types.String `tfsdk:"ip_addr"`
	Area           types.String `tfsdk:"area"`
	NetMask        types.String `tfsdk:"net_mask"`
	AreaRange      types.String `tfsdk:"area_range"`
	PrefixListIn   types.String `tfsdk:"prefix_list_in"`
	PrefixListOut  types.String `tfsdk:"prefix_list_out"`

	Id types.String `tfsdk:"id"`
}

func ospfNetworkResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure OSPF networks.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this OSPF network. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"ip_addr": schema.StringAttribute{
				MarkdownDescription: "The network IP address.",
				Required:            true,
			},
			"area": schema.StringAttribute{
				MarkdownDescription: "The OSPF area in IPv4 dotted notation (e.g. `0.0.0.0`).",
				Required:            true,
			},
			"net_mask": schema.StringAttribute{
				MarkdownDescription: "The network mask.",
				Required:            true,
			},
			"area_range": schema.StringAttribute{
				MarkdownDescription: "Area range for summarization. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"prefix_list_in": schema.StringAttribute{
				MarkdownDescription: "UUID of the prefix list for inbound direction. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"prefix_list_out": schema.StringAttribute{
				MarkdownDescription: "UUID of the prefix list for outbound direction. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the OSPF network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func ospfNetworkDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure OSPF networks.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this OSPF network.",
				Computed:            true,
			},
			"ip_addr": dschema.StringAttribute{
				MarkdownDescription: "The network IP address.",
				Computed:            true,
			},
			"area": dschema.StringAttribute{
				MarkdownDescription: "The OSPF area in IPv4 dotted notation.",
				Computed:            true,
			},
			"net_mask": dschema.StringAttribute{
				MarkdownDescription: "The network mask.",
				Computed:            true,
			},
			"area_range": dschema.StringAttribute{
				MarkdownDescription: "Area range for summarization.",
				Computed:            true,
			},
			"prefix_list_in": dschema.StringAttribute{
				MarkdownDescription: "UUID of the prefix list for inbound direction.",
				Computed:            true,
			},
			"prefix_list_out": dschema.StringAttribute{
				MarkdownDescription: "UUID of the prefix list for outbound direction.",
				Computed:            true,
			},
		},
	}
}

func convertOSPFNetworkSchemaToStruct(d *ospfNetworkResourceModel) (*quagga.OSPFNetwork, error) {
	return &quagga.OSPFNetwork{
		Enabled:             tools.BoolToString(d.Enabled.ValueBool()),
		IPAddr:              d.IPAddr.ValueString(),
		Area:                d.Area.ValueString(),
		NetMask:             d.NetMask.ValueString(),
		AreaRange:           d.AreaRange.ValueString(),
		LinkedPrefixListIn:  api.SelectedMap(d.PrefixListIn.ValueString()),
		LinkedPrefixListOut: api.SelectedMap(d.PrefixListOut.ValueString()),
	}, nil
}

func convertOSPFNetworkStructToSchema(d *quagga.OSPFNetwork) (*ospfNetworkResourceModel, error) {
	return &ospfNetworkResourceModel{
		Enabled:       types.BoolValue(tools.StringToBool(d.Enabled)),
		IPAddr:        types.StringValue(d.IPAddr),
		Area:          types.StringValue(d.Area),
		NetMask:       types.StringValue(d.NetMask),
		AreaRange:     types.StringValue(d.AreaRange),
		PrefixListIn:  types.StringValue(d.LinkedPrefixListIn.String()),
		PrefixListOut: types.StringValue(d.LinkedPrefixListOut.String()),
	}, nil
}
