package mdnsrepeater_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const plugin = "os-mdns-repeater"

// The CI VM only has a wan interface, and the repeater needs two to run, so
// these tests leave it disabled.
func TestAccMdnsrepeaterSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PluginPreCheck(t, plugin) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create is blocked, so start by importing.
			{
				Config:             testAccSettingsConfig(false, `[]`),
				ResourceName:       "opnsense_mdnsrepeater_settings.test",
				ImportState:        true,
				ImportStateId:      "mdnsrepeater_settings",
				ImportStatePersist: true,
			},
			{
				Config: testAccSettingsConfig(false, `[]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_mdnsrepeater_settings.test", "id", "mdnsrepeater_settings"),
					resource.TestCheckResourceAttr("opnsense_mdnsrepeater_settings.test", "enabled", "false"),
					resource.TestCheckResourceAttr("opnsense_mdnsrepeater_settings.test", "enable_carp", "false"),
					resource.TestCheckTypeSetElemAttr("opnsense_mdnsrepeater_settings.test", "interfaces.*", "wan"),
					resource.TestCheckResourceAttr("opnsense_mdnsrepeater_settings.test", "blocklist.#", "0"),
				),
			},
			{
				Config: testAccSettingsConfig(true, `["10.0.0.0/8", "192.168.1.10/32"]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_mdnsrepeater_settings.test", "enable_carp", "true"),
					resource.TestCheckResourceAttr("opnsense_mdnsrepeater_settings.test", "blocklist.#", "2"),
					resource.TestCheckTypeSetElemAttr("opnsense_mdnsrepeater_settings.test", "blocklist.*", "10.0.0.0/8"),
					resource.TestCheckTypeSetElemAttr("opnsense_mdnsrepeater_settings.test", "blocklist.*", "192.168.1.10/32"),
				),
			},
			{
				ResourceName:      "opnsense_mdnsrepeater_settings.test",
				ImportState:       true,
				ImportStateId:     "mdnsrepeater_settings",
				ImportStateVerify: true,
			},
			// Restore defaults.
			{
				Config: testAccSettingsConfig(false, `[]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_mdnsrepeater_settings.test", "blocklist.#", "0"),
				),
			},
		},
	})
}

func TestAccMdnsrepeaterSettingsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PluginPreCheck(t, plugin) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "opnsense_mdnsrepeater_settings" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.opnsense_mdnsrepeater_settings.test", "id", "mdnsrepeater_settings"),
					resource.TestCheckResourceAttrSet("data.opnsense_mdnsrepeater_settings.test", "enabled"),
				),
			},
		},
	})
}

func TestAccMdnsrepeaterSettingsResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PluginPreCheck(t, plugin) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSettingsConfig(false, `[]`),
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

func TestAccMdnsrepeaterSettingsResource_EnabledNeedsTwoInterfaces(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PluginPreCheck(t, plugin) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "opnsense_mdnsrepeater_settings" "test" {
  enabled    = true
  interfaces = ["wan"]
}
`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("requires at least 2 interfaces"),
			},
		},
	})
}

func testAccSettingsConfig(enableCarp bool, blocklist string) string {
	carp := "false"
	if enableCarp {
		carp = "true"
	}
	return `
resource "opnsense_mdnsrepeater_settings" "test" {
  enabled     = false
  enable_carp = ` + carp + `
  interfaces  = ["wan"]
  blocklist   = ` + blocklist + `
}
`
}
