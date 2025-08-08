package service

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallAliasResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallAliasResourceConfig("testalias", "Test alias description", "host", "192.168.1.100"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "name", "testalias"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "description", "Test alias description"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "type", "host"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "content.#", "1"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "content.0", "192.168.1.100"),
					resource.TestCheckResourceAttrSet("opnsense_firewall_alias.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_firewall_alias.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFirewallAliasResourceConfig("testaliasupdated", "Updated alias description", "host", "192.168.1.101"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "name", "testaliasupdated"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "description", "Updated alias description"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "type", "host"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "content.#", "1"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "content.0", "192.168.1.101"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccFirewallAliasResource_MultipleHosts(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallAliasResourceConfigMultiple("multihostalias", "Multiple hosts", "host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "name", "multihostalias"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "description", "Multiple hosts"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "type", "host"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "content.#", "3"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_alias.test", "content.*", "192.168.1.100"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_alias.test", "content.*", "192.168.1.101"),
					resource.TestCheckTypeSetElemAttr("opnsense_firewall_alias.test", "content.*", "192.168.1.102"),
				),
			},
		},
	})
}

func TestAccFirewallAliasResource_Network(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallAliasResourceConfig("networkalias", "Network alias", "network", "192.168.1.0/24"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "name", "networkalias"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "description", "Network alias"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "type", "network"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "content.#", "1"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "content.0", "192.168.1.0/24"),
				),
			},
		},
	})
}

func TestAccFirewallAliasResource_Port(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallAliasResourceConfig("portalias", "Port alias", "port", "80"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "name", "portalias"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "description", "Port alias"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "type", "port"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "content.#", "1"),
					resource.TestCheckResourceAttr("opnsense_firewall_alias.test", "content.0", "80"),
				),
			},
		},
	})
}

func testAccFirewallAliasResourceConfig(name, description, aliasType, content string) string {
	return fmt.Sprintf(`
resource "opnsense_firewall_alias" "test" {
  name        = %[1]q
  description = %[2]q
  type        = %[3]q
  content     = [%[4]q]
}
`, name, description, aliasType, content)
}

func testAccFirewallAliasResourceConfigMultiple(name, description, aliasType string) string {
	return fmt.Sprintf(`
resource "opnsense_firewall_alias" "test" {
  name        = %[1]q
  description = %[2]q
  type        = %[3]q
  content     = ["192.168.1.100", "192.168.1.101", "192.168.1.102"]
}
`, name, description, aliasType)
}
