package service

import (
	"regexp"
	"terraform-provider-opnsense/internal/tools"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/interfaces"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// InterfacesVipResourceModel describes the resource data model.
type InterfacesVipResourceModel struct {
	Mode        types.String `tfsdk:"mode"`
	Interface   types.String `tfsdk:"interface"`
	Network     types.String `tfsdk:"network"`
	Gateway     types.String `tfsdk:"gateway"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
}

var ipOrCidrValidator = stringvalidator.RegexMatches(
	regexp.MustCompile(`^(([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))?|([0-9a-fA-F:]+)(\/([0-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8]))?)$`),
	"must be a valid IPv4 or IPv6 address or CIDR (e.g. 192.168.0.1, 192.168.0.0/24, 2001:db8::1, 2001:db8::/64)",
)

var cidrValidator = stringvalidator.RegexMatches(
	regexp.MustCompile(`^(([0-9]{1,3}\.){3}[0-9]{1,3}\/(3[0-2]|[1-2]?[0-9]))$|^(([0-9a-fA-F:]+)\/(12[0-8]|1[0-1][0-9]|[1-9]?[0-9]))$`),
	"must be a valid IPv4 or IPv6 CIDR (e.g. 192.168.0.0/24, 2001:db8::/64)",
)

func InterfacesVipResourceSchema() schema.Schema {
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
					cidrValidator,
				},
			},
			"gateway": schema.StringAttribute{
				MarkdownDescription: "For some interface types a gateway is required to configure an IP Alias (ppp/pppoe/tun), leave this field empty for all other interface types.",
				Optional:            true,
				Validators: []validator.String{
					ipOrCidrValidator,
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

func InterfacesVipDataSourceSchema() dschema.Schema {
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

func convertInterfacesVipSchemaToStruct(d *InterfacesVipResourceModel) (*interfaces.Vip, error) {
	return &interfaces.Vip{
		Description: d.Description.ValueString(),
		Mode:        api.SelectedMap(d.Mode.ValueString()),
		Interface:   api.SelectedMap(d.Interface.ValueString()),
		Network:     d.Network.ValueString(),
		Gateway:     d.Gateway.ValueString(),
	}, nil
}

func convertInterfacesVipStructToSchema(d *interfaces.Vip) (*InterfacesVipResourceModel, error) {
	return &InterfacesVipResourceModel{
		Mode:        types.StringValue(d.Mode.String()),
		Interface:   types.StringValue(d.Interface.String()),
		Network:     types.StringValue(d.Network),
		Gateway:     tools.StringOrNull(d.Gateway),
		Description: tools.StringOrNull(d.Description),
	}, nil
}
