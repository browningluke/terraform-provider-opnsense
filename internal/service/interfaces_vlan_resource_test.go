package service

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfacesVlanResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccInterfacesVlanResourceConfig(100, "High VLAN ID test", 4, "vtnet0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_vlan.test", "tag", "100"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vlan.test", "description", "High VLAN ID test"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vlan.test", "priority", "4"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vlan.test", "parent", "vtnet0"),
					resource.TestCheckResourceAttrSet("opnsense_interfaces_vlan.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_interfaces_vlan.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccInterfacesVlanResourceConfig(100, "Updated VLAN 100", 6, "vtnet0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_vlan.test", "tag", "100"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vlan.test", "description", "Updated VLAN 100"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vlan.test", "priority", "6"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vlan.test", "parent", "vtnet0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccInterfacesVlanResource_HighVlanId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfacesVlanResourceConfig(4093, "High VLAN ID test", 6, "vtnet0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_vlan.test", "tag", "4093"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vlan.test", "description", "High VLAN ID test"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vlan.test", "priority", "6"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vlan.test", "parent", "vtnet0"),
				),
			},
		},
	})
}

func testAccInterfacesVlanResourceConfig(tag int, description string, priority int, parent string) string {
	return fmt.Sprintf(`
resource "opnsense_interfaces_vlan" "test" {
  tag         = %[1]d
  description = %[2]q
  priority    = %[3]d
  parent      = %[4]q
}
`, tag, description, priority, parent)
}
