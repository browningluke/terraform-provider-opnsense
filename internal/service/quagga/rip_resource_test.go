package quagga_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaRIPResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccQuaggaRIPResourceMinimalConfig(),
				ResourceName:       "opnsense_quagga_rip.test",
				ImportState:        true,
				ImportStateId:      "quagga_rip",
				ImportStatePersist: true,
			},
			{
				Config: testAccQuaggaRIPResourceMinimalConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_rip.test", "id", "quagga_rip"),
					resource.TestCheckResourceAttr("opnsense_quagga_rip.test", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccQuaggaRIPResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `resource "opnsense_quagga_rip" "test" {}`,
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

func testAccQuaggaRIPResourceMinimalConfig() string {
	return `
resource "opnsense_quagga_rip" "test" {
}
`
}
