package monit_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccMonitSettingsResource tests the singleton monit settings resource.
// Because this resource blocks creation, the test begins with an import step.
func TestAccMonitSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccMonitSettingsResourceConfig("120"),
				ResourceName:       "opnsense_monit_settings.test",
				ImportState:        true,
				ImportStateId:      "monit_settings",
				ImportStatePersist: true,
			},
			{
				Config: testAccMonitSettingsResourceConfig("120"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_monit_settings.test", "id", "monit_settings"),
					resource.TestCheckResourceAttr("opnsense_monit_settings.test", "startdelay", "120"),
				),
			},
			{
				Config: testAccMonitSettingsResourceConfig("90"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_monit_settings.test", "startdelay", "90"),
				),
			},
			// Restore original state
			{
				Config: testAccMonitSettingsResourceConfig("120"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_monit_settings.test", "startdelay", "120"),
				),
			},
		},
	})
}

// TestAccMonitSettingsResource_CreateBlocked verifies that attempting to
// create this singleton resource without importing it first returns a clear
// error.
func TestAccMonitSettingsResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `resource "opnsense_monit_settings" "test" {}`,
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

func testAccMonitSettingsResourceConfig(startdelay string) string {
	return `
resource "opnsense_monit_settings" "test" {
  enabled    = false
  interval   = "120"
  startdelay = "` + startdelay + `"
}
`
}
