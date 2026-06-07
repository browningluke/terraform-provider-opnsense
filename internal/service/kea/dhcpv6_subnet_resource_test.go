package kea_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKeaDhcpv6SubnetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDhcpv6SubnetResourceConfig(
					"fd00:200::/64",
					"fd00:200::100-fd00:200::200",
					"wan",
					"Test Kea DHCPv6 Subnet",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_subnet.test", "subnet", "fd00:200::/64"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_subnet.test", "pools.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_kea_dhcpv6_subnet.test", "pools.*", "fd00:200::100-fd00:200::200"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_subnet.test", "interface", "wan"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_subnet.test", "description", "Test Kea DHCPv6 Subnet"),
					resource.TestCheckResourceAttrSet("opnsense_kea_dhcpv6_subnet.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_kea_dhcpv6_subnet.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDhcpv6SubnetResourceConfig(
					"fd00:200::/64",
					"fd00:200::100-fd00:200::150",
					"wan",
					"Test Kea DHCPv6 Subnet Updated",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_subnet.test", "pools.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_kea_dhcpv6_subnet.test", "pools.*", "fd00:200::100-fd00:200::150"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_subnet.test", "description", "Test Kea DHCPv6 Subnet Updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDhcpv6SubnetResourceConfig(subnet, pool, iface, description string) string {
	return fmt.Sprintf(`
resource "opnsense_kea_dhcpv6_subnet" "test" {
  subnet      = %[1]q
  pools       = [%[2]q]
  interface   = %[3]q
  description = %[4]q
}
`, subnet, pool, iface, description)
}
