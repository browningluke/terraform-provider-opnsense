package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaOSPFRouteMapResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaOSPFRouteMapResourceConfig("TEST_OSPF_RM", "permit", "10", "metric 100"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_route_map.test", "name", "TEST_OSPF_RM"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_route_map.test", "action", "permit"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_route_map.test", "route_map_id", "10"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_route_map.test", "set", "metric 100"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_route_map.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_ospf_route_map.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_ospf_route_map.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaOSPFRouteMapResourceConfig("TEST_OSPF_RM", "deny", "20", "metric 200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_route_map.test", "action", "deny"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_route_map.test", "set", "metric 200"),
				),
			},
		},
	})
}

func testAccQuaggaOSPFRouteMapResourceConfig(name, action, id, set string) string {
	return `
resource "opnsense_quagga_ospf_route_map" "test" {
  name         = "` + name + `"
  action       = "` + action + `"
  route_map_id = "` + id + `"
  set          = "` + set + `"
}
`
}
