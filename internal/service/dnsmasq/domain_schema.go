package dnsmasq

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/browningluke/terraform-provider-opnsense/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type domainResourceModel struct {
	Sequence      types.Int64  `tfsdk:"sequence"`
	Domain        types.String `tfsdk:"domain"`
	FirewallAlias types.String `tfsdk:"firewall_alias"`
	SourceIp      types.String `tfsdk:"srcip"`
	Port          types.Int64  `tfsdk:"port"`
	Ip            types.String `tfsdk:"ip"`
	Description   types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func domainResourceSchema() schema.Schema {
	return schema.Schema{
		Version:             1,
		MarkdownDescription: "Configure domain information for dnsmasq.",
		Attributes: map[string]schema.Attribute{
			"sequence": schema.Int64Attribute{
				MarkdownDescription: "Sort with a sequence number, e.g., for strict processing order when using the \"Query DNS servers sequentially\" option in general settings.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 99999),
				},
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Domain to override (NOTE: this does not have to be a valid TLD!)",
				Required:            true,
			},
			"firewall_alias": schema.StringAttribute{
				MarkdownDescription: "Choose an \"external (advanced)\" type alias from \"Firewall - Aliases\". Whenever a client successfully resolves the domain, the resolved IP addresses will be automatically added to the chosen alias. Adding a domain will also add all IP addresses of resolved subdomains. Please note that DNS record TTL is not evaluated; once an IP address is added, it will stay permanently, or until manually flushed in \"Firewall - Diagnostics - Aliases\", or until removed automatically when setting an expiration on the alias.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"srcip": schema.StringAttribute{
				MarkdownDescription: "Source IP address for queries to the DNS server for the override domain. Best to leave empty",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					validators.IP(),
				},
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Specify a non standard port number here, leave blank for default",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"ip": schema.StringAttribute{
				MarkdownDescription: "IP address of the authoritative DNS server for this domain, leave empty to prevent lookups for this domain",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					validators.IP(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "You may enter a description here for your reference (not parsed)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID of the domain.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func domainDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure domain information for dnsmasq.",
		Attributes: map[string]dschema.Attribute{
			"sequence": schema.Int64Attribute{
				MarkdownDescription: "Sort with a sequence number, e.g., for strict processing order when using the \"Query DNS servers sequentially\" option in general settings.",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Domain to override (NOTE: this does not have to be a valid TLD!)",
				Computed:            true,
			},
			"firewall_alias": schema.StringAttribute{
				MarkdownDescription: "Choose an \"external (advanced)\" type alias from \"Firewall - Aliases\". Whenever a client successfully resolves the domain, the resolved IP addresses will be automatically added to the chosen alias. Adding a domain will also add all IP addresses of resolved subdomains. Please note that DNS record TTL is not evaluated; once an IP address is added, it will stay permanently, or until manually flushed in \"Firewall - Diagnostics - Aliases\", or until removed automatically when setting an expiration on the alias.",
				Computed:            true,
			},
			"srcip": schema.StringAttribute{
				MarkdownDescription: "Source IP address for queries to the DNS server for the override domain. Best to leave empty",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Specify a non standard port number here, leave blank for default",
				Computed:            true,
			},
			"ip": schema.StringAttribute{
				MarkdownDescription: "IP address of the authoritative DNS server for this domain, leave empty to prevent lookups for this domain",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "You may enter a description here for your reference (not parsed)",
				Computed:            true,
			},
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the domain.",
				Required:            true,
			},
		},
	}
}

func convertDomainSchemaToStruct(d *domainResourceModel) (*dnsmasq.Domain, error) {
	domain := &dnsmasq.Domain{
		Sequence:      d.Sequence.String(),
		Domain:        d.Domain.ValueString(),
		FirewallAlias: api.SelectedMap(d.FirewallAlias.ValueString()),
		SourceIp:      d.SourceIp.ValueString(),
		Ip:            d.Ip.ValueString(),
		Description:   d.Description.ValueString(),
	}
	if !d.Port.IsNull() {
		portValue := d.Port.ValueInt64()
		domain.Port = tools.Int64ToString(portValue)
	}
	return domain, nil
}

func convertDomainStructToSchema(d *dnsmasq.Domain) (*domainResourceModel, error) {
	domain := &domainResourceModel{
		Domain:        types.StringValue(d.Domain),
		FirewallAlias: types.StringValue(d.FirewallAlias.String()),
		SourceIp:      types.StringValue(d.SourceIp),
		Ip:            types.StringValue(d.Ip),
		Description:   types.StringValue(d.Description),
	}

	if d.Sequence != "" {
		domain.Sequence = types.Int64Value(tools.StringToInt64(d.Sequence))
	} else {
		domain.Sequence = types.Int64Null()
	}
	if d.Port != "" {
		domain.Port = types.Int64Value(tools.StringToInt64(d.Port))
	} else {
		domain.Port = types.Int64Null()
	}

	return domain, nil
}
