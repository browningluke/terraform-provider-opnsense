package quagga_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaStaticResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccQuaggaStaticResourceConfig(false),
				ResourceName:       "opnsense_quagga_static.test",
				ImportState:        true,
				ImportStateId:      "quagga_static",
				ImportStatePersist: true,
			},
			{
				Config: testAccQuaggaStaticResourceConfig(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_static.test", "id", "quagga_static"),
					resource.TestCheckResourceAttr("opnsense_quagga_static.test", "enabled", "false"),
				),
			},
			{
				Config: testAccQuaggaStaticResourceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_static.test", "enabled", "true"),
				),
			},
			{
				Config: testAccQuaggaStaticResourceConfig(false),
			},
		},
	})
}

func TestAccQuaggaStaticResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `resource "opnsense_quagga_static" "test" {}`,
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

func testAccQuaggaStaticResourceConfig(enabled bool) string {
	return `
resource "opnsense_quagga_static" "test" {
  enabled = ` + boolStr(enabled) + `
}
`
}
