package openvpn

import (
	"context"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/openvpn"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// instanceResourceModel describes the resource data model.
type instanceResourceModel struct {
	Enabled              types.Bool   `tfsdk:"enabled"`
	Role                 types.String `tfsdk:"role"`
	VPNID                types.Int64  `tfsdk:"vpn_id"`
	Description          types.String `tfsdk:"description"`
	DevType              types.String `tfsdk:"dev_type"`
	Protocol             types.String `tfsdk:"protocol"`
	Port                 types.Int64  `tfsdk:"port"`
	PortShare            types.String `tfsdk:"port_share"`
	Local                types.String `tfsdk:"local"`
	Remote               types.Set    `tfsdk:"remote"`
	Topology             types.String `tfsdk:"topology"`
	Server               types.String `tfsdk:"server"`
	ServerIPv6           types.String `tfsdk:"server_ipv6"`
	NoPool               types.Bool   `tfsdk:"no_pool"`
	BridgeGateway        types.String `tfsdk:"bridge_gateway"`
	BridgePool           types.String `tfsdk:"bridge_pool"`
	Route                types.Set    `tfsdk:"route"`
	PushRoute            types.Set    `tfsdk:"push_route"`
	PushExcludedRoutes   types.Set    `tfsdk:"push_excluded_routes"`
	Certificate          types.String `tfsdk:"certificate"`
	CRL                  types.String `tfsdk:"crl"`
	CertificateAuthority types.String `tfsdk:"certificate_authority"`
	CertDepth            types.String `tfsdk:"cert_depth"`
	RemoteCertTLS        types.Bool   `tfsdk:"remote_cert_tls"`
	VerifyClientCert     types.String `tfsdk:"verify_client_cert"`
	UseOCSP              types.Bool   `tfsdk:"use_ocsp"`
	AuthDigest           types.String `tfsdk:"auth_digest"`
	DataCiphers          types.Set    `tfsdk:"data_ciphers"`
	DataCiphersFallback  types.String `tfsdk:"data_ciphers_fallback"`
	TLSKey               types.String `tfsdk:"tls_key"`
	AuthMode             types.Set    `tfsdk:"auth_mode"`
	LocalGroup           types.String `tfsdk:"local_group"`
	VariousFlags         types.Set    `tfsdk:"various_flags"`
	VariousPushFlags     types.Set    `tfsdk:"various_push_flags"`
	PushInactive         types.Int64  `tfsdk:"push_inactive"`
	UsernameAsCommonName types.Bool   `tfsdk:"username_as_common_name"`
	StrictUserCN         types.String `tfsdk:"strict_user_cn"`
	Username             types.String `tfsdk:"username"`
	Password             types.String `tfsdk:"password"`
	MaxClients           types.Int64  `tfsdk:"max_clients"`
	KeepaliveInterval    types.Int64  `tfsdk:"keepalive_interval"`
	KeepaliveTimeout     types.Int64  `tfsdk:"keepalive_timeout"`
	RenegSec             types.Int64  `tfsdk:"reneg_sec"`
	AuthGenToken         types.Int64  `tfsdk:"auth_gen_token"`
	AuthGenTokenRenewal  types.Int64  `tfsdk:"auth_gen_token_renewal"`
	AuthGenTokenSecret   types.String `tfsdk:"auth_gen_token_secret"`
	ProvisionExclusive   types.Bool   `tfsdk:"provision_exclusive"`
	RedirectGateway      types.Set    `tfsdk:"redirect_gateway"`
	RouteMetric          types.Int64  `tfsdk:"route_metric"`
	RegisterDNS          types.Bool   `tfsdk:"register_dns"`
	DNSDomain            types.Set    `tfsdk:"dns_domain"`
	DNSDomainSearch      types.Set    `tfsdk:"dns_domain_search"`
	DNSServers           types.Set    `tfsdk:"dns_servers"`
	NTPServers           types.Set    `tfsdk:"ntp_servers"`
	TunMTU               types.Int64  `tfsdk:"tun_mtu"`
	Fragment             types.Int64  `tfsdk:"fragment"`
	MSSFix               types.Bool   `tfsdk:"mss_fix"`
	CARPDependOn         types.String `tfsdk:"carp_depend_on"`
	CompressMigrate      types.Bool   `tfsdk:"compress_migrate"`
	IfConfigPoolPersist  types.Bool   `tfsdk:"ifconfig_pool_persist"`
	HTTPProxy            types.String `tfsdk:"http_proxy"`
	VerifyX509Name       types.String `tfsdk:"verify_x509_name"`

	Id types.String `tfsdk:"id"`
}

func instanceResourceSchema() schema.Schema {
	emptySet := setdefault.StaticValue(tools.EmptySetValue(types.StringType))

	return schema.Schema{
		MarkdownDescription: "OpenVPN instances (servers or clients) configured under `VPN > OpenVPN > Instances`. Each instance corresponds to one OpenVPN daemon process on the OPNsense host.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this instance. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "Whether this instance acts as a `server` or a `client`. Defaults to `server`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("server"),
				Validators: []validator.String{
					stringvalidator.OneOf("server", "client"),
				},
			},
			"vpn_id": schema.Int64Attribute{
				MarkdownDescription: "Numeric VPN ID. Must be unique across all OpenVPN instances. The OPNsense web UI auto-assigns one; when managing instances via Terraform you should set this explicitly to keep ordering deterministic. Defaults to `-1` (auto).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of this instance. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"dev_type": schema.StringAttribute{
				MarkdownDescription: "Tunnel device type. One of `tun`, `tap`, or `ovpn` (DCO). Defaults to `tun`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("tun"),
				Validators: []validator.String{
					stringvalidator.OneOf("tun", "tap", "ovpn"),
				},
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Network protocol. One of `udp`, `udp4`, `udp6`, `tcp`, `tcp4`, `tcp6`. Defaults to `udp`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("udp"),
				Validators: []validator.String{
					stringvalidator.OneOf("udp", "udp4", "udp6", "tcp", "tcp4", "tcp6"),
				},
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Listening port (server) or remote port (client). Defaults to `-1` (no fixed port).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"port_share": schema.StringAttribute{
				MarkdownDescription: "Share the OpenVPN port with another service using `host:port`. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"local": schema.StringAttribute{
				MarkdownDescription: "Local IP to bind to. Defaults to `\"\"` (all interfaces).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"remote": schema.SetAttribute{
				MarkdownDescription: "Remote host(s) the client should connect to. Use `host` or `host:port`. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"topology": schema.StringAttribute{
				MarkdownDescription: "Tunnel topology. One of `subnet`, `net30`, or `p2p`. Defaults to `subnet`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("subnet"),
				Validators: []validator.String{
					stringvalidator.OneOf("subnet", "net30", "p2p"),
				},
			},
			"server": schema.StringAttribute{
				MarkdownDescription: "IPv4 tunnel network in CIDR notation (server mode). Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"server_ipv6": schema.StringAttribute{
				MarkdownDescription: "IPv6 tunnel network in CIDR notation (server mode). Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"no_pool": schema.BoolAttribute{
				MarkdownDescription: "When `true`, disables the dynamic IP pool. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"bridge_gateway": schema.StringAttribute{
				MarkdownDescription: "Bridge gateway IP when bridging the OpenVPN tap interface. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"bridge_pool": schema.StringAttribute{
				MarkdownDescription: "Bridge DHCP pool range when bridging the OpenVPN tap interface (`start-end`). Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"route": schema.SetAttribute{
				MarkdownDescription: "Local routes to install (CIDR). Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"push_route": schema.SetAttribute{
				MarkdownDescription: "Routes to push to clients (CIDR). Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"push_excluded_routes": schema.SetAttribute{
				MarkdownDescription: "Routes excluded from being pushed to clients (CIDR). Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"certificate": schema.StringAttribute{
				MarkdownDescription: "ID (refid) of the server/client certificate to use. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"crl": schema.StringAttribute{
				MarkdownDescription: "ID (refid) of the certificate revocation list. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"certificate_authority": schema.StringAttribute{
				MarkdownDescription: "ID (refid) of the certificate authority. Leave blank to derive from the server/client certificate. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"cert_depth": schema.StringAttribute{
				MarkdownDescription: "Maximum certificate chain depth. One of `\"\"` (do not check), `1`, `2`, `3`, `4`, `5`. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf("", "1", "2", "3", "4", "5"),
				},
			},
			"remote_cert_tls": schema.BoolAttribute{
				MarkdownDescription: "Require the remote peer to have a certificate with the correct extended key usage. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"verify_client_cert": schema.StringAttribute{
				MarkdownDescription: "Client certificate verification mode. One of `require` or `none`. Defaults to `require`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("require"),
				Validators: []validator.String{
					stringvalidator.OneOf("require", "none"),
				},
			},
			"use_ocsp": schema.BoolAttribute{
				MarkdownDescription: "Use OCSP to verify the peer certificate. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"auth_digest": schema.StringAttribute{
				MarkdownDescription: "Authentication digest. Set to `\"\"` to use the OpenVPN default. Valid values include `SHA1`, `SHA256`, `SHA512`, etc. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"data_ciphers": schema.SetAttribute{
				MarkdownDescription: "Accepted data-channel ciphers, in priority order (e.g. `[\"AES-256-GCM\", \"CHACHA20-POLY1305\"]`). Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"data_ciphers_fallback": schema.StringAttribute{
				MarkdownDescription: "Fallback cipher for legacy clients that don't support cipher negotiation. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tls_key": schema.StringAttribute{
				MarkdownDescription: "UUID of an `opnsense_openvpn_static_key` to use as TLS key. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"auth_mode": schema.SetAttribute{
				MarkdownDescription: "Authentication backends (e.g. `[\"Local Database\"]`). Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"local_group": schema.StringAttribute{
				MarkdownDescription: "Required local group membership for authenticating clients. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"various_flags": schema.SetAttribute{
				MarkdownDescription: "Miscellaneous instance flags. Any of `block-ipv6`, `client-to-client`, `duplicate-cn`, `float`, `passtos`, `persist-remote-ip`, `remote-random`, `route-noexec`, `route-nopull`, `explicit-exit-notify`, `fast-io`. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"various_push_flags": schema.SetAttribute{
				MarkdownDescription: "Miscellaneous push flags. Any of `block-ipv6`, `block-outside-dns`, `register-dns`, `explicit-exit-notify`. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"push_inactive": schema.Int64Attribute{
				MarkdownDescription: "Disconnect clients after this many seconds of inactivity. Defaults to `-1` (disabled).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"username_as_common_name": schema.BoolAttribute{
				MarkdownDescription: "Use the authenticated username as the common name. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"strict_user_cn": schema.StringAttribute{
				MarkdownDescription: "Whether to enforce that the username matches the certificate common name. One of `0` (no), `1` (yes), or `2` (yes, case-insensitive). Defaults to `0`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("0"),
				Validators: []validator.String{
					stringvalidator.OneOf("0", "1", "2"),
				},
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Client mode: username for the remote server. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Client mode: password for the remote server. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Default:             stringdefault.StaticString(""),
			},
			"max_clients": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of concurrent clients. Defaults to `-1` (unset).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"keepalive_interval": schema.Int64Attribute{
				MarkdownDescription: "Keepalive ping interval in seconds. Defaults to `-1` (unset).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"keepalive_timeout": schema.Int64Attribute{
				MarkdownDescription: "Keepalive timeout in seconds (peer is declared dead). Defaults to `-1` (unset).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"reneg_sec": schema.Int64Attribute{
				MarkdownDescription: "Renegotiate the data channel key after this many seconds. `0` disables renegotiation. Defaults to `-1` (unset).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"auth_gen_token": schema.Int64Attribute{
				MarkdownDescription: "Generate auth tokens valid for this many seconds. Defaults to `-1` (unset).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"auth_gen_token_renewal": schema.Int64Attribute{
				MarkdownDescription: "Renew the auth token every N seconds. Defaults to `-1` (unset).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"auth_gen_token_secret": schema.StringAttribute{
				MarkdownDescription: "Secret used to sign auth tokens. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Default:             stringdefault.StaticString(""),
			},
			"provision_exclusive": schema.BoolAttribute{
				MarkdownDescription: "Only allow the most recently authenticated session for a given common name. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"redirect_gateway": schema.SetAttribute{
				MarkdownDescription: "OpenVPN `redirect-gateway` flags. Any of `local`, `autolocal`, `def1`, `bypass-dhcp`, `bypass-dns`, `block-local`, `ipv6`, `!ipv4`. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"route_metric": schema.Int64Attribute{
				MarkdownDescription: "Metric for installed routes. Defaults to `-1` (unset).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"register_dns": schema.BoolAttribute{
				MarkdownDescription: "Run `register-dns` on Windows clients. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"dns_domain": schema.SetAttribute{
				MarkdownDescription: "DNS domains pushed to clients. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"dns_domain_search": schema.SetAttribute{
				MarkdownDescription: "DNS search domains pushed to clients. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"dns_servers": schema.SetAttribute{
				MarkdownDescription: "DNS servers pushed to clients. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"ntp_servers": schema.SetAttribute{
				MarkdownDescription: "NTP servers pushed to clients. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				Default:             emptySet,
				ElementType:         types.StringType,
			},
			"tun_mtu": schema.Int64Attribute{
				MarkdownDescription: "MTU for the tunnel interface. Defaults to `-1` (unset).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"fragment": schema.Int64Attribute{
				MarkdownDescription: "Enable internal datagram fragmentation. Defaults to `-1` (disabled).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"mss_fix": schema.BoolAttribute{
				MarkdownDescription: "Enable `mssfix` to clamp TCP MSS to the tunnel MTU. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"carp_depend_on": schema.StringAttribute{
				MarkdownDescription: "Only run this instance when the named CARP VIP is master. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"compress_migrate": schema.BoolAttribute{
				MarkdownDescription: "Enable compression migration (legacy clients). Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ifconfig_pool_persist": schema.BoolAttribute{
				MarkdownDescription: "Persist client-IP assignments across restarts. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"http_proxy": schema.StringAttribute{
				MarkdownDescription: "Connect through an HTTP proxy (`host:port`). Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"verify_x509_name": schema.StringAttribute{
				MarkdownDescription: "Verify the peer certificate against this X509 name. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func instanceDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Read an existing OpenVPN instance (server or client) by UUID.",
		Attributes: map[string]dschema.Attribute{
			"id":                      dschema.StringAttribute{MarkdownDescription: "UUID of the resource.", Required: true},
			"enabled":                 dschema.BoolAttribute{MarkdownDescription: "Whether this instance is enabled.", Computed: true},
			"role":                    dschema.StringAttribute{MarkdownDescription: "Role: `server` or `client`.", Computed: true},
			"vpn_id":                  dschema.Int64Attribute{MarkdownDescription: "Numeric VPN ID.", Computed: true},
			"description":             dschema.StringAttribute{MarkdownDescription: "Description.", Computed: true},
			"dev_type":                dschema.StringAttribute{MarkdownDescription: "Tunnel device type.", Computed: true},
			"protocol":                dschema.StringAttribute{MarkdownDescription: "Network protocol.", Computed: true},
			"port":                    dschema.Int64Attribute{MarkdownDescription: "Port.", Computed: true},
			"port_share":              dschema.StringAttribute{MarkdownDescription: "Port-share target.", Computed: true},
			"local":                   dschema.StringAttribute{MarkdownDescription: "Local bind IP.", Computed: true},
			"remote":                  dschema.SetAttribute{MarkdownDescription: "Remote hosts.", Computed: true, ElementType: types.StringType},
			"topology":                dschema.StringAttribute{MarkdownDescription: "Tunnel topology.", Computed: true},
			"server":                  dschema.StringAttribute{MarkdownDescription: "IPv4 tunnel network.", Computed: true},
			"server_ipv6":             dschema.StringAttribute{MarkdownDescription: "IPv6 tunnel network.", Computed: true},
			"no_pool":                 dschema.BoolAttribute{MarkdownDescription: "Dynamic IP pool disabled.", Computed: true},
			"bridge_gateway":          dschema.StringAttribute{MarkdownDescription: "Bridge gateway IP.", Computed: true},
			"bridge_pool":             dschema.StringAttribute{MarkdownDescription: "Bridge DHCP pool.", Computed: true},
			"route":                   dschema.SetAttribute{MarkdownDescription: "Local routes.", Computed: true, ElementType: types.StringType},
			"push_route":              dschema.SetAttribute{MarkdownDescription: "Pushed routes.", Computed: true, ElementType: types.StringType},
			"push_excluded_routes":    dschema.SetAttribute{MarkdownDescription: "Excluded pushed routes.", Computed: true, ElementType: types.StringType},
			"certificate":             dschema.StringAttribute{MarkdownDescription: "Certificate refid.", Computed: true},
			"crl":                     dschema.StringAttribute{MarkdownDescription: "CRL refid.", Computed: true},
			"certificate_authority":   dschema.StringAttribute{MarkdownDescription: "CA refid.", Computed: true},
			"cert_depth":              dschema.StringAttribute{MarkdownDescription: "Certificate chain depth.", Computed: true},
			"remote_cert_tls":         dschema.BoolAttribute{MarkdownDescription: "Remote cert TLS enforcement.", Computed: true},
			"verify_client_cert":      dschema.StringAttribute{MarkdownDescription: "Client cert verification mode.", Computed: true},
			"use_ocsp":                dschema.BoolAttribute{MarkdownDescription: "OCSP enforcement.", Computed: true},
			"auth_digest":             dschema.StringAttribute{MarkdownDescription: "Authentication digest.", Computed: true},
			"data_ciphers":            dschema.SetAttribute{MarkdownDescription: "Data ciphers.", Computed: true, ElementType: types.StringType},
			"data_ciphers_fallback":   dschema.StringAttribute{MarkdownDescription: "Fallback data cipher.", Computed: true},
			"tls_key":                 dschema.StringAttribute{MarkdownDescription: "Static key UUID.", Computed: true},
			"auth_mode":               dschema.SetAttribute{MarkdownDescription: "Authentication backends.", Computed: true, ElementType: types.StringType},
			"local_group":             dschema.StringAttribute{MarkdownDescription: "Required local group.", Computed: true},
			"various_flags":           dschema.SetAttribute{MarkdownDescription: "Misc flags.", Computed: true, ElementType: types.StringType},
			"various_push_flags":      dschema.SetAttribute{MarkdownDescription: "Misc push flags.", Computed: true, ElementType: types.StringType},
			"push_inactive":           dschema.Int64Attribute{MarkdownDescription: "Inactivity disconnect (seconds).", Computed: true},
			"username_as_common_name": dschema.BoolAttribute{MarkdownDescription: "Use auth username as CN.", Computed: true},
			"strict_user_cn":          dschema.StringAttribute{MarkdownDescription: "Strict username/CN check.", Computed: true},
			"username":                dschema.StringAttribute{MarkdownDescription: "Client-mode username.", Computed: true},
			"password":                dschema.StringAttribute{MarkdownDescription: "Client-mode password.", Computed: true, Sensitive: true},
			"max_clients":             dschema.Int64Attribute{MarkdownDescription: "Maximum concurrent clients.", Computed: true},
			"keepalive_interval":      dschema.Int64Attribute{MarkdownDescription: "Keepalive interval.", Computed: true},
			"keepalive_timeout":       dschema.Int64Attribute{MarkdownDescription: "Keepalive timeout.", Computed: true},
			"reneg_sec":               dschema.Int64Attribute{MarkdownDescription: "Renegotiation interval.", Computed: true},
			"auth_gen_token":          dschema.Int64Attribute{MarkdownDescription: "Auth token lifetime.", Computed: true},
			"auth_gen_token_renewal":  dschema.Int64Attribute{MarkdownDescription: "Auth token renewal interval.", Computed: true},
			"auth_gen_token_secret":   dschema.StringAttribute{MarkdownDescription: "Auth token signing secret.", Computed: true, Sensitive: true},
			"provision_exclusive":     dschema.BoolAttribute{MarkdownDescription: "Provision exclusive sessions.", Computed: true},
			"redirect_gateway":        dschema.SetAttribute{MarkdownDescription: "Redirect gateway flags.", Computed: true, ElementType: types.StringType},
			"route_metric":            dschema.Int64Attribute{MarkdownDescription: "Route metric.", Computed: true},
			"register_dns":            dschema.BoolAttribute{MarkdownDescription: "register-dns on Windows.", Computed: true},
			"dns_domain":              dschema.SetAttribute{MarkdownDescription: "DNS domains pushed.", Computed: true, ElementType: types.StringType},
			"dns_domain_search":       dschema.SetAttribute{MarkdownDescription: "DNS search domains pushed.", Computed: true, ElementType: types.StringType},
			"dns_servers":             dschema.SetAttribute{MarkdownDescription: "DNS servers pushed.", Computed: true, ElementType: types.StringType},
			"ntp_servers":             dschema.SetAttribute{MarkdownDescription: "NTP servers pushed.", Computed: true, ElementType: types.StringType},
			"tun_mtu":                 dschema.Int64Attribute{MarkdownDescription: "Tunnel MTU.", Computed: true},
			"fragment":                dschema.Int64Attribute{MarkdownDescription: "Fragment size.", Computed: true},
			"mss_fix":                 dschema.BoolAttribute{MarkdownDescription: "Enable mssfix.", Computed: true},
			"carp_depend_on":          dschema.StringAttribute{MarkdownDescription: "CARP VIP dependency.", Computed: true},
			"compress_migrate":        dschema.BoolAttribute{MarkdownDescription: "Compress migrate.", Computed: true},
			"ifconfig_pool_persist":   dschema.BoolAttribute{MarkdownDescription: "Persist client IP assignments.", Computed: true},
			"http_proxy":              dschema.StringAttribute{MarkdownDescription: "HTTP proxy target.", Computed: true},
			"verify_x509_name":        dschema.StringAttribute{MarkdownDescription: "Verify peer X509 name.", Computed: true},
		},
	}
}

func convertInstanceSchemaToStruct(d *instanceResourceModel) (*openvpn.Instance, error) {
	ctx := context.Background()
	toList := func(s types.Set) []string {
		var out []string
		s.ElementsAs(ctx, &out, false)
		return out
	}

	return &openvpn.Instance{
		Enabled:              tools.BoolToString(d.Enabled.ValueBool()),
		Role:                 api.SelectedMap(d.Role.ValueString()),
		VPNID:                tools.Int64ToStringNegative(d.VPNID.ValueInt64()),
		Description:          d.Description.ValueString(),
		DevType:              api.SelectedMap(d.DevType.ValueString()),
		Protocol:             api.SelectedMap(d.Protocol.ValueString()),
		Port:                 tools.Int64ToStringNegative(d.Port.ValueInt64()),
		PortShare:            d.PortShare.ValueString(),
		Local:                d.Local.ValueString(),
		Remote:               api.SelectedMapList(toList(d.Remote)),
		Topology:             api.SelectedMap(d.Topology.ValueString()),
		Server:               d.Server.ValueString(),
		ServerIPv6:           d.ServerIPv6.ValueString(),
		NoPool:               tools.BoolToString(d.NoPool.ValueBool()),
		BridgeGateway:        d.BridgeGateway.ValueString(),
		BridgePool:           d.BridgePool.ValueString(),
		Route:                api.SelectedMapList(toList(d.Route)),
		PushRoute:            api.SelectedMapList(toList(d.PushRoute)),
		PushExcludedRoutes:   api.SelectedMapList(toList(d.PushExcludedRoutes)),
		Certificate:          api.SelectedMap(d.Certificate.ValueString()),
		CRL:                  api.SelectedMap(d.CRL.ValueString()),
		CertificateAuthority: api.SelectedMap(d.CertificateAuthority.ValueString()),
		CertDepth:            api.SelectedMap(d.CertDepth.ValueString()),
		RemoteCertTLS:        tools.BoolToString(d.RemoteCertTLS.ValueBool()),
		VerifyClientCert:     api.SelectedMap(d.VerifyClientCert.ValueString()),
		UseOCSP:              tools.BoolToString(d.UseOCSP.ValueBool()),
		AuthDigest:           api.SelectedMap(d.AuthDigest.ValueString()),
		DataCiphers:          api.SelectedMapList(toList(d.DataCiphers)),
		DataCiphersFallback:  api.SelectedMap(d.DataCiphersFallback.ValueString()),
		TLSKey:               api.SelectedMap(d.TLSKey.ValueString()),
		AuthMode:             api.SelectedMapList(toList(d.AuthMode)),
		LocalGroup:           api.SelectedMap(d.LocalGroup.ValueString()),
		VariousFlags:         api.SelectedMapList(toList(d.VariousFlags)),
		VariousPushFlags:     api.SelectedMapList(toList(d.VariousPushFlags)),
		PushInactive:         tools.Int64ToStringNegative(d.PushInactive.ValueInt64()),
		UsernameAsCommonName: tools.BoolToString(d.UsernameAsCommonName.ValueBool()),
		StrictUserCN:         api.SelectedMap(d.StrictUserCN.ValueString()),
		Username:             d.Username.ValueString(),
		Password:             d.Password.ValueString(),
		MaxClients:           tools.Int64ToStringNegative(d.MaxClients.ValueInt64()),
		KeepaliveInterval:    tools.Int64ToStringNegative(d.KeepaliveInterval.ValueInt64()),
		KeepaliveTimeout:     tools.Int64ToStringNegative(d.KeepaliveTimeout.ValueInt64()),
		RenegSec:             tools.Int64ToStringNegative(d.RenegSec.ValueInt64()),
		AuthGenToken:         tools.Int64ToStringNegative(d.AuthGenToken.ValueInt64()),
		AuthGenTokenRenewal:  tools.Int64ToStringNegative(d.AuthGenTokenRenewal.ValueInt64()),
		AuthGenTokenSecret:   d.AuthGenTokenSecret.ValueString(),
		ProvisionExclusive:   tools.BoolToString(d.ProvisionExclusive.ValueBool()),
		RedirectGateway:      api.SelectedMapList(toList(d.RedirectGateway)),
		RouteMetric:          tools.Int64ToStringNegative(d.RouteMetric.ValueInt64()),
		RegisterDNS:          tools.BoolToString(d.RegisterDNS.ValueBool()),
		DNSDomain:            api.SelectedMapList(toList(d.DNSDomain)),
		DNSDomainSearch:      api.SelectedMapList(toList(d.DNSDomainSearch)),
		DNSServers:           api.SelectedMapList(toList(d.DNSServers)),
		NTPServers:           api.SelectedMapList(toList(d.NTPServers)),
		TunMTU:               tools.Int64ToStringNegative(d.TunMTU.ValueInt64()),
		Fragment:             tools.Int64ToStringNegative(d.Fragment.ValueInt64()),
		MSSFix:               tools.BoolToString(d.MSSFix.ValueBool()),
		CARPDependOn:         api.SelectedMap(d.CARPDependOn.ValueString()),
		CompressMigrate:      tools.BoolToString(d.CompressMigrate.ValueBool()),
		IfConfigPoolPersist:  tools.BoolToString(d.IfConfigPoolPersist.ValueBool()),
		HTTPProxy:            d.HTTPProxy.ValueString(),
		VerifyX509Name:       d.VerifyX509Name.ValueString(),
	}, nil
}

func convertInstanceStructToSchema(d *openvpn.Instance) (*instanceResourceModel, error) {
	return &instanceResourceModel{
		Enabled:              types.BoolValue(tools.StringToBool(d.Enabled)),
		Role:                 types.StringValue(d.Role.String()),
		VPNID:                types.Int64Value(tools.StringToInt64(d.VPNID)),
		Description:          types.StringValue(d.Description),
		DevType:              types.StringValue(d.DevType.String()),
		Protocol:             types.StringValue(d.Protocol.String()),
		Port:                 types.Int64Value(tools.StringToInt64(d.Port)),
		PortShare:            types.StringValue(d.PortShare),
		Local:                types.StringValue(d.Local),
		Remote:               tools.StringSliceToSet(d.Remote),
		Topology:             types.StringValue(d.Topology.String()),
		Server:               types.StringValue(d.Server),
		ServerIPv6:           types.StringValue(d.ServerIPv6),
		NoPool:               types.BoolValue(tools.StringToBool(d.NoPool)),
		BridgeGateway:        types.StringValue(d.BridgeGateway),
		BridgePool:           types.StringValue(d.BridgePool),
		Route:                tools.StringSliceToSet(d.Route),
		PushRoute:            tools.StringSliceToSet(d.PushRoute),
		PushExcludedRoutes:   tools.StringSliceToSet(d.PushExcludedRoutes),
		Certificate:          types.StringValue(d.Certificate.String()),
		CRL:                  types.StringValue(d.CRL.String()),
		CertificateAuthority: types.StringValue(d.CertificateAuthority.String()),
		CertDepth:            types.StringValue(d.CertDepth.String()),
		RemoteCertTLS:        types.BoolValue(tools.StringToBool(d.RemoteCertTLS)),
		VerifyClientCert:     types.StringValue(d.VerifyClientCert.String()),
		UseOCSP:              types.BoolValue(tools.StringToBool(d.UseOCSP)),
		AuthDigest:           types.StringValue(d.AuthDigest.String()),
		DataCiphers:          tools.StringSliceToSet(d.DataCiphers),
		DataCiphersFallback:  types.StringValue(d.DataCiphersFallback.String()),
		TLSKey:               types.StringValue(d.TLSKey.String()),
		AuthMode:             tools.StringSliceToSet(d.AuthMode),
		LocalGroup:           types.StringValue(d.LocalGroup.String()),
		VariousFlags:         tools.StringSliceToSet(d.VariousFlags),
		VariousPushFlags:     tools.StringSliceToSet(d.VariousPushFlags),
		PushInactive:         types.Int64Value(tools.StringToInt64(d.PushInactive)),
		UsernameAsCommonName: types.BoolValue(tools.StringToBool(d.UsernameAsCommonName)),
		StrictUserCN:         types.StringValue(d.StrictUserCN.String()),
		Username:             types.StringValue(d.Username),
		Password:             types.StringValue(d.Password),
		MaxClients:           types.Int64Value(tools.StringToInt64(d.MaxClients)),
		KeepaliveInterval:    types.Int64Value(tools.StringToInt64(d.KeepaliveInterval)),
		KeepaliveTimeout:     types.Int64Value(tools.StringToInt64(d.KeepaliveTimeout)),
		RenegSec:             types.Int64Value(tools.StringToInt64(d.RenegSec)),
		AuthGenToken:         types.Int64Value(tools.StringToInt64(d.AuthGenToken)),
		AuthGenTokenRenewal:  types.Int64Value(tools.StringToInt64(d.AuthGenTokenRenewal)),
		AuthGenTokenSecret:   types.StringValue(d.AuthGenTokenSecret),
		ProvisionExclusive:   types.BoolValue(tools.StringToBool(d.ProvisionExclusive)),
		RedirectGateway:      tools.StringSliceToSet(d.RedirectGateway),
		RouteMetric:          types.Int64Value(tools.StringToInt64(d.RouteMetric)),
		RegisterDNS:          types.BoolValue(tools.StringToBool(d.RegisterDNS)),
		DNSDomain:            tools.StringSliceToSet(d.DNSDomain),
		DNSDomainSearch:      tools.StringSliceToSet(d.DNSDomainSearch),
		DNSServers:           tools.StringSliceToSet(d.DNSServers),
		NTPServers:           tools.StringSliceToSet(d.NTPServers),
		TunMTU:               types.Int64Value(tools.StringToInt64(d.TunMTU)),
		Fragment:             types.Int64Value(tools.StringToInt64(d.Fragment)),
		MSSFix:               types.BoolValue(tools.StringToBool(d.MSSFix)),
		CARPDependOn:         types.StringValue(d.CARPDependOn.String()),
		CompressMigrate:      types.BoolValue(tools.StringToBool(d.CompressMigrate)),
		IfConfigPoolPersist:  types.BoolValue(tools.StringToBool(d.IfConfigPoolPersist)),
		HTTPProxy:            types.StringValue(d.HTTPProxy),
		VerifyX509Name:       types.StringValue(d.VerifyX509Name),
	}, nil
}
