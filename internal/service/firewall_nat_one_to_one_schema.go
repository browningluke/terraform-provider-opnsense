package service

import (
	"context"
	"regexp"
	"terraform-provider-opnsense/internal/tools"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/firewall"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type firewallLocationOneToOne struct {
	Net    types.String `tfsdk:"net"`
	Invert types.Bool   `tfsdk:"invert"`
}

// FirewallNATOneToOneResourceModel describes the resource data model.
type FirewallNATOneToOneResourceModel struct {
	Enabled       types.Bool                `tfsdk:"enabled"`
	Log           types.Bool                `tfsdk:"log"`
	Sequence      types.Int64               `tfsdk:"sequence"`
	Interface     types.String              `tfsdk:"interface"`
	Type          types.String              `tfsdk:"type"`
	Source        *firewallLocationOneToOne `tfsdk:"source"`
	Destination   *firewallLocationOneToOne `tfsdk:"destination"`
	ExternalNet   types.String              `tfsdk:"external_net"`
	NatReflection types.String              `tfsdk:"nat_reflection"`
	Categories    types.Set                 `tfsdk:"categories"`
	Description   types.String              `tfsdk:"description"`
	Id            types.String              `tfsdk:"id"`
}

func FirewallNATOneToOneResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "1:1 NAT maps a public IP or subnet to an internal private IP or subnet. All traffic to the public address is forwarded to the internal host or network. Unlike port forwarding, it exposes the full internal system, useful for servers behind a firewall. BINAT rules enable bidirectional translation for consistent incoming and outgoing connections.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this firewall NAT rule. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"log": schema.BoolAttribute{
				MarkdownDescription: "Log packets that are handled by this rule. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"sequence": schema.Int64Attribute{
				MarkdownDescription: "Specify the order of this NAT rule. Defaults to `1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"interface": schema.StringAttribute{
				MarkdownDescription: "Choose which interface this rule applies to. Defaults to `wan`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("wan"),
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Select `binat` (default) or `nat` here, when nets are equally sized `binat` is usually the best option. Using `nat` we can also map unequal sized networks. A `binat` rule specifies a bidirectional mapping between an external and internal network and can be used from both ends, `nat` only applies in one direction.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("binat"),
				Validators: []validator.String{
					stringvalidator.OneOf(
						"binat", "nat",
					),
				},
			},
			"source": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"net": schema.StringAttribute{
						MarkdownDescription: "Enter the internal IP address or CIDR for the 1:1 mapping. Aliases are only allowed in nat, not in binat type!",
						Required:            true,
					},
					"invert": schema.BoolAttribute{
						MarkdownDescription: "Use this option to invert the sense of the match. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
			"destination": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						map[string]attr.Type{
							"net":    types.StringType,
							"invert": types.BoolType,
						},
						map[string]attr.Value{
							"net":    types.StringValue("any"),
							"invert": types.BoolValue(false),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"net": schema.StringAttribute{
						MarkdownDescription: "The 1:1 mapping will only be used for connections to or from the specified destination. Hint: this is usually 'any'. Defaults to `any`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("any"),
					},
					"invert": schema.BoolAttribute{
						MarkdownDescription: "Use this option to invert the sense of the match. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
			"external_net": schema.StringAttribute{
				MarkdownDescription: "Enter the external subnet's starting address for the 1:1 mapping or network. This is the address or network the traffic will translate to/from.",
				Required:            true,
				Validators: []validator.String{
					ipOrCidrValidator,
				},
			},
			"nat_reflection": schema.StringAttribute{
				MarkdownDescription: "NAT reflection mode. One of `default`, `enable`, or `disable`. `default` means OPNsense uses the global firewall NAT reflection setting.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("default"),
				Validators: []validator.String{
					stringvalidator.OneOf("default", "enable", "disable"),
				},
			},
			"categories": schema.SetAttribute{
				MarkdownDescription: "Set of category IDs to apply. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed). Must be between 0 and 255 characters. Must be a character in set `[a-zA-Z0-9 .]`.",
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

func FirewallNATOneToOneDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "1:1 NAT maps a public IP or subnet to an internal private IP or subnet. All traffic to the public address is forwarded to the internal host or network. Unlike port forwarding, it exposes the full internal system, useful for servers behind a firewall. BINAT rules enable bidirectional translation for consistent incoming and outgoing connections.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this firewall NAT rule.",
				Computed:            true,
			},
			"log": dschema.BoolAttribute{
				MarkdownDescription: "Log packets that are handled by this rule.",
				Computed:            true,
			},
			"sequence": dschema.Int64Attribute{
				MarkdownDescription: "Specify the order of this NAT rule.",
				Computed:            true,
			},
			"interface": dschema.StringAttribute{
				MarkdownDescription: "Choose which interface this rule applies to.",
				Computed:            true,
			},
			"type": dschema.StringAttribute{
				MarkdownDescription: "Select `binat` (default) or `nat` here, when nets are equally sized `binat` is usually the best option. Using `nat` we can also map unequal sized networks. A `binat` rule specifies a bidirectional mapping between an external and internal network and can be used from both ends, `nat` only applies in one direction.",
				Computed:            true,
			},
			"source": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"net": dschema.StringAttribute{
						MarkdownDescription: "Enter the internal IP address, CIDR or alias for the 1:1 mapping.",
						Computed:            true,
					},
					"invert": dschema.BoolAttribute{
						MarkdownDescription: "Use this option to invert the sense of the match.",
						Computed:            true,
					},
				},
			},
			"destination": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"net": dschema.StringAttribute{
						MarkdownDescription: "The 1:1 mapping will only be used for connections to or from the specified destination.",
						Computed:            true,
					},
					"invert": dschema.BoolAttribute{
						MarkdownDescription: "Use this option to invert the sense of the match.",
						Computed:            true,
					},
				},
			},
			"external_net": dschema.StringAttribute{
				MarkdownDescription: "Enter the external subnet's starting address for the 1:1 mapping or network. This is the address or network the traffic will translate to/from.",
				Computed:            true,
			},
			"nat_reflection": dschema.StringAttribute{
				MarkdownDescription: "NAT reflection mode. One of `default`, `enable`, or `disable`. `default` means OPNsense uses the global firewall NAT reflection setting.",
				Computed:            true,
			},
			"categories": dschema.SetAttribute{
				MarkdownDescription: "Set of category IDs to apply.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Computed:            true,
			},
		},
	}
}

func convertFirewallNATOneToOneSchemaToStruct(d *FirewallNATOneToOneResourceModel) (*firewall.NatOneToOne, error) {

	// Parse 'Categories'
	var categoriesList []string
	d.Categories.ElementsAs(context.Background(), &categoriesList, false)

	// map Terraform "default" → API ""
	natReflection := d.NatReflection.ValueString()
	if natReflection == "default" {
		natReflection = ""
	}

	return &firewall.NatOneToOne{
		Enabled:           tools.BoolToString(d.Enabled.ValueBool()),
		Sequence:          tools.Int64ToString(d.Sequence.ValueInt64()),
		Interface:         api.SelectedMap(d.Interface.ValueString()),
		Type:              api.SelectedMap(d.Type.ValueString()),
		SourceNet:         d.Source.Net.ValueString(),
		SourceInvert:      tools.BoolToString(d.Source.Invert.ValueBool()),
		DestinationNet:    d.Destination.Net.ValueString(),
		DestinationInvert: tools.BoolToString(d.Destination.Invert.ValueBool()),
		ExternalNet:       d.ExternalNet.ValueString(),
		NatReflection:     api.SelectedMap(natReflection),
		Categories:        categoriesList,
		Log:               tools.BoolToString(d.Log.ValueBool()),
		Description:       d.Description.ValueString(),
	}, nil
}

func convertFirewallNATOneToOneStructToSchema(d *firewall.NatOneToOne) (*FirewallNATOneToOneResourceModel, error) {

	// map API "" → Terraform "default"
	natReflection := d.NatReflection.String()
	if natReflection == "" {
		natReflection = "default"
	}

	model := &FirewallNATOneToOneResourceModel{
		Enabled:   types.BoolValue(tools.StringToBool(d.Enabled)),
		Log:       types.BoolValue(tools.StringToBool(d.Log)),
		Sequence:  tools.StringToInt64Null(d.Sequence),
		Interface: types.StringValue(d.Interface.String()),
		Type:      types.StringValue(d.Type.String()),
		Source: &firewallLocationOneToOne{
			Net:    types.StringValue(d.SourceNet),
			Invert: types.BoolValue(tools.StringToBool(d.SourceInvert)),
		},
		Destination: &firewallLocationOneToOne{
			Net:    types.StringValue(d.DestinationNet),
			Invert: types.BoolValue(tools.StringToBool(d.DestinationInvert)),
		},
		ExternalNet:   types.StringValue(d.ExternalNet),
		NatReflection: types.StringValue(natReflection),
		Categories:    types.SetNull(types.StringType),
		Description:   tools.StringOrNull(d.Description),
	}

	// Parse 'Categories'
	var categoriesList []attr.Value
	for _, i := range d.Categories {
		categoriesList = append(categoriesList, basetypes.NewStringValue(i))
	}
	categoriesTypeList, _ := types.SetValue(types.StringType, categoriesList)
	model.Categories = categoriesTypeList

	return model, nil
}
