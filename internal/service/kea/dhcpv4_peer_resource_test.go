package kea_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKeaDhcpv4PeerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDhcpv4PeerResourceConfig("test-peer-v4", "http://192.168.1.2:647/", "primary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_peer.test", "name", "test-peer-v4"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_peer.test", "url", "http://192.168.1.2:647/"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_peer.test", "role", "primary"),
					resource.TestCheckResourceAttrSet("opnsense_kea_dhcpv4_peer.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_kea_dhcpv4_peer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDhcpv4PeerResourceConfig("test-peer-v4-upd", "http://192.168.1.3:647/", "standby"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_peer.test", "name", "test-peer-v4-upd"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_peer.test", "url", "http://192.168.1.3:647/"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv4_peer.test", "role", "standby"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDhcpv4PeerResourceConfig(name, url, role string) string {
	return fmt.Sprintf(`
resource "opnsense_kea_dhcpv4_peer" "test" {
  name = %[1]q
  url  = %[2]q
  role = %[3]q
}
`, name, url, role)
}
