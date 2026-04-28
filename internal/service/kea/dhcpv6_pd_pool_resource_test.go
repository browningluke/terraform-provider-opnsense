package kea_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKeaDhcpv6PdPoolResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDhcpv6PdPoolResourceConfig(
					"fd00:202::/64",
					"fd00:203::",
					"64",
					"80",
					"Test Kea DHCPv6 PD Pool",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("opnsense_kea_dhcpv6_pd_pool.test", "subnet_id", "opnsense_kea_dhcpv6_subnet.test", "id"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_pd_pool.test", "prefix", "fd00:203::"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_pd_pool.test", "prefix_len", "64"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_pd_pool.test", "delegated_len", "80"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_pd_pool.test", "description", "Test Kea DHCPv6 PD Pool"),
					resource.TestCheckResourceAttrSet("opnsense_kea_dhcpv6_pd_pool.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_kea_dhcpv6_pd_pool.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDhcpv6PdPoolResourceConfig(
					"fd00:202::/64",
					"fd00:203::",
					"64",
					"96",
					"Test Kea DHCPv6 PD Pool Updated",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_pd_pool.test", "delegated_len", "96"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_pd_pool.test", "description", "Test Kea DHCPv6 PD Pool Updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDhcpv6PdPoolResourceConfig(subnetCIDR, prefix, prefixLen, delegatedLen, description string) string {
	return fmt.Sprintf(`
resource "opnsense_kea_dhcpv6_subnet" "test" {
  subnet      = %[1]q
  interface   = "wan"
  description = "Test subnet for DHCPv6 PD pool"
}

resource "opnsense_kea_dhcpv6_pd_pool" "test" {
  subnet_id     = opnsense_kea_dhcpv6_subnet.test.id
  prefix        = %[2]q
  prefix_len    = %[3]q
  delegated_len = %[4]q
  description   = %[5]q
}
`, subnetCIDR, prefix, prefixLen, delegatedLen, description)
}
