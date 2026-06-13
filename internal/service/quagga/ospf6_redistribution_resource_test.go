package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaOSPF6RedistributionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaOSPF6RedistributionResourceConfig("connected", "test-ospf6-redistribution"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_redistribution.test", "redistribute", "connected"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_redistribution.test", "description", "test-ospf6-redistribution"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_redistribution.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_ospf6_redistribution.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_ospf6_redistribution.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaOSPF6RedistributionResourceConfig("static", "test-ospf6-redistribution-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_redistribution.test", "redistribute", "static"),
				),
			},
		},
	})
}

func testAccQuaggaOSPF6RedistributionResourceConfig(redistribute, description string) string {
	return `
resource "opnsense_quagga_ospf6_redistribution" "test" {
  redistribute = "` + redistribute + `"
  description  = "` + description + `"
}
`
}
