package service

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/opnsense"
)

// RouteResourceModel describes the resource data model.
type RouteResourceModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Description types.String `tfsdk:"description"`
	Gateway     types.String `tfsdk:"gateway"`
	Network     types.String `tfsdk:"network"`

	Id types.String `tfsdk:"id"`
}

func RouteResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Routes can be used to teach your firewall which path it should take when forwarding packets to a specific network.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this route.  Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Optional:            true,
			},
			"gateway": schema.StringAttribute{
				MarkdownDescription: "Which gateway this route applies, e.g. `WAN`. Must be an existing gateway.",
				Required:            true,
			},
			"network": schema.StringAttribute{
				MarkdownDescription: "Destination network for this static route.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the route.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func convertRouteSchemaToStruct(d *RouteResourceModel) (*opnsense.Route, error) {
	// Convert 'Enabled' to 'Disabled'
	var disabled string
	if d.Enabled.ValueBool() {
		disabled = "0"
	} else {
		disabled = "1"
	}

	return &opnsense.Route{
		Disabled:    disabled,
		Description: d.Description.ValueString(),
		Gateway:     opnsense.SelectedMap(d.Gateway.ValueString()),
		Network:     d.Network.ValueString(),
	}, nil
}

func convertRouteStructToSchema(d *opnsense.Route) (*RouteResourceModel, error) {
	model := &RouteResourceModel{
		Enabled:     types.BoolValue(true),
		Description: types.StringNull(),
		Gateway:     types.StringValue(d.Gateway.String()),
		Network:     types.StringValue(d.Network),
	}

	// Parse 'Disabled'
	if d.Disabled == "1" {
		model.Enabled = types.BoolValue(false)
	}

	// Parse 'Description'
	if d.Description != "" {
		model.Description = types.StringValue(d.Description)
	}

	return model, nil
}
