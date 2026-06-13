package quagga

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// bfdNeighborResourceModel describes the resource data model.
type bfdNeighborResourceModel struct {
	Enabled           types.Bool   `tfsdk:"enabled"`
	Description       types.String `tfsdk:"description"`
	Address           types.String `tfsdk:"address"`
	MultiHop          types.Bool   `tfsdk:"multi_hop"`
	LocalAddress      types.String `tfsdk:"local_address"`
	InterfaceName     types.String `tfsdk:"interface_name"`
	DetectMultiplier  types.Int64  `tfsdk:"detect_multiplier"`
	ReceiveInterval   types.Int64  `tfsdk:"receive_interval"`
	TransmitInterval  types.Int64  `tfsdk:"transmit_interval"`

	Id types.String `tfsdk:"id"`
}

func bfdNeighborResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure BFD neighbors.",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this BFD neighbor. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "An optional description for this BFD neighbor. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "The IPv4 or IPv6 address of the BFD neighbor.",
				Required:            true,
			},
			"multi_hop": schema.BoolAttribute{
				MarkdownDescription: "Enable multi-hop BFD. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"local_address": schema.StringAttribute{
				MarkdownDescription: "The local address used for the BFD session. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"interface_name": schema.StringAttribute{
				MarkdownDescription: "The interface used for the BFD session. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"detect_multiplier": schema.Int64Attribute{
				MarkdownDescription: "The detection multiplier. Defaults to `3`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(3),
			},
			"receive_interval": schema.Int64Attribute{
				MarkdownDescription: "The minimum receive interval in milliseconds. Defaults to `300`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(300),
			},
			"transmit_interval": schema.Int64Attribute{
				MarkdownDescription: "The minimum transmit interval in milliseconds. Defaults to `300`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(300),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the BFD neighbor.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func bfdNeighborDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure BFD neighbors.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this BFD neighbor.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "An optional description for this BFD neighbor.",
				Computed:            true,
			},
			"address": dschema.StringAttribute{
				MarkdownDescription: "The IPv4 or IPv6 address of the BFD neighbor.",
				Computed:            true,
			},
			"multi_hop": dschema.BoolAttribute{
				MarkdownDescription: "Enable multi-hop BFD.",
				Computed:            true,
			},
			"local_address": dschema.StringAttribute{
				MarkdownDescription: "The local address used for the BFD session.",
				Computed:            true,
			},
			"interface_name": dschema.StringAttribute{
				MarkdownDescription: "The interface used for the BFD session.",
				Computed:            true,
			},
			"detect_multiplier": dschema.Int64Attribute{
				MarkdownDescription: "The detection multiplier.",
				Computed:            true,
			},
			"receive_interval": dschema.Int64Attribute{
				MarkdownDescription: "The minimum receive interval in milliseconds.",
				Computed:            true,
			},
			"transmit_interval": dschema.Int64Attribute{
				MarkdownDescription: "The minimum transmit interval in milliseconds.",
				Computed:            true,
			},
		},
	}
}

func convertBFDNeighborSchemaToStruct(d *bfdNeighborResourceModel) (*quagga.BFDNeighbor, error) {
	return &quagga.BFDNeighbor{
		Enabled:          tools.BoolToString(d.Enabled.ValueBool()),
		Description:      d.Description.ValueString(),
		Address:          d.Address.ValueString(),
		MultiHop:         tools.BoolToString(d.MultiHop.ValueBool()),
		LocalAddress:     d.LocalAddress.ValueString(),
		InterfaceName:    api.SelectedMap(d.InterfaceName.ValueString()),
		DetectMultiplier: tools.Int64ToString(d.DetectMultiplier.ValueInt64()),
		ReceiveInterval:  tools.Int64ToString(d.ReceiveInterval.ValueInt64()),
		TransmitInterval: tools.Int64ToString(d.TransmitInterval.ValueInt64()),
	}, nil
}

func convertBFDNeighborStructToSchema(d *quagga.BFDNeighbor) (*bfdNeighborResourceModel, error) {
	return &bfdNeighborResourceModel{
		Enabled:          types.BoolValue(tools.StringToBool(d.Enabled)),
		Description:      types.StringValue(d.Description),
		Address:          types.StringValue(d.Address),
		MultiHop:         types.BoolValue(tools.StringToBool(d.MultiHop)),
		LocalAddress:     types.StringValue(d.LocalAddress),
		InterfaceName:    types.StringValue(d.InterfaceName.String()),
		DetectMultiplier: types.Int64Value(tools.StringToInt64(d.DetectMultiplier)),
		ReceiveInterval:  types.Int64Value(tools.StringToInt64(d.ReceiveInterval)),
		TransmitInterval: types.Int64Value(tools.StringToInt64(d.TransmitInterval)),
	}, nil
}
