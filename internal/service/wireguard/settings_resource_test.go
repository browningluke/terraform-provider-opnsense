package wireguard_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccWireguardSettingsResource tests the singleton wireguard settings resource.
//
// Because this resource blocks creation (terraform import must be used instead),
// the test begins with an import step rather than an apply step.
func TestAccWireguardSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import the singleton resource. Create is blocked, so we must import first.
			{
				Config:             testAccWireguardSettingsResourceConfig(true),
				ResourceName:       "opnsense_wireguard_settings.settings",
				ImportState:        true,
				ImportStateId:      "wireguard_settings",
				ImportStatePersist: true,
			},
			// Apply the baseline config and verify key attributes round-trip correctly.
			{
				Config: testAccWireguardSettingsResourceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_wireguard_settings.settings", "id", "wireguard_settings"),
					resource.TestCheckResourceAttr("opnsense_wireguard_settings.settings", "enabled", "true"),
				),
			},
			// Disable wireguard and verify.
			{
				Config: testAccWireguardSettingsResourceConfig(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_wireguard_settings.settings", "id", "wireguard_settings"),
					resource.TestCheckResourceAttr("opnsense_wireguard_settings.settings", "enabled", "false"),
				),
			},
			// Re-enable wireguard and verify.
			{
				Config: testAccWireguardSettingsResourceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_wireguard_settings.settings", "enabled", "true"),
				),
			},
			// Delete testing: automatically removes from state only (no upstream change).
		},
	})
}

// TestAccWireguardSettingsResource_CreateBlocked verifies that attempting to create
// this singleton resource without importing it first returns a clear error.
func TestAccWireguardSettingsResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccWireguardSettingsResourceConfigMinimal(),
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

func testAccWireguardSettingsResourceConfig(enabled bool) string {
	return `
resource "opnsense_wireguard_settings" "settings" {
  enabled = ` + boolStr(enabled) + `
}
`
}

func testAccWireguardSettingsResourceConfigMinimal() string {
	return `
resource "opnsense_wireguard_settings" "settings" {
}
`
}
