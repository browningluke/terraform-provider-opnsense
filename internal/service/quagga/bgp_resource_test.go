package quagga_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccQuaggaBGPResource tests the singleton BGP resource.
// The test uses enabled=false/true toggle since as_number is a required field
// and toggling it in tests can leave the upstream in an unpredictable state
// between test runs.
func TestAccQuaggaBGPResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import the singleton with its current upstream state.
			{
				Config:             testAccQuaggaBGPResourceMinimalConfig(),
				ResourceName:       "opnsense_quagga_bgp.test",
				ImportState:        true,
				ImportStateId:      "quagga_bgp",
				ImportStatePersist: true,
			},
			// Verify we can read the resource without error.
			{
				Config: testAccQuaggaBGPResourceMinimalConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp.test", "id", "quagga_bgp"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_bgp.test", "as_number"),
				),
			},
		},
	})
}

func TestAccQuaggaBGPResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `resource "opnsense_quagga_bgp" "test" { as_number = "65001" }`,
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

func testAccQuaggaBGPResourceMinimalConfig() string {
	return `
resource "opnsense_quagga_bgp" "test" {
}
`
}
