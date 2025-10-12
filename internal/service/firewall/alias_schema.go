package firewall

import (
	"context"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/firewall"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// aliasResourceModel describes the resource data model.
type aliasResourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled"`
	Name    types.String `tfsdk:"name"`
	Type    types.String `tfsdk:"type"`

	IPProtocol types.String `tfsdk:"ip_protocol"`
	Interface  types.String `tfsdk:"interface"`

	Content    types.Set `tfsdk:"content"`
	Categories types.Set `tfsdk:"categories"`

	UpdateFreq types.Float64 `tfsdk:"update_freq"`

	Statistics  types.Bool   `tfsdk:"stats"`
	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func aliasResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Aliases are named lists of networks, hosts or ports that can be used as one entity by selecting the alias name in the various supported sections of the firewall. These aliases are particularly useful to condense firewall rules and minimize changes.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this firewall alias. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name must start with a letter or single underscore, be less than 32 characters and only consist of alphanumeric characters or underscores. Aliases can be nested using this name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of alias.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"host", "network", "port", "url", "urltable", "geoip", "networkgroup",
						"mac", "asn", "dynipv6host", "authgroup", "internal", "external",
					),
				},
			},
			"ip_protocol": schema.StringAttribute{
				MarkdownDescription: "Select the Internet Protocol version this alias applies to. Available values: `IPv4`, `IPv6`. Only applies when `type = \"asn\"`, `type = \"geoip\"`, or `type = \"external\"`. Defaults to `IPv4`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("IPv4", "IPv6"),
				},
				Default: stringdefault.StaticString("IPv4"),
			},
			"interface": schema.StringAttribute{
				MarkdownDescription: "Choose on which interface this alias applies. Only applies (and must be set) when `type = \"dynipv6host\"`. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"content": schema.SetAttribute{
				MarkdownDescription: "The content of the alias. Enter ISO 3166-1 country codes when `type = \"geoip\"` (e.g. `[\"CA\", \"FR\"]`). Enter `__<int>_network`, or alias when `type = \"networkgroup\"` (e.g. `[\"__wan_network\", \"otheralias\"]`). Enter OpenVPN group when `type = \"authgroup\"` (e.g. `[\"admins\"]`). Set to `[]` when `type = \"external\"`. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"categories": schema.SetAttribute{
				MarkdownDescription: "Set of category IDs to apply. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"update_freq": schema.Float64Attribute{
				MarkdownDescription: "The frequency that the list will be refreshed, in days (e.g. for 30 hours, enter `1.25`). Only applies (and must be set) when `type = \"urltable\"`. Defaults to `-1`.",
				Optional:            true,
				Computed:            true,
				Default:             float64default.StaticFloat64(-1),
			},
			"stats": schema.BoolAttribute{
				MarkdownDescription: "Whether to maintain a set of counters for each table entry.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Optional:            true,
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

func aliasDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Aliases are named lists of networks, hosts or ports that can be used as one entity by selecting the alias name in the various supported sections of the firewall. These aliases are particularly useful to condense firewall rules and minimize changes.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this firewall alias.",
				Computed:            true,
			},
			"name": dschema.StringAttribute{
				MarkdownDescription: "The name must start with a letter or single underscore, be less than 32 characters and only consist of alphanumeric characters or underscores. Aliases can be nested using this name.",
				Computed:            true,
			},
			"type": dschema.StringAttribute{
				MarkdownDescription: "The type of alias.",
				Computed:            true,
			},
			"ip_protocol": dschema.StringAttribute{
				MarkdownDescription: "Select the Internet Protocol version this alias applies to. Available values: `IPv4`, `IPv6`. Only applies when `type = \"asn\"`, `type = \"geoip\"`, or `type = \"external\"`.",
				Computed:            true,
			},
			"interface": dschema.StringAttribute{
				MarkdownDescription: "Choose on which interface this alias applies. Only applies (and must be set) when `type = \"dynipv6host\"`.",
				Computed:            true,
			},
			"content": dschema.SetAttribute{
				MarkdownDescription: "The content of the alias. Enter ISO 3166-1 country codes when `type = \"geoip\"` (e.g. `[\"CA\", \"FR\"]`). Enter `__<int>_network`, or alias when `type = \"networkgroup\"` (e.g. `[\"__wan_network\", \"otheralias\"]`). Enter OpenVPN group when `type = \"authgroup\"` (e.g. `[\"admins\"]`). Set to `[]` when `type = \"external\"`.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"categories": dschema.SetAttribute{
				MarkdownDescription: "Set of category IDs to apply.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"update_freq": dschema.Float64Attribute{
				MarkdownDescription: "The frequency that the list will be refreshed, in days (e.g. for 30 hours, enter `1.25`). Only applies (and must be set) when `type = \"urltable\"`.",
				Computed:            true,
			},
			"stats": dschema.BoolAttribute{
				MarkdownDescription: "Whether to maintain a set of counters for each table entry.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Computed:            true,
			},
		},
	}
}

func convertAliasSchemaToStruct(d *aliasResourceModel) (*firewall.Alias, error) {
	// Parse 'Content'
	var contentList []string
	d.Content.ElementsAs(context.Background(), &contentList, false)

	// Parse 'Categories'
	var categoriesList []string
	d.Categories.ElementsAs(context.Background(), &categoriesList, false)

	return &firewall.Alias{
		Enabled:     tools.BoolToString(d.Enabled.ValueBool()),
		Name:        d.Name.ValueString(),
		Type:        api.SelectedMap(d.Type.ValueString()),
		IPProtocol:  api.SelectedMap(d.IPProtocol.ValueString()),
		Interface:   api.SelectedMap(d.Interface.ValueString()),
		Content:     contentList,
		Categories:  categoriesList,
		UpdateFreq:  tools.Float64ToStringNegative(d.UpdateFreq.ValueFloat64()),
		Statistics:  tools.BoolToString(d.Statistics.ValueBool()),
		Description: d.Description.ValueString(),
	}, nil
}

func convertAliasStructToSchema(d *firewall.Alias) (*aliasResourceModel, error) {
	model := &aliasResourceModel{
		Enabled:     types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:        types.StringValue(d.Name),
		Type:        types.StringValue(d.Type.String()),
		IPProtocol:  types.StringValue(d.IPProtocol.String()),
		Interface:   types.StringValue(d.Interface.String()),
		Content:     types.SetNull(types.StringType),
		Categories:  types.SetNull(types.StringType),
		UpdateFreq:  types.Float64Value(tools.StringToFloat64(d.UpdateFreq)),
		Statistics:  types.BoolValue(tools.StringToBool(d.Statistics)),
		Description: tools.StringOrNull(d.Description),
	}

	// Parse 'Content'
	var contentList []attr.Value
	for _, i := range d.Content {
		// OPNsense API always returns empty string in list of content, skip it.
		if i == "" {
			continue
		}
		contentList = append(contentList, basetypes.NewStringValue(i))
	}
	contentTypeList, _ := types.SetValue(types.StringType, contentList)
	model.Content = contentTypeList

	// Parse 'Categories'
	var categoriesList []attr.Value
	for _, i := range d.Categories {
		categoriesList = append(categoriesList, basetypes.NewStringValue(i))
	}
	categoriesTypeList, _ := types.SetValue(types.StringType, categoriesList)
	model.Categories = categoriesTypeList

	return model, nil
}
