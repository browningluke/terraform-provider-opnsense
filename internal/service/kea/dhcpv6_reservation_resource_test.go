package kea_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKeaDhcpv6ReservationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDhcpv6ReservationResourceConfig(
					"fd00:201::/64",
					"fd00:201::150",
					"00:03:00:01:aa:bb:cc:dd:ee:ff",
					"test-host-v6",
					"Test Kea DHCPv6 Reservation",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("opnsense_kea_dhcpv6_reservation.test", "subnet_id", "opnsense_kea_dhcpv6_subnet.test", "id"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_reservation.test", "ip_address", "fd00:201::150"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_reservation.test", "duid", "00:03:00:01:aa:bb:cc:dd:ee:ff"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_reservation.test", "hostname", "test-host-v6"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_reservation.test", "description", "Test Kea DHCPv6 Reservation"),
					resource.TestCheckResourceAttrSet("opnsense_kea_dhcpv6_reservation.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_kea_dhcpv6_reservation.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDhcpv6ReservationResourceConfig(
					"fd00:201::/64",
					"fd00:201::150",
					"00:03:00:01:aa:bb:cc:dd:ee:ff",
					"test-host-v6-upd",
					"Test Kea DHCPv6 Reservation Updated",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_reservation.test", "hostname", "test-host-v6-upd"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_reservation.test", "description", "Test Kea DHCPv6 Reservation Updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDhcpv6ReservationResourceConfig(subnetCIDR, ipAddress, duid, hostname, description string) string {
	return fmt.Sprintf(`
resource "opnsense_kea_dhcpv6_subnet" "test" {
  subnet      = %[1]q
  interface   = "wan"
  description = "Test subnet for DHCPv6 reservation"
}

resource "opnsense_kea_dhcpv6_reservation" "test" {
  subnet_id   = opnsense_kea_dhcpv6_subnet.test.id
  ip_address  = %[2]q
  duid        = %[3]q
  hostname    = %[4]q
  description = %[5]q
}
`, subnetCIDR, ipAddress, duid, hostname, description)
}
