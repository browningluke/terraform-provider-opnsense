package quagga_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaGeneralResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccQuaggaGeneralResourceConfig(false, "traditional", "notifications"),
				ResourceName:       "opnsense_quagga_general.test",
				ImportState:        true,
				ImportStateId:      "quagga_general",
				ImportStatePersist: true,
			},
			{
				Config: testAccQuaggaGeneralResourceConfig(false, "traditional", "notifications"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_general.test", "id", "quagga_general"),
					resource.TestCheckResourceAttr("opnsense_quagga_general.test", "enabled", "false"),
					resource.TestCheckResourceAttr("opnsense_quagga_general.test", "profile", "traditional"),
					resource.TestCheckResourceAttr("opnsense_quagga_general.test", "syslog_level", "notifications"),
				),
			},
			{
				Config: testAccQuaggaGeneralResourceConfig(false, "datacenter", "informational"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_general.test", "profile", "datacenter"),
					resource.TestCheckResourceAttr("opnsense_quagga_general.test", "syslog_level", "informational"),
				),
			},
			{
				Config: testAccQuaggaGeneralResourceConfig(false, "traditional", "notifications"),
			},
		},
	})
}

func TestAccQuaggaGeneralResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccQuaggaGeneralResourceMinimalConfig(),
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

func testAccQuaggaGeneralResourceConfig(enabled bool, profile, syslogLevel string) string {
	return `
resource "opnsense_quagga_general" "test" {
  enabled      = ` + boolStr(enabled) + `
  profile      = "` + profile + `"
  syslog_level = "` + syslogLevel + `"
}
`
}

func testAccQuaggaGeneralResourceMinimalConfig() string {
	return `
resource "opnsense_quagga_general" "test" {
}
`
}
