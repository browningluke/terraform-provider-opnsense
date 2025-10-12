package routes

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/routes"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// routeResourceModel describes the resource data model.
type routeResourceModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Description types.String `tfsdk:"description"`
	Gateway     types.String `tfsdk:"gateway"`
	Network     types.String `tfsdk:"network"`

	Id types.String `tfsdk:"id"`
}

func routeResourceSchema() schema.Schema {
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

func routeDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Routes can be used to teach your firewall which path it should take when forwarding packets to a specific network.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether this route is enabled.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Computed:            true,
			},
			"gateway": dschema.StringAttribute{
				MarkdownDescription: "Which gateway this route applies, e.g. `WAN`.",
				Computed:            true,
			},
			"network": dschema.StringAttribute{
				MarkdownDescription: "Destination network for this static route.",
				Computed:            true,
			},
		},
	}
}

func convertRouteSchemaToStruct(d *routeResourceModel) (*routes.Route, error) {
	return &routes.Route{
		Disabled:    tools.BoolToString(!d.Enabled.ValueBool()),
		Description: d.Description.ValueString(),
		Gateway:     api.SelectedMap(d.Gateway.ValueString()),
		Network:     d.Network.ValueString(),
	}, nil
}

func convertRouteStructToSchema(d *routes.Route) (*routeResourceModel, error) {
	return &routeResourceModel{
		Enabled:     types.BoolValue(!tools.StringToBool(d.Disabled)),
		Description: tools.StringOrNull(d.Description),
		Gateway:     types.StringValue(d.Gateway.String()),
		Network:     types.StringValue(d.Network),
	}, nil
}
