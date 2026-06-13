package quagga_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaOSPFResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccQuaggaOSPFResourceMinimalConfig(),
				ResourceName:       "opnsense_quagga_ospf.test",
				ImportState:        true,
				ImportStateId:      "quagga_ospf",
				ImportStatePersist: true,
			},
			{
				Config: testAccQuaggaOSPFResourceMinimalConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf.test", "id", "quagga_ospf"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf.test", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccQuaggaOSPFResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `resource "opnsense_quagga_ospf" "test" {}`,
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

func testAccQuaggaOSPFResourceMinimalConfig() string {
	return `
resource "opnsense_quagga_ospf" "test" {
}
`
}
