package quagga

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
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
)

// ospf6InterfaceResourceModel describes the resource data model.
type ospf6InterfaceResourceModel struct {
	Enabled            types.Bool   `tfsdk:"enabled"`
	InterfaceName      types.String `tfsdk:"interface_name"`
	Area               types.String `tfsdk:"area"`
	Passive            types.Bool   `tfsdk:"passive"`
	Cost               types.Int64  `tfsdk:"cost"`
	CostDemoted        types.Int64  `tfsdk:"cost_demoted"`
	CARPDependOn       types.String `tfsdk:"carp_depend_on"`
	HelloInterval      types.Int64  `tfsdk:"hello_interval"`
	DeadInterval       types.Int64  `tfsdk:"dead_interval"`
	RetransmitInterval types.Int64  `tfsdk:"retransmit_interval"`
	TransmitDelay      types.Int64  `tfsdk:"transmit_delay"`
	Priority           types.Int64  `tfsdk:"priority"`
	BFD                types.Bool   `tfsdk:"bfd"`
	NetworkType        types.String `tfsdk:"network_type"`

	Id types.String `tfsdk:"id"`
}

func ospf6InterfaceResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure OSPFv3 interfaces.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this OSPFv3 interface. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"interface_name": schema.StringAttribute{
				MarkdownDescription: "The interface name. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"area": schema.StringAttribute{
				MarkdownDescription: "The OSPFv3 area this interface belongs to. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"passive": schema.BoolAttribute{
				MarkdownDescription: "Enable passive mode for this interface. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"cost": schema.Int64Attribute{
				MarkdownDescription: "Interface cost. Use `-1` for unset. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"cost_demoted": schema.Int64Attribute{
				MarkdownDescription: "Interface cost when demoted. Use `-1` for unset. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"carp_depend_on": schema.StringAttribute{
				MarkdownDescription: "CARP VHID to depend on. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"hello_interval": schema.Int64Attribute{
				MarkdownDescription: "Hello interval in seconds. Use `-1` for unset. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"dead_interval": schema.Int64Attribute{
				MarkdownDescription: "Dead interval in seconds. Use `-1` for unset. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"retransmit_interval": schema.Int64Attribute{
				MarkdownDescription: "Retransmit interval in seconds. Use `-1` for unset. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"transmit_delay": schema.Int64Attribute{
				MarkdownDescription: "Transmit delay in seconds. Use `-1` for unset. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Interface priority. Use `-1` for unset. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"bfd": schema.BoolAttribute{
				MarkdownDescription: "Enable BFD support for this interface. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"network_type": schema.StringAttribute{
				MarkdownDescription: "Network type. One of `\"\"`, `\"broadcast\"`, `\"point-to-point\"`. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf("", "broadcast", "point-to-point"),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the OSPFv3 interface.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func ospf6InterfaceDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure OSPFv3 interfaces.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this OSPFv3 interface.",
				Computed:            true,
			},
			"interface_name": dschema.StringAttribute{
				MarkdownDescription: "The interface name.",
				Computed:            true,
			},
			"area": dschema.StringAttribute{
				MarkdownDescription: "The OSPFv3 area this interface belongs to.",
				Computed:            true,
			},
			"passive": dschema.BoolAttribute{
				MarkdownDescription: "Enable passive mode for this interface.",
				Computed:            true,
			},
			"cost": dschema.Int64Attribute{
				MarkdownDescription: "Interface cost.",
				Computed:            true,
			},
			"cost_demoted": dschema.Int64Attribute{
				MarkdownDescription: "Interface cost when demoted.",
				Computed:            true,
			},
			"carp_depend_on": dschema.StringAttribute{
				MarkdownDescription: "CARP VHID to depend on.",
				Computed:            true,
			},
			"hello_interval": dschema.Int64Attribute{
				MarkdownDescription: "Hello interval in seconds.",
				Computed:            true,
			},
			"dead_interval": dschema.Int64Attribute{
				MarkdownDescription: "Dead interval in seconds.",
				Computed:            true,
			},
			"retransmit_interval": dschema.Int64Attribute{
				MarkdownDescription: "Retransmit interval in seconds.",
				Computed:            true,
			},
			"transmit_delay": dschema.Int64Attribute{
				MarkdownDescription: "Transmit delay in seconds.",
				Computed:            true,
			},
			"priority": dschema.Int64Attribute{
				MarkdownDescription: "Interface priority.",
				Computed:            true,
			},
			"bfd": dschema.BoolAttribute{
				MarkdownDescription: "Enable BFD support for this interface.",
				Computed:            true,
			},
			"network_type": dschema.StringAttribute{
				MarkdownDescription: "Network type.",
				Computed:            true,
			},
		},
	}
}

func convertOSPF6InterfaceSchemaToStruct(d *ospf6InterfaceResourceModel) (*quagga.OSPF6Interface, error) {
	return &quagga.OSPF6Interface{
		Enabled:            tools.BoolToString(d.Enabled.ValueBool()),
		InterfaceName:      api.SelectedMap(d.InterfaceName.ValueString()),
		Area:               d.Area.ValueString(),
		Passive:            tools.BoolToString(d.Passive.ValueBool()),
		Cost:               tools.Int64ToStringNegative(d.Cost.ValueInt64()),
		CostDemoted:        tools.Int64ToStringNegative(d.CostDemoted.ValueInt64()),
		CARPDependOn:       api.SelectedMap(d.CARPDependOn.ValueString()),
		HelloInterval:      tools.Int64ToStringNegative(d.HelloInterval.ValueInt64()),
		DeadInterval:       tools.Int64ToStringNegative(d.DeadInterval.ValueInt64()),
		RetransmitInterval: tools.Int64ToStringNegative(d.RetransmitInterval.ValueInt64()),
		TransmitDelay:      tools.Int64ToStringNegative(d.TransmitDelay.ValueInt64()),
		Priority:           tools.Int64ToStringNegative(d.Priority.ValueInt64()),
		BFD:                tools.BoolToString(d.BFD.ValueBool()),
		NetworkType:        api.SelectedMap(d.NetworkType.ValueString()),
	}, nil
}

func convertOSPF6InterfaceStructToSchema(d *quagga.OSPF6Interface) (*ospf6InterfaceResourceModel, error) {
	return &ospf6InterfaceResourceModel{
		Enabled:            types.BoolValue(tools.StringToBool(d.Enabled)),
		InterfaceName:      types.StringValue(d.InterfaceName.String()),
		Area:               types.StringValue(d.Area),
		Passive:            types.BoolValue(tools.StringToBool(d.Passive)),
		Cost:               types.Int64Value(tools.StringToInt64(d.Cost)),
		CostDemoted:        types.Int64Value(tools.StringToInt64(d.CostDemoted)),
		CARPDependOn:       types.StringValue(d.CARPDependOn.String()),
		HelloInterval:      types.Int64Value(tools.StringToInt64(d.HelloInterval)),
		DeadInterval:       types.Int64Value(tools.StringToInt64(d.DeadInterval)),
		RetransmitInterval: types.Int64Value(tools.StringToInt64(d.RetransmitInterval)),
		TransmitDelay:      types.Int64Value(tools.StringToInt64(d.TransmitDelay)),
		Priority:           types.Int64Value(tools.StringToInt64(d.Priority)),
		BFD:                types.BoolValue(tools.StringToBool(d.BFD)),
		NetworkType:        types.StringValue(d.NetworkType.String()),
	}, nil
}
