package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaOSPFRedistributionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaOSPFRedistributionResourceConfig("connected", "test-ospf-redistribution"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_redistribution.test", "redistribute", "connected"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_redistribution.test", "description", "test-ospf-redistribution"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_redistribution.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_ospf_redistribution.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_ospf_redistribution.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaOSPFRedistributionResourceConfig("static", "test-ospf-redistribution-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_redistribution.test", "redistribute", "static"),
				),
			},
		},
	})
}

func testAccQuaggaOSPFRedistributionResourceConfig(redistribute, description string) string {
	return `
resource "opnsense_quagga_ospf_redistribution" "test" {
  redistribute = "` + redistribute + `"
  description  = "` + description + `"
}
`
}
