package dnsmasq_test

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDnsmasqSettingsResource tests the singleton dnsmasq settings resource.
//
// Because this resource blocks creation (terraform import must be used instead),
// the test begins with an import step rather than an apply step.
func TestAccDnsmasqSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// This singleton already exists on OPNsense. Its resource-level
			// acceptance contract is therefore import plus read; changing global
			// DNS/DHCP settings in a test would be destructive and environment
			// dependent. The schema round-trip unit tests cover the mapping used by
			// Update.
			{
				Config:             testAccSettingsResourceConfig(false, 0, 0),
				ResourceName:       "opnsense_dnsmasq_settings.settings",
				ImportState:        true,
				ImportStateId:      "dnsmasq_settings",
				ImportStatePersist: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_settings.settings", "id", "dnsmasq_settings"),
				),
			},
		},
	})
}

// TestAccDnsmasqSettingsResource_CreateBlocked verifies that attempting to create
// this singleton resource without importing it first returns a clear error.
func TestAccDnsmasqSettingsResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSettingsResourceConfigMinimal(),
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

// testAccSettingsResourceConfig returns a resource config used for the
// import, baseline, update and restore steps.
func testAccSettingsResourceConfig(enabled bool, cacheSize int, maxConcurrentQueries int) string {
	return `
resource "opnsense_dnsmasq_settings" "settings" {
  enabled = ` + boolStr(enabled) + `

  interface = []
  strict_interface_binding = false

  dns = {
    port                   = 53
    dnssec                 = false
    no_host_lookup         = false
    log_queries            = false
    max_concurrent_queries = ` + intStr(maxConcurrentQueries) + `
    cache_size             = ` + intStr(cacheSize) + `
    local_entry_ttl        = 0
    no_ident               = true
  }

  dns_query_forwarding = {
    query_sequentially             = false
    require_domain                 = false
    do_not_forward_system_dns      = false
    do_not_forward_private_reverse = false
    add_mac                        = ""
    add_subnet                     = false
    strip_subnet                   = false
  }

  dhcp = {
    no_dhcp_interfaces      = []
    fqdn                    = true
    default_domain          = ""
    local_domain            = false
    max_leases              = ` + intStr(cacheSize) + `
    authoritative           = false
    reply_delay             = ` + intStr(maxConcurrentQueries) + `
    register_firewall_rules = true
    router_advertisements   = false
    host_ping               = false
    disable_ha_sync         = false
    log_dhcp                = false
    log_quiet               = false
  }

  legacy = {
    register_isc_dhcp4_leases     = false
    dhcp_domain_override          = ""
    register_dhcp_static_mappings = false
    prefer_dhcp                   = false
  }
}
`
}

// testAccSettingsResourceConfigMinimal is used only by
// TestAccDnsmasqSettingsResource_CreateBlocked to verify that attempting
// to create (rather than import) the singleton fails.
func testAccSettingsResourceConfigMinimal() string {
	// Populate nested blocks so Terraform can decode the plan before Create
	// returns the singleton-specific error.
	return testAccSettingsResourceConfig(false, 0, 0)
}

func boolStr(b bool) string {
	if b {
		return "true"
	}

	return "false"
}

func intStr(i int) string {
	return strconv.Itoa(i)
}
