package crowdsec_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccCrowdsecSettingsDataSource verifies the singleton crowdsec settings
// data source reads the upstream configuration.
//
// The data source takes no arguments, so every attribute is computed from the
// upstream system. The test asserts the fixed `id` and that each attribute is
// populated, rather than asserting values that vary per-instance.
func TestAccCrowdsecSettingsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCrowdsecSettingsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.opnsense_crowdsec_settings.test", "id", "crowdsec_settings"),
					resource.TestCheckResourceAttrSet("data.opnsense_crowdsec_settings.test", "agent_enabled"),
					resource.TestCheckResourceAttrSet("data.opnsense_crowdsec_settings.test", "lapi_enabled"),
					resource.TestCheckResourceAttrSet("data.opnsense_crowdsec_settings.test", "firewall_bouncer_enabled"),
					resource.TestCheckResourceAttrSet("data.opnsense_crowdsec_settings.test", "lapi_manual_configuration"),
					resource.TestCheckResourceAttrSet("data.opnsense_crowdsec_settings.test", "lapi_listen_address"),
					resource.TestCheckResourceAttrSet("data.opnsense_crowdsec_settings.test", "lapi_listen_port"),
					resource.TestCheckResourceAttrSet("data.opnsense_crowdsec_settings.test", "rules_enabled"),
					resource.TestCheckResourceAttrSet("data.opnsense_crowdsec_settings.test", "rules_log"),
					resource.TestCheckResourceAttrSet("data.opnsense_crowdsec_settings.test", "crowdsec_firewall_verbose"),
				),
			},
		},
	})
}

// TestAccCrowdsecSettingsDataSource_MatchesResource verifies the data source
// reports the same values as the managed resource after an apply.
func TestAccCrowdsecSettingsDataSource_MatchesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import the singleton resource. Create is blocked, so we must import first.
			{
				Config:             testAccCrowdsecSettingsDataSourceConfigWithResource(),
				ResourceName:       "opnsense_crowdsec_settings.settings",
				ImportState:        true,
				ImportStateId:      "crowdsec_settings",
				ImportStatePersist: true,
			},
			{
				Config: testAccCrowdsecSettingsDataSourceConfigWithResource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.opnsense_crowdsec_settings.test", "agent_enabled",
						"opnsense_crowdsec_settings.settings", "agent_enabled",
					),
					resource.TestCheckResourceAttrPair(
						"data.opnsense_crowdsec_settings.test", "rules_log",
						"opnsense_crowdsec_settings.settings", "rules_log",
					),
					resource.TestCheckResourceAttrPair(
						"data.opnsense_crowdsec_settings.test", "lapi_listen_address",
						"opnsense_crowdsec_settings.settings", "lapi_listen_address",
					),
					resource.TestCheckResourceAttrPair(
						"data.opnsense_crowdsec_settings.test", "lapi_listen_port",
						"opnsense_crowdsec_settings.settings", "lapi_listen_port",
					),
				),
			},
		},
	})
}

func testAccCrowdsecSettingsDataSourceConfig() string {
	return `
data "opnsense_crowdsec_settings" "test" {}
`
}

func testAccCrowdsecSettingsDataSourceConfigWithResource() string {
	return `
resource "opnsense_crowdsec_settings" "settings" {
  agent_enabled = true
  rules_log     = false
}

data "opnsense_crowdsec_settings" "test" {
  depends_on = [opnsense_crowdsec_settings.settings]
}
`
}
