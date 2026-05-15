package unbound

import (
	"context"
	"regexp"
	"sort"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/unbound"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Nested block structs

type settingsGeneralBlock struct {
	Enabled                    types.Bool   `tfsdk:"enabled"`
	Port                       types.Int64  `tfsdk:"port"`
	ListenInterfaces           types.Set    `tfsdk:"listen_interfaces"`
	EnableDNSSEC               types.Bool   `tfsdk:"enable_dnssec"`
	EnableDNS64                types.Bool   `tfsdk:"enable_dns64"`
	DNS64Prefix                types.String `tfsdk:"dns64_prefix"`
	NoARecords                 types.Bool   `tfsdk:"no_a_records"`
	RegisterDHCPLeases         types.Bool   `tfsdk:"register_dhcp_leases"`
	DHCPDomainOverride         types.String `tfsdk:"dhcp_domain_override"`
	RegisterDHCPStaticMappings types.Bool   `tfsdk:"register_dhcp_static_mappings"`
	RegisterIPv6LinkLocal      types.Bool   `tfsdk:"register_ipv6_link_local"`
	RegisterSystemsRecords     types.Bool   `tfsdk:"register_systems_records"`
	EnableTXTSupport           types.Bool   `tfsdk:"enable_txt_support"`
	EnableReloadCacheFlush     types.Bool   `tfsdk:"enable_reload_cache_flush"`
	LocalZoneType              types.String `tfsdk:"local_zone_type"`
	OutgoingInterfaces         types.Set    `tfsdk:"outgoing_interfaces"`
	EnableWPAD                 types.Bool   `tfsdk:"enable_wpad"`
}

type settingsAdvancedBlock struct {
	HideIdentity        types.Bool  `tfsdk:"hide_identity"`
	HideVersion         types.Bool  `tfsdk:"hide_version"`
	EnablePrefetchKey   types.Bool  `tfsdk:"enable_prefetch_key"`
	DNSSECStripped      types.Bool  `tfsdk:"dnssec_stripped"`
	AggressiveNSEC      types.Bool  `tfsdk:"aggressive_nsec"`
	QnameMinStrict      types.Bool  `tfsdk:"qname_min_strict"`
	OutgoingNumTCP      types.Int64 `tfsdk:"outgoing_num_tcp"`
	IncomingNumTCP      types.Int64 `tfsdk:"incoming_num_tcp"`
	NumQueriesPerThread types.Int64 `tfsdk:"num_queries_per_thread"`
	OutgoingRange       types.Int64 `tfsdk:"outgoing_range"`
	JostleTimeout       types.Int64 `tfsdk:"jostle_timeout"`
	DiscardTimeout      types.Int64 `tfsdk:"discard_timeout"`
	PrivateDomains      types.Set   `tfsdk:"private_domains"`
	PrivateAddresses    types.Set   `tfsdk:"private_addresses"`
	InsecureDomains     types.Set   `tfsdk:"insecure_domains"`

	ServeExpired *settingsAdvancedServeExpiredBlock `tfsdk:"serve_expired"`
	Logging      *settingsAdvancedLoggingBlock      `tfsdk:"logging"`
	Cache        *settingsAdvancedCacheBlock        `tfsdk:"cache"`
}

type settingsAdvancedServeExpiredBlock struct {
	Enabled               types.Bool   `tfsdk:"enabled"`
	RecordReplyTTL        types.String `tfsdk:"record_reply_ttl"`
	TTL                   types.String `tfsdk:"ttl"`
	ResetTTL              types.Bool   `tfsdk:"reset_ttl"`
	ClientResponseTimeout types.String `tfsdk:"client_response_timeout"`
}

type settingsAdvancedLoggingBlock struct {
	ExtendedStatistics types.Bool  `tfsdk:"extended_statistics"`
	LogQueries         types.Bool  `tfsdk:"log_queries"`
	LogReplies         types.Bool  `tfsdk:"log_replies"`
	TagQueryReply      types.Bool  `tfsdk:"tag_query_reply"`
	LogLocalActions    types.Bool  `tfsdk:"log_local_actions"`
	LogServFail        types.Bool  `tfsdk:"log_servfail"`
	VerbosityLevel     types.Int64 `tfsdk:"verbosity_level"`
	ValidationLevel    types.Int64 `tfsdk:"validation_level"`
}

type settingsAdvancedCacheBlock struct {
	EnablePrefetch         types.Bool   `tfsdk:"enable_prefetch"`
	UnwantedReplyThreshold types.Int64  `tfsdk:"unwanted_reply_threshold"`
	MsgCacheSize           types.String `tfsdk:"msg_cache_size"`
	RRSetCacheSize         types.String `tfsdk:"rrset_cache_size"`
	MaxTTL                 types.Int64  `tfsdk:"max_ttl"`
	MaxNegativeTTL         types.Int64  `tfsdk:"max_negative_ttl"`
	MinTTL                 types.Int64  `tfsdk:"min_ttl"`
	HostTTL                types.Int64  `tfsdk:"host_ttl"`
	KeepProbingHosts       types.Bool   `tfsdk:"keep_probing_hosts"`
	NumHosts               types.Int64  `tfsdk:"num_hosts"`
}

type settingsACLsBlock struct {
	DefaultAction types.String `tfsdk:"default_action"`
}

type settingsDNSBLBlock struct {
	Enabled            types.Bool   `tfsdk:"enabled"`
	ForceSafeSearch    types.Bool   `tfsdk:"force_safe_search"`
	Type               types.Set    `tfsdk:"type"`
	Blocklists         types.Set    `tfsdk:"blocklists"`
	WhitelistDomains   types.Set    `tfsdk:"whitelist_domains"`
	BlocklistDomains   types.Set    `tfsdk:"blocklist_domains"`
	WildcardDomains    types.Set    `tfsdk:"wildcard_domains"`
	DestinationAddress types.String `tfsdk:"destination_address"`
	ReturnNXDomain     types.Bool   `tfsdk:"return_nxdomain"`
}

type settingsForwardingBlock struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

// settingsResourceModel describes the resource data model.
// This is a SINGLETON resource - it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type settingsResourceModel struct {
	Id         types.String             `tfsdk:"id"`
	General    *settingsGeneralBlock    `tfsdk:"general"`
	Advanced   *settingsAdvancedBlock   `tfsdk:"advanced"`
	ACLs       *settingsACLsBlock       `tfsdk:"acls"`
	DNSBL      *settingsDNSBLBlock      `tfsdk:"dnsbl"`
	Forwarding *settingsForwardingBlock `tfsdk:"forwarding"`
}

func settingsResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages Unbound DNS resolver settings. This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_unbound_settings.settings unbound_settings\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `unbound_settings`. Use this value when importing: `terraform import opnsense_unbound_settings.settings unbound_settings`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"general": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					map[string]attr.Type{
						"enabled":                       types.BoolType,
						"port":                          types.Int64Type,
						"listen_interfaces":             types.SetType{ElemType: types.StringType},
						"enable_dnssec":                 types.BoolType,
						"enable_dns64":                  types.BoolType,
						"dns64_prefix":                  types.StringType,
						"no_a_records":                  types.BoolType,
						"register_dhcp_leases":          types.BoolType,
						"dhcp_domain_override":          types.StringType,
						"register_dhcp_static_mappings": types.BoolType,
						"register_ipv6_link_local":      types.BoolType,
						"register_systems_records":      types.BoolType,
						"enable_txt_support":            types.BoolType,
						"enable_reload_cache_flush":     types.BoolType,
						"local_zone_type":               types.StringType,
						"outgoing_interfaces":           types.SetType{ElemType: types.StringType},
						"enable_wpad":                   types.BoolType,
					},
					map[string]attr.Value{
						"enabled":                       types.BoolValue(false),
						"port":                          types.Int64Value(53),
						"listen_interfaces":             types.SetValueMust(types.StringType, []attr.Value{}),
						"enable_dnssec":                 types.BoolValue(false),
						"enable_dns64":                  types.BoolValue(false),
						"dns64_prefix":                  types.StringValue(""),
						"no_a_records":                  types.BoolValue(false),
						"register_dhcp_leases":          types.BoolValue(false),
						"dhcp_domain_override":          types.StringValue(""),
						"register_dhcp_static_mappings": types.BoolValue(false),
						"register_ipv6_link_local":      types.BoolValue(true),
						"register_systems_records":      types.BoolValue(true),
						"enable_txt_support":            types.BoolValue(false),
						"enable_reload_cache_flush":     types.BoolValue(false),
						"local_zone_type":               types.StringValue(""),
						"outgoing_interfaces":           types.SetValueMust(types.StringType, []attr.Value{}),
						"enable_wpad":                   types.BoolValue(false),
					},
				)),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable Unbound DNS resolver. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"port": schema.Int64Attribute{
						MarkdownDescription: "The TCP/UDP port used for responding to DNS queries. Defaults to `53`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(53),
					},
					"listen_interfaces": schema.SetAttribute{
						MarkdownDescription: "Interface IP addresses used for responding to queries from clients. If an interface has both IPv4 and IPv6 IPs, both are used. Queries to other interface IPs not selected below are discarded. The default behavior is to respond to queries on every available IPv4 and IPv6 address. Defaults to `[]`.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
					"enable_dnssec": schema.BoolAttribute{
						MarkdownDescription: "Enable DNSSEC validation. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"enable_dns64": schema.BoolAttribute{
						MarkdownDescription: "When set, Unbound will synthesize AAAA records from A records if no actual AAAA records are present. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"dns64_prefix": schema.StringAttribute{
						MarkdownDescription: "Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"no_a_records": schema.BoolAttribute{
						MarkdownDescription: "When set, Unbound will remove all A records from the answer section of all responses. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"register_dhcp_leases": schema.BoolAttribute{
						MarkdownDescription: "When set, then machines that specify their hostname when requesting a DHCP lease will be registered in Unbound, so that their name can be resolved. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"dhcp_domain_override": schema.StringAttribute{
						MarkdownDescription: "The default domain name to use for DHCP lease registration. If empty, the system domain is used. Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"register_dhcp_static_mappings": schema.BoolAttribute{
						MarkdownDescription: "When set, then DHCP static mappings will be registered in Unbound, so that their name can be resolved. You should also set the domain in `System: Settings: General` to the proper value. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"register_ipv6_link_local": schema.BoolAttribute{
						MarkdownDescription: "When set, then IPv6 link-local addresses will be registered in Unbound, allowing return of unreachable address when more than one listen interface is configured. Defaults to `true`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
					},
					"register_systems_records": schema.BoolAttribute{
						MarkdownDescription: "When set, then A/AAAA records for the configured listen interfaces will be generated. Disable this to control which interface IP addresses are mapped to the system host/domain name as well as to restrict the amount of information exposed in replies to queries for the system host/domain name. Defaults to `true`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
					},
					"enable_txt_support": schema.BoolAttribute{
						MarkdownDescription: "When set, then any descriptions associated with Host entries and DHCP Static mappings will create a corresponding TXT record. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"enable_reload_cache_flush": schema.BoolAttribute{
						MarkdownDescription: "When set, the DNS cache will be flushed during each daemon reload. This is the default behavior for Unbound, but may be undesired when multiple dynamic interfaces require frequent reloading. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"local_zone_type": schema.StringAttribute{
						MarkdownDescription: "The local zone type used for the system domain. Type descriptions are available under \"local-zone:\" in the `unbound.conf(5)` manual page. Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.OneOf(
								"always_nxdomain",
								"always_refuse",
								"always_transparent",
								"deny",
								"inform",
								"inform_deny",
								"nodefault",
								"refuse",
								"static",
								"transparent",
								"typetransparent",
								"",
							),
						},
					},
					"outgoing_interfaces": schema.SetAttribute{
						MarkdownDescription: "Utilize different network interfaces that Unbound will use to send queries to authoritative servers and receive their replies. By default all interfaces are used. Note that setting explicit outgoing interfaces only works when they are statically configured. Defaults to `[]`.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
					"enable_wpad": schema.BoolAttribute{
						MarkdownDescription: "When set, CNAME records for the WPAD host of all configured domains will be automatically added as well as overrides for TXT records for domains. This allows automatic proxy configuration in your network but you should not enable it if you are not using WPAD or if you want to configure it by yourself. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
			"advanced": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					map[string]attr.Type{
						"hide_identity":          types.BoolType,
						"hide_version":           types.BoolType,
						"enable_prefetch_key":    types.BoolType,
						"dnssec_stripped":        types.BoolType,
						"aggressive_nsec":        types.BoolType,
						"qname_min_strict":       types.BoolType,
						"outgoing_num_tcp":       types.Int64Type,
						"incoming_num_tcp":       types.Int64Type,
						"num_queries_per_thread": types.Int64Type,
						"outgoing_range":         types.Int64Type,
						"jostle_timeout":         types.Int64Type,
						"discard_timeout":        types.Int64Type,
						"private_domains":        types.SetType{ElemType: types.StringType},
						"private_addresses":      types.SetType{ElemType: types.StringType},
						"insecure_domains":       types.SetType{ElemType: types.StringType},
						"serve_expired": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"enabled":                 types.BoolType,
								"record_reply_ttl":        types.StringType,
								"ttl":                     types.StringType,
								"reset_ttl":               types.BoolType,
								"client_response_timeout": types.StringType,
							},
						},
						"logging": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"extended_statistics": types.BoolType,
								"log_queries":         types.BoolType,
								"log_replies":         types.BoolType,
								"tag_query_reply":     types.BoolType,
								"log_local_actions":   types.BoolType,
								"log_servfail":        types.BoolType,
								"verbosity_level":     types.Int64Type,
								"validation_level":    types.Int64Type,
							},
						},
						"cache": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"enable_prefetch":          types.BoolType,
								"unwanted_reply_threshold": types.Int64Type,
								"msg_cache_size":           types.StringType,
								"rrset_cache_size":         types.StringType,
								"max_ttl":                  types.Int64Type,
								"max_negative_ttl":         types.Int64Type,
								"min_ttl":                  types.Int64Type,
								"host_ttl":                 types.Int64Type,
								"keep_probing_hosts":       types.BoolType,
								"num_hosts":                types.Int64Type,
							},
						},
					},
					map[string]attr.Value{
						"hide_identity":          types.BoolValue(false),
						"hide_version":           types.BoolValue(false),
						"enable_prefetch_key":    types.BoolValue(false),
						"dnssec_stripped":        types.BoolValue(false),
						"aggressive_nsec":        types.BoolValue(true),
						"qname_min_strict":       types.BoolValue(false),
						"outgoing_num_tcp":       types.Int64Value(-1),
						"incoming_num_tcp":       types.Int64Value(-1),
						"num_queries_per_thread": types.Int64Value(-1),
						"outgoing_range":         types.Int64Value(-1),
						"jostle_timeout":         types.Int64Value(-1),
						"discard_timeout":        types.Int64Value(-1),
						"private_domains":        types.SetValueMust(types.StringType, []attr.Value{}),
						"private_addresses": types.SetValueMust(types.StringType, []attr.Value{
							types.StringValue("0.0.0.0/8"),
							types.StringValue("10.0.0.0/8"),
							types.StringValue("100.64.0.0/10"),
							types.StringValue("169.254.0.0/16"),
							types.StringValue("172.16.0.0/12"),
							types.StringValue("192.0.2.0/24"),
							types.StringValue("192.168.0.0/16"),
							types.StringValue("198.18.0.0/15"),
							types.StringValue("198.51.100.0/24"),
							types.StringValue("2001:db8::/32"),
							types.StringValue("203.0.113.0/24"),
							types.StringValue("233.252.0.0/24"),
							types.StringValue("::1/128"),
							types.StringValue("fc00::/8"),
							types.StringValue("fd00::/8"),
							types.StringValue("fe80::/10"),
						}),
						"insecure_domains": types.SetValueMust(types.StringType, []attr.Value{}),
						"serve_expired": types.ObjectValueMust(
							map[string]attr.Type{
								"enabled":                 types.BoolType,
								"record_reply_ttl":        types.StringType,
								"ttl":                     types.StringType,
								"reset_ttl":               types.BoolType,
								"client_response_timeout": types.StringType,
							},
							map[string]attr.Value{
								"enabled":                 types.BoolValue(false),
								"record_reply_ttl":        types.StringValue(""),
								"ttl":                     types.StringValue(""),
								"reset_ttl":               types.BoolValue(false),
								"client_response_timeout": types.StringValue(""),
							},
						),
						"logging": types.ObjectValueMust(
							map[string]attr.Type{
								"extended_statistics": types.BoolType,
								"log_queries":         types.BoolType,
								"log_replies":         types.BoolType,
								"tag_query_reply":     types.BoolType,
								"log_local_actions":   types.BoolType,
								"log_servfail":        types.BoolType,
								"verbosity_level":     types.Int64Type,
								"validation_level":    types.Int64Type,
							},
							map[string]attr.Value{
								"extended_statistics": types.BoolValue(false),
								"log_queries":         types.BoolValue(false),
								"log_replies":         types.BoolValue(false),
								"tag_query_reply":     types.BoolValue(false),
								"log_local_actions":   types.BoolValue(false),
								"log_servfail":        types.BoolValue(false),
								"verbosity_level":     types.Int64Value(1),
								"validation_level":    types.Int64Value(0),
							},
						),
						"cache": types.ObjectValueMust(
							map[string]attr.Type{
								"enable_prefetch":          types.BoolType,
								"unwanted_reply_threshold": types.Int64Type,
								"msg_cache_size":           types.StringType,
								"rrset_cache_size":         types.StringType,
								"max_ttl":                  types.Int64Type,
								"max_negative_ttl":         types.Int64Type,
								"min_ttl":                  types.Int64Type,
								"host_ttl":                 types.Int64Type,
								"keep_probing_hosts":       types.BoolType,
								"num_hosts":                types.Int64Type,
							},
							map[string]attr.Value{
								"enable_prefetch":          types.BoolValue(false),
								"unwanted_reply_threshold": types.Int64Value(-1),
								"msg_cache_size":           types.StringValue(""),
								"rrset_cache_size":         types.StringValue(""),
								"max_ttl":                  types.Int64Value(-1),
								"max_negative_ttl":         types.Int64Value(-1),
								"min_ttl":                  types.Int64Value(-1),
								"host_ttl":                 types.Int64Value(-1),
								"keep_probing_hosts":       types.BoolValue(false),
								"num_hosts":                types.Int64Value(-1),
							},
						),
					},
				)),
				Attributes: map[string]schema.Attribute{
					"hide_identity": schema.BoolAttribute{
						MarkdownDescription: "When enabled, id.server and hostname.bind queries are refused. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"hide_version": schema.BoolAttribute{
						MarkdownDescription: "When enabled, version.server and version.bind queries are refused. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"enable_prefetch_key": schema.BoolAttribute{
						MarkdownDescription: "When enabled, DNSKEYs are fetched earlier in the validation process when a Delegation signer is encountered. This helps lower the latency of requests but does utilize a little more CPU. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"dnssec_stripped": schema.BoolAttribute{
						MarkdownDescription: "DNSSEC data is required for trust-anchored zones. If such data is absent, the zone becomes bogus. If this is disabled and no DNSSEC data is received, then the zone is made insecure. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"aggressive_nsec": schema.BoolAttribute{
						MarkdownDescription: "Whether to enable RFC8198-based aggressive use of the DNSSEC-Validated cache. Helps to reduce the query rate towards targets but may lead to false negative responses if there are errors in the zone config. Defaults to `true`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
					},
					"qname_min_strict": schema.BoolAttribute{
						MarkdownDescription: "Whether to send minimum amount of information to upstream servers to enhance privacy. Do not fall-back to sending full QNAME to potentially broken nameservers. A lot of domains will not be resolvable when this option in enabled. Only use if you know what you are doing. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"outgoing_num_tcp": schema.Int64Attribute{
						MarkdownDescription: "The number of outgoing TCP buffers to allocate per thread. If 0 is selected then no TCP queries, to authoritative servers, are done. Defaults to `-1`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(-1),
					},
					"incoming_num_tcp": schema.Int64Attribute{
						MarkdownDescription: "The number of incoming TCP buffers to allocate per thread. If 0 is selected then no TCP queries, from clients, are accepted. Defaults to `-1`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(-1),
					},
					"num_queries_per_thread": schema.Int64Attribute{
						MarkdownDescription: "The number of queries that every thread will service simultaneously. If more queries arrive that need to be serviced, and no queries can be jostled out (see \"jostle_timeout\"), then these queries are dropped. This forces the client to resend after a timeout, allowing the server time to work on the existing queries. Defaults to `-1`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(-1),
					},
					"outgoing_range": schema.Int64Attribute{
						MarkdownDescription: "The number of ports to open. This number of file descriptors can be opened per thread. Larger numbers need extra resources from the operating system. For performance a very large value is best. For reference, usually double the amount of queries per thread is used. Defaults to `-1`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(-1),
					},
					"jostle_timeout": schema.Int64Attribute{
						MarkdownDescription: "This timeout is used for when the server is very busy. Set to a value that usually results in one round-trip to the authority servers. If too many queries arrive, then 50% of the queries are allowed to run to completion, and the other 50% are replaced with the new incoming query if they have already spent more than their allowed time. This protects against denial of service by slow queries or high query rates. Defaults to `-1`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(-1),
					},
					"discard_timeout": schema.Int64Attribute{
						MarkdownDescription: "The wait time in msec where recursion requests are dropped. This is to stop a large number of replies from accumulating. If 'Serve Expired Responses' is enabled this field should be set greater than 'Client Expired Response Timeout', otherwise, these late responses will not update the cache. The value 0 disables it. Default 1900. This setting may increase the \"request queue exceeded\" counter. Defaults to `-1`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(-1),
					},
					"private_domains": schema.SetAttribute{
						MarkdownDescription: "List of domains to mark as private. These domains and all its subdomains are allowed to contain private addresses. Defaults to `[]`.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
					"private_addresses": schema.SetAttribute{
						MarkdownDescription: "These are addresses on your private network, and are not allowed to be returned for public internet names. Any occurrence of such addresses are removed from DNS answers. Additionally, the DNSSEC validator may mark the answers bogus. This protects against so-called DNS Rebinding (Only applicable when DNS rebind check is enabled in `System->Settings->Administration`). Defaults to a list of RFC1918 and other private IP ranges.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default: setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{
							types.StringValue("0.0.0.0/8"),
							types.StringValue("10.0.0.0/8"),
							types.StringValue("100.64.0.0/10"),
							types.StringValue("169.254.0.0/16"),
							types.StringValue("172.16.0.0/12"),
							types.StringValue("192.0.2.0/24"),
							types.StringValue("192.168.0.0/16"),
							types.StringValue("198.18.0.0/15"),
							types.StringValue("198.51.100.0/24"),
							types.StringValue("2001:db8::/32"),
							types.StringValue("203.0.113.0/24"),
							types.StringValue("233.252.0.0/24"),
							types.StringValue("::1/128"),
							types.StringValue("fc00::/8"),
							types.StringValue("fd00::/8"),
							types.StringValue("fe80::/10"),
						})),
					},
					"insecure_domains": schema.SetAttribute{
						MarkdownDescription: "List of domains to mark as insecure. DNSSEC chain of trust is ignored towards the domain name. Defaults to `[]`.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
					"serve_expired": schema.SingleNestedAttribute{
						MarkdownDescription: "Serve expired cache configuration.",
						Optional:            true,
						Computed:            true,
						Default: objectdefault.StaticValue(types.ObjectValueMust(
							map[string]attr.Type{
								"enabled":                 types.BoolType,
								"record_reply_ttl":        types.StringType,
								"ttl":                     types.StringType,
								"reset_ttl":               types.BoolType,
								"client_response_timeout": types.StringType,
							},
							map[string]attr.Value{
								"enabled":                 types.BoolValue(false),
								"record_reply_ttl":        types.StringValue(""),
								"ttl":                     types.StringValue(""),
								"reset_ttl":               types.BoolValue(false),
								"client_response_timeout": types.StringValue(""),
							},
						)),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Serve expired responses from the cache with a TTL of 0 without waiting for the actual resolution to finish. The TTL can be modified with \"record_reply_ttl\" value. Defaults to `false`.",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"record_reply_ttl": schema.StringAttribute{
								MarkdownDescription: "TTL value to use when replying with expired data. If \"Client Expired Response Timeout\" is also used then it is recommended to use 30 as the value as per RFC 8767. Defaults to `\"\"`.",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString(""),
							},
							"ttl": schema.StringAttribute{
								MarkdownDescription: "Limits the serving of expired responses to the configured amount of seconds after expiration. A value of 0 disables the limit. A suggested value per RFC 8767 is between 86400 (1 day) and 259200 (3 days). Defaults to `\"\"`.",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString(""),
							},
							"reset_ttl": schema.BoolAttribute{
								MarkdownDescription: "Whether to set the TTL of expired records to the \"TTL for Expired Responses\" value after a failed attempt to retrieve the record from an upstream server. This makes sure that the expired records will be served as long as there are queries for it. Defaults to `false`.",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"client_response_timeout": schema.StringAttribute{
								MarkdownDescription: "Time in milliseconds before replying to the client with expired data. This essentially enables the serve- stable behavior as specified in RFC 8767 that first tries to resolve before immediately responding with expired data. A recommended value per RFC 8767 is 1800. Setting this to 0 will disable this behavior. Defaults to `\"\"`.",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString(""),
							},
						},
					},
					"logging": schema.SingleNestedAttribute{
						MarkdownDescription: "Logging configuration.",
						Optional:            true,
						Computed:            true,
						Default: objectdefault.StaticValue(types.ObjectValueMust(
							map[string]attr.Type{
								"extended_statistics": types.BoolType,
								"log_queries":         types.BoolType,
								"log_replies":         types.BoolType,
								"tag_query_reply":     types.BoolType,
								"log_local_actions":   types.BoolType,
								"log_servfail":        types.BoolType,
								"verbosity_level":     types.Int64Type,
								"validation_level":    types.Int64Type,
							},
							map[string]attr.Value{
								"extended_statistics": types.BoolValue(false),
								"log_queries":         types.BoolValue(false),
								"log_replies":         types.BoolValue(false),
								"tag_query_reply":     types.BoolValue(false),
								"log_local_actions":   types.BoolValue(false),
								"log_servfail":        types.BoolValue(false),
								"verbosity_level":     types.Int64Value(1),
								"validation_level":    types.Int64Value(0),
							},
						)),
						Attributes: map[string]schema.Attribute{
							"extended_statistics": schema.BoolAttribute{
								MarkdownDescription: "When enabled, extended statistics are printed. Defaults to `false`.",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"log_queries": schema.BoolAttribute{
								MarkdownDescription: "When enabled, prints one line per query to the log, with the log timestamp and IP address, name, type and class. Note that it takes time to print these lines, which makes the server (significantly) slower. Odd (non-printable) characters in names are printed as '?'. Defaults to `false`.",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"log_replies": schema.BoolAttribute{
								MarkdownDescription: "When enabled, prints one line per reply to the log, with the log timestamp and IP address, name, type, class, return code, time to resolve, whether the reply is from the cache and the response size. Note that it takes time to print these lines, which makes the server (significantly) slower. Odd (non-printable) characters in names are printed as '?'. Defaults to `false`.",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"tag_query_reply": schema.BoolAttribute{
								MarkdownDescription: "When enabled, prints the word 'query: ' and 'reply: ' with logged queries and replies. This makes filtering logs easier. Defaults to `false`.",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"log_local_actions": schema.BoolAttribute{
								MarkdownDescription: "When enabled, log lines to inform about local zone actions. These lines are like the local-zone type inform prints out, but they are also printed for the other types of local zones. Defaults to `false`.",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"log_servfail": schema.BoolAttribute{
								MarkdownDescription: "When enabled, log lines that say why queries return SERVFAIL to clients. This is separate from the verbosity debug logs, much smaller, and printed at the error level, not the info level of debug info from verbosity. Defaults to `false`.",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"verbosity_level": schema.Int64Attribute{
								MarkdownDescription: "Select the log verbosity. Level 0 means no verbosity, only errors. Level 1 gives operational information. Level 2 gives detailed operational information. Level 3 gives query level information, output per query. Level 4 gives algorithm level information. Level 5 logs client identification for cache misses. Defaults to `1`.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(1),
								Validators: []validator.Int64{
									int64validator.OneOf(0, 1, 2, 3, 4, 5),
								},
							},
							"validation_level": schema.Int64Attribute{
								MarkdownDescription: "Have the validator print validation failures to the log. Regardless of the verbosity setting. Default is 0, off. At 1, for every user query that fails a line is printed to the logs. This way you can monitor what happens with validation. Use a diagnosis tool, such as dig or drill, to find out why validation is failing for these queries. At 2, not only the query that failed is printed but also the reason why Unbound thought it was wrong and which server sent the faulty data. Defaults to `0`.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(0),
								Validators: []validator.Int64{
									int64validator.OneOf(0, 1, 2),
								},
							},
						},
					},
					"cache": schema.SingleNestedAttribute{
						MarkdownDescription: "Cache configuration.",
						Optional:            true,
						Computed:            true,
						Default: objectdefault.StaticValue(types.ObjectValueMust(
							map[string]attr.Type{
								"enable_prefetch":          types.BoolType,
								"unwanted_reply_threshold": types.Int64Type,
								"msg_cache_size":           types.StringType,
								"rrset_cache_size":         types.StringType,
								"max_ttl":                  types.Int64Type,
								"max_negative_ttl":         types.Int64Type,
								"min_ttl":                  types.Int64Type,
								"host_ttl":                 types.Int64Type,
								"keep_probing_hosts":       types.BoolType,
								"num_hosts":                types.Int64Type,
							},
							map[string]attr.Value{
								"enable_prefetch":          types.BoolValue(false),
								"unwanted_reply_threshold": types.Int64Value(-1),
								"msg_cache_size":           types.StringValue(""),
								"rrset_cache_size":         types.StringValue(""),
								"max_ttl":                  types.Int64Value(-1),
								"max_negative_ttl":         types.Int64Value(-1),
								"min_ttl":                  types.Int64Value(-1),
								"host_ttl":                 types.Int64Value(-1),
								"keep_probing_hosts":       types.BoolValue(false),
								"num_hosts":                types.Int64Value(-1),
							},
						)),
						Attributes: map[string]schema.Attribute{
							"enable_prefetch": schema.BoolAttribute{
								MarkdownDescription: "Message cache elements are prefetched before they expire to help keep the cache up to date. When enabled, this option can cause an increase of around 10% more DNS traffic and load on the server, but frequently requested items will not expire from the cache. Defaults to `false`.",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"unwanted_reply_threshold": schema.Int64Attribute{
								MarkdownDescription: "When enabled, a total number of unwanted replies is kept track of in every thread. When it reaches the threshold, a defensive action is taken and a warning is printed to the log file. This defensive action is to clear the RRSet and message caches, hopefully flushing away any poison. Defaults to `-1`.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(-1),
							},
							"msg_cache_size": schema.StringAttribute{
								MarkdownDescription: "Size of the message cache. The message cache stores DNS rcodes and validation statuses. Valid input is plain bytes, optionally appended with 'k', 'm', or 'g' for kilobytes, megabytes or gigabytes respectively. Defaults to `\"\"`.",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString(""),
								Validators: []validator.String{
									stringvalidator.RegexMatches(
										regexp.MustCompile(`^(\d+[kmg]?)?$`),
										"must be a number optionally followed by 'k', 'm', or 'g' (e.g., '4m', '512k', '1g')",
									),
								},
							},
							"rrset_cache_size": schema.StringAttribute{
								MarkdownDescription: "Size of the RRset cache. Contains the actual RR data. Valid input is plain bytes, optionally appended with 'k', 'm', or 'g' for kilobytes, megabytes or gigabytes respectively. Defaults to `\"\"`.",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString(""),
								Validators: []validator.String{
									stringvalidator.RegexMatches(
										regexp.MustCompile(`^(\d+[kmg]?)?$`),
										"must be a number optionally followed by 'k', 'm', or 'g' (e.g., '4m', '512k', '1g')",
									),
								},
							},
							"max_ttl": schema.Int64Attribute{
								MarkdownDescription: "Configure a maximum Time to live in seconds for RRsets and messages in the cache. When the internal TTL expires the cache item is expired. This can be configured to force the resolver to query for data more often and not trust (very large) TTL values. Defaults to `-1`.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(-1),
							},
							"max_negative_ttl": schema.Int64Attribute{
								MarkdownDescription: "Configure a maximum Negative Time to live in seconds for RRsets and messages in the cache. When the internal TTL expires the negative response cache item is expired. This can be configured to force the resolver to query for data more often in case you wont get a valid answer. Defaults to `-1`.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(-1),
							},
							"min_ttl": schema.Int64Attribute{
								MarkdownDescription: "Configure a minimum Time to live in seconds for RRsets and messages in the cache. If the minimum value kicks in, the data is cached for longer than the domain owner intended, and thus fewer queries are made to look up the data. The 0 value ensures the data in the cache is as the domain owner intended. High values can lead to trouble as the data in the cache might not match up with the actual data anymore. Defaults to `-1`.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(-1),
							},
							"host_ttl": schema.Int64Attribute{
								MarkdownDescription: "Time to live in seconds for entries in the host cache. The host cache contains round-trip timing, lameness and EDNS support information. Defaults to `-1`.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(-1),
							},
							"keep_probing_hosts": schema.BoolAttribute{
								MarkdownDescription: "Whether to keep probing hosts that are down in the infrastructure host cache. Hosts that are down are probed about every 120 seconds with an exponential backoff. If hosts do not respond within this time period, they are marked as down for the duration of the host cache TTL. Defaults to `false`.",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"num_hosts": schema.Int64Attribute{
								MarkdownDescription: "Number of hosts for which information is cached. Defaults to `-1`.",
								Optional:            true,
								Computed:            true,
								Default:             int64default.StaticInt64(-1),
							},
						},
					},
				},
			},
			"acls": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					map[string]attr.Type{
						"default_action": types.StringType,
					},
					map[string]attr.Value{
						"default_action": types.StringValue("allow"),
					},
				)),
				Attributes: map[string]schema.Attribute{
					"default_action": schema.StringAttribute{
						MarkdownDescription: "By default, Unbound will allow queries from all networks. Use this setting to change this behaviour. Since the most specific net block is used, the ACLs as defined in the grid below have preference and will therefore override the behaviour of this setting for the specified networks. Defaults to `\"allow\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("allow"),
						Validators: []validator.String{
							stringvalidator.OneOf("allow", "deny", "refuse"),
						},
					},
				},
			},
			"dnsbl": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					map[string]attr.Type{
						"enabled":             types.BoolType,
						"force_safe_search":   types.BoolType,
						"type":                types.SetType{ElemType: types.StringType},
						"blocklists":          types.SetType{ElemType: types.StringType},
						"whitelist_domains":   types.SetType{ElemType: types.StringType},
						"blocklist_domains":   types.SetType{ElemType: types.StringType},
						"wildcard_domains":    types.SetType{ElemType: types.StringType},
						"destination_address": types.StringType,
						"return_nxdomain":     types.BoolType,
					},
					map[string]attr.Value{
						"enabled":             types.BoolValue(false),
						"force_safe_search":   types.BoolValue(false),
						"type":                types.SetValueMust(types.StringType, []attr.Value{}),
						"blocklists":          types.SetValueMust(types.StringType, []attr.Value{}),
						"whitelist_domains":   types.SetValueMust(types.StringType, []attr.Value{}),
						"blocklist_domains":   types.SetValueMust(types.StringType, []attr.Value{}),
						"wildcard_domains":    types.SetValueMust(types.StringType, []attr.Value{}),
						"destination_address": types.StringValue(""),
						"return_nxdomain":     types.BoolValue(false),
					},
				)),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable the usage of DNS blocklists. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"force_safe_search": schema.BoolAttribute{
						MarkdownDescription: "Force the usage of SafeSearch on Google, DuckDuckGo, Bing, Qwant, PixaBay and YouTube. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"type": schema.SetAttribute{
						MarkdownDescription: "Select which kind of DNSBL you want to use. Defaults to `[]`.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
					"blocklists": schema.SetAttribute{
						MarkdownDescription: "List of domains (URLs) from where blocklist will be downloaded. Defaults to `[]`.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
					"whitelist_domains": schema.SetAttribute{
						MarkdownDescription: "List of domains to whitelist. You can use regular expressions. Defaults to `[]`.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
					"blocklist_domains": schema.SetAttribute{
						MarkdownDescription: "List of domains to blocklist. Only exact matches are supported. Defaults to `[]`.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
					"wildcard_domains": schema.SetAttribute{
						MarkdownDescription: "List of wildcard domains to blocklist. All subdomains of the given domain will be blocked. Blocking first-level domains is not supported. Defaults to `[]`.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
					},
					"destination_address": schema.StringAttribute{
						MarkdownDescription: "Destination ip address for entries in the blocklist (leave as `\"\"` to use default: `0.0.0.0`). Not used when \"return_nxdomain\" is set. Defaults to `\"\"`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"return_nxdomain": schema.BoolAttribute{
						MarkdownDescription: "Whether to use the DNS response code NXDOMAIN instead of a destination address. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
			"forwarding": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					map[string]attr.Type{
						"enabled": types.BoolType,
					},
					map[string]attr.Value{
						"enabled": types.BoolValue(false),
					},
				)),
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "The configured system nameservers will be used to forward queries to. This will override any entry in the grid below, except for entries with a specific domain. DNS over TLS will never be used for any query bound for a system nameserver. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
		},
	}
}

func settingsDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Reads Unbound DNS resolver settings from the upstream system.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `singleton`.",
			},
			"general": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"enabled": dschema.BoolAttribute{
						MarkdownDescription: "Whether Unbound DNS resolver is enabled.",
						Computed:            true,
					},
					"port": dschema.Int64Attribute{
						MarkdownDescription: "The TCP/UDP port used for responding to DNS queries.",
						Computed:            true,
					},
					"listen_interfaces": dschema.SetAttribute{
						MarkdownDescription: "Interface IP addresses used for responding to queries from clients. If an interface has both IPv4 and IPv6 IPs, both are used. Queries to other interface IPs not selected below are discarded. The default behavior is to respond to queries on every available IPv4 and IPv6 address.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"enable_dnssec": dschema.BoolAttribute{
						MarkdownDescription: "Whether DNSSEC validation is enabled.",
						Computed:            true,
					},
					"enable_dns64": dschema.BoolAttribute{
						MarkdownDescription: "When set, Unbound will synthesize AAAA records from A records if no actual AAAA records are present.",
						Computed:            true,
					},
					"dns64_prefix": dschema.StringAttribute{
						MarkdownDescription: "",
						Computed:            true,
					},
					"no_a_records": dschema.BoolAttribute{
						MarkdownDescription: "When set, Unbound will remove all A records from the answer section of all responses.",
						Computed:            true,
					},
					"register_dhcp_leases": dschema.BoolAttribute{
						MarkdownDescription: "When set, then machines that specify their hostname when requesting a DHCP lease will be registered in Unbound, so that their name can be resolved.",
						Computed:            true,
					},
					"dhcp_domain_override": dschema.StringAttribute{
						MarkdownDescription: "The default domain name to use for DHCP lease registration. If empty, the system domain is used.",
						Computed:            true,
					},
					"register_dhcp_static_mappings": dschema.BoolAttribute{
						MarkdownDescription: "When set, then DHCP static mappings will be registered in Unbound, so that their name can be resolved.",
						Computed:            true,
					},
					"register_ipv6_link_local": dschema.BoolAttribute{
						MarkdownDescription: "When set, then IPv6 link-local addresses will also be registered.",
						Computed:            true,
					},
					"register_systems_records": dschema.BoolAttribute{
						MarkdownDescription: "When set, then entries configured under System: Settings: General will be added to Unbound's DNS resolver.",
						Computed:            true,
					},
					"enable_txt_support": dschema.BoolAttribute{
						MarkdownDescription: "When set, then DNS queries for TXT records will also be registered with the source address of the query.",
						Computed:            true,
					},
					"enable_reload_cache_flush": dschema.BoolAttribute{
						MarkdownDescription: "When set, the DNS cache is flushed every time the service is restarted.",
						Computed:            true,
					},
					"local_zone_type": dschema.StringAttribute{
						MarkdownDescription: "The local zone type controls how the local zone is handled. By default unbound will be a 'transparent' forwarding server for the local zone. Other options like 'static' will make it directly authoritative for the local zone.",
						Computed:            true,
					},
					"outgoing_interfaces": dschema.SetAttribute{
						MarkdownDescription: "Interface IP addresses used for outgoing queries to authoritative servers and receiving their replies. If an interface has both IPv4 and IPv6 IPs, both are used. The default behavior is to use all available IPv4 and IPv6 addresses.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"enable_wpad": dschema.BoolAttribute{
						MarkdownDescription: "When set, Unbound will respond to queries for wpad.yourdomain.tld.",
						Computed:            true,
					},
				},
			},
			"advanced": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"hide_identity": dschema.BoolAttribute{
						MarkdownDescription: "When enabled, id.server and hostname.bind queries are refused.",
						Computed:            true,
					},
					"hide_version": dschema.BoolAttribute{
						MarkdownDescription: "When enabled, version.server and version.bind queries are refused.",
						Computed:            true,
					},
					"enable_prefetch_key": dschema.BoolAttribute{
						MarkdownDescription: "When enabled, DNSKEYs are fetched earlier in the validation process when a Delegation signer is encountered. This helps lower the latency of requests but does utilize a little more CPU.",
						Computed:            true,
					},
					"dnssec_stripped": dschema.BoolAttribute{
						MarkdownDescription: "DNSSEC data is required for trust-anchored zones. If such data is absent, the zone becomes bogus. If this is disabled and no DNSSEC data is received, then the zone is made insecure.",
						Computed:            true,
					},
					"aggressive_nsec": dschema.BoolAttribute{
						MarkdownDescription: "Whether to enable RFC8198-based aggressive use of the DNSSEC-Validated cache. Helps to reduce the query rate towards targets but may lead to false negative responses if there are errors in the zone config.",
						Computed:            true,
					},
					"qname_min_strict": dschema.BoolAttribute{
						MarkdownDescription: "Whether to send minimum amount of information to upstream servers to enhance privacy. Do not fall-back to sending full QNAME to potentially broken nameservers. A lot of domains will not be resolvable when this option in enabled. Only use if you know what you are doing.",
						Computed:            true,
					},
					"outgoing_num_tcp": dschema.Int64Attribute{
						MarkdownDescription: "The number of outgoing TCP buffers to allocate per thread. If 0 is selected then no TCP queries, to authoritative servers, are done.",
						Computed:            true,
					},
					"incoming_num_tcp": dschema.Int64Attribute{
						MarkdownDescription: "The number of incoming TCP buffers to allocate per thread. If 0 is selected then no TCP queries, from clients, are accepted.",
						Computed:            true,
					},
					"num_queries_per_thread": dschema.Int64Attribute{
						MarkdownDescription: "The number of queries that every thread will service simultaneously. If more queries arrive that need to be serviced, and no queries can be jostled out (see \"jostle_timeout\"), then these queries are dropped. This forces the client to resend after a timeout, allowing the server time to work on the existing queries.",
						Computed:            true,
					},
					"outgoing_range": dschema.Int64Attribute{
						MarkdownDescription: "The number of ports to open. This number of file descriptors can be opened per thread. Larger numbers need extra resources from the operating system. For performance a very large value is best. For reference, usually double the amount of queries per thread is used.",
						Computed:            true,
					},
					"jostle_timeout": dschema.Int64Attribute{
						MarkdownDescription: "This timeout is used for when the server is very busy. Set to a value that usually results in one round-trip to the authority servers. If too many queries arrive, then 50% of the queries are allowed to run to completion, and the other 50% are replaced with the new incoming query if they have already spent more than their allowed time. This protects against denial of service by slow queries or high query rates.",
						Computed:            true,
					},
					"discard_timeout": dschema.Int64Attribute{
						MarkdownDescription: "The wait time in msec where recursion requests are dropped. This is to stop a large number of replies from accumulating. If 'Serve Expired Responses' is enabled this field should be set greater than 'Client Expired Response Timeout', otherwise, these late responses will not update the cache. The value 0 disables it. Default 1900. This setting may increase the \"request queue exceeded\" counter.",
						Computed:            true,
					},
					"private_domains": dschema.SetAttribute{
						MarkdownDescription: "List of domains to mark as private. These domains and all its subdomains are allowed to contain private addresses.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"private_addresses": dschema.SetAttribute{
						MarkdownDescription: "These are addresses on your private network, and are not allowed to be returned for public internet names. Any occurrence of such addresses are removed from DNS answers. Additionally, the DNSSEC validator may mark the answers bogus. This protects against so-called DNS Rebinding (Only applicable when DNS rebind check is enabled in `System->Settings->Administration`).",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"insecure_domains": dschema.SetAttribute{
						MarkdownDescription: "List of domains to mark as insecure. DNSSEC chain of trust is ignored towards the domain name.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"serve_expired": dschema.SingleNestedAttribute{
						MarkdownDescription: "Serve expired cache configuration.",
						Computed:            true,
						Attributes: map[string]dschema.Attribute{
							"enabled": dschema.BoolAttribute{
								MarkdownDescription: "Serve expired responses from the cache with a TTL of 0 without waiting for the actual resolution to finish. The TTL can be modified with \"record_reply_ttl\" value.",
								Computed:            true,
							},
							"record_reply_ttl": dschema.StringAttribute{
								MarkdownDescription: "TTL value to use when replying with expired data. If \"Client Expired Response Timeout\" is also used then it is recommended to use 30 as the value as per RFC 8767.",
								Computed:            true,
							},
							"ttl": dschema.StringAttribute{
								MarkdownDescription: "Limits the serving of expired responses to the configured amount of seconds after expiration. A value of 0 disables the limit. A suggested value per RFC 8767 is between 86400 (1 day) and 259200 (3 days).",
								Computed:            true,
							},
							"reset_ttl": dschema.BoolAttribute{
								MarkdownDescription: "Whether to set the TTL of expired records to the \"TTL for Expired Responses\" value after a failed attempt to retrieve the record from an upstream server. This makes sure that the expired records will be served as long as there are queries for it.",
								Computed:            true,
							},
							"client_response_timeout": dschema.StringAttribute{
								MarkdownDescription: "Time in milliseconds before replying to the client with expired data. This essentially enables the serve- stable behavior as specified in RFC 8767 that first tries to resolve before immediately responding with expired data. A recommended value per RFC 8767 is 1800. Setting this to 0 will disable this behavior.",
								Computed:            true,
							},
						},
					},
					"logging": dschema.SingleNestedAttribute{
						MarkdownDescription: "Logging configuration.",
						Computed:            true,
						Attributes: map[string]dschema.Attribute{
							"extended_statistics": dschema.BoolAttribute{
								MarkdownDescription: "When enabled, extended statistics are printed.",
								Computed:            true,
							},
							"log_queries": dschema.BoolAttribute{
								MarkdownDescription: "When enabled, prints one line per query to the log, with the log timestamp and IP address, name, type and class. Note that it takes time to print these lines, which makes the server (significantly) slower. Odd (non-printable) characters in names are printed as '?'.",
								Computed:            true,
							},
							"log_replies": dschema.BoolAttribute{
								MarkdownDescription: "When enabled, prints one line per reply to the log, with the log timestamp and IP address, name, type, class, return code, time to resolve, whether the reply is from the cache and the response size. Note that it takes time to print these lines, which makes the server (significantly) slower. Odd (non-printable) characters in names are printed as '?'.",
								Computed:            true,
							},
							"tag_query_reply": dschema.BoolAttribute{
								MarkdownDescription: "When enabled, prints the word 'query: ' and 'reply: ' with logged queries and replies. This makes filtering logs easier.",
								Computed:            true,
							},
							"log_local_actions": dschema.BoolAttribute{
								MarkdownDescription: "When enabled, log lines to inform about local zone actions. These lines are like the local-zone type inform prints out, but they are also printed for the other types of local zones.",
								Computed:            true,
							},
							"log_servfail": dschema.BoolAttribute{
								MarkdownDescription: "When enabled, log lines that say why queries return SERVFAIL to clients. This is separate from the verbosity debug logs, much smaller, and printed at the error level, not the info level of debug info from verbosity.",
								Computed:            true,
							},
							"verbosity_level": dschema.Int64Attribute{
								MarkdownDescription: "Select the log verbosity. Level 0 means no verbosity, only errors. Level 1 gives operational information. Level 2 gives detailed operational information. Level 3 gives query level information, output per query. Level 4 gives algorithm level information. Level 5 logs client identification for cache misses.",
								Computed:            true,
							},
							"validation_level": dschema.Int64Attribute{
								MarkdownDescription: "Have the validator print validation failures to the log. Regardless of the verbosity setting. Default is 0, off. At 1, for every user query that fails a line is printed to the logs. This way you can monitor what happens with validation. Use a diagnosis tool, such as dig or drill, to find out why validation is failing for these queries. At 2, not only the query that failed is printed but also the reason why Unbound thought it was wrong and which server sent the faulty data.",
								Computed:            true,
							},
						},
					},
					"cache": dschema.SingleNestedAttribute{
						MarkdownDescription: "Cache configuration.",
						Computed:            true,
						Attributes: map[string]dschema.Attribute{
							"enable_prefetch": dschema.BoolAttribute{
								MarkdownDescription: "Message cache elements are prefetched before they expire to help keep the cache up to date. When enabled, this option can cause an increase of around 10% more DNS traffic and load on the server, but frequently requested items will not expire from the cache.",
								Computed:            true,
							},
							"unwanted_reply_threshold": dschema.Int64Attribute{
								MarkdownDescription: "When enabled, a total number of unwanted replies is kept track of in every thread. When it reaches the threshold, a defensive action is taken and a warning is printed to the log file. This defensive action is to clear the RRSet and message caches, hopefully flushing away any poison.",
								Computed:            true,
							},
							"msg_cache_size": dschema.StringAttribute{
								MarkdownDescription: "Size of the message cache. The message cache stores DNS rcodes and validation statuses. Valid input is plain bytes, optionally appended with 'k', 'm', or 'g' for kilobytes, megabytes or gigabytes respectively.",
								Computed:            true,
							},
							"rrset_cache_size": dschema.StringAttribute{
								MarkdownDescription: "Size of the RRset cache. Contains the actual RR data. Valid input is plain bytes, optionally appended with 'k', 'm', or 'g' for kilobytes, megabytes or gigabytes respectively.",
								Computed:            true,
							},
							"max_ttl": dschema.Int64Attribute{
								MarkdownDescription: "Configure a maximum Time to live in seconds for RRsets and messages in the cache. When the internal TTL expires the cache item is expired. This can be configured to force the resolver to query for data more often and not trust (very large) TTL values.",
								Computed:            true,
							},
							"max_negative_ttl": dschema.Int64Attribute{
								MarkdownDescription: "Configure a maximum Negative Time to live in seconds for RRsets and messages in the cache. When the internal TTL expires the negative response cache item is expired. This can be configured to force the resolver to query for data more often in case you wont get a valid answer.",
								Computed:            true,
							},
							"min_ttl": dschema.Int64Attribute{
								MarkdownDescription: "Configure a minimum Time to live in seconds for RRsets and messages in the cache. If the minimum value kicks in, the data is cached for longer than the domain owner intended, and thus fewer queries are made to look up the data. The 0 value ensures the data in the cache is as the domain owner intended. High values can lead to trouble as the data in the cache might not match up with the actual data anymore.",
								Computed:            true,
							},
							"host_ttl": dschema.Int64Attribute{
								MarkdownDescription: "Time to live in seconds for entries in the host cache. The host cache contains round-trip timing, lameness and EDNS support information.",
								Computed:            true,
							},
							"keep_probing_hosts": dschema.BoolAttribute{
								MarkdownDescription: "Whether to keep probing hosts that are down in the infrastructure host cache. Hosts that are down are probed about every 120 seconds with an exponential backoff. If hosts do not respond within this time period, they are marked as down for the duration of the host cache TTL.",
								Computed:            true,
							},
							"num_hosts": dschema.Int64Attribute{
								MarkdownDescription: "Number of hosts for which information is cached.",
								Computed:            true,
							},
						},
					},
				},
			},
			"acls": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"default_action": dschema.StringAttribute{
						MarkdownDescription: "By default, Unbound will allow queries from all networks. Use this setting to change this behaviour. Since the most specific net block is used, the ACLs as defined in the grid below have preference and will therefore override the behaviour of this setting for the specified networks.",
						Computed:            true,
					},
				},
			},
			"dnsbl": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"enabled": dschema.BoolAttribute{
						MarkdownDescription: "Enable the usage of DNS blocklists.",
						Computed:            true,
					},
					"force_safe_search": dschema.BoolAttribute{
						MarkdownDescription: "Force the usage of SafeSearch on Google, DuckDuckGo, Bing, Qwant, PixaBay and YouTube.",
						Computed:            true,
					},
					"type": dschema.SetAttribute{
						MarkdownDescription: "Select which kind of DNSBL you want to use.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"blocklists": dschema.SetAttribute{
						MarkdownDescription: "List of domains (URLs) from where blocklist will be downloaded.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"whitelist_domains": dschema.SetAttribute{
						MarkdownDescription: "List of domains to whitelist. You can use regular expressions.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"blocklist_domains": dschema.SetAttribute{
						MarkdownDescription: "List of domains to blocklist. Only exact matches are supported.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"wildcard_domains": dschema.SetAttribute{
						MarkdownDescription: "List of wildcard domains to blocklist. All subdomains of the given domain will be blocked. Blocking first-level domains is not supported.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"destination_address": dschema.StringAttribute{
						MarkdownDescription: "Destination ip address for entries in the blocklist (leave as `\"\"` to use default: `0.0.0.0`). Not used when \"return_nxdomain\" is set.",
						Computed:            true,
					},
					"return_nxdomain": dschema.BoolAttribute{
						MarkdownDescription: "Whether to use the DNS response code NXDOMAIN instead of a destination address.",
						Computed:            true,
					},
				},
			},
			"forwarding": dschema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]dschema.Attribute{
					"enabled": dschema.BoolAttribute{
						MarkdownDescription: "The configured system nameservers will be used to forward queries to. This will override any entry in the grid below, except for entries with a specific domain. DNS over TLS will never be used for any query bound for a system nameserver.",
						Computed:            true,
					},
				},
			},
		},
	}
}

// convertSettingsSchemaToStruct converts TF schema to upstream API struct
func convertSettingsSchemaToStruct(d *settingsResourceModel) (*unbound.Settings, error) {
	result := &unbound.Settings{}

	// Parse 'General' block
	if d.General != nil {
		result.General.Enabled = tools.BoolToString(d.General.Enabled.ValueBool())
		result.General.Port = tools.Int64ToString(d.General.Port.ValueInt64())
		result.General.DNSSEC = tools.BoolToString(d.General.EnableDNSSEC.ValueBool())
		result.General.DNS64 = tools.BoolToString(d.General.EnableDNS64.ValueBool())
		result.General.DNS64Prefix = d.General.DNS64Prefix.ValueString()
		result.General.NoARecords = tools.BoolToString(d.General.NoARecords.ValueBool())
		result.General.RegDHCP = tools.BoolToString(d.General.RegisterDHCPLeases.ValueBool())
		result.General.RegDHCPDomain = d.General.DHCPDomainOverride.ValueString()
		result.General.RegDHCPStatic = tools.BoolToString(d.General.RegisterDHCPStaticMappings.ValueBool())
		result.General.NoRegLLAddr6 = tools.BoolToString(!d.General.RegisterIPv6LinkLocal.ValueBool())
		result.General.NoRegRecords = tools.BoolToString(!d.General.RegisterSystemsRecords.ValueBool())
		result.General.TXTSupport = tools.BoolToString(d.General.EnableTXTSupport.ValueBool())
		result.General.CacheFlush = tools.BoolToString(d.General.EnableReloadCacheFlush.ValueBool())
		result.General.LocalZoneType = api.SelectedMap(d.General.LocalZoneType.ValueString())
		result.General.EnableWPAD = tools.BoolToString(d.General.EnableWPAD.ValueBool())

		// ListenInterfaces
		if !d.General.ListenInterfaces.IsNull() && !d.General.ListenInterfaces.IsUnknown() {
			var activeInterfaceList []string
			d.General.ListenInterfaces.ElementsAs(context.Background(), &activeInterfaceList, false)
			sort.Strings(activeInterfaceList)
			result.General.ActiveInterface = api.SelectedMapList(activeInterfaceList)
		}

		// OutgoingInterfaces
		if !d.General.OutgoingInterfaces.IsNull() && !d.General.OutgoingInterfaces.IsUnknown() {
			var outgoingInterfaceList []string
			d.General.OutgoingInterfaces.ElementsAs(context.Background(), &outgoingInterfaceList, false)
			sort.Strings(outgoingInterfaceList)
			result.General.OutgoingInterface = api.SelectedMapList(outgoingInterfaceList)
		}
	}

	// Parse 'Advanced' block
	if d.Advanced != nil {
		result.Advanced.HideIdentity = tools.BoolToString(d.Advanced.HideIdentity.ValueBool())
		result.Advanced.HideVersion = tools.BoolToString(d.Advanced.HideVersion.ValueBool())
		result.Advanced.PrefetchKey = tools.BoolToString(d.Advanced.EnablePrefetchKey.ValueBool())
		result.Advanced.DNSSECStripped = tools.BoolToString(d.Advanced.DNSSECStripped.ValueBool())
		result.Advanced.AggressiveNSEC = tools.BoolToString(d.Advanced.AggressiveNSEC.ValueBool())
		result.Advanced.QnameMinStrict = tools.BoolToString(d.Advanced.QnameMinStrict.ValueBool())
		result.Advanced.OutgoingNumTCP = tools.Int64ToStringNegative(d.Advanced.OutgoingNumTCP.ValueInt64())
		result.Advanced.IncomingNumTCP = tools.Int64ToStringNegative(d.Advanced.IncomingNumTCP.ValueInt64())
		result.Advanced.NumQueriesPerThread = tools.Int64ToStringNegative(d.Advanced.NumQueriesPerThread.ValueInt64())
		result.Advanced.OutgoingRange = tools.Int64ToStringNegative(d.Advanced.OutgoingRange.ValueInt64())
		result.Advanced.JostleTimeout = tools.Int64ToStringNegative(d.Advanced.JostleTimeout.ValueInt64())
		result.Advanced.DiscardTimeout = tools.Int64ToStringNegative(d.Advanced.DiscardTimeout.ValueInt64())

		// Parse 'ServeExpired' nested block
		if d.Advanced.ServeExpired != nil {
			result.Advanced.ServeExpired = tools.BoolToString(d.Advanced.ServeExpired.Enabled.ValueBool())
			result.Advanced.ServeExpiredReplyTTL = d.Advanced.ServeExpired.RecordReplyTTL.ValueString()
			result.Advanced.ServeExpiredTTL = d.Advanced.ServeExpired.TTL.ValueString()
			result.Advanced.ServeExpiredTTLReset = tools.BoolToString(d.Advanced.ServeExpired.ResetTTL.ValueBool())
			result.Advanced.ServeExpiredClientTimeout = d.Advanced.ServeExpired.ClientResponseTimeout.ValueString()
		}

		// Parse 'Logging' nested block
		if d.Advanced.Logging != nil {
			result.Advanced.ExtendedStatistics = tools.BoolToString(d.Advanced.Logging.ExtendedStatistics.ValueBool())
			result.Advanced.LogQueries = tools.BoolToString(d.Advanced.Logging.LogQueries.ValueBool())
			result.Advanced.LogReplies = tools.BoolToString(d.Advanced.Logging.LogReplies.ValueBool())
			result.Advanced.LogTagQueryReply = tools.BoolToString(d.Advanced.Logging.TagQueryReply.ValueBool())
			result.Advanced.LogServFail = tools.BoolToString(d.Advanced.Logging.LogServFail.ValueBool())
			result.Advanced.LogLocalActions = tools.BoolToString(d.Advanced.Logging.LogLocalActions.ValueBool())
			result.Advanced.LogVerbosity = api.SelectedMap(tools.Int64ToStringNegative(d.Advanced.Logging.VerbosityLevel.ValueInt64()))
			result.Advanced.ValLogLevel = api.SelectedMap(tools.Int64ToStringNegative(d.Advanced.Logging.ValidationLevel.ValueInt64()))
		}

		// Parse 'Cache' nested block
		if d.Advanced.Cache != nil {
			result.Advanced.Prefetch = tools.BoolToString(d.Advanced.Cache.EnablePrefetch.ValueBool())
			result.Advanced.MsgCacheSize = d.Advanced.Cache.MsgCacheSize.ValueString()
			result.Advanced.RRSetCacheSize = d.Advanced.Cache.RRSetCacheSize.ValueString()
			result.Advanced.CacheMaxTTL = tools.Int64ToStringNegative(d.Advanced.Cache.MaxTTL.ValueInt64())
			result.Advanced.CacheMaxNegativeTTL = tools.Int64ToStringNegative(d.Advanced.Cache.MaxNegativeTTL.ValueInt64())
			result.Advanced.CacheMinTTL = tools.Int64ToStringNegative(d.Advanced.Cache.MinTTL.ValueInt64())
			result.Advanced.InfraHostTTL = tools.Int64ToStringNegative(d.Advanced.Cache.HostTTL.ValueInt64())
			result.Advanced.InfraKeepProbing = tools.BoolToString(d.Advanced.Cache.KeepProbingHosts.ValueBool())
			result.Advanced.InfraCacheNumHosts = tools.Int64ToStringNegative(d.Advanced.Cache.NumHosts.ValueInt64())
			result.Advanced.UnwantedReplyThreshold = tools.Int64ToStringNegative(d.Advanced.Cache.UnwantedReplyThreshold.ValueInt64())
		}

		// PrivateDomains
		if !d.Advanced.PrivateDomains.IsNull() && !d.Advanced.PrivateDomains.IsUnknown() {
			var privateDomainList []string
			d.Advanced.PrivateDomains.ElementsAs(context.Background(), &privateDomainList, false)
			sort.Strings(privateDomainList)
			result.Advanced.PrivateDomain = api.SelectedMapList(privateDomainList)
		}

		// PrivateAddresses
		if !d.Advanced.PrivateAddresses.IsNull() && !d.Advanced.PrivateAddresses.IsUnknown() {
			var privateAddressList []string
			d.Advanced.PrivateAddresses.ElementsAs(context.Background(), &privateAddressList, false)
			sort.Strings(privateAddressList)
			result.Advanced.PrivateAddress = api.SelectedMapList(privateAddressList)
		}

		// InsecureDomains
		if !d.Advanced.InsecureDomains.IsNull() && !d.Advanced.InsecureDomains.IsUnknown() {
			var insecureDomainList []string
			d.Advanced.InsecureDomains.ElementsAs(context.Background(), &insecureDomainList, false)
			sort.Strings(insecureDomainList)
			result.Advanced.InsecureDomain = api.SelectedMapList(insecureDomainList)
		}
	}

	// Parse 'ACLs' block
	if d.ACLs != nil {
		result.ACLs.DefaultAction = api.SelectedMap(d.ACLs.DefaultAction.ValueString())
	}

	// Parse 'DNSBL' block
	if d.DNSBL != nil {
		result.DNSBL.Enabled = tools.BoolToString(d.DNSBL.Enabled.ValueBool())
		result.DNSBL.SafeSearch = tools.BoolToString(d.DNSBL.ForceSafeSearch.ValueBool())
		result.DNSBL.Address = d.DNSBL.DestinationAddress.ValueString()
		result.DNSBL.NXDomain = tools.BoolToString(d.DNSBL.ReturnNXDomain.ValueBool())

		// Type
		if !d.DNSBL.Type.IsNull() && !d.DNSBL.Type.IsUnknown() {
			var typeList []string
			d.DNSBL.Type.ElementsAs(context.Background(), &typeList, false)
			sort.Strings(typeList)
			result.DNSBL.Type = api.SelectedMapList(typeList)
		}

		// Blocklists
		if !d.DNSBL.Blocklists.IsNull() && !d.DNSBL.Blocklists.IsUnknown() {
			var listsList []string
			d.DNSBL.Blocklists.ElementsAs(context.Background(), &listsList, false)
			sort.Strings(listsList)
			result.DNSBL.Lists = api.SelectedMapList(listsList)
		}

		// WhitelistDomains
		if !d.DNSBL.WhitelistDomains.IsNull() && !d.DNSBL.WhitelistDomains.IsUnknown() {
			var whitelistsList []string
			d.DNSBL.WhitelistDomains.ElementsAs(context.Background(), &whitelistsList, false)
			sort.Strings(whitelistsList)
			result.DNSBL.Whitelists = api.SelectedMapList(whitelistsList)
		}

		// BlocklistDomains
		if !d.DNSBL.BlocklistDomains.IsNull() && !d.DNSBL.BlocklistDomains.IsUnknown() {
			var blocklistsList []string
			d.DNSBL.BlocklistDomains.ElementsAs(context.Background(), &blocklistsList, false)
			sort.Strings(blocklistsList)
			result.DNSBL.Blocklists = api.SelectedMapList(blocklistsList)
		}

		// WildcardDomains
		if !d.DNSBL.WildcardDomains.IsNull() && !d.DNSBL.WildcardDomains.IsUnknown() {
			var wildcardsList []string
			d.DNSBL.WildcardDomains.ElementsAs(context.Background(), &wildcardsList, false)
			sort.Strings(wildcardsList)
			result.DNSBL.Wildcards = api.SelectedMapList(wildcardsList)
		}
	}

	// Parse 'Forwarding' block
	if d.Forwarding != nil {
		result.Forwarding.Enabled = tools.BoolToString(d.Forwarding.Enabled.ValueBool())
	}

	return result, nil
}

// convertSettingsStructToSchema converts upstream API struct to TF schema
func convertSettingsStructToSchema(dRaw *unbound.SettingsMonad) (*settingsResourceModel, error) {
	// Unpack monad
	d := dRaw.Unbound

	model := &settingsResourceModel{
		Id: types.StringValue("unbound_settings"),
	}

	// Parse 'General' block
	model.General = &settingsGeneralBlock{
		Enabled:                    types.BoolValue(tools.StringToBool(d.General.Enabled)),
		Port:                       tools.StringToInt64Null(d.General.Port),
		EnableDNSSEC:               types.BoolValue(tools.StringToBool(d.General.DNSSEC)),
		EnableDNS64:                types.BoolValue(tools.StringToBool(d.General.DNS64)),
		DNS64Prefix:                types.StringValue(d.General.DNS64Prefix),
		NoARecords:                 types.BoolValue(tools.StringToBool(d.General.NoARecords)),
		RegisterDHCPLeases:         types.BoolValue(tools.StringToBool(d.General.RegDHCP)),
		DHCPDomainOverride:         types.StringValue(d.General.RegDHCPDomain),
		RegisterDHCPStaticMappings: types.BoolValue(tools.StringToBool(d.General.RegDHCPStatic)),
		RegisterIPv6LinkLocal:      types.BoolValue(!tools.StringToBool(d.General.NoRegLLAddr6)),
		RegisterSystemsRecords:     types.BoolValue(!tools.StringToBool(d.General.NoRegRecords)),
		EnableTXTSupport:           types.BoolValue(tools.StringToBool(d.General.TXTSupport)),
		EnableReloadCacheFlush:     types.BoolValue(tools.StringToBool(d.General.CacheFlush)),
		LocalZoneType:              types.StringValue(d.General.LocalZoneType.String()),
		EnableWPAD:                 types.BoolValue(tools.StringToBool(d.General.EnableWPAD)),
		ListenInterfaces:           tools.StringSliceToSet(d.General.ActiveInterface),
		OutgoingInterfaces:         tools.StringSliceToSet(d.General.OutgoingInterface),
	}

	// Parse 'Advanced' block
	model.Advanced = &settingsAdvancedBlock{
		HideIdentity:        types.BoolValue(tools.StringToBool(d.Advanced.HideIdentity)),
		HideVersion:         types.BoolValue(tools.StringToBool(d.Advanced.HideVersion)),
		EnablePrefetchKey:   types.BoolValue(tools.StringToBool(d.Advanced.PrefetchKey)),
		DNSSECStripped:      types.BoolValue(tools.StringToBool(d.Advanced.DNSSECStripped)),
		AggressiveNSEC:      types.BoolValue(tools.StringToBool(d.Advanced.AggressiveNSEC)),
		QnameMinStrict:      types.BoolValue(tools.StringToBool(d.Advanced.QnameMinStrict)),
		OutgoingNumTCP:      types.Int64Value(tools.StringToInt64(d.Advanced.OutgoingNumTCP)),
		IncomingNumTCP:      types.Int64Value(tools.StringToInt64(d.Advanced.IncomingNumTCP)),
		NumQueriesPerThread: types.Int64Value(tools.StringToInt64(d.Advanced.NumQueriesPerThread)),
		OutgoingRange:       types.Int64Value(tools.StringToInt64(d.Advanced.OutgoingRange)),
		JostleTimeout:       types.Int64Value(tools.StringToInt64(d.Advanced.JostleTimeout)),
		DiscardTimeout:      types.Int64Value(tools.StringToInt64(d.Advanced.DiscardTimeout)),
		PrivateDomains:      tools.StringSliceToSet(d.Advanced.PrivateDomain),
		PrivateAddresses:    tools.StringSliceToSet(d.Advanced.PrivateAddress),
		InsecureDomains:     tools.StringSliceToSet(d.Advanced.InsecureDomain),
		// Initialize nested blocks
		ServeExpired: &settingsAdvancedServeExpiredBlock{
			Enabled:               types.BoolValue(tools.StringToBool(d.Advanced.ServeExpired)),
			RecordReplyTTL:        types.StringValue(d.Advanced.ServeExpiredReplyTTL),
			TTL:                   types.StringValue(d.Advanced.ServeExpiredTTL),
			ResetTTL:              types.BoolValue(tools.StringToBool(d.Advanced.ServeExpiredTTLReset)),
			ClientResponseTimeout: types.StringValue(d.Advanced.ServeExpiredClientTimeout),
		},
		Logging: &settingsAdvancedLoggingBlock{
			ExtendedStatistics: types.BoolValue(tools.StringToBool(d.Advanced.ExtendedStatistics)),
			LogQueries:         types.BoolValue(tools.StringToBool(d.Advanced.LogQueries)),
			LogReplies:         types.BoolValue(tools.StringToBool(d.Advanced.LogReplies)),
			TagQueryReply:      types.BoolValue(tools.StringToBool(d.Advanced.LogTagQueryReply)),
			LogServFail:        types.BoolValue(tools.StringToBool(d.Advanced.LogServFail)),
			LogLocalActions:    types.BoolValue(tools.StringToBool(d.Advanced.LogLocalActions)),
			VerbosityLevel:     types.Int64Value(tools.StringToInt64(d.Advanced.LogVerbosity.String())),
			ValidationLevel:    types.Int64Value(tools.StringToInt64(d.Advanced.ValLogLevel.String())),
		},
		Cache: &settingsAdvancedCacheBlock{
			EnablePrefetch:         types.BoolValue(tools.StringToBool(d.Advanced.Prefetch)),
			UnwantedReplyThreshold: types.Int64Value(tools.StringToInt64(d.Advanced.UnwantedReplyThreshold)),
			MsgCacheSize:           types.StringValue(d.Advanced.MsgCacheSize),
			RRSetCacheSize:         types.StringValue(d.Advanced.RRSetCacheSize),
			MaxTTL:                 types.Int64Value(tools.StringToInt64(d.Advanced.CacheMaxTTL)),
			MaxNegativeTTL:         types.Int64Value(tools.StringToInt64(d.Advanced.CacheMaxNegativeTTL)),
			MinTTL:                 types.Int64Value(tools.StringToInt64(d.Advanced.CacheMinTTL)),
			HostTTL:                types.Int64Value(tools.StringToInt64(d.Advanced.InfraHostTTL)),
			KeepProbingHosts:       types.BoolValue(tools.StringToBool(d.Advanced.InfraKeepProbing)),
			NumHosts:               types.Int64Value(tools.StringToInt64(d.Advanced.InfraCacheNumHosts)),
		},
	}

	// Parse 'ACLs' block
	model.ACLs = &settingsACLsBlock{
		DefaultAction: types.StringValue(d.ACLs.DefaultAction.String()),
	}

	// Parse 'DNSBL' block
	model.DNSBL = &settingsDNSBLBlock{
		Enabled:            types.BoolValue(tools.StringToBool(d.DNSBL.Enabled)),
		ForceSafeSearch:    types.BoolValue(tools.StringToBool(d.DNSBL.SafeSearch)),
		DestinationAddress: types.StringValue(d.DNSBL.Address),
		ReturnNXDomain:     types.BoolValue(tools.StringToBool(d.DNSBL.NXDomain)),
		Type:               tools.StringSliceToSet(d.DNSBL.Type),
		Blocklists:         tools.StringSliceToSet(d.DNSBL.Lists),
		WhitelistDomains:   tools.StringSliceToSet(d.DNSBL.Whitelists),
		BlocklistDomains:   tools.StringSliceToSet(d.DNSBL.Blocklists),
		WildcardDomains:    tools.StringSliceToSet(d.DNSBL.Wildcards),
	}

	// Parse 'Forwarding' block
	model.Forwarding = &settingsForwardingBlock{
		Enabled: types.BoolValue(tools.StringToBool(d.Forwarding.Enabled)),
	}

	return model, nil
}
