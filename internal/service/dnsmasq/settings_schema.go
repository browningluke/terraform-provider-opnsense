package dnsmasq

import (
	"context"
	"fmt"
	"sort"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// settingsDNSBlock contains the DNS-related dnsmasq settings.
type settingsDNSBlock struct {
	Port                 types.Int64 `tfsdk:"port"`
	DNSSEC               types.Bool  `tfsdk:"dnssec"`
	NoHostLookup         types.Bool  `tfsdk:"no_host_lookup"`
	LogQueries           types.Bool  `tfsdk:"log_queries"`
	MaxConcurrentQueries types.Int64 `tfsdk:"max_concurrent_queries"`
	CacheSize            types.Int64 `tfsdk:"cache_size"`
	LocalEntryTTL        types.Int64 `tfsdk:"local_entry_ttl"`
	NoIdent              types.Bool  `tfsdk:"no_ident"`
}

// settingsDNSQueryForwardingBlock contains DNS query forwarding settings.
type settingsDNSQueryForwardingBlock struct {
	QuerySequentially          types.Bool   `tfsdk:"query_sequentially"`
	RequireDomain              types.Bool   `tfsdk:"require_domain"`
	DoNotForwardSystemDNS      types.Bool   `tfsdk:"do_not_forward_system_dns"`
	DoNotForwardPrivateReverse types.Bool   `tfsdk:"do_not_forward_private_reverse"`
	AddMAC                     types.String `tfsdk:"add_mac"`
	AddSubnet                  types.Bool   `tfsdk:"add_subnet"`
	StripSubnet                types.Bool   `tfsdk:"strip_subnet"`
}

// settingsDHCPBlock contains the general DHCP settings.
type settingsDHCPBlock struct {
	NoDHCPInterfaces      types.Set    `tfsdk:"no_dhcp_interfaces"`
	FQDN                  types.Bool   `tfsdk:"fqdn"`
	DefaultDomain         types.String `tfsdk:"default_domain"`
	LocalDomain           types.Bool   `tfsdk:"local_domain"`
	MaxLeases             types.Int64  `tfsdk:"max_leases"`
	Authoritative         types.Bool   `tfsdk:"authoritative"`
	ReplyDelay            types.Int64  `tfsdk:"reply_delay"`
	RegisterFirewallRules types.Bool   `tfsdk:"register_firewall_rules"`
	RouterAdvertisements  types.Bool   `tfsdk:"router_advertisements"`
	HostPing              types.Bool   `tfsdk:"host_ping"`
	DisableHASync         types.Bool   `tfsdk:"disable_ha_sync"`
	LogDhcp               types.Bool   `tfsdk:"log_dhcp"`
	LogQuiet              types.Bool   `tfsdk:"log_quiet"`
}

// settingsLegacyBlock contains legacy dnsmasq integration settings.
type settingsLegacyBlock struct {
	RegisterISCDHCP4Leases    types.Bool   `tfsdk:"register_isc_dhcp4_leases"`
	DHCPDomainOverride        types.String `tfsdk:"dhcp_domain_override"`
	RegisterDHCPStaticMapping types.Bool   `tfsdk:"register_dhcp_static_mappings"`
	PreferDHCP                types.Bool   `tfsdk:"prefer_dhcp"`
}

// settingsResourceModel describes the resource data model.
//
// This is a SINGLETON resource — it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type settingsResourceModel struct {
	Id types.String `tfsdk:"id"`

	Enabled                types.Bool `tfsdk:"enabled"`
	Interface              types.Set  `tfsdk:"interface"`
	StrictInterfaceBinding types.Bool `tfsdk:"strict_interface_binding"`

	DNS                *settingsDNSBlock                `tfsdk:"dns"`
	DNSQueryForwarding *settingsDNSQueryForwardingBlock `tfsdk:"dns_query_forwarding"`
	DHCP               *settingsDHCPBlock               `tfsdk:"dhcp"`
	Legacy             *settingsLegacyBlock             `tfsdk:"legacy"`
}

func settingsResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages the general dnsmasq settings. This is a singleton resource that manages existing upstream configuration and cannot be created or destroyed.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `dnsmasq_settings`. Use this value when importing: `terraform import opnsense_dnsmasq_settings.settings dnsmasq_settings`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("dnsmasq_settings"),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable the dnsmasq service.",
				Optional:            true,
				Computed:            true,
			},
			"interface": schema.SetAttribute{
				MarkdownDescription: "Interface IPs used to respond to queries from clients. If no interfaces are selected, Dnsmasq will listen on all available IPv4 and IPv6 addresses by default. However, DHCP related firewall rules will only be added for explicitly selected interfaces, never for all interfaces.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"strict_interface_binding": schema.BoolAttribute{
				MarkdownDescription: "By default we bind the wildcard address, even when listening on some interfaces. Requests that should not be handled are discarded, this has the advantage of working even when interfaces come and go and change address. This option forces binding to only the interfaces we are listening on, which is less stable in non static environments.",
				Optional:            true,
				Computed:            true,
			},

			"dns": schema.SingleNestedAttribute{
				MarkdownDescription: "DNS server settings.",
				Optional:            true,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"port": schema.Int64Attribute{
						MarkdownDescription: "The port used for responding to DNS queries. It should normally be left blank unless another service needs to bind to TCP/UDP port 53. Setting this to zero (0) completely disables DNS function",
						Optional:            true,
						Computed:            true,
						Validators:          []validator.Int64{int64validator.Between(0, 65535)},
					},
					"dnssec": schema.BoolAttribute{
						MarkdownDescription: "Enable DNSSEC validation.",
						Optional:            true,
						Computed:            true,
					},
					"no_host_lookup": schema.BoolAttribute{
						MarkdownDescription: "Do not use `/etc/hosts` for DNS lookups.",
						Optional:            true,
						Computed:            true,
					},
					"log_queries": schema.BoolAttribute{
						MarkdownDescription: "Append the configured domain to simple hostnames read from hosts files.",
						Optional:            true,
						Computed:            true,
					},
					"max_concurrent_queries": schema.Int64Attribute{
						MarkdownDescription: "Set the maximum number of concurrent DNS queries. On configurations with tight resources, this value may need to be reduced.",
						Optional:            true,
						Computed:            true,
						Validators:          []validator.Int64{int64validator.AtLeast(0)},
					},
					"cache_size": schema.Int64Attribute{
						MarkdownDescription: "Set the size of the cache. Setting the cache size to zero disables caching. Please note that huge cache size impacts performance.",
						Optional:            true,
						Computed:            true,
						Validators:          []validator.Int64{int64validator.AtLeast(0)},
					},
					"local_entry_ttl": schema.Int64Attribute{
						MarkdownDescription: "This option allows a time-to-live (in seconds) to be given for local DNS entries, i.e. /etc/hosts or DHCP leases. This will reduce the load on the server at the expense of clients using stale data under some circumstances. A value of zero will disable client-side caching.",
						Optional:            true,
						Computed:            true,
						Validators:          []validator.Int64{int64validator.AtLeast(0)},
					},
					"no_ident": schema.BoolAttribute{
						MarkdownDescription: "Do not respond to class CHAOS and type TXT in domain bind queries. Without this option being set, the cache statistics are also available in the DNS as answers to queries of class CHAOS and type TXT in domain bind.",
						Optional:            true,
						Computed:            true,
					},
				},
			},

			"dns_query_forwarding": schema.SingleNestedAttribute{
				MarkdownDescription: "DNS query forwarding settings.",
				Optional:            true,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"query_sequentially": schema.BoolAttribute{
						MarkdownDescription: "If this option is set, we will query the DNS servers sequentially in the order specified (System: General Setup: DNS Servers), rather than all at once in parallel.",
						Optional:            true,
						Computed:            true,
					},
					"require_domain": schema.BoolAttribute{
						MarkdownDescription: "If this option is set, we will not forward A or AAAA queries for plain names, without dots or domain parts, to upstream name servers. If the name is not known from /etc/hosts or DHCP then a \"not found\" answer is returned.",
						Optional:            true,
						Computed:            true,
					},
					"do_not_forward_system_dns": schema.BoolAttribute{
						MarkdownDescription: "If this option is set, DNS forwarding to system nameservers (defined in System: General Setup: DNS Servers) will be disabled. Upstream servers defined in Services: Dnsmasq DNS & DHCP: Domains will still be used. This option is recommended when Unbound forwards local domain queries to Dnsmasq, so that all queries terminate without further lookups if they are unknown.",
						Optional:            true,
						Computed:            true,
					},
					"do_not_forward_private_reverse": schema.BoolAttribute{
						MarkdownDescription: "If this option is set, we will not forward reverse DNS lookups (PTR) for private addresses (RFC 1918) to upstream name servers. Any entries in the Domain Overrides section forwarding private \"n.n.n.in-addr.arpa\" names to a specific server are still forwarded. If the IP to name is not known from /etc/hosts, DHCP or a specific domain override then a \"not found\" answer is immediately returned.",
						Optional:            true,
						Computed:            true,
					},
					"add_mac": schema.StringAttribute{
						MarkdownDescription: "Add the MAC address of the requestor to DNS queries which are forwarded upstream. The MAC address will only be added if the upstream DNS Server is in the same subnet as the requestor. Since this is not standardized, it should be considered experimental. This is useful for selective DNS filtering on the upstream DNS server.",
						Optional:            true,
						Computed:            true,
						Validators:          []validator.String{stringvalidator.OneOf("", "standard", "base64", "text")},
					},
					"add_subnet": schema.BoolAttribute{
						MarkdownDescription: "Add the real client IPv4 and IPv6 addresses (add-subnet=32,128) to DNS queries which are forwarded upstream. Be careful setting this option as it can undermine privacy. This is useful for selective DNS filtering on the upstream DNS server.",
						Optional:            true,
						Computed:            true,
					},
					"strip_subnet": schema.BoolAttribute{
						MarkdownDescription: "Strip the subnet received by a downstream DNS server. If \"Add subnet\" is used and the downstream DNS server already added a subnet, Dnsmasq will not replace it without setting \"Strip subnet\".",
						Optional:            true,
						Computed:            true,
					},
				},
			},

			"dhcp": schema.SingleNestedAttribute{
				MarkdownDescription: "General DHCP settings.",
				Optional:            true,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"no_dhcp_interfaces": schema.SetAttribute{
						MarkdownDescription: "Do not provide DHCP, TFTP or router advertisement on the specified interfaces, but do provide DNS service. Please note that Dnsmasq continues to listen on the default DHCP ports as long as any DHCP ranges are configured; setting this option only ignores these packets on the selected interfaces.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
					},
					"fqdn": schema.BoolAttribute{
						MarkdownDescription: "If disabled, registers the unqualified names of DHCP clients into the DNS (e.g. 'smartphone'), in which case they should be unique. If enabled, the qualified name (e.g. 'smartphone.lan.internal') is registered. This option must be enabled if you are forwarding from Unbound to Dnsmasq for specific local domains.",
						Optional:            true,
						Computed:            true,
					},
					"default_domain": schema.StringAttribute{
						MarkdownDescription: "To ensure that all names have a domain part, there must be a default domain specified when DHCP FQDN is set. Leave empty to use the system domain.",
						Optional:            true,
						Computed:            true,
					},
					"local_domain": schema.BoolAttribute{
						MarkdownDescription: "Sets all DHCP domains as local. This will configure this DNS server as authoritative; it will not forward queries to any upstream servers for these domains.",
						Optional:            true,
						Computed:            true,
					},
					"max_leases": schema.Int64Attribute{
						MarkdownDescription: "Limits Dnsmasq to the specified maximum number of DHCP leases. This limit is to prevent DoS attacks from hosts which create thousands of leases and use lots of memory in the Dnsmasq process.",
						Optional:            true,
						Computed:            true,
						Validators:          []validator.Int64{int64validator.AtLeast(0)},
					},
					"authoritative": schema.BoolAttribute{
						MarkdownDescription: "Should be set when Dnsmasq is definitely the only DHCP server on a network. For DHCPv4, it changes the behaviour from strict RFC compliance so that DHCP requests on unknown leases from unknown hosts are not ignored.",
						Optional:            true,
						Computed:            true,
					},
					"reply_delay": schema.Int64Attribute{
						MarkdownDescription: "Delays sending DHCPOFFER and PROXYDHCP replies for at least the specified number of seconds. This can be practical for split DHCP solutions, to make sure the secondary server answers slower than the primary.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
					"register_firewall_rules": schema.BoolAttribute{
						MarkdownDescription: "Automatically register firewall rules to allow DHCP traffic for all explicitly selected interfaces, can be disabled for more fine grained control if needed.",
						Optional:            true,
						Computed:            true,
					},
					"router_advertisements": schema.BoolAttribute{
						MarkdownDescription: "Setting this will enable Router Advertisements for all configured DHCPv6 ranges with the managed address bits set, and the use SLAAC bit reset. To change this default, select a combination of the possible options in the individual DHCPv6 ranges. Keep in mind that this is a global option; if there are configured DHCPv6 ranges, RAs will be sent unconditionally and cannot be deactivated selectively. Setting Router Advertisement modes in DHCPv6 ranges will have no effect without this global option enabled.",
						Optional:            true,
						Computed:            true,
					},
					"host_ping": schema.BoolAttribute{
						MarkdownDescription: "By default, the DHCP server will use a ping to ensure that an address is not in use before allocating it to a host.",
						Optional:            true,
						Computed:            true,
					},
					"disable_ha_sync": schema.BoolAttribute{
						MarkdownDescription: "Ignore the DHCP general settings from being updated using HA sync.",
						Optional:            true,
						Computed:            true,
					},
					"log_dhcp": schema.BoolAttribute{
						MarkdownDescription: "Extra logging for DHCP, log all the options sent to DHCP clients and the tags used to determine them.",
						Optional:            true,
						Computed:            true,
					},
					"log_quiet": schema.BoolAttribute{
						MarkdownDescription: "Suppress logging of the routine operation of DHCP, RA and TFTP. Errors and problems will still be logged.",
						Optional:            true,
						Computed:            true,
					},
				},
			},

			"legacy": schema.SingleNestedAttribute{
				MarkdownDescription: "Legacy DHCP integration settings.",
				Optional:            true,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"register_isc_dhcp4_leases": schema.BoolAttribute{
						MarkdownDescription: "If this option is set, then machines that specify their hostname when requesting a DHCP lease will be registered, so that their name can be resolved.",
						Optional:            true,
						Computed:            true,
					},
					"dhcp_domain_override": schema.StringAttribute{
						MarkdownDescription: "The domain name to use for DHCP hostname registration. If empty, the default system domain is used. Note that all DHCP leases will be assigned to the same domain. If this is undesired, static DHCP lease registration is able to provide coherent mappings.",
						Optional:            true,
						Computed:            true,
					},
					"register_dhcp_static_mappings": schema.BoolAttribute{
						MarkdownDescription: "If this option is set, then DHCP static mappings will be registered, so that their name can be resolved.",
						Optional:            true,
						Computed:            true,
					},
					"prefer_dhcp": schema.BoolAttribute{
						MarkdownDescription: "If this option is set, then DHCP mappings will be resolved before the manual list of names below. This only affects the name given for a reverse lookup (PTR).",
						Optional:            true,
						Computed:            true,
					},
				},
			},
		},
	}
}

func settingsDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Reads the general dnsmasq settings from OPNsense.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Always set to `dnsmasq_settings`.",
				Validators: []validator.String{
					stringvalidator.OneOf("dnsmasq_settings"),
				},
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable the dnsmasq service.",
				Computed:            true,
			},
			"interface": dschema.SetAttribute{
				MarkdownDescription: "Interface IPs used to respond to queries from clients. If no interfaces are selected, Dnsmasq will listen on all available IPv4 and IPv6 addresses by default. However, DHCP related firewall rules will only be added for explicitly selected interfaces, never for all interfaces.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"strict_interface_binding": dschema.BoolAttribute{
				MarkdownDescription: "By default we bind the wildcard address, even when listening on some interfaces. Requests that should not be handled are discarded, this has the advantage of working even when interfaces come and go and change address. This option forces binding to only the interfaces we are listening on, which is less stable in non static environments.",
				Computed:            true,
			},

			"dns": dschema.SingleNestedAttribute{
				MarkdownDescription: "DNS server settings.",
				Computed:            true,

				Attributes: map[string]dschema.Attribute{
					"port": dschema.Int64Attribute{
						MarkdownDescription: "The port used for responding to DNS queries. It should normally be left blank unless another service needs to bind to TCP/UDP port 53. Setting this to zero (0) completely disables DNS function",
						Computed:            true,
					},
					"dnssec": dschema.BoolAttribute{
						MarkdownDescription: "Enable DNSSEC validation.",
						Computed:            true,
					},
					"no_host_lookup": dschema.BoolAttribute{
						MarkdownDescription: "Do not use `/etc/hosts` for DNS lookups.",
						Computed:            true,
					},
					"log_queries": dschema.BoolAttribute{
						MarkdownDescription: "Append the configured domain to simple hostnames read from hosts files.",
						Computed:            true,
					},
					"max_concurrent_queries": dschema.Int64Attribute{
						MarkdownDescription: "Set the maximum number of concurrent DNS queries. On configurations with tight resources, this value may need to be reduced.",
						Computed:            true,
					},
					"cache_size": dschema.Int64Attribute{
						MarkdownDescription: "Set the size of the cache. Setting the cache size to zero disables caching. Please note that huge cache size impacts performance.",
						Computed:            true,
					},
					"local_entry_ttl": dschema.Int64Attribute{
						MarkdownDescription: "This option allows a time-to-live (in seconds) to be given for local DNS entries, i.e. /etc/hosts or DHCP leases. This will reduce the load on the server at the expense of clients using stale data under some circumstances. A value of zero will disable client-side caching.",
						Computed:            true,
					},
					"no_ident": dschema.BoolAttribute{
						MarkdownDescription: "Do not respond to class CHAOS and type TXT in domain bind queries. Without this option being set, the cache statistics are also available in the DNS as answers to queries of class CHAOS and type TXT in domain bind.",
						Computed:            true,
					},
				},
			},

			"dns_query_forwarding": dschema.SingleNestedAttribute{
				MarkdownDescription: "DNS query forwarding settings.",
				Computed:            true,

				Attributes: map[string]dschema.Attribute{
					"query_sequentially": dschema.BoolAttribute{
						MarkdownDescription: "If this option is set, we will query the DNS servers sequentially in the order specified (System: General Setup: DNS Servers), rather than all at once in parallel.",
						Computed:            true,
					},
					"require_domain": dschema.BoolAttribute{
						MarkdownDescription: "If this option is set, we will not forward A or AAAA queries for plain names, without dots or domain parts, to upstream name servers. If the name is not known from /etc/hosts or DHCP then a \"not found\" answer is returned.",
						Computed:            true,
					},
					"do_not_forward_system_dns": dschema.BoolAttribute{
						MarkdownDescription: "If this option is set, DNS forwarding to system nameservers (defined in System: General Setup: DNS Servers) will be disabled. Upstream servers defined in Services: Dnsmasq DNS & DHCP: Domains will still be used. This option is recommended when Unbound forwards local domain queries to Dnsmasq, so that all queries terminate without further lookups if they are unknown.",
						Computed:            true,
					},
					"do_not_forward_private_reverse": dschema.BoolAttribute{
						MarkdownDescription: "If this option is set, we will not forward reverse DNS lookups (PTR) for private addresses (RFC 1918) to upstream name servers. Any entries in the Domain Overrides section forwarding private \"n.n.n.in-addr.arpa\" names to a specific server are still forwarded. If the IP to name is not known from /etc/hosts, DHCP or a specific domain override then a \"not found\" answer is immediately returned.",
						Computed:            true,
					},
					"add_mac": dschema.StringAttribute{
						MarkdownDescription: "Add the MAC address of the requestor to DNS queries which are forwarded upstream. The MAC address will only be added if the upstream DNS Server is in the same subnet as the requestor. Since this is not standardized, it should be considered experimental. This is useful for selective DNS filtering on the upstream DNS server.",
						Computed:            true,
					},
					"add_subnet": dschema.BoolAttribute{
						MarkdownDescription: "Add the real client IPv4 and IPv6 addresses (add-subnet=32,128) to DNS queries which are forwarded upstream. Be careful setting this option as it can undermine privacy. This is useful for selective DNS filtering on the upstream DNS server.",
						Computed:            true,
					},
					"strip_subnet": dschema.BoolAttribute{
						MarkdownDescription: "Strip the subnet received by a downstream DNS server. If \"Add subnet\" is used and the downstream DNS server already added a subnet, Dnsmasq will not replace it without setting \"Strip subnet\".",
						Computed:            true,
					},
				},
			},

			"dhcp": dschema.SingleNestedAttribute{
				MarkdownDescription: "General DHCP settings.",
				Computed:            true,

				Attributes: map[string]dschema.Attribute{
					"no_dhcp_interfaces": dschema.SetAttribute{
						MarkdownDescription: "Do not provide DHCP, TFTP or router advertisement on the specified interfaces, but do provide DNS service. Please note that Dnsmasq continues to listen on the default DHCP ports as long as any DHCP ranges are configured; setting this option only ignores these packets on the selected interfaces.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"fqdn": dschema.BoolAttribute{
						MarkdownDescription: "If disabled, registers the unqualified names of DHCP clients into the DNS (e.g. 'smartphone'), in which case they should be unique. If enabled, the qualified name (e.g. 'smartphone.lan.internal') is registered. This option must be enabled if you are forwarding from Unbound to Dnsmasq for specific local domains.",
						Computed:            true,
					},
					"default_domain": dschema.StringAttribute{
						MarkdownDescription: "To ensure that all names have a domain part, there must be a default domain specified when DHCP FQDN is set. Leave empty to use the system domain.",
						Computed:            true,
					},
					"local_domain": dschema.BoolAttribute{
						MarkdownDescription: "Sets all DHCP domains as local. This will configure this DNS server as authoritative; it will not forward queries to any upstream servers for these domains.",
						Computed:            true,
					},
					"max_leases": dschema.Int64Attribute{
						MarkdownDescription: "Limits Dnsmasq to the specified maximum number of DHCP leases. This limit is to prevent DoS attacks from hosts which create thousands of leases and use lots of memory in the Dnsmasq process.",
						Computed:            true,
					},
					"authoritative": dschema.BoolAttribute{
						MarkdownDescription: "Should be set when Dnsmasq is definitely the only DHCP server on a network. For DHCPv4, it changes the behaviour from strict RFC compliance so that DHCP requests on unknown leases from unknown hosts are not ignored.",
						Computed:            true,
					},
					"reply_delay": dschema.Int64Attribute{
						MarkdownDescription: "Delays sending DHCPOFFER and PROXYDHCP replies for at least the specified number of seconds. This can be practical for split DHCP solutions, to make sure the secondary server answers slower than the primary.",
						Computed:            true,
					},
					"register_firewall_rules": dschema.BoolAttribute{
						MarkdownDescription: "Automatically register firewall rules to allow DHCP traffic for all explicitly selected interfaces, can be disabled for more fine grained control if needed.",
						Computed:            true,
					},
					"router_advertisements": dschema.BoolAttribute{
						MarkdownDescription: "Setting this will enable Router Advertisements for all configured DHCPv6 ranges with the managed address bits set, and the use SLAAC bit reset. To change this default, select a combination of the possible options in the individual DHCPv6 ranges. Keep in mind that this is a global option; if there are configured DHCPv6 ranges, RAs will be sent unconditionally and cannot be deactivated selectively. Setting Router Advertisement modes in DHCPv6 ranges will have no effect without this global option enabled.",
						Computed:            true,
					},
					"host_ping": dschema.BoolAttribute{
						MarkdownDescription: "By default, the DHCP server will use a ping to ensure that an address is not in use before allocating it to a host.",
						Computed:            true,
					},
					"disable_ha_sync": dschema.BoolAttribute{
						MarkdownDescription: "Ignore the DHCP general settings from being updated using HA sync.",
						Computed:            true,
					},
					"log_dhcp": dschema.BoolAttribute{
						MarkdownDescription: "Extra logging for DHCP, log all the options sent to DHCP clients and the tags used to determine them.",
						Computed:            true,
					},
					"log_quiet": dschema.BoolAttribute{
						MarkdownDescription: "Suppress logging of the routine operation of DHCP, RA and TFTP. Errors and problems will still be logged.",
						Computed:            true,
					},
				},
			},

			"legacy": dschema.SingleNestedAttribute{
				MarkdownDescription: "Legacy DHCP integration settings.",
				Computed:            true,

				Attributes: map[string]dschema.Attribute{
					"register_isc_dhcp4_leases": dschema.BoolAttribute{
						MarkdownDescription: "If this option is set, then machines that specify their hostname when requesting a DHCP lease will be registered, so that their name can be resolved.",
						Computed:            true,
					},
					"dhcp_domain_override": dschema.StringAttribute{
						MarkdownDescription: "The domain name to use for DHCP hostname registration. If empty, the default system domain is used. Note that all DHCP leases will be assigned to the same domain. If this is undesired, static DHCP lease registration is able to provide coherent mappings.",
						Computed:            true,
					},
					"register_dhcp_static_mappings": dschema.BoolAttribute{
						MarkdownDescription: "If this option is set, then DHCP static mappings will be registered, so that their name can be resolved.",
						Computed:            true,
					},
					"prefer_dhcp": dschema.BoolAttribute{
						MarkdownDescription: "If this option is set, then DHCP mappings will be resolved before the manual list of names below. This only affects the name given for a reverse lookup (PTR).",
						Computed:            true,
					},
				},
			},
		},
	}
}

func convertSettingsSchemaToStruct(d *settingsResourceModel) (*dnsmasq.GeneralSettingsWrapper, error) {
	result := &dnsmasq.GeneralSettingsWrapper{}

	result.Dnsmasq.IsEnabled =
		tools.BoolToString(d.Enabled.ValueBool())

	result.Dnsmasq.StrictInterfaceBinding =
		tools.BoolToString(d.StrictInterfaceBinding.ValueBool())

	// Interface
	if !d.Interface.IsNull() && !d.Interface.IsUnknown() {
		var interfaces []string

		diags := d.Interface.ElementsAs(
			context.Background(),
			&interfaces,
			false,
		)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse interfaces")
		}

		sort.Strings(interfaces)

		result.Dnsmasq.Interface =
			api.SelectedMapList(interfaces)
	}

	// DNS
	if d.DNS != nil {
		result.Dnsmasq.DNS_Port = tools.Int64ToStringNegative(d.DNS.Port.ValueInt64())
		result.Dnsmasq.DNS_DNSSEC = tools.BoolToString(d.DNS.DNSSEC.ValueBool())
		result.Dnsmasq.DNS_NoHostLookup = tools.BoolToString(d.DNS.NoHostLookup.ValueBool())
		result.Dnsmasq.DNS_LogDnsQueries = tools.BoolToString(d.DNS.LogQueries.ValueBool())
		result.Dnsmasq.DNS_MaxConcurrentQueries = tools.Int64ToStringNegative(d.DNS.MaxConcurrentQueries.ValueInt64())
		result.Dnsmasq.DNS_CacheSize = tools.Int64ToStringNegative(d.DNS.CacheSize.ValueInt64())
		result.Dnsmasq.DNS_LocalEntryTTL = tools.Int64ToStringNegative(d.DNS.LocalEntryTTL.ValueInt64())
		result.Dnsmasq.DNS_NoIdent = tools.BoolToString(d.DNS.NoIdent.ValueBool())
	}

	// DNS query forwarding
	if d.DNSQueryForwarding != nil {
		result.Dnsmasq.DNS_QF_QuerySequentially = tools.BoolToString(d.DNSQueryForwarding.QuerySequentially.ValueBool())
		result.Dnsmasq.DNS_QF_RequireDomain = tools.BoolToString(d.DNSQueryForwarding.RequireDomain.ValueBool())
		result.Dnsmasq.DNS_QF_DoNotForwardSystemDNS = tools.BoolToString(d.DNSQueryForwarding.DoNotForwardSystemDNS.ValueBool())
		result.Dnsmasq.DNS_QF_DoNotForwardPrivateReverseLookup = tools.BoolToString(d.DNSQueryForwarding.DoNotForwardPrivateReverse.ValueBool())
		result.Dnsmasq.DNS_QF_AddMac = api.SelectedMap(d.DNSQueryForwarding.AddMAC.ValueString())
		result.Dnsmasq.DNS_QF_AddSubnet = tools.BoolToString(d.DNSQueryForwarding.AddSubnet.ValueBool())
		result.Dnsmasq.DNS_QF_StripSubnet = tools.BoolToString(d.DNSQueryForwarding.StripSubnet.ValueBool())
	}

	// DHCP
	if d.DHCP != nil {
		result.Dnsmasq.DHCPSettings.FQDN = tools.BoolToString(d.DHCP.FQDN.ValueBool())
		result.Dnsmasq.DHCPSettings.DefaultDomain = d.DHCP.DefaultDomain.ValueString()
		result.Dnsmasq.DHCPSettings.LocalDomain = tools.BoolToString(d.DHCP.LocalDomain.ValueBool())
		result.Dnsmasq.DHCPSettings.MaxLeases = tools.Int64ToStringNegative(d.DHCP.MaxLeases.ValueInt64())
		result.Dnsmasq.DHCPSettings.Authoritative = tools.BoolToString(d.DHCP.Authoritative.ValueBool())
		result.Dnsmasq.DHCPSettings.ReplyDelay = tools.Int64ToStringNegative(d.DHCP.ReplyDelay.ValueInt64())
		result.Dnsmasq.DHCPSettings.RegisterFirewallRules = tools.BoolToString(d.DHCP.RegisterFirewallRules.ValueBool())
		result.Dnsmasq.DHCPSettings.RouterAdvertisements = tools.BoolToString(d.DHCP.RouterAdvertisements.ValueBool())
		result.Dnsmasq.DHCPSettings.HostPing = tools.BoolToString(d.DHCP.HostPing.ValueBool())
		result.Dnsmasq.DHCPSettings.DisableHASync = tools.BoolToString(d.DHCP.DisableHASync.ValueBool())
		result.Dnsmasq.DHCPSettings.LogDhcp = tools.BoolToString(d.DHCP.LogDhcp.ValueBool())
		result.Dnsmasq.DHCPSettings.LogQuiet = tools.BoolToString(d.DHCP.LogQuiet.ValueBool())

		if !d.DHCP.NoDHCPInterfaces.IsNull() &&
			!d.DHCP.NoDHCPInterfaces.IsUnknown() {

			var interfaces []string

			diags := d.DHCP.NoDHCPInterfaces.ElementsAs(
				context.Background(),
				&interfaces,
				false,
			)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to parse no_dhcp_interfaces")
			}

			sort.Strings(interfaces)

			result.Dnsmasq.DHCPSettings.InterfaceNoDhcp = api.SelectedMapList(interfaces)
		}
	}

	// Legacy
	if d.Legacy != nil {
		result.Dnsmasq.Legacy_RegisterISCDhcp4Leases = tools.BoolToString(d.Legacy.RegisterISCDHCP4Leases.ValueBool())
		result.Dnsmasq.Legacy_DhcpDomainOverride = d.Legacy.DHCPDomainOverride.ValueString()
		result.Dnsmasq.Legacy_RegisterDhcpStaticMappings = tools.BoolToString(d.Legacy.RegisterDHCPStaticMapping.ValueBool())
		result.Dnsmasq.Legacy_PreferDhcp = tools.BoolToString(d.Legacy.PreferDHCP.ValueBool())
	}

	return result, nil
}

func convertSettingsStructToSchema(dRaw *dnsmasq.GeneralSettingsWrapper) (*settingsResourceModel, error) {
	d := dRaw.Dnsmasq

	model := &settingsResourceModel{
		Id: types.StringValue("dnsmasq_settings"),
	}

	// DNS
	model.DNS = &settingsDNSBlock{
		Port:                 tools.StringToInt64Null(d.DNS_Port),
		DNSSEC:               types.BoolValue(tools.StringToBool(d.DNS_DNSSEC)),
		NoHostLookup:         types.BoolValue(tools.StringToBool(d.DNS_NoHostLookup)),
		LogQueries:           types.BoolValue(tools.StringToBool(d.DNS_LogDnsQueries)),
		MaxConcurrentQueries: tools.StringToInt64Null(d.DNS_MaxConcurrentQueries),
		CacheSize:            tools.StringToInt64Null(d.DNS_CacheSize),
		LocalEntryTTL:        tools.StringToInt64Null(d.DNS_LocalEntryTTL),
		NoIdent:              types.BoolValue(tools.StringToBool(d.DNS_NoIdent)),
	}

	// DNS query forwarding
	model.DNSQueryForwarding = &settingsDNSQueryForwardingBlock{
		QuerySequentially:          types.BoolValue(tools.StringToBool(d.DNS_QF_QuerySequentially)),
		RequireDomain:              types.BoolValue(tools.StringToBool(d.DNS_QF_RequireDomain)),
		DoNotForwardSystemDNS:      types.BoolValue(tools.StringToBool(d.DNS_QF_DoNotForwardSystemDNS)),
		DoNotForwardPrivateReverse: types.BoolValue(tools.StringToBool(d.DNS_QF_DoNotForwardPrivateReverseLookup)),
		AddMAC:                     types.StringValue(d.DNS_QF_AddMac.String()),
		AddSubnet:                  types.BoolValue(tools.StringToBool(d.DNS_QF_AddSubnet)),
		StripSubnet:                types.BoolValue(tools.StringToBool(d.DNS_QF_StripSubnet)),
	}

	// DHCP
	model.DHCP = &settingsDHCPBlock{
		FQDN:                  types.BoolValue(tools.StringToBool(d.DHCPSettings.FQDN)),
		DefaultDomain:         types.StringValue(d.DHCPSettings.DefaultDomain),
		LocalDomain:           types.BoolValue(tools.StringToBool(d.DHCPSettings.LocalDomain)),
		MaxLeases:             tools.StringToInt64Null(d.DHCPSettings.MaxLeases),
		Authoritative:         types.BoolValue(tools.StringToBool(d.DHCPSettings.Authoritative)),
		ReplyDelay:            tools.StringToInt64Null(d.DHCPSettings.ReplyDelay),
		RegisterFirewallRules: types.BoolValue(tools.StringToBool(d.DHCPSettings.RegisterFirewallRules)),
		RouterAdvertisements:  types.BoolValue(tools.StringToBool(d.DHCPSettings.RouterAdvertisements)),
		HostPing:              types.BoolValue(tools.StringToBool(d.DHCPSettings.HostPing)),
		DisableHASync:         types.BoolValue(tools.StringToBool(d.DHCPSettings.DisableHASync)),
		LogDhcp:               types.BoolValue(tools.StringToBool(d.DHCPSettings.LogDhcp)),
		LogQuiet:              types.BoolValue(tools.StringToBool(d.DHCPSettings.LogQuiet)),
	}

	// DHCP interfaces
	model.DHCP.NoDHCPInterfaces = tools.StringSliceToSet([]string(d.DHCPSettings.InterfaceNoDhcp))

	// Legacy
	model.Legacy = &settingsLegacyBlock{
		RegisterISCDHCP4Leases:    types.BoolValue(tools.StringToBool(d.Legacy_RegisterISCDhcp4Leases)),
		DHCPDomainOverride:        types.StringValue(d.Legacy_DhcpDomainOverride),
		RegisterDHCPStaticMapping: types.BoolValue(tools.StringToBool(d.Legacy_RegisterDhcpStaticMappings)),
		PreferDHCP:                types.BoolValue(tools.StringToBool(d.Legacy_PreferDhcp)),
	}

	// General settings
	model.Enabled = types.BoolValue(tools.StringToBool(d.IsEnabled))
	model.StrictInterfaceBinding = types.BoolValue(tools.StringToBool(d.StrictInterfaceBinding))
	model.Interface = tools.StringSliceToSet([]string(d.Interface))

	return model, nil
}
