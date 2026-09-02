package dnsmasq_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingsDataSource(t *testing.T) {
	ctx := context.Background()

	acctest.AccPreCheck(t)
	client := acctest.Client(t)

	settings, err := client.Dnsmasq().GeneralSettingsGet(ctx)
	if err != nil {
		t.Fatalf("failed to read dnsmasq settings: %v", err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccSettingsDataSourceConfig(),

				Check: resource.ComposeAggregateTestCheckFunc(
					// General
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_settings.test", "id", "dnsmasq_settings"),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"enabled",
						boolString(settings.Dnsmasq.IsEnabled),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"strict_interface_binding",
						boolString(settings.Dnsmasq.StrictInterfaceBinding),
					),

					// Interface: SelectedMapList -> Terraform Set
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"interface.#",
						fmt.Sprintf("%d", len(settings.Dnsmasq.Interface)),
					),

					// DNS
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_settings.test", "dns.port", settings.Dnsmasq.DNS_Port),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dns.dnssec",
						boolString(settings.Dnsmasq.DNS_DNSSEC),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dns.no_host_lookup",
						boolString(settings.Dnsmasq.DNS_NoHostLookup),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dns.log_queries",
						boolString(settings.Dnsmasq.DNS_LogDnsQueries),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dns.no_ident",
						boolString(settings.Dnsmasq.DNS_NoIdent),
					),

					// DNS query forwarding
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dns_query_forwarding.query_sequentially",
						boolString(settings.Dnsmasq.DNS_QF_QuerySequentially),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dns_query_forwarding.require_domain",
						boolString(settings.Dnsmasq.DNS_QF_RequireDomain),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dns_query_forwarding.do_not_forward_system_dns",
						boolString(settings.Dnsmasq.DNS_QF_DoNotForwardSystemDNS),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dns_query_forwarding.do_not_forward_private_reverse",
						boolString(settings.Dnsmasq.DNS_QF_DoNotForwardPrivateReverseLookup),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dns_query_forwarding.add_mac",
						settings.Dnsmasq.DNS_QF_AddMac.String(),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dns_query_forwarding.add_subnet",
						boolString(settings.Dnsmasq.DNS_QF_AddSubnet),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dns_query_forwarding.strip_subnet",
						boolString(settings.Dnsmasq.DNS_QF_StripSubnet),
					),

					// DHCP
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dhcp.fqdn",
						boolString(settings.Dnsmasq.DHCPSettings.FQDN),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dhcp.default_domain",
						settings.Dnsmasq.DHCPSettings.DefaultDomain,
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dhcp.local_domain",
						boolString(settings.Dnsmasq.DHCPSettings.LocalDomain),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dhcp.authoritative",
						boolString(settings.Dnsmasq.DHCPSettings.Authoritative),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dhcp.register_firewall_rules",
						boolString(settings.Dnsmasq.DHCPSettings.RegisterFirewallRules),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dhcp.router_advertisements",
						boolString(settings.Dnsmasq.DHCPSettings.RouterAdvertisements),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dhcp.host_ping",
						boolString(settings.Dnsmasq.DHCPSettings.HostPing),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dhcp.disable_ha_sync",
						boolString(settings.Dnsmasq.DHCPSettings.DisableHASync),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dhcp.log_dhcp",
						boolString(settings.Dnsmasq.DHCPSettings.LogDhcp),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dhcp.log_quiet",
						boolString(settings.Dnsmasq.DHCPSettings.LogQuiet),
					),

					// DHCP interfaces: SelectedMapList -> Terraform Set
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"dhcp.no_dhcp_interfaces.#",
						fmt.Sprintf(
							"%d",
							len(settings.Dnsmasq.DHCPSettings.InterfaceNoDhcp),
						),
					),

					// Legacy
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"legacy.register_isc_dhcp4_leases",
						boolString(settings.Dnsmasq.Legacy_RegisterISCDhcp4Leases),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"legacy.dhcp_domain_override",
						settings.Dnsmasq.Legacy_DhcpDomainOverride,
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"legacy.register_dhcp_static_mappings",
						boolString(settings.Dnsmasq.Legacy_RegisterDhcpStaticMappings),
					),
					resource.TestCheckResourceAttr(
						"data.opnsense_dnsmasq_settings.test",
						"legacy.prefer_dhcp",
						boolString(settings.Dnsmasq.Legacy_PreferDhcp),
					),
				),
			},
		},
	})

}

func testAccSettingsDataSourceConfig() string {
	return `data "opnsense_dnsmasq_settings" "test" {
  id = "dnsmasq_settings"
}`
}

func boolString(value string) string {
	switch value {
	case "1", "true":
		return "true"
	case "0", "false":
		return "false"
	case "":
		// The API omits unset checkboxes; the provider exposes those as false.
		return "false"
	default:
		return value
	}
}
