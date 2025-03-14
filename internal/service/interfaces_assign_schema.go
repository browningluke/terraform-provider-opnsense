package service

import (
	"terraform-provider-opnsense/internal/tools"

	"github.com/browningluke/opnsense-go/pkg/interfaces"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// InterfacesAssignResourceModel describes the resource data model.
type InterfacesAssignResourceModel struct {
	Interface        types.String `tfsdk:"interface"`
	Device           types.String `tfsdk:"device"`
	Ip               types.String `tfsdk:"ipaddr"`
	Gateway          types.String `tfsdk:"gateway"`
	Description      types.String `tfsdk:"description"`
	Enable           types.Bool   `tfsdk:"enable"`
	Subnet           types.String `tfsdk:"subnet"`
	GatewayInterface types.Bool   `tfsdk:"gateway_interface"`

	Id types.String `tfsdk:"id"`
}

func InterfacesAssignResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "OPNsense interfaces assign resource",

		Attributes: map[string]schema.Attribute{
			"interface": schema.StringAttribute{
				MarkdownDescription: "Interface name, e.g. `lan`, `wan`, `opt1`, `opt2`. (Automated next free interface name opt+i)",
				Optional:            true,
			},
			"device": schema.StringAttribute{
				MarkdownDescription: "Device name, e.g. `igb0`, `igb1`, `igb2`, `igb3`, `igb4`, `igb5`.",
				Required:            true,
			},
			"ipaddr": schema.StringAttribute{
				MarkdownDescription: "IP address, e.g. ``",
				Optional:            true,
			},
			"gateway": schema.StringAttribute{
				MarkdownDescription: "Gateway address, e.g. ``.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Automated Upper Case of the interface name.",
				Optional:            true,
			},
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable the interface.",
				Optional:            true,
			},
			"subnet": schema.StringAttribute{
				MarkdownDescription: "Subnet mask, e.g. ``.",
				Optional:            true,
			},
			"gateway_interface": schema.BoolAttribute{
				MarkdownDescription: "Gateway interface",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func convertInterfacesAssignSchemaToStruct(d *InterfacesAssignResourceModel) (*interfaces.Assign, error) {
	return &interfaces.Assign{
		Interface:        d.Interface.ValueString(),
		Device:           d.Device.ValueString(),
		Ip:               d.Ip.ValueString(),
		Gateway:          d.Gateway.ValueString(),
		Description:      d.Description.ValueString(),
		Enable:           tools.BoolToString(d.Enable.ValueBool()),
		Subnet:           d.Subnet.ValueString(),
		GatewayInterface: tools.BoolToString(d.GatewayInterface.ValueBool()),
	}, nil
}

func convertInterfacesAssignStructToSchema(d *interfaces.Assign) (*InterfacesAssignResourceModel, error) {
	return &InterfacesAssignResourceModel{
		Interface:        types.StringValue(d.Interface),
		Device:           types.StringValue(d.Device),
		Ip:               types.StringValue(d.Ip),
		Gateway:          types.StringValue(d.Gateway),
		Description:      tools.StringOrNull(d.Description),
		Enable:           types.BoolValue(tools.StringToBool(d.Enable)),
		Subnet:           types.StringValue(d.Subnet),
		GatewayInterface: types.BoolValue(tools.StringToBool(d.GatewayInterface)),
	}, nil
}
