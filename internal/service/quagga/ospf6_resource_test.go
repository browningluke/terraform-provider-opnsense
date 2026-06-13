package quagga_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaOSPF6Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccQuaggaOSPF6ResourceConfig(false, ""),
				ResourceName:       "opnsense_quagga_ospf6.test",
				ImportState:        true,
				ImportStateId:      "quagga_ospf6",
				ImportStatePersist: true,
			},
			{
				Config: testAccQuaggaOSPF6ResourceConfig(false, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6.test", "id", "quagga_ospf6"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6.test", "enabled", "false"),
				),
			},
			{
				Config: testAccQuaggaOSPF6ResourceConfig(false, "4.3.2.1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6.test", "router_id", "4.3.2.1"),
				),
			},
			{
				Config: testAccQuaggaOSPF6ResourceConfig(false, ""),
			},
		},
	})
}

func TestAccQuaggaOSPF6Resource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `resource "opnsense_quagga_ospf6" "test" {}`,
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

func testAccQuaggaOSPF6ResourceConfig(enabled bool, routerID string) string {
	return `
resource "opnsense_quagga_ospf6" "test" {
  enabled   = ` + boolStr(enabled) + `
  router_id = "` + routerID + `"
}
`
}
