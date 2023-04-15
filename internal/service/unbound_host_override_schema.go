package service

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-opnsense/internal/opnsense"
)

// UnboundHostOverrideResourceModel describes the resource data model.
type UnboundHostOverrideResourceModel struct {
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

func unboundHostOverrideResourceSchema() schema.Schema {
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

func convertUnboundHostOverrideSchemaToStruct(d *UnboundHostOverrideResourceModel) (*opnsense.UnboundHostOverride, error) {
	// Parse 'Enabled'
	var enabled string
	if d.Enabled.ValueBool() {
		enabled = "1"
	} else {
		enabled = "0"
	}

	// Parse 'MXPriority'
	mxPriority := fmt.Sprintf("%d", d.MXPriority.ValueInt64())
	if d.MXPriority.ValueInt64() == -1 {
		mxPriority = ""
	}

	return &opnsense.UnboundHostOverride{
		Enabled:     enabled,
		Hostname:    d.Hostname.ValueString(),
		Domain:      d.Domain.ValueString(),
		Type:        d.Type.ValueString(),
		Server:      d.Server.ValueString(),
		MXDomain:    d.MXDomain.ValueString(),
		MXPriority:  mxPriority,
		Description: d.Description.ValueString(),
	}, nil
}

func convertUnboundHostOverrideStructToSchema(d *opnsense.UnboundHostOverride) (*UnboundHostOverrideResourceModel, error) {
	model := &UnboundHostOverrideResourceModel{
		Enabled:     types.BoolValue(false),
		Hostname:    types.StringValue(d.Hostname),
		Domain:      types.StringValue(d.Domain),
		Type:        types.StringValue(d.Type),
		Server:      types.StringValue(d.Server),
		MXPriority:  types.Int64Value(-1),
		MXDomain:    types.StringValue(d.MXDomain),
		Description: types.StringNull(),
	}

	// Parse 'Enabled'
	if d.Enabled == "1" {
		model.Enabled = types.BoolValue(true)
	}

	// Parse 'MXPriority'
	mxPriority, err := strconv.ParseInt(d.MXPriority, 10, 64)
	if err == nil {
		model.MXPriority = types.Int64Value(mxPriority)
	}

	// Parse 'Description'
	if d.Description != "" {
		model.Description = types.StringValue(d.Description)
	}

	return model, nil
}
