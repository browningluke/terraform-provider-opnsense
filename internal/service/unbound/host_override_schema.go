package unbound

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/unbound"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// hostOverrideResourceModel describes the resource data model.
type hostOverrideResourceModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Hostname    types.String `tfsdk:"hostname"`
	Domain      types.String `tfsdk:"domain"`
	Type        types.String `tfsdk:"type"`
	Server      types.String `tfsdk:"server"`
	Description types.String `tfsdk:"description"`

	MXPriority types.Int64  `tfsdk:"mx_priority"`
	MXDomain   types.String `tfsdk:"mx_host"`

	Id types.String `tfsdk:"id"`
}

func hostOverrideResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Host overrides can be used to change DNS results from client queries or to add custom DNS records.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable the override for this host. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Name of the host, without the domain part. Use `*` to create a wildcard entry.",
				Required:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Domain of the host, e.g. example.com",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of resource record. Available values: `A`, `AAAA`, `MX`. Defaults to `A`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("A"),
				Validators: []validator.String{
					stringvalidator.OneOf("A", "AAAA", "MX"),
				},
			},
			"server": schema.StringAttribute{
				MarkdownDescription: "IP address of the host, e.g. 192.168.100.100 or fd00:abcd::1. Must be set when `type` is `A` or `AAAA`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"mx_priority": schema.Int64Attribute{
				MarkdownDescription: "Priority of MX record, e.g. 10. Must be set when `type` is `MX`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.Expressions{
						path.MatchRoot("mx_host"),
					}...),
				},
			},
			"mx_host": schema.StringAttribute{
				MarkdownDescription: "Host name of MX host, e.g. mail.example.com. Must be set when `type` is `MX`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("mx_priority"),
					}...),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the host override.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func hostOverrideDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Host overrides can be used to change DNS results from client queries or to add custom DNS records.",

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
				MarkdownDescription: "Name of the host, without the domain part. Use `*` to create a wildcard entry.",
				Computed:            true,
			},
			"domain": dschema.StringAttribute{
				MarkdownDescription: "Domain of the host, e.g. example.com",
				Computed:            true,
			},
			"type": dschema.StringAttribute{
				MarkdownDescription: "Type of resource record. Available values: `A`, `AAAA`, `MX`.",
				Computed:            true,
			},
			"server": dschema.StringAttribute{
				MarkdownDescription: "IP address of the host, e.g. 192.168.100.100 or fd00:abcd::1.",
				Computed:            true,
			},
			"mx_priority": dschema.Int64Attribute{
				MarkdownDescription: "Priority of MX record, e.g. 10.",
				Computed:            true,
			},
			"mx_host": dschema.StringAttribute{
				MarkdownDescription: "Host name of MX host, e.g. mail.example.com.",
				Computed:            true,
			},
		},
	}
}

func convertHostOverrideSchemaToStruct(d *hostOverrideResourceModel) (*unbound.HostOverride, error) {
	return &unbound.HostOverride{
		Enabled:     tools.BoolToString(d.Enabled.ValueBool()),
		Hostname:    d.Hostname.ValueString(),
		Domain:      d.Domain.ValueString(),
		Type:        api.SelectedMap(d.Type.ValueString()),
		Server:      d.Server.ValueString(),
		MXDomain:    d.MXDomain.ValueString(),
		MXPriority:  tools.Int64ToStringNegative(d.MXPriority.ValueInt64()),
		Description: d.Description.ValueString(),
	}, nil
}

func convertHostOverrideStructToSchema(d *unbound.HostOverride) (*hostOverrideResourceModel, error) {
	return &hostOverrideResourceModel{
		Enabled:     types.BoolValue(tools.StringToBool(d.Enabled)),
		Hostname:    types.StringValue(d.Hostname),
		Domain:      types.StringValue(d.Domain),
		Type:        types.StringValue(d.Type.String()),
		Server:      types.StringValue(d.Server),
		MXPriority:  types.Int64Value(tools.StringToInt64(d.MXPriority)),
		MXDomain:    types.StringValue(d.MXDomain),
		Description: tools.StringOrNull(d.Description),
	}, nil
}
