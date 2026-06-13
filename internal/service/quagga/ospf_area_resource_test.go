package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaOSPFAreaResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaOSPFAreaResourceConfig("0.0.0.1", "stub"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_area.test", "area_id", "0.0.0.1"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_area.test", "type", "stub"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_area.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_ospf_area.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_ospf_area.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaOSPFAreaResourceConfig("0.0.0.1", "nssa"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_area.test", "type", "nssa"),
				),
			},
		},
	})
}

func testAccQuaggaOSPFAreaResourceConfig(areaID, areaType string) string {
	return `
resource "opnsense_quagga_ospf_area" "test" {
  area_id = "` + areaID + `"
  type    = "` + areaType + `"
}
`
}
