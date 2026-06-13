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

// bgpRedistributionResourceModel describes the resource data model.
type bgpRedistributionResourceModel struct {
	Enabled        types.Bool   `tfsdk:"enabled"`
	Description    types.String `tfsdk:"description"`
	Redistribute   types.String `tfsdk:"redistribute"`
	LinkedRouteMap types.String `tfsdk:"linked_route_map"`

	Id types.String `tfsdk:"id"`
}

func bgpRedistributionResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure redistribution rules for BGP.",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this redistribution rule. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "An optional description for this redistribution rule. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"redistribute": schema.StringAttribute{
				MarkdownDescription: "The protocol to redistribute into BGP.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ospf", "connected", "kernel", "rip", "static"),
				},
			},
			"linked_route_map": schema.StringAttribute{
				MarkdownDescription: "The route map UUID to apply to redistributed routes. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the redistribution rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func bgpRedistributionDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure redistribution rules for BGP.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this redistribution rule.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "An optional description for this redistribution rule.",
				Computed:            true,
			},
			"redistribute": dschema.StringAttribute{
				MarkdownDescription: "The protocol to redistribute into BGP.",
				Computed:            true,
			},
			"linked_route_map": dschema.StringAttribute{
				MarkdownDescription: "The route map UUID to apply to redistributed routes.",
				Computed:            true,
			},
		},
	}
}

func convertBGPRedistributionSchemaToStruct(d *bgpRedistributionResourceModel) (*quagga.BGPRedistribution, error) {
	return &quagga.BGPRedistribution{
		Enabled:        tools.BoolToString(d.Enabled.ValueBool()),
		Description:    d.Description.ValueString(),
		Redistribute:   api.SelectedMap(d.Redistribute.ValueString()),
		LinkedRouteMap: api.SelectedMap(d.LinkedRouteMap.ValueString()),
	}, nil
}

func convertBGPRedistributionStructToSchema(d *quagga.BGPRedistribution) (*bgpRedistributionResourceModel, error) {
	return &bgpRedistributionResourceModel{
		Enabled:        types.BoolValue(tools.StringToBool(d.Enabled)),
		Description:    types.StringValue(d.Description),
		Redistribute:   types.StringValue(d.Redistribute.String()),
		LinkedRouteMap: types.StringValue(d.LinkedRouteMap.String()),
	}, nil
}
