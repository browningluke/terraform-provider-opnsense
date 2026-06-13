package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaBGPRedistributionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaBGPRedistributionResourceConfig("connected", "test-bgp-redistribution"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_redistribution.test", "redistribute", "connected"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_redistribution.test", "description", "test-bgp-redistribution"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_redistribution.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_bgp_redistribution.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_bgp_redistribution.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaBGPRedistributionResourceConfig("static", "test-bgp-redistribution-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_redistribution.test", "redistribute", "static"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_redistribution.test", "description", "test-bgp-redistribution-updated"),
				),
			},
		},
	})
}

func testAccQuaggaBGPRedistributionResourceConfig(redistribute, description string) string {
	return `
resource "opnsense_quagga_bgp_redistribution" "test" {
  redistribute = "` + redistribute + `"
  description  = "` + description + `"
}
`
}
