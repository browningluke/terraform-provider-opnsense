package unbound_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUnboundHostAliasResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundHostAliasResourceConfig("testalias", "example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_host_alias.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_unbound_host_alias.test", "hostname", "testalias"),
					resource.TestCheckResourceAttr("opnsense_unbound_host_alias.test", "domain", "example.com"),
					resource.TestCheckResourceAttrPair(
						"opnsense_unbound_host_alias.test", "override",
						"opnsense_unbound_host_override.test", "id",
					),
					resource.TestCheckResourceAttrSet("opnsense_unbound_host_alias.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_host_alias.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccUnboundHostAliasResourceConfigUpdated("testalias-upd", "example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_host_alias.test", "hostname", "testalias-upd"),
				),
			},
		},
	})
}

func TestAccUnboundHostAliasResource_Disabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundHostAliasResourceConfigDisabled("disabledalias", "example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_host_alias.test", "enabled", "false"),
					resource.TestCheckResourceAttr("opnsense_unbound_host_alias.test", "hostname", "disabledalias"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_host_alias.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_host_alias.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccUnboundHostAliasResource_WithDescription(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundHostAliasResourceConfigWithDescription("descalias", "example.com", "Test host alias"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_host_alias.test", "hostname", "descalias"),
					resource.TestCheckResourceAttr("opnsense_unbound_host_alias.test", "description", "Test host alias"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_host_alias.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_host_alias.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccUnboundHostAliasBase() string {
	return `
resource "opnsense_unbound_host_override" "test" {
  hostname = "basehost"
  domain   = "example.com"
  server   = "192.168.1.100"
}
`
}

func testAccUnboundHostAliasResourceConfig(hostname, domain string) string {
	return testAccUnboundHostAliasBase() + fmt.Sprintf(`
resource "opnsense_unbound_host_alias" "test" {
  override = opnsense_unbound_host_override.test.id
  hostname = %[1]q
  domain   = %[2]q
}
`, hostname, domain)
}

func testAccUnboundHostAliasResourceConfigUpdated(hostname, domain string) string {
	return testAccUnboundHostAliasBase() + fmt.Sprintf(`
resource "opnsense_unbound_host_alias" "test" {
  override = opnsense_unbound_host_override.test.id
  hostname = %[1]q
  domain   = %[2]q
}
`, hostname, domain)
}

func testAccUnboundHostAliasResourceConfigDisabled(hostname, domain string) string {
	return testAccUnboundHostAliasBase() + fmt.Sprintf(`
resource "opnsense_unbound_host_alias" "test" {
  override = opnsense_unbound_host_override.test.id
  enabled  = false
  hostname = %[1]q
  domain   = %[2]q
}
`, hostname, domain)
}

func testAccUnboundHostAliasResourceConfigWithDescription(hostname, domain, description string) string {
	return testAccUnboundHostAliasBase() + fmt.Sprintf(`
resource "opnsense_unbound_host_alias" "test" {
  override    = opnsense_unbound_host_override.test.id
  hostname    = %[1]q
  domain      = %[2]q
  description = %[3]q
}
`, hostname, domain, description)
}
