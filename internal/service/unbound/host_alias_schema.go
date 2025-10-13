package unbound

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/unbound"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// hostAliasResourceModel describes the resource data model.
type hostAliasResourceModel struct {
	Override    types.String `tfsdk:"override"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Hostname    types.String `tfsdk:"hostname"`
	Domain      types.String `tfsdk:"domain"`
	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func hostAliasResourceSchema() schema.Schema {
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

func hostAliasDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Host aliases can be used to create alternative names for a Host.",

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
			"hostname": dschema.StringAttribute{
				MarkdownDescription: "Name of the host, without the domain part.",
				Computed:            true,
			},
			"domain": dschema.StringAttribute{
				MarkdownDescription: "Domain of the host, e.g. example.com",
				Computed:            true,
			},
			"override": dschema.StringAttribute{
				MarkdownDescription: "The associated host override to apply this alias on.",
				Computed:            true,
			},
		},
	}
}

func convertHostAliasSchemaToStruct(d *hostAliasResourceModel) (*unbound.HostAlias, error) {
	return &unbound.HostAlias{
		Enabled:     tools.BoolToString(d.Enabled.ValueBool()),
		Host:        api.SelectedMap(d.Override.ValueString()),
		Hostname:    d.Hostname.ValueString(),
		Domain:      d.Domain.ValueString(),
		Description: d.Description.ValueString(),
	}, nil
}

func convertHostAliasStructToSchema(d *unbound.HostAlias) (*hostAliasResourceModel, error) {
	return &hostAliasResourceModel{
		Enabled:     types.BoolValue(tools.StringToBool(d.Enabled)),
		Hostname:    types.StringValue(d.Hostname),
		Domain:      types.StringValue(d.Domain),
		Description: tools.StringOrNull(d.Description),
		Override:    types.StringValue(d.Host.String()),
	}, nil
}
