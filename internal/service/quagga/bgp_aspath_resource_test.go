package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaBGPASPathResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaBGPASPathResourceConfig("10", "permit", "^65100$"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_aspath.test", "number", "10"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_aspath.test", "action", "permit"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_aspath.test", "as", "^65100$"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_aspath.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_bgp_aspath.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_bgp_aspath.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaBGPASPathResourceConfig("10", "deny", "^65200$"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_aspath.test", "action", "deny"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_aspath.test", "as", "^65200$"),
				),
			},
		},
	})
}

func testAccQuaggaBGPASPathResourceConfig(number, action, as string) string {
	return `
resource "opnsense_quagga_bgp_aspath" "test" {
  number = ` + number + `
  action = "` + action + `"
  as     = "` + as + `"
}
`
}
