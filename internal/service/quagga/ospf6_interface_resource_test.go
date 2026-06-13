package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaOSPF6InterfaceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaOSPF6InterfaceResourceConfig("wan", "0.0.0.0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_interface.test", "interface_name", "wan"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_interface.test", "area", "0.0.0.0"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_interface.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_ospf6_interface.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_ospf6_interface.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaOSPF6InterfaceResourceConfig("wan", "0.0.0.1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_interface.test", "area", "0.0.0.1"),
				),
			},
		},
	})
}

func testAccQuaggaOSPF6InterfaceResourceConfig(interfaceName, area string) string {
	return `
resource "opnsense_quagga_ospf6_interface" "test" {
  interface_name = "` + interfaceName + `"
  area           = "` + area + `"
}
`
}
