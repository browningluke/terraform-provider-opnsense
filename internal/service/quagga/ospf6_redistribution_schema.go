package quagga

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ospf6RedistributionResourceModel describes the resource data model.
type ospf6RedistributionResourceModel struct {
	Enabled        types.Bool   `tfsdk:"enabled"`
	Description    types.String `tfsdk:"description"`
	Redistribute   types.String `tfsdk:"redistribute"`
	LinkedRouteMap types.String `tfsdk:"linked_route_map"`

	Id types.String `tfsdk:"id"`
}

func ospf6RedistributionResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure OSPFv3 redistribution.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this redistribution. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "An optional description for this redistribution. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"redistribute": schema.StringAttribute{
				MarkdownDescription: "The protocol to redistribute. One of `\"bgp\"`, `\"connected\"`, `\"kernel\"`, `\"rip\"`, `\"static\"`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("bgp", "connected", "kernel", "rip", "static"),
				},
			},
			"linked_route_map": schema.StringAttribute{
				MarkdownDescription: "UUID of the route map to apply. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the OSPFv3 redistribution.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func ospf6RedistributionDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure OSPFv3 redistribution.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this redistribution.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "An optional description for this redistribution.",
				Computed:            true,
			},
			"redistribute": dschema.StringAttribute{
				MarkdownDescription: "The protocol to redistribute.",
				Computed:            true,
			},
			"linked_route_map": dschema.StringAttribute{
				MarkdownDescription: "UUID of the route map to apply.",
				Computed:            true,
			},
		},
	}
}

func convertOSPF6RedistributionSchemaToStruct(d *ospf6RedistributionResourceModel) (*quagga.OSPF6Redistribution, error) {
	return &quagga.OSPF6Redistribution{
		Enabled:        tools.BoolToString(d.Enabled.ValueBool()),
		Description:    d.Description.ValueString(),
		Redistribute:   api.SelectedMap(d.Redistribute.ValueString()),
		LinkedRouteMap: api.SelectedMap(d.LinkedRouteMap.ValueString()),
	}, nil
}

func convertOSPF6RedistributionStructToSchema(d *quagga.OSPF6Redistribution) (*ospf6RedistributionResourceModel, error) {
	return &ospf6RedistributionResourceModel{
		Enabled:        types.BoolValue(tools.StringToBool(d.Enabled)),
		Description:    types.StringValue(d.Description),
		Redistribute:   types.StringValue(d.Redistribute.String()),
		LinkedRouteMap: types.StringValue(d.LinkedRouteMap.String()),
	}, nil
}
