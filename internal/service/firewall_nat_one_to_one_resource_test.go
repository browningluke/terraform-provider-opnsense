package service

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallNatOneToOneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFirewallNatOneToOneResourceConfig(false, false, "10.10.10.22/32", "nat", "192.168.3.22/32", false, "default", "Tesging NAT One-to-One"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "enabled", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "log", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "external_net", "10.10.10.22/32"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "type", "nat"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "source.net", "192.168.3.22/32"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "source.invert", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "nat_reflection", "default"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "description", "Tesging NAT One-to-One"),
					resource.TestCheckResourceAttrSet("opnsense_firewall_nat_one_to_one.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_firewall_nat_one_to_one.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccFirewallNatOneToOneResourceConfig(true, true, "10.10.10.23/32", "nat", "192.168.3.23/32", false, "default", "Updated NAT One-to-One"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "log", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "external_net", "10.10.10.23/32"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "type", "nat"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "source.net", "192.168.3.23/32"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "source.invert", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "nat_reflection", "default"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "description", "Updated NAT One-to-One"),
					resource.TestCheckResourceAttrSet("opnsense_firewall_nat_one_to_one.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccFirewallNatOneToOneBinatResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFirewallNatOneToOneResourceConfig(false, false, "10.10.10.22/32", "binat", "192.168.3.22/32", false, "default", "Tesging BNAT One-to-One"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "enabled", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "log", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "external_net", "10.10.10.22/32"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "type", "binat"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "source.net", "192.168.3.22/32"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "source.invert", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "nat_reflection", "default"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "description", "Tesging BNAT One-to-One"),
					resource.TestCheckResourceAttrSet("opnsense_firewall_nat_one_to_one.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_firewall_nat_one_to_one.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccFirewallNatOneToOneResourceConfig(true, true, "10.10.10.23/32", "binat", "192.168.3.23/32", false, "default", "Updated BNAT One-to-One"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "log", "true"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "external_net", "10.10.10.23/32"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "type", "binat"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "source.net", "192.168.3.23/32"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "source.invert", "false"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "nat_reflection", "default"),
					resource.TestCheckResourceAttr("opnsense_firewall_nat_one_to_one.test", "description", "Updated BNAT One-to-One"),
					resource.TestCheckResourceAttrSet("opnsense_firewall_nat_one_to_one.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccFirewallNatOneToOneResourceConfig(enabled, log bool, external_net, typ, source_net string, source_invert bool, nat_reflection, description string) string {
	return fmt.Sprintf(`
resource "opnsense_firewall_nat_one_to_one" "test" {
  enabled         = %[1]t
  log             = %[2]t
  external_net    = %[3]q
  type            = %[4]q
  source = {
    net           = %[5]q
	invert        = %[6]t
  }
  nat_reflection  = %[7]q
  description     = %[8]q
}
`, enabled, log, external_net, typ, source_net, source_invert, nat_reflection, description)
}
