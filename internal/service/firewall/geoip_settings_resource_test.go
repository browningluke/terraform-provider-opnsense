package firewall_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccFirewallGeoIPSettingsResource tests the singleton firewall GeoIP
// settings resource. Because this resource blocks creation, the test begins
// with an import step.
//
// NOTE: applying this resource triggers a real GeoIP database download on
// the OPNsense host (POST /firewall/alias/update/geoip), so the configured
// URL must be reachable from the test runner/VM.
func TestAccFirewallGeoIPSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccFirewallGeoIPSettingsResourceConfig(testAccGeoIPURLPrimary),
				ResourceName:       "opnsense_firewall_geoip_settings.test",
				ImportState:        true,
				ImportStateId:      "firewall_geoip_settings",
				ImportStatePersist: true,
			},
			{
				Config: testAccFirewallGeoIPSettingsResourceConfig(testAccGeoIPURLPrimary),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_geoip_settings.test", "id", "firewall_geoip_settings"),
					resource.TestCheckResourceAttr("opnsense_firewall_geoip_settings.test", "url", testAccGeoIPURLPrimary),
				),
			},
			{
				Config: testAccFirewallGeoIPSettingsResourceConfig(testAccGeoIPURLUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_geoip_settings.test", "url", testAccGeoIPURLUpdated),
				),
			},
		},
	})
}

// TestAccFirewallGeoIPSettingsResource_CreateBlocked verifies that attempting
// to create this singleton resource without importing it first returns a
// clear error.
func TestAccFirewallGeoIPSettingsResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `resource "opnsense_firewall_geoip_settings" "test" { url = "https://example.com/geoip.zip" }`,
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

const (
	testAccGeoIPURLPrimary = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country-CSV&license_key=test&suffix=zip"
	testAccGeoIPURLUpdated = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country-CSV&license_key=test2&suffix=zip"
)

func testAccFirewallGeoIPSettingsResourceConfig(url string) string {
	return `
resource "opnsense_firewall_geoip_settings" "test" {
  url = "` + url + `"
}
`
}
