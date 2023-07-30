package service

import (
	"github.com/browningluke/opnsense-go/pkg/firewall"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/tools"
)

// FirewallCategoryResourceModel describes the resource data model.
type FirewallCategoryResourceModel struct {
	Automatic types.Bool   `tfsdk:"auto"`
	Name      types.String `tfsdk:"name"`
	Color     types.String `tfsdk:"color"`

	Id types.String `tfsdk:"id"`
}

func FirewallCategoryResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "To ease maintenance of larger rulesets, OPNsense includes categories for the firewall. Each rule can contain one or more categories.",

		Attributes: map[string]schema.Attribute{
			"auto": schema.BoolAttribute{
				MarkdownDescription: "If set, this category will be removed when unused. This is included for completeness, but will result in constant recreations if not attached to any rules, and thus it is advised to leave it as default. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Enter a name for this category.",
				Required:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Pick a color to use. Must be a hex color in format `rrggbb` (e.g. `ff0000`). Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func FirewallCategoryDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "To ease maintenance of larger rulesets, OPNsense includes categories for the firewall. Each rule can contain one or more categories.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"auto": dschema.BoolAttribute{
				MarkdownDescription: "If set, this category will be removed when unused.",
				Computed:            true,
			},
			"name": dschema.StringAttribute{
				MarkdownDescription: "The name for this category.",
				Computed:            true,
			},
			"color": dschema.StringAttribute{
				MarkdownDescription: "The color to use. Must be a hex color in format `rrggbb` (e.g. `ff0000`).",
				Computed:            true,
			},
		},
	}
}

func convertFirewallCategorySchemaToStruct(d *FirewallCategoryResourceModel) (*firewall.Category, error) {
	return &firewall.Category{
		Automatic: tools.BoolToString(d.Automatic.ValueBool()),
		Name:      d.Name.ValueString(),
		Color:     d.Color.ValueString(),
	}, nil
}

func convertFirewallCategoryStructToSchema(d *firewall.Category) (*FirewallCategoryResourceModel, error) {
	return &FirewallCategoryResourceModel{
		Automatic: types.BoolValue(tools.StringToBool(d.Automatic)),
		Name:      types.StringValue(d.Name),
		Color:     types.StringValue(d.Color),
	}, nil
}
