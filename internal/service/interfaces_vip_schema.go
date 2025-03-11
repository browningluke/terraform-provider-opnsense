package service

import (
	"terraform-provider-opnsense/internal/tools"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/interfaces"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// InterfacesVipResourceModel describes the resource data model.
type InterfacesVipResourceModel struct {
	Interface   types.String `tfsdk:"interface"`
	Mode        types.String `tfsdk:"mode"`
	Network     types.String `tfsdk:"network"`
	Gateway     types.String `tfsdk:"gateway"`
	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func InterfacesVipResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "VIPs (Virtual IPs) can be used to segment a single physical network into multiple virtual networks.",

		Attributes: map[string]schema.Attribute{
			"interface": schema.StringAttribute{
				MarkdownDescription: "Interface to assign the VIP to, e.g. `lan`.",
				Required:            true,
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "VIP mode, e.g. `ipalias`, `carp`, `proxyarp`, `other`.",
				Required:            true,
			},
			"network": schema.StringAttribute{
				MarkdownDescription: "network address, e.g. ``",
				Required:            true,
			},
			"gateway": schema.StringAttribute{
				MarkdownDescription: "Gateway address, e.g. ``.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the VIP.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func convertInterfacesVipSchemaToStruct(d *InterfacesVipResourceModel) (*interfaces.Vip, error) {
	return &interfaces.Vip{
		Interface:   api.SelectedMap(d.Interface.ValueString()),
		Mode:        api.SelectedMap(d.Mode.ValueString()),
		Network:     d.Network.ValueString(),
		Gateway:     d.Gateway.ValueString(),
		Description: d.Description.ValueString(),
	}, nil
}

func convertInterfacesVipStructToSchema(d *interfaces.Vip) (*InterfacesVipResourceModel, error) {
	return &InterfacesVipResourceModel{
		Interface:   types.StringValue(d.Interface.String()),
		Mode:        types.StringValue(d.Mode.String()),
		Network:     types.StringValue(d.Network),
		Gateway:     types.StringValue(d.Gateway),
		Description: tools.StringOrNull(d.Description),
	}, nil
}
