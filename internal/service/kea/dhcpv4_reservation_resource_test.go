package kea_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKeaDhcpv4ReservationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDhcpv4ReservationResourceConfig(
					"192.168.201.0/24",
					"192.168.201.150",
					"aa:bb:cc:dd:ee:ff",
					"test-host",
					"Test Kea DHCPv4 Reservation",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("opnsense_kea_dhcpv4_reservation.test", "subnet_id", "opnsense_kea_dhcpv4_subnet.test", "id"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_reservation.test", "ip_address", "192.168.201.150"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_reservation.test", "mac_address", "aa:bb:cc:dd:ee:ff"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_reservation.test", "hostname", "test-host"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_reservation.test", "description", "Test Kea DHCPv4 Reservation"),
					resource.TestCheckResourceAttrSet("opnsense_kea_dhcpv4_reservation.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_kea_dhcpv4_reservation.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDhcpv4ReservationResourceConfig(
					"192.168.201.0/24",
					"192.168.201.150",
					"aa:bb:cc:dd:ee:ff",
					"test-host-upd",
					"Test Kea DHCPv4 Reservation Updated",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_reservation.test", "hostname", "test-host-upd"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_reservation.test", "description", "Test Kea DHCPv4 Reservation Updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDhcpv4ReservationResourceConfig(subnetCIDR, ipAddress, macAddress, hostname, description string) string {
	return fmt.Sprintf(`
resource "opnsense_kea_dhcpv4_subnet" "test" {
  subnet      = %[1]q
  description = "Test subnet for DHCPv4 reservation"
}

resource "opnsense_kea_dhcpv4_reservation" "test" {
  subnet_id   = opnsense_kea_dhcpv4_subnet.test.id
  ip_address  = %[2]q
  mac_address = %[3]q
  hostname    = %[4]q
  description = %[5]q
}
`, subnetCIDR, ipAddress, macAddress, hostname, description)
}
