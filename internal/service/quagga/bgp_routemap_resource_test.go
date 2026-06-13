package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaBGPRouteMapResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaBGPRouteMapResourceConfig("TEST_RM", "permit", "10", "local-preference 200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_routemap.test", "name", "TEST_RM"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_routemap.test", "action", "permit"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_routemap.test", "route_map_id", "10"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_routemap.test", "set", "local-preference 200"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_bgp_routemap.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_bgp_routemap.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaBGPRouteMapResourceConfig("TEST_RM", "deny", "20", "local-preference 300"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_routemap.test", "action", "deny"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_routemap.test", "set", "local-preference 300"),
				),
			},
		},
	})
}

func testAccQuaggaBGPRouteMapResourceConfig(name, action, id, set string) string {
	return `
resource "opnsense_quagga_bgp_routemap" "test" {
  name         = "` + name + `"
  action       = "` + action + `"
  route_map_id = "` + id + `"
  set          = "` + set + `"
}
`
}
