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

// ospf6NetworkResourceModel describes the resource data model.
type ospf6NetworkResourceModel struct {
	Enabled       types.Bool   `tfsdk:"enabled"`
	IPAddr        types.String `tfsdk:"ip_addr"`
	NetMask       types.String `tfsdk:"net_mask"`
	Area          types.String `tfsdk:"area"`
	AreaRange     types.String `tfsdk:"area_range"`
	PrefixListIn  types.String `tfsdk:"prefix_list_in"`
	PrefixListOut types.String `tfsdk:"prefix_list_out"`

	Id types.String `tfsdk:"id"`
}

func ospf6NetworkResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure OSPFv3 networks.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this OSPFv3 network. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"ip_addr": schema.StringAttribute{
				MarkdownDescription: "The IPv6 network address.",
				Required:            true,
			},
			"net_mask": schema.StringAttribute{
				MarkdownDescription: "The prefix length (0-128).",
				Required:            true,
			},
			"area": schema.StringAttribute{
				MarkdownDescription: "The OSPFv3 area.",
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
				MarkdownDescription: "UUID of the OSPFv3 network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func ospf6NetworkDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure OSPFv3 networks.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this OSPFv3 network.",
				Computed:            true,
			},
			"ip_addr": dschema.StringAttribute{
				MarkdownDescription: "The IPv6 network address.",
				Computed:            true,
			},
			"net_mask": dschema.StringAttribute{
				MarkdownDescription: "The prefix length (0-128).",
				Computed:            true,
			},
			"area": dschema.StringAttribute{
				MarkdownDescription: "The OSPFv3 area.",
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

func convertOSPF6NetworkSchemaToStruct(d *ospf6NetworkResourceModel) (*quagga.OSPF6Network, error) {
	return &quagga.OSPF6Network{
		Enabled:             tools.BoolToString(d.Enabled.ValueBool()),
		IPAddr:              d.IPAddr.ValueString(),
		NetMask:             d.NetMask.ValueString(),
		Area:                d.Area.ValueString(),
		AreaRange:           d.AreaRange.ValueString(),
		LinkedPrefixListIn:  api.SelectedMap(d.PrefixListIn.ValueString()),
		LinkedPrefixListOut: api.SelectedMap(d.PrefixListOut.ValueString()),
	}, nil
}

func convertOSPF6NetworkStructToSchema(d *quagga.OSPF6Network) (*ospf6NetworkResourceModel, error) {
	return &ospf6NetworkResourceModel{
		Enabled:       types.BoolValue(tools.StringToBool(d.Enabled)),
		IPAddr:        types.StringValue(d.IPAddr),
		NetMask:       types.StringValue(d.NetMask),
		Area:          types.StringValue(d.Area),
		AreaRange:     types.StringValue(d.AreaRange),
		PrefixListIn:  types.StringValue(d.LinkedPrefixListIn.String()),
		PrefixListOut: types.StringValue(d.LinkedPrefixListOut.String()),
	}, nil
}
