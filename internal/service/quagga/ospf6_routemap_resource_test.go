package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaOSPF6RouteMapResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaOSPF6RouteMapResourceConfig("TEST_OSPF6_RM", "permit", "10", "metric 100"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_route_map.test", "name", "TEST_OSPF6_RM"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_route_map.test", "action", "permit"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_route_map.test", "route_map_id", "10"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_route_map.test", "set", "metric 100"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_route_map.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_ospf6_route_map.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_ospf6_route_map.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaOSPF6RouteMapResourceConfig("TEST_OSPF6_RM", "deny", "20", "metric 200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_route_map.test", "action", "deny"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_route_map.test", "set", "metric 200"),
				),
			},
		},
	})
}

func testAccQuaggaOSPF6RouteMapResourceConfig(name, action, id, set string) string {
	return `
resource "opnsense_quagga_ospf6_route_map" "test" {
  name         = "` + name + `"
  action       = "` + action + `"
  route_map_id = "` + id + `"
  set          = "` + set + `"
}
`
}
