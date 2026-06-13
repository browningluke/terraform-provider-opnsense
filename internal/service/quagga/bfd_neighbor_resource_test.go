package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaBFDNeighborResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaBFDNeighborResourceConfig("192.0.2.1", "test-bfd-neighbor"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bfd_neighbor.test", "address", "192.0.2.1"),
					resource.TestCheckResourceAttr("opnsense_quagga_bfd_neighbor.test", "description", "test-bfd-neighbor"),
					resource.TestCheckResourceAttr("opnsense_quagga_bfd_neighbor.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_quagga_bfd_neighbor.test", "detect_multiplier", "3"),
					resource.TestCheckResourceAttr("opnsense_quagga_bfd_neighbor.test", "receive_interval", "300"),
					resource.TestCheckResourceAttr("opnsense_quagga_bfd_neighbor.test", "transmit_interval", "300"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_bfd_neighbor.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_bfd_neighbor.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaBFDNeighborResourceConfig("192.0.2.2", "test-bfd-neighbor-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bfd_neighbor.test", "address", "192.0.2.2"),
					resource.TestCheckResourceAttr("opnsense_quagga_bfd_neighbor.test", "description", "test-bfd-neighbor-updated"),
				),
			},
		},
	})
}

func testAccQuaggaBFDNeighborResourceConfig(address, description string) string {
	return `
resource "opnsense_quagga_bfd_neighbor" "test" {
  address     = "` + address + `"
  description = "` + description + `"
}
`
}
