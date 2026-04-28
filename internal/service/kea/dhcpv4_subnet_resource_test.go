package kea_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKeaDhcpv4SubnetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDhcpv4SubnetResourceConfig(
					"192.168.200.0/24",
					"192.168.200.100 - 192.168.200.200",
					false,
					"192.168.200.1",
					"8.8.8.8",
					"Test Kea DHCPv4 Subnet",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_subnet.test", "subnet", "192.168.200.0/24"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_subnet.test", "pools.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_kea_dhcpv4_subnet.test", "pools.*", "192.168.200.100 - 192.168.200.200"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_subnet.test", "match_client_id", "true"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_subnet.test", "auto_collect", "false"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_subnet.test", "routers.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_kea_dhcpv4_subnet.test", "routers.*", "192.168.200.1"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_subnet.test", "dns_servers.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_kea_dhcpv4_subnet.test", "dns_servers.*", "8.8.8.8"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_subnet.test", "description", "Test Kea DHCPv4 Subnet"),
					resource.TestCheckResourceAttrSet("opnsense_kea_dhcpv4_subnet.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_kea_dhcpv4_subnet.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDhcpv4SubnetResourceConfig(
					"192.168.200.0/24",
					"192.168.200.100 - 192.168.200.150",
					false,
					"192.168.200.1",
					"8.8.8.8",
					"Test Kea DHCPv4 Subnet Updated",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_subnet.test", "pools.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_kea_dhcpv4_subnet.test", "pools.*", "192.168.200.100 - 192.168.200.150"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_subnet.test", "description", "Test Kea DHCPv4 Subnet Updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDhcpv4SubnetResourceConfig(subnet, pool string, autoCollect bool, router, dns, description string) string {
	return fmt.Sprintf(`
resource "opnsense_kea_dhcpv4_subnet" "test" {
  subnet       = %[1]q
  pools        = [%[2]q]
  auto_collect = %[3]t
  routers      = [%[4]q]
  dns_servers  = [%[5]q]
  description  = %[6]q
}
`, subnet, pool, autoCollect, router, dns, description)
}
