package dnsmasq

import (
	"regexp"
	"strconv"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/browningluke/terraform-provider-opnsense/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	LeaseTime       types.Int64  `tfsdk:"lease_time"`
	Tag             types.String `tfsdk:"tag"`
	IsIgnored       types.Bool   `tfsdk:"is_ignored"`
	Description     types.String `tfsdk:"description"`
	Comment         types.String `tfsdk:"comment"`

	Id types.String `tfsdk:"id"`
}

func hostResourceSchema() schema.Schema {
	return schema.Schema{
		Version:             1,
		MarkdownDescription: "Configure hosts override for dnsmasq.",
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Name of the host, without the domain part. Use \"*\" to create a wildcard entry.",
				Required:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Domain of the host, e.g. example.com",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"is_local_domain": schema.BoolAttribute{
				MarkdownDescription: "Set the above domain as local. This will configure this DNS server as authoritative; it will not forward queries to any upstream servers for this domain.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ip_addresses": schema.SetAttribute{
				MarkdownDescription: "IP addresses of the host, e.g. 192.168.100.100 or fd00:abcd::1. Can be multiple IPv4 and IPv6 addresses for dual stack configurations. Setting multiple addresses will automatically assign the best match based on the subnet of the interface receiving the DHCP Discover.",
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(validators.IP()),
					setvalidator.SizeAtLeast(1),
				},
			},
			"alias_records": schema.SetAttribute{
				MarkdownDescription: "Adds additional static A, AAAA and PTR records for the given alternative names (FQDN). Please note that these records are only created if IP addresses are configured in this host entry.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"cname_records": schema.SetAttribute{
				MarkdownDescription: "Adds additional CNAME records for the given alternative names (FQDN). Useful if this host entry has dynamic IPv4 and partial IPv6 addresses, as the CNAME record will point to the name instead of static IP addresses.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,62}$`), "Labels must be 1–63 characters, start with a letter or digit, and may contain only letters, digits, '-' or '_'.")),
				},
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Match the identifier of the client, e.g., DUID for DHCPv6. Setting the special character \"*\" will ignore the client identifier for DHCPv4 leases if a client offers both as choice.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^(?:\*|(?:[0-9A-Fa-f]{2}(?::[0-9A-Fa-f]{2})+))$`), "Value must be a colon-separated hexadecimal sequence (e.g., 01:02:f3) or \"*\"."),
				},
			},
			"hardware_addresses": schema.SetAttribute{
				MarkdownDescription: "Match the hardware address of the client. Can be multiple addresses, e.g., if the client has multiple network cards. Though keep in mind that Dnsmasq cannot assume which address is the correct one when multiple send DHCP Discover at the same time.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"lease_time": schema.Int64Attribute{
				MarkdownDescription: "Defines how long the addresses (leases) given out by the server are valid (in seconds). Set 0 for infinite.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "Optional tag to set for requests matching this range which can be used to selectively match DHCP options. Can be left empty if options with an interface tag exist, since the client automatically receives this tag based on the interface receiving the DHCP Discover.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"is_ignored": schema.BoolAttribute{
				MarkdownDescription: "Ignore any DHCP packets of this host. Useful if it should get served by a different DHCP server.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "You may enter a description here for your reference (not parsed)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"comment": schema.StringAttribute{
				MarkdownDescription: "You may enter a comment here for your reference (not parsed)",
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
				MarkdownDescription: "Name of the host, without the domain part. Use \"*\" to create a wildcard entry.",
				Computed:            true,
			},
			"domain": dschema.StringAttribute{
				MarkdownDescription: "Domain of the host, e.g. example.com",
				Computed:            true,
			},
			"is_local_domain": dschema.BoolAttribute{
				MarkdownDescription: "Set the above domain as local. This will configure this DNS server as authoritative; it will not forward queries to any upstream servers for this domain.",
				Computed:            true,
			},
			"ip_addresses": dschema.SetAttribute{
				MarkdownDescription: "IP addresses of the host, e.g. 192.168.100.100 or fd00:abcd::1. Can be multiple IPv4 and IPv6 addresses for dual stack configurations. Setting multiple addresses will automatically assign the best match based on the subnet of the interface receiving the DHCP Discover.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"alias_records": dschema.SetAttribute{
				MarkdownDescription: "Adds additional static A, AAAA and PTR records for the given alternative names (FQDN). Please note that these records are only created if IP addresses are configured in this host entry.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"cname_records": dschema.SetAttribute{
				MarkdownDescription: "Adds additional CNAME records for the given alternative names (FQDN). Useful if this host entry has dynamic IPv4 and partial IPv6 addresses, as the CNAME record will point to the name instead of static IP addresses.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"client_id": dschema.StringAttribute{
				MarkdownDescription: "Match the identifier of the client, e.g., DUID for DHCPv6. Setting the special character \"*\" will ignore the client identifier for DHCPv4 leases if a client offers both as choice.",
				Computed:            true,
			},
			"hardware_addresses": dschema.SetAttribute{
				MarkdownDescription: "Match the hardware address of the client. Can be multiple addresses, e.g., if the client has multiple network cards. Though keep in mind that Dnsmasq cannot assume which address is the correct one when multiple send DHCP Discover at the same time.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"lease_time": dschema.Int64Attribute{
				MarkdownDescription: "Defines how long the addresses (leases) given out by the server are valid (in seconds). Set 0 for infinite.",
				Computed:            true,
			},
			"tag": dschema.StringAttribute{
				MarkdownDescription: "Optional tag to set for requests matching this range which can be used to selectively match DHCP options. Can be left empty if options with an interface tag exist, since the client automatically receives this tag based on the interface receiving the DHCP Discover.",
				Computed:            true,
			},
			"is_ignored": dschema.BoolAttribute{
				MarkdownDescription: "Ignore any DHCP packets of this host. Useful if it should get served by a different DHCP server.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "You may enter a description here for your reference (not parsed)",
				Computed:            true,
			},
			"comment": dschema.StringAttribute{
				MarkdownDescription: "You may enter a comment here for your reference (not parsed)",
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
	leaseTime := ""
	if !d.LeaseTime.IsNull() {
		leaseTime = strconv.FormatInt(d.LeaseTime.ValueInt64(), 10)
	}
	return &dnsmasq.Host{
		Hostname:          d.Hostname.ValueString(),
		Domain:            d.Domain.ValueString(),
		IsLocalDomain:     tools.BoolToString(d.IsLocalDomain.ValueBool()),
		IpAddresses:       tools.SetToStringSlice(d.IpAddresses),
		AliasRecords:      tools.SetToStringSlice(d.AliasRecords),
		CnameRecords:      tools.SetToStringSlice(d.CnameRecords),
		ClientId:          d.ClientID.ValueString(),
		HardwareAddresses: tools.SetToStringSlice(d.HarwareAdresses),
		LeaseTime:         leaseTime,
		Tag:               api.SelectedMap(d.Tag.ValueString()),
		IsIgnored:         tools.BoolToString(d.IsIgnored.ValueBool()),
		Description:       d.Description.ValueString(),
		Comments:          d.Comment.ValueString(),
	}, nil
}

func convertHostStructToSchema(d *dnsmasq.Host) (*hostResourceModel, error) {
	var leaseTime types.Int64

	if d.LeaseTime == "" {
		leaseTime = types.Int64Null()
	} else {
		leaseTime = types.Int64Value(tools.StringToInt64(d.LeaseTime))
	}
	return &hostResourceModel{
		Hostname:        types.StringValue(d.Hostname),
		Domain:          types.StringValue(d.Domain),
		IsLocalDomain:   types.BoolValue(tools.StringToBool(d.IsLocalDomain)),
		IpAddresses:     tools.StringSliceToSet(d.IpAddresses),
		AliasRecords:    tools.StringSliceToSet(d.AliasRecords),
		CnameRecords:    tools.StringSliceToSet(d.CnameRecords),
		ClientID:        types.StringValue(d.ClientId),
		HarwareAdresses: tools.StringSliceToSet(d.HardwareAddresses),
		LeaseTime:       leaseTime,
		Tag:             types.StringValue(d.Tag.String()),
		IsIgnored:       types.BoolValue(tools.StringToBool(d.IsIgnored)),
		Description:     types.StringValue(d.Description),
		Comment:         types.StringValue(d.Comments),
	}, nil
}
