package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaOSPFInterfaceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaOSPFInterfaceResourceConfig("wan", "0.0.0.0", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_interface.test", "interface_name", "wan"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_interface.test", "area", "0.0.0.0"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_interface.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_ospf_interface.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_ospf_interface.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaOSPFInterfaceResourceConfig("wan", "0.0.0.0", "broadcast"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_interface.test", "network_type", "broadcast"),
				),
			},
		},
	})
}

func testAccQuaggaOSPFInterfaceResourceConfig(interfaceName, area, networkType string) string {
	config := `
resource "opnsense_quagga_ospf_interface" "test" {
  interface_name = "` + interfaceName + `"
  area           = "` + area + `"
`
	if networkType != "" {
		config += `  network_type   = "` + networkType + `"` + "\n"
	}
	config += `}
`
	return config
}
