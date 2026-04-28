package kea_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKeaDhcpv6PeerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDhcpv6PeerResourceConfig("test-peer-v6", "http://[2001:db8::2]:647/", "primary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_peer.test", "name", "test-peer-v6"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_peer.test", "url", "http://[2001:db8::2]:647/"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_peer.test", "role", "primary"),
					resource.TestCheckResourceAttrSet("opnsense_kea_dhcpv6_peer.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_kea_dhcpv6_peer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDhcpv6PeerResourceConfig("test-peer-v6-upd", "http://[2001:db8::3]:647/", "standby"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_peer.test", "name", "test-peer-v6-upd"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_peer.test", "url", "http://[2001:db8::3]:647/"),
					resource.TestCheckResourceAttr("opnsense_kea_dhcpv6_peer.test", "role", "standby"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDhcpv6PeerResourceConfig(name, url, role string) string {
	return fmt.Sprintf(`
resource "opnsense_kea_dhcpv6_peer" "test" {
  name = %[1]q
  url  = %[2]q
  role = %[3]q
}
`, name, url, role)
}
