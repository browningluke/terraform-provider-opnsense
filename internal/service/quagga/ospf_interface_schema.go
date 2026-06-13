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

// ospfInterfaceResourceModel describes the resource data model.
type ospfInterfaceResourceModel struct {
	Enabled            types.Bool   `tfsdk:"enabled"`
	InterfaceName      types.String `tfsdk:"interface_name"`
	AuthType           types.String `tfsdk:"auth_type"`
	AuthKey            types.String `tfsdk:"auth_key"`
	AuthKeyID          types.String `tfsdk:"auth_key_id"`
	Area               types.String `tfsdk:"area"`
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
	P2MPOptions        types.String `tfsdk:"p2mp_options"`

	Id types.String `tfsdk:"id"`
}

func ospfInterfaceResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure OSPF interfaces.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this OSPF interface. Defaults to `true`.",
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
			"auth_type": schema.StringAttribute{
				MarkdownDescription: "Authentication type. One of `\"\"`, `\"message-digest\"`, `\"plain\"`. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf("", "message-digest", "plain"),
				},
			},
			"auth_key": schema.StringAttribute{
				MarkdownDescription: "Authentication key. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"auth_key_id": schema.StringAttribute{
				MarkdownDescription: "Authentication key ID. Defaults to `\"1\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("1"),
			},
			"area": schema.StringAttribute{
				MarkdownDescription: "The OSPF area this interface belongs to. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
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
				MarkdownDescription: "Network type. One of `\"\"`, `\"broadcast\"`, `\"non-broadcast\"`, `\"point-to-multipoint\"`, `\"point-to-point\"`. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf("", "broadcast", "non-broadcast", "point-to-multipoint", "point-to-point"),
				},
			},
			"p2mp_options": schema.StringAttribute{
				MarkdownDescription: "Point-to-multipoint options. One of `\"\"`, `\"delay-reflood\"`, `\"non-broadcast\"`. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf("", "delay-reflood", "non-broadcast"),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the OSPF interface.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func ospfInterfaceDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure OSPF interfaces.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this OSPF interface.",
				Computed:            true,
			},
			"interface_name": dschema.StringAttribute{
				MarkdownDescription: "The interface name.",
				Computed:            true,
			},
			"auth_type": dschema.StringAttribute{
				MarkdownDescription: "Authentication type.",
				Computed:            true,
			},
			"auth_key": dschema.StringAttribute{
				MarkdownDescription: "Authentication key.",
				Computed:            true,
			},
			"auth_key_id": dschema.StringAttribute{
				MarkdownDescription: "Authentication key ID.",
				Computed:            true,
			},
			"area": dschema.StringAttribute{
				MarkdownDescription: "The OSPF area this interface belongs to.",
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
			"p2mp_options": dschema.StringAttribute{
				MarkdownDescription: "Point-to-multipoint options.",
				Computed:            true,
			},
		},
	}
}

func convertOSPFInterfaceSchemaToStruct(d *ospfInterfaceResourceModel) (*quagga.OSPFInterface, error) {
	return &quagga.OSPFInterface{
		Enabled:            tools.BoolToString(d.Enabled.ValueBool()),
		InterfaceName:      api.SelectedMap(d.InterfaceName.ValueString()),
		AuthType:           api.SelectedMap(d.AuthType.ValueString()),
		AuthKey:            d.AuthKey.ValueString(),
		AuthKeyID:          d.AuthKeyID.ValueString(),
		Area:               d.Area.ValueString(),
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
		P2MPOptions:        api.SelectedMap(d.P2MPOptions.ValueString()),
	}, nil
}

func convertOSPFInterfaceStructToSchema(d *quagga.OSPFInterface) (*ospfInterfaceResourceModel, error) {
	return &ospfInterfaceResourceModel{
		Enabled:            types.BoolValue(tools.StringToBool(d.Enabled)),
		InterfaceName:      types.StringValue(d.InterfaceName.String()),
		AuthType:           types.StringValue(d.AuthType.String()),
		AuthKey:            types.StringValue(d.AuthKey),
		AuthKeyID:          types.StringValue(d.AuthKeyID),
		Area:               types.StringValue(d.Area),
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
		P2MPOptions:        types.StringValue(d.P2MPOptions.String()),
	}, nil
}
