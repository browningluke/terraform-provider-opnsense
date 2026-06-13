package quagga_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaBFDResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccQuaggaBFDResourceConfig(false),
				ResourceName:       "opnsense_quagga_bfd.test",
				ImportState:        true,
				ImportStateId:      "quagga_bfd",
				ImportStatePersist: true,
			},
			{
				Config: testAccQuaggaBFDResourceConfig(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bfd.test", "id", "quagga_bfd"),
					resource.TestCheckResourceAttr("opnsense_quagga_bfd.test", "enabled", "false"),
				),
			},
			{
				Config: testAccQuaggaBFDResourceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bfd.test", "enabled", "true"),
				),
			},
			{
				Config: testAccQuaggaBFDResourceConfig(false),
			},
		},
	})
}

func TestAccQuaggaBFDResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `resource "opnsense_quagga_bfd" "test" {}`,
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

func testAccQuaggaBFDResourceConfig(enabled bool) string {
	return `
resource "opnsense_quagga_bfd" "test" {
  enabled = ` + boolStr(enabled) + `
}
`
}
