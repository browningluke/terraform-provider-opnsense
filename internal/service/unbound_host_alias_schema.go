package service

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
