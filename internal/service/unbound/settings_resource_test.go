package unbound_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccUnboundSettingsResource tests the singleton unbound settings resource.
//
// Because this resource blocks creation (terraform import must be used instead),
// the test begins with an import step rather than an apply step.
//
// The general block is always included explicitly in test configs with the actual
// upstream values (enabled=true, local_zone_type="transparent"). This prevents
// the schema's objectdefault from overriding the imported state with Default values
// that differ from the live OPNsense configuration.
func TestAccUnboundSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import the singleton resource. Create is blocked, so we must import first.
			{
				Config:             testAccSettingsResourceConfig(false, false, false),
				ResourceName:       "opnsense_unbound_settings.settings",
				ImportState:        true,
				ImportStateId:      "unbound_settings",
				ImportStatePersist: true,
			},
			// Apply the baseline config and verify key attributes round-trip correctly.
			{
				Config: testAccSettingsResourceConfig(false, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_settings.settings", "id", "unbound_settings"),
					resource.TestCheckResourceAttr("opnsense_unbound_settings.settings", "general.enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_unbound_settings.settings", "advanced.hide_identity", "false"),
					resource.TestCheckResourceAttr("opnsense_unbound_settings.settings", "advanced.hide_version", "false"),
					resource.TestCheckResourceAttr("opnsense_unbound_settings.settings", "advanced.logging.log_queries", "false"),
				),
			},
			// Update safe settings and verify changes are applied upstream.
			{
				Config: testAccSettingsResourceConfig(true, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_settings.settings", "id", "unbound_settings"),
					resource.TestCheckResourceAttr("opnsense_unbound_settings.settings", "advanced.hide_identity", "true"),
					resource.TestCheckResourceAttr("opnsense_unbound_settings.settings", "advanced.hide_version", "true"),
					resource.TestCheckResourceAttr("opnsense_unbound_settings.settings", "advanced.logging.log_queries", "true"),
				),
			},
			// Restore settings and verify.
			{
				Config: testAccSettingsResourceConfig(false, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_settings.settings", "advanced.hide_identity", "false"),
					resource.TestCheckResourceAttr("opnsense_unbound_settings.settings", "advanced.hide_version", "false"),
					resource.TestCheckResourceAttr("opnsense_unbound_settings.settings", "advanced.logging.log_queries", "false"),
				),
			},
			// Delete testing: automatically removes from state only (no upstream change).
		},
	})
}

// TestAccUnboundSettingsResource_CreateBlocked verifies that attempting to create
// this singleton resource without importing it first returns a clear error.
func TestAccUnboundSettingsResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSettingsResourceConfigMinimal(),
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

// testAccSettingsResourceConfig returns a resource config that always explicitly
// sets the general block with the actual upstream values on the test VM
// (enabled=true, local_zone_type="transparent"). This prevents objectdefault from
// overriding imported state with schema defaults that differ from the live config.
//
// The advanced block sets safe-to-toggle fields so the test can verify round-trip
// updates without disrupting the DNS service.
func testAccSettingsResourceConfig(hideIdentity, hideVersion, logQueries bool) string {
	return `
resource "opnsense_unbound_settings" "settings" {
  general = {
    enabled         = true
    local_zone_type = "transparent"
  }

  advanced = {
    hide_identity = ` + boolStr(hideIdentity) + `
    hide_version  = ` + boolStr(hideVersion) + `

    logging = {
      log_queries = ` + boolStr(logQueries) + `
    }
  }
}
`
}

// testAccSettingsResourceConfigMinimal is used only by TestAccUnboundSettingsResource_CreateBlocked
// to verify that attempting to create (rather than import) the singleton fails.
func testAccSettingsResourceConfigMinimal() string {
	return `
resource "opnsense_unbound_settings" "settings" {}
`
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
