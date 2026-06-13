package trust_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccTrustSettingsResource tests the singleton trust settings resource.
// Because this resource blocks creation, the test begins with an import step.
func TestAccTrustSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccTrustSettingsResourceConfig(false, false),
				ResourceName:       "opnsense_trust_settings.test",
				ImportState:        true,
				ImportStateId:      "trust_settings",
				ImportStatePersist: true,
			},
			{
				Config: testAccTrustSettingsResourceConfig(false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_trust_settings.test", "id", "trust_settings"),
					resource.TestCheckResourceAttr("opnsense_trust_settings.test", "install_crls", "false"),
					resource.TestCheckResourceAttr("opnsense_trust_settings.test", "fetch_crls", "false"),
				),
			},
			{
				Config: testAccTrustSettingsResourceConfig(true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_trust_settings.test", "store_intermediate_certs", "true"),
				),
			},
			// Restore original state
			{
				Config: testAccTrustSettingsResourceConfig(false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_trust_settings.test", "store_intermediate_certs", "false"),
				),
			},
		},
	})
}

// TestAccTrustSettingsResource_CreateBlocked verifies that attempting to create
// this singleton resource without importing it first returns a clear error.
func TestAccTrustSettingsResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `resource "opnsense_trust_settings" "test" {}`,
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

func testAccTrustSettingsResourceConfig(storeIntermediateCerts, fetchCrls bool) string {
	return `
resource "opnsense_trust_settings" "test" {
  store_intermediate_certs  = ` + boolStr(storeIntermediateCerts) + `
  install_crls              = false
  fetch_crls                = ` + boolStr(fetchCrls) + `
  enable_legacy_sect        = true
  enable_config_constraints = false
}
`
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
