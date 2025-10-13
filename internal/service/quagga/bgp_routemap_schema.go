package quagga

import (
	"context"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// bgpRouteMapResourceModel describes the resource data model.
type bgpRouteMapResourceModel struct {
	Enabled       types.Bool   `tfsdk:"enabled"`
	Description   types.String `tfsdk:"description"`
	Name          types.String `tfsdk:"name"`
	Action        types.String `tfsdk:"action"`
	RouteMapID    types.Int64  `tfsdk:"route_map_id"`
	ASPathList    types.Set    `tfsdk:"aspaths"`
	PrefixList    types.Set    `tfsdk:"prefix_lists"`
	CommunityList types.Set    `tfsdk:"community_lists"`
	Set           types.String `tfsdk:"set"`

	Id types.String `tfsdk:"id"`
}

func bgpRouteMapResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure route maps for BGP.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this route map. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "An optional description for this route map. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of this route map.",
				Required:            true,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Set permit for match or deny to negate the rule. Defaults to `\"permit\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("permit"),
				Validators: []validator.String{
					stringvalidator.OneOf("permit", "deny"),
				},
			},
			"route_map_id": schema.Int64Attribute{
				MarkdownDescription: "The Route-map ID between 1 and 65535. Be aware that the sorting will be done under the hood, so when you add an entry between it gets to the right position.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"aspaths": schema.SetAttribute{
				MarkdownDescription: "Set the AS Path list IDs to use. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
				ElementType:         types.StringType,
			},
			"prefix_lists": schema.SetAttribute{
				MarkdownDescription: "Set the prefix list IDs to use. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
				ElementType:         types.StringType,
			},
			"community_lists": schema.SetAttribute{
				MarkdownDescription: "Set the community list IDs to use. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
				ElementType:         types.StringType,
			},
			"set": schema.StringAttribute{
				MarkdownDescription: "Free text field for your set, please be careful! You can set e.g. `local-preference 300` or `community 1:1` (http://www.nongnu.org/quagga/docs/docs-multi/Route-Map-Set-Command.html#Route-Map-Set-Command). Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the route map.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func bgpRouteMapDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure route maps for BGP.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this route map.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "An optional description for this route map.",
				Computed:            true,
			},
			"name": dschema.StringAttribute{
				MarkdownDescription: "The name of this route map.",
				Computed:            true,
			},
			"action": dschema.StringAttribute{
				MarkdownDescription: "Set permit for match or deny to negate the rule.",
				Computed:            true,
			},
			"route_map_id": dschema.Int64Attribute{
				MarkdownDescription: "The Route-map ID between 1 and 65535. Be aware that the sorting will be done under the hood, so when you add an entry between it gets to the right position.",
				Computed:            true,
			},
			"aspaths": dschema.SetAttribute{
				MarkdownDescription: "Set the AS Path list IDs to use.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"prefix_lists": dschema.SetAttribute{
				MarkdownDescription: "Set the prefix list IDs to use.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"community_lists": dschema.SetAttribute{
				MarkdownDescription: "Set the community list IDs to use.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"set": dschema.StringAttribute{
				MarkdownDescription: "Free text field for your set, please be careful! You can set e.g. `local-preference 300` or `community 1:1` (http://www.nongnu.org/quagga/docs/docs-multi/Route-Map-Set-Command.html#Route-Map-Set-Command). Defaults to `\"\"`.",
				Computed:            true,
			},
		},
	}
}

func convertBGPRouteMapSchemaToStruct(d *bgpRouteMapResourceModel) (*quagga.BGPRouteMap, error) {
	// Parse 'ASPathList'
	var asPathList []string
	d.ASPathList.ElementsAs(context.Background(), &asPathList, false)

	// Parse 'PrefixList'
	var prefixList []string
	d.PrefixList.ElementsAs(context.Background(), &prefixList, false)

	// Parse 'CommunityList'
	var communityList []string
	d.CommunityList.ElementsAs(context.Background(), &communityList, false)

	return &quagga.BGPRouteMap{
		Enabled:       tools.BoolToString(d.Enabled.ValueBool()),
		Description:   d.Description.ValueString(),
		Name:          d.Name.ValueString(),
		Action:        api.SelectedMap(d.Action.ValueString()),
		RouteMapID:    tools.Int64ToString(d.RouteMapID.ValueInt64()),
		ASPathList:    asPathList,
		PrefixList:    prefixList,
		CommunityList: communityList,
		Set:           d.Set.ValueString(),
	}, nil
}

func convertBGPRouteMapStructToSchema(d *quagga.BGPRouteMap) (*bgpRouteMapResourceModel, error) {
	model := &bgpRouteMapResourceModel{
		Enabled:     types.BoolValue(tools.StringToBool(d.Enabled)),
		Description: types.StringValue(d.Description),
		Name:        types.StringValue(d.Name),
		Action:      types.StringValue(d.Action.String()),
		RouteMapID:  types.Int64Value(tools.StringToInt64(d.RouteMapID)),
		Set:         types.StringValue(d.Set),
	}

	// Parse 'ASPathList'
	model.ASPathList = tools.StringSliceToSet(d.ASPathList)

	// Parse 'PrefixList'
	model.PrefixList = tools.StringSliceToSet(d.PrefixList)

	// Parse 'CommunityList'
	model.CommunityList = tools.StringSliceToSet(d.CommunityList)

	return model, nil
}
