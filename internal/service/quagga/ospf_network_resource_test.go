package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaOSPFNetworkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaOSPFNetworkResourceConfig("10.0.0.0", "0.0.0.0", "24"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_network.test", "ip_addr", "10.0.0.0"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_network.test", "area", "0.0.0.0"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_network.test", "net_mask", "24"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_network.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_ospf_network.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_ospf_network.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaOSPFNetworkResourceConfig("10.0.0.0", "0.0.0.0", "16"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_network.test", "net_mask", "16"),
				),
			},
		},
	})
}

func testAccQuaggaOSPFNetworkResourceConfig(ipAddr, area, netMask string) string {
	return `
resource "opnsense_quagga_ospf_network" "test" {
  ip_addr  = "` + ipAddr + `"
  area     = "` + area + `"
  net_mask = "` + netMask + `"
}
`
}
