package service

import (
	"github.com/browningluke/opnsense-go/pkg/unbound"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UnboundDomainOverrideResourceModel describes the resource data model.
type UnboundDomainOverrideResourceModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Domain      types.String `tfsdk:"domain"`
	Server      types.String `tfsdk:"server"`
	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func unboundDomainOverrideResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Domain overrides can be used to forward queries for specific domains (and subsequent subdomains) to local or remote DNS servers.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this domain override. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Domain to override (NOTE: this does not have to be a valid TLD!), e.g. `test` or `mycompany.localdomain` or `1.168.192.in-addr.arpa`.",
				Required:            true,
			},
			"server": schema.StringAttribute{
				MarkdownDescription: "IP address of the authoritative DNS server for this domain, e.g. `192.168.100.100`. To use a nondefault port for communication, append an `@` with the port number.",
				Required:            true,
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

func UnboundDomainOverrideDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Domain overrides can be used to forward queries for specific domains (and subsequent subdomains) to local or remote DNS servers.",

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
			"domain": dschema.StringAttribute{
				MarkdownDescription: "Domain to override (NOTE: this does not have to be a valid TLD!), e.g. `test` or `mycompany.localdomain` or `1.168.192.in-addr.arpa`.",
				Computed:            true,
			},
			"server": dschema.StringAttribute{
				MarkdownDescription: "IP address of the authoritative DNS server for this domain, e.g. `192.168.100.100`.",
				Computed:            true,
			},
		},
	}
}

func convertUnboundDomainOverrideSchemaToStruct(d *UnboundDomainOverrideResourceModel) (*unbound.DomainOverride, error) {
	// Parse 'Enabled'
	var enabled string
	if d.Enabled.ValueBool() {
		enabled = "1"
	} else {
		enabled = "0"
	}

	return &unbound.DomainOverride{
		Enabled:     enabled,
		Domain:      d.Domain.ValueString(),
		Server:      d.Server.ValueString(),
		Description: d.Description.ValueString(),
	}, nil
}

func convertUnboundDomainOverrideStructToSchema(d *unbound.DomainOverride) (*UnboundDomainOverrideResourceModel, error) {
	model := &UnboundDomainOverrideResourceModel{
		Enabled:     types.BoolValue(false),
		Domain:      types.StringValue(d.Domain),
		Server:      types.StringValue(d.Server),
		Description: types.StringNull(),
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
