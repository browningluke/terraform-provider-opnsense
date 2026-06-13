package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaOSPF6NetworkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaOSPF6NetworkResourceConfig("2001:db8::", "32", "0.0.0.0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_network.test", "ip_addr", "2001:db8::"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_network.test", "net_mask", "32"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_network.test", "area", "0.0.0.0"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_network.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_ospf6_network.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_ospf6_network.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaOSPF6NetworkResourceConfig("2001:db8::", "48", "0.0.0.0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_network.test", "net_mask", "48"),
				),
			},
		},
	})
}

func testAccQuaggaOSPF6NetworkResourceConfig(ipAddr, netMask, area string) string {
	return `
resource "opnsense_quagga_ospf6_network" "test" {
  ip_addr  = "` + ipAddr + `"
  net_mask = "` + netMask + `"
  area     = "` + area + `"
}
`
}
