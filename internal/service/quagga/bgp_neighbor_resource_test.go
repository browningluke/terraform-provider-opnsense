package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaBGPNeighborResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaBGPNeighborResourceConfig("192.0.2.2", "65100"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_neighbor.test", "peer_ip", "192.0.2.2"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_neighbor.test", "remote_as", "65100"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_neighbor.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_bgp_neighbor.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_bgp_neighbor.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaBGPNeighborResourceConfig("192.0.2.3", "65200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_neighbor.test", "peer_ip", "192.0.2.3"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_neighbor.test", "remote_as", "65200"),
				),
			},
		},
	})
}

func testAccQuaggaBGPNeighborResourceConfig(peerIP, remoteAS string) string {
	return `
resource "opnsense_quagga_bgp_neighbor" "test" {
  peer_ip   = "` + peerIP + `"
  remote_as = ` + remoteAS + `
}
`
}
