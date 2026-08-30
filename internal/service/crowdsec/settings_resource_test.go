package crowdsec_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccCrowdsecSettingsResource tests the singleton crowdsec settings resource.
//
// Because this resource blocks creation (terraform import must be used instead),
// the test begins with an import step rather than an apply step.
func TestAccCrowdsecSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import the singleton resource. Create is blocked, so we must import first.
			{
				Config:             testAccCrowdsecSettingsResourceConfig(true, false),
				ResourceName:       "opnsense_crowdsec_settings.settings",
				ImportState:        true,
				ImportStateId:      "crowdsec_settings",
				ImportStatePersist: true,
			},
			// Apply the baseline config and verify key attributes round-trip correctly.
			{
				Config: testAccCrowdsecSettingsResourceConfig(true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_crowdsec_settings.settings", "id", "crowdsec_settings"),
					resource.TestCheckResourceAttr("opnsense_crowdsec_settings.settings", "agent_enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_crowdsec_settings.settings", "rules_log", "false"),
					resource.TestCheckResourceAttr("opnsense_crowdsec_settings.settings", "lapi_listen_address", "127.0.0.1"),
					resource.TestCheckResourceAttr("opnsense_crowdsec_settings.settings", "lapi_listen_port", "8080"),
				),
			},
			// Update: disable the agent, enable rules_log.
			{
				Config: testAccCrowdsecSettingsResourceConfig(false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_crowdsec_settings.settings", "agent_enabled", "false"),
					resource.TestCheckResourceAttr("opnsense_crowdsec_settings.settings", "rules_log", "true"),
				),
			},
			// Restore original values.
			{
				Config: testAccCrowdsecSettingsResourceConfig(true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_crowdsec_settings.settings", "agent_enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_crowdsec_settings.settings", "rules_log", "false"),
				),
			},
			// ImportState verification.
			{
				ResourceName:      "opnsense_crowdsec_settings.settings",
				ImportState:       true,
				ImportStateId:     "crowdsec_settings",
				ImportStateVerify: true,
			},
			// Delete testing: automatically removes from state only (no upstream change).
		},
	})
}

// TestAccCrowdsecSettingsResource_CreateBlocked verifies that attempting to create
// this singleton resource without importing it first returns a clear error.
func TestAccCrowdsecSettingsResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCrowdsecSettingsResourceConfigMinimal(),
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

func testAccCrowdsecSettingsResourceConfig(agentEnabled, rulesLog bool) string {
	return `
resource "opnsense_crowdsec_settings" "settings" {
  agent_enabled    = ` + boolStr(agentEnabled) + `
  rules_log        = ` + boolStr(rulesLog) + `
}
`
}

func testAccCrowdsecSettingsResourceConfigMinimal() string {
	return `
resource "opnsense_crowdsec_settings" "settings" {
}
`
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
