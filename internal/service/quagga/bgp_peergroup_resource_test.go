package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaBGPPeerGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaBGPPeerGroupResourceConfig("test-peergroup", "65002", "ipv4"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_peer_group.test", "name", "test-peergroup"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_peer_group.test", "remote_as", "65002"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_peer_group.test", "family", "ipv4"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_peer_group.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_bgp_peer_group.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_bgp_peer_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaBGPPeerGroupResourceConfig("test-peergroup", "65003", "ipv6"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_peer_group.test", "remote_as", "65003"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_peer_group.test", "family", "ipv6"),
				),
			},
		},
	})
}

func testAccQuaggaBGPPeerGroupResourceConfig(name, remoteAS, family string) string {
	return `
resource "opnsense_quagga_bgp_peer_group" "test" {
  name      = "` + name + `"
  remote_as = "` + remoteAS + `"
  family    = "` + family + `"
}
`
}
