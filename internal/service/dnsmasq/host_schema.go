package dnsmasq

import (
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type hostResourceModel struct {
	Hostname        types.String `tfsdk:"hostname"`
	Domain          types.String `tfsdk:"domain"`
	IsLocalDomain   types.Bool   `tfsdk:"is_local_domain"`
	IpAddresses     types.Set    `tfsdk:"ip_addresses"`
	AliasRecords    types.Set    `tfsdk:"alias_records"`
	CnameRecords    types.Set    `tfsdk:"cname_records"`
	ClientID        types.String `tfsdk:"client_id"`
	HarwareAdresses types.Set    `tfsdk:"hardware_addresses"`
	IsIgnored       types.Bool   `tfsdk:"is_ignored"`
	Description     types.String `tfsdk:"description"`
	Comment         types.String `tfsdk:"comment"`

	Id types.String `tfsdk:"id"`
}

func hostResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure hosts override for dnsmasq.",
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Name of the host, without the domain part.",
				Required:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Domain of the host.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"is_local_domain": schema.BoolAttribute{
				MarkdownDescription: "Whether this is a local domain.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ip_addresses": schema.SetAttribute{
				MarkdownDescription: "IP addresses of the host.",
				ElementType:         types.StringType,
				Required:            true,
			},
			"alias_records": schema.SetAttribute{
				MarkdownDescription: "Alias records of the host.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"cname_records": schema.SetAttribute{
				MarkdownDescription: "CNAME records of the host.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID of the host.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"hardware_addresses": schema.SetAttribute{
				MarkdownDescription: "Hardware addresses of the host.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"is_ignored": schema.BoolAttribute{
				MarkdownDescription: "Whether DHCP packet is ignored for this host.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"comment": schema.StringAttribute{
				MarkdownDescription: "Optional comment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID of the host.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func hostDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure hosts override for dnsmasq.",
		Attributes: map[string]dschema.Attribute{
			"hostname": dschema.StringAttribute{
				MarkdownDescription: "Name of the host, without the domain part.",
				Required:            true,
			},
			"domain": dschema.StringAttribute{
				MarkdownDescription: "Domain of the host.",
				Computed:            true,
			},
			"is_local_domain": dschema.BoolAttribute{
				MarkdownDescription: "Whether this is a local domain.",
				Computed:            true,
			},
			"ip_addresses": dschema.SetAttribute{
				MarkdownDescription: "IP addresses of the host.",
				ElementType:         types.StringType,
				Required:            true,
			},
			"alias_records": dschema.SetAttribute{
				MarkdownDescription: "Alias records of the host.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"cname_records": dschema.SetAttribute{
				MarkdownDescription: "CNAME records of the host.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"client_id": dschema.StringAttribute{
				MarkdownDescription: "Client identifier of the host.",
				Computed:            true,
			},
			"hardware_addresses": dschema.SetAttribute{
				MarkdownDescription: "Hardware addresses of the host.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"is_ignored": dschema.BoolAttribute{
				MarkdownDescription: "Whether DHCP packet is ignored for this host.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description.",
				Computed:            true,
			},
			"comment": dschema.StringAttribute{
				MarkdownDescription: "Optional comment.",
				Computed:            true,
			},
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the host.",
				Required:            true,
			},
		},
	}
}

func convertHostSchemaToStruct(d *hostResourceModel) (*dnsmasq.Host, error) {

	return &dnsmasq.Host{
		Hostname:          d.Hostname.ValueString(),
		Domain:            d.Domain.ValueString(),
		IsLocalDomain:     tools.BoolToString(d.IsLocalDomain.ValueBool()),
		IpAddresses:       tools.SetToStringSlice(d.IpAddresses),
		AliasRecords:      tools.SetToStringSlice(d.AliasRecords),
		CnameRecords:      tools.SetToStringSlice(d.CnameRecords),
		ClientId:          d.ClientID.ValueString(),
		HardwareAddresses: tools.SetToStringSlice(d.HarwareAdresses),
		IsIgnored:         tools.BoolToString(d.IsIgnored.ValueBool()),
		Description:       d.Description.ValueString(),
		Comments:          d.Comment.ValueString(),
	}, nil
}

func convertHostStructToSchema(d *dnsmasq.Host) (*hostResourceModel, error) {
	return &hostResourceModel{
		Hostname:        types.StringValue(d.Hostname),
		Domain:          types.StringValue(d.Domain),
		IsLocalDomain:   types.BoolValue(tools.StringToBool(d.IsLocalDomain)),
		IpAddresses:     tools.StringSliceToSet(d.IpAddresses),
		AliasRecords:    tools.StringSliceToSet(d.AliasRecords),
		CnameRecords:    tools.StringSliceToSet(d.CnameRecords),
		ClientID:        types.StringValue(d.ClientId),
		HarwareAdresses: tools.StringSliceToSet(d.HardwareAddresses),
		IsIgnored:       types.BoolValue(tools.StringToBool(d.IsIgnored)),
		Description:     types.StringValue(d.Description),
		Comment:         types.StringValue(d.Comments),
	}, nil
}
