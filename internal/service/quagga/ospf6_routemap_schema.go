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

// ospf6RouteMapResourceModel describes the resource data model.
type ospf6RouteMapResourceModel struct {
	Enabled    types.Bool   `tfsdk:"enabled"`
	Name       types.String `tfsdk:"name"`
	Action     types.String `tfsdk:"action"`
	RouteMapID types.String `tfsdk:"route_map_id"`
	PrefixList types.String `tfsdk:"prefix_list"`
	Set        types.String `tfsdk:"set"`

	Id types.String `tfsdk:"id"`
}

func ospf6RouteMapResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure OSPFv3 route maps.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this route map. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of this route map.",
				Required:            true,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Set permit for match or deny to negate the rule. One of `\"permit\"`, `\"deny\"`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("permit", "deny"),
				},
			},
			"route_map_id": schema.StringAttribute{
				MarkdownDescription: "The route map ID (sequence number).",
				Required:            true,
			},
			"prefix_list": schema.StringAttribute{
				MarkdownDescription: "UUID of the prefix list to match. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"set": schema.StringAttribute{
				MarkdownDescription: "Free text field for set commands. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the OSPFv3 route map.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func ospf6RouteMapDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure OSPFv3 route maps.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this route map.",
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
			"route_map_id": dschema.StringAttribute{
				MarkdownDescription: "The route map ID (sequence number).",
				Computed:            true,
			},
			"prefix_list": dschema.StringAttribute{
				MarkdownDescription: "UUID of the prefix list to match.",
				Computed:            true,
			},
			"set": dschema.StringAttribute{
				MarkdownDescription: "Free text field for set commands.",
				Computed:            true,
			},
		},
	}
}

func convertOSPF6RouteMapSchemaToStruct(d *ospf6RouteMapResourceModel) (*quagga.OSPF6RouteMap, error) {
	return &quagga.OSPF6RouteMap{
		Enabled:    tools.BoolToString(d.Enabled.ValueBool()),
		Name:       d.Name.ValueString(),
		Action:     api.SelectedMap(d.Action.ValueString()),
		RouteMapID: d.RouteMapID.ValueString(),
		PrefixList: api.SelectedMap(d.PrefixList.ValueString()),
		Set:        d.Set.ValueString(),
	}, nil
}

func convertOSPF6RouteMapStructToSchema(d *quagga.OSPF6RouteMap) (*ospf6RouteMapResourceModel, error) {
	return &ospf6RouteMapResourceModel{
		Enabled:    types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:       types.StringValue(d.Name),
		Action:     types.StringValue(d.Action.String()),
		RouteMapID: types.StringValue(d.RouteMapID),
		PrefixList: types.StringValue(d.PrefixList.String()),
		Set:        types.StringValue(d.Set),
	}, nil
}
