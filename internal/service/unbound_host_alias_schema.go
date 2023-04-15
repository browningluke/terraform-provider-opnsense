package service

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/opnsense"
)

// UnboundHostAliasResourceModel describes the resource data model.
type UnboundHostAliasResourceModel struct {
	Override    types.String `tfsdk:"override"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Hostname    types.String `tfsdk:"hostname"`
	Domain      types.String `tfsdk:"domain"`
	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func unboundHostAliasResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Host aliases can be used to create alternative names for a Host",

		Attributes: map[string]schema.Attribute{
			"override": schema.StringAttribute{
				MarkdownDescription: "The associated host override to apply this alias on.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this alias for the selected host. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Name of the host, without the domain part. Use `*` to create a wildcard entry.",
				Required:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Domain of the host, e.g. example.com.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the host alias.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func convertUnboundHostAliasSchemaToStruct(d *UnboundHostAliasResourceModel) (*opnsense.UnboundHostAlias, error) {
	// Parse 'Enabled'
	var enabled string
	if d.Enabled.ValueBool() {
		enabled = "1"
	} else {
		enabled = "0"
	}

	return &opnsense.UnboundHostAlias{
		Enabled:     enabled,
		Host:        d.Override.ValueString(),
		Hostname:    d.Hostname.ValueString(),
		Domain:      d.Domain.ValueString(),
		Description: d.Description.ValueString(),
	}, nil
}

func convertUnboundHostAliasStructToSchema(d *opnsense.UnboundHostAlias) (*UnboundHostAliasResourceModel, error) {
	model := &UnboundHostAliasResourceModel{
		Enabled:     types.BoolValue(false),
		Hostname:    types.StringValue(d.Hostname),
		Domain:      types.StringValue(d.Domain),
		Description: types.StringValue(d.Description),
	}

	// Parse 'Enabled'
	if d.Enabled == "1" {
		model.Enabled = types.BoolValue(true)
	}

	// Parse 'Description'
	if d.Description != "" {
		model.Description = types.StringValue(d.Description)
	}

	return model, nil
}
