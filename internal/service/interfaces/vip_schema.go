package interfaces

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/interfaces"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/browningluke/terraform-provider-opnsense/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// vipResourceModel describes the resource data model.
type vipResourceModel struct {
	Mode        types.String `tfsdk:"mode"`
	Interface   types.String `tfsdk:"interface"`
	Network     types.String `tfsdk:"network"`
	Gateway     types.String `tfsdk:"gateway"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
}

func vipResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Virtual IPs allow an OPNsense firewall to assign multiple IP addresses to the same network interface.",

		Attributes: map[string]schema.Attribute{
			"mode": schema.StringAttribute{
				MarkdownDescription: "Mode of the VIP. One of `ipalias` or `proxyarp`. `proxyarp` cannot be bound to by anything running on the firewall, such as IPsec, OpenVPN, etc. In most cases an `ipalias` should be used.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("ipalias"),
				Validators: []validator.String{
					stringvalidator.OneOf("ipalias", "proxyarp"),
				},
			},
			"interface": schema.StringAttribute{
				MarkdownDescription: "Choose which interface this VIP applies to.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("wan"),
			},
			"network": schema.StringAttribute{
				MarkdownDescription: "Provide an address and subnet to use. (e.g 192.168.0.1/24)",
				Required:            true,
				Validators: []validator.String{
					validators.CIDR(),
				},
			},
			"gateway": schema.StringAttribute{
				MarkdownDescription: "For some interface types a gateway is required to configure an IP Alias (ppp/pppoe/tun), leave this field empty for all other interface types.",
				Optional:            true,
				Validators: []validator.String{
					validators.IpOrCIDR(),
				},
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

func vipDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Virtual IPs allow an OPNsense firewall to assign multiple IP addresses to the same network interface.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "Mode of the VIP. One of `ipalias`, `carp`, or `proxyarp`. `proxyarp` cannot be bound to by anything running on the firewall, such as IPsec, OpenVPN, etc. In most cases an `ipalias` should be used.",
				Computed:            true,
			},
			"interface": schema.StringAttribute{
				MarkdownDescription: "Choose which interface this VIP applies to.",
				Computed:            true,
			},
			"network": schema.StringAttribute{
				MarkdownDescription: "Provide an address and subnet to use. (e.g 192.168.0.1/24)",
				Computed:            true,
			},
			"gateway": schema.StringAttribute{
				MarkdownDescription: "For some interface types a gateway is required to configure an IP Alias (ppp/pppoe/tun), leave this field empty for all other interface types.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Computed:            true,
			},
		},
	}
}

func convertVipSchemaToStruct(d *vipResourceModel) (*interfaces.Vip, error) {
	return &interfaces.Vip{
		Description: d.Description.ValueString(),
		Mode:        api.SelectedMap(d.Mode.ValueString()),
		Interface:   api.SelectedMap(d.Interface.ValueString()),
		Network:     d.Network.ValueString(),
		Gateway:     d.Gateway.ValueString(),
	}, nil
}

func convertVipStructToSchema(d *interfaces.Vip) (*vipResourceModel, error) {
	return &vipResourceModel{
		Mode:        types.StringValue(d.Mode.String()),
		Interface:   types.StringValue(d.Interface.String()),
		Network:     types.StringValue(d.Network),
		Gateway:     tools.StringOrNull(d.Gateway),
		Description: tools.StringOrNull(d.Description),
	}, nil
}
