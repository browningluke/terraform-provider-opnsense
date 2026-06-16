package unbound_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUnboundAclResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundAclResourceConfig("lan-allow", "allow", []string{"10.0.0.0/24", "10.0.1.0/24"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_acl.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_unbound_acl.test", "name", "lan-allow"),
					resource.TestCheckResourceAttr("opnsense_unbound_acl.test", "action", "allow"),
					resource.TestCheckResourceAttr("opnsense_unbound_acl.test", "networks.#", "2"),
					resource.TestCheckTypeSetElemAttr("opnsense_unbound_acl.test", "networks.*", "10.0.0.0/24"),
					resource.TestCheckTypeSetElemAttr("opnsense_unbound_acl.test", "networks.*", "10.0.1.0/24"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_acl.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_acl.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccUnboundAclResourceConfig("lan-allow", "refuse", []string{"10.0.0.0/24"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_acl.test", "action", "refuse"),
					resource.TestCheckResourceAttr("opnsense_unbound_acl.test", "networks.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_unbound_acl.test", "networks.*", "10.0.0.0/24"),
				),
			},
		},
	})
}

func TestAccUnboundAclResource_Disabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundAclResourceConfigDisabled("corp-deny", "deny", []string{"172.16.0.0/12"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_acl.test", "enabled", "false"),
					resource.TestCheckResourceAttr("opnsense_unbound_acl.test", "action", "deny"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_acl.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_acl.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccUnboundAclResource_WithDescription(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundAclResourceConfigWithDescription("mgmt-snoop", "allow_snoop", []string{"192.168.0.0/24"}, "Management network"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_acl.test", "description", "Management network"),
					resource.TestCheckResourceAttr("opnsense_unbound_acl.test", "action", "allow_snoop"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_acl.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_acl.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccUnboundAclResourceConfig(name, action string, networks []string) string {
	networksStr := ""
	for _, n := range networks {
		networksStr += fmt.Sprintf("%q,\n    ", n)
	}
	return fmt.Sprintf(`
resource "opnsense_unbound_acl" "test" {
  name     = %[1]q
  action   = %[2]q
  networks = [
    %[3]s]
}
`, name, action, networksStr)
}

func testAccUnboundAclResourceConfigDisabled(name, action string, networks []string) string {
	networksStr := ""
	for _, n := range networks {
		networksStr += fmt.Sprintf("%q,\n    ", n)
	}
	return fmt.Sprintf(`
resource "opnsense_unbound_acl" "test" {
  enabled  = false
  name     = %[1]q
  action   = %[2]q
  networks = [
    %[3]s]
}
`, name, action, networksStr)
}

func testAccUnboundAclResourceConfigWithDescription(name, action string, networks []string, description string) string {
	networksStr := ""
	for _, n := range networks {
		networksStr += fmt.Sprintf("%q,\n    ", n)
	}
	return fmt.Sprintf(`
resource "opnsense_unbound_acl" "test" {
  name        = %[1]q
  action      = %[2]q
  networks    = [
    %[3]s]
  description = %[4]q
}
`, name, action, networksStr, description)
}
