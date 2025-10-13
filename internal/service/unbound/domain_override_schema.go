package unbound

import (
	"github.com/browningluke/opnsense-go/pkg/unbound"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// domainOverrideResourceModel describes the resource data model.
type domainOverrideResourceModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Domain      types.String `tfsdk:"domain"`
	Server      types.String `tfsdk:"server"`
	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func domainOverrideResourceSchema() schema.Schema {
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

func domainOverrideDataSourceSchema() dschema.Schema {
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

func convertDomainOverrideSchemaToStruct(d *domainOverrideResourceModel) (*unbound.DomainOverride, error) {
	return &unbound.DomainOverride{
		Enabled:     tools.BoolToString(d.Enabled.ValueBool()),
		Domain:      d.Domain.ValueString(),
		Server:      d.Server.ValueString(),
		Description: d.Description.ValueString(),
	}, nil
}

func convertDomainOverrideStructToSchema(d *unbound.DomainOverride) (*domainOverrideResourceModel, error) {
	return &domainOverrideResourceModel{
		Enabled:     types.BoolValue(tools.StringToBool(d.Enabled)),
		Domain:      types.StringValue(d.Domain),
		Server:      types.StringValue(d.Server),
		Description: tools.StringOrNull(d.Description),
	}, nil
}
