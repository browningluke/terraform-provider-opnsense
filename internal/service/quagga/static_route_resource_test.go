package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaStaticRouteResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaStaticRouteResourceConfig("10.100.0.0/24", "192.168.1.1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_static_route.test", "network", "10.100.0.0/24"),
					resource.TestCheckResourceAttr("opnsense_quagga_static_route.test", "gateway", "192.168.1.1"),
					resource.TestCheckResourceAttr("opnsense_quagga_static_route.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_static_route.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_static_route.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaStaticRouteResourceConfig("10.100.0.0/24", "192.168.1.254"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_static_route.test", "gateway", "192.168.1.254"),
				),
			},
		},
	})
}

func testAccQuaggaStaticRouteResourceConfig(network, gateway string) string {
	return `
resource "opnsense_quagga_static_route" "test" {
  network = "` + network + `"
  gateway = "` + gateway + `"
}
`
}
