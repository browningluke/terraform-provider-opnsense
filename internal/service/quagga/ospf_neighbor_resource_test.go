package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaOSPFNeighborResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaOSPFNeighborResourceConfig("192.0.2.1", "test-ospf-neighbor"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_neighbor.test", "address", "192.0.2.1"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_neighbor.test", "description", "test-ospf-neighbor"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_neighbor.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_ospf_neighbor.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_ospf_neighbor.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaOSPFNeighborResourceConfig("192.0.2.2", "test-ospf-neighbor-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_neighbor.test", "address", "192.0.2.2"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_neighbor.test", "description", "test-ospf-neighbor-updated"),
				),
			},
		},
	})
}

func testAccQuaggaOSPFNeighborResourceConfig(address, description string) string {
	return `
resource "opnsense_quagga_ospf_neighbor" "test" {
  address     = "` + address + `"
  description = "` + description + `"
}
`
}
