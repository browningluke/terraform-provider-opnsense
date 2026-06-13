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

// staticRouteResourceModel describes the resource data model.
type staticRouteResourceModel struct {
	Enabled       types.Bool   `tfsdk:"enabled"`
	Network       types.String `tfsdk:"network"`
	Gateway       types.String `tfsdk:"gateway"`
	InterfaceName types.String `tfsdk:"interface_name"`
	BFD           types.Bool   `tfsdk:"bfd"`

	Id types.String `tfsdk:"id"`
}

func staticRouteResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure static routes for Quagga.",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this static route. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"network": schema.StringAttribute{
				MarkdownDescription: "The destination network in CIDR notation (e.g. `192.168.0.0/24`).",
				Required:            true,
			},
			"gateway": schema.StringAttribute{
				MarkdownDescription: "The gateway IP address for this static route.",
				Required:            true,
			},
			"interface_name": schema.StringAttribute{
				MarkdownDescription: "The interface to use for this static route. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"bfd": schema.BoolAttribute{
				MarkdownDescription: "Enable BFD tracking for this static route. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the static route.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func staticRouteDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure static routes for Quagga.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this static route.",
				Computed:            true,
			},
			"network": dschema.StringAttribute{
				MarkdownDescription: "The destination network in CIDR notation.",
				Computed:            true,
			},
			"gateway": dschema.StringAttribute{
				MarkdownDescription: "The gateway IP address for this static route.",
				Computed:            true,
			},
			"interface_name": dschema.StringAttribute{
				MarkdownDescription: "The interface used for this static route.",
				Computed:            true,
			},
			"bfd": dschema.BoolAttribute{
				MarkdownDescription: "Enable BFD tracking for this static route.",
				Computed:            true,
			},
		},
	}
}

func convertStaticRouteSchemaToStruct(d *staticRouteResourceModel) (*quagga.StaticRoute, error) {
	return &quagga.StaticRoute{
		Enabled:       tools.BoolToString(d.Enabled.ValueBool()),
		Network:       d.Network.ValueString(),
		Gateway:       d.Gateway.ValueString(),
		InterfaceName: api.SelectedMap(d.InterfaceName.ValueString()),
		BFD:           tools.BoolToString(d.BFD.ValueBool()),
	}, nil
}

func convertStaticRouteStructToSchema(d *quagga.StaticRoute) (*staticRouteResourceModel, error) {
	return &staticRouteResourceModel{
		Enabled:       types.BoolValue(tools.StringToBool(d.Enabled)),
		Network:       types.StringValue(d.Network),
		Gateway:       types.StringValue(d.Gateway),
		InterfaceName: types.StringValue(d.InterfaceName.String()),
		BFD:           types.BoolValue(tools.StringToBool(d.BFD)),
	}, nil
}
