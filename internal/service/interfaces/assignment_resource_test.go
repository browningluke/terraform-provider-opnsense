package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccInterfacesAssignmentResource exercises the full create/update/import/
// delete lifecycle. Rather than requiring a spare physical NIC, it creates a
// throwaway VLAN on the CI VM's NIC and assigns that VLAN's device to a new
// optional interface, so the test runs on a single-NIC box without touching
// live traffic.
//
// Skipped unless OPNSENSE_TEST_ASSIGNMENT=1: the underlying
// AssignmentController endpoint only exists on OPNsense 26.7+, and the CI VM
// image is 26.1.
func TestAccInterfacesAssignmentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AssignmentPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAssignmentResourceConfig("AssignmentTest"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_assignment.test", "device", testAccAssignmentVlanDevice),
					resource.TestCheckResourceAttr("opnsense_interfaces_assignment.test", "description", "AssignmentTest"),
					resource.TestCheckResourceAttr("opnsense_interfaces_assignment.test", "lock", "false"),
					resource.TestCheckResourceAttrSet("opnsense_interfaces_assignment.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_interfaces_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAssignmentResourceConfig("UpdatedAssignmentTest"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_assignment.test", "device", testAccAssignmentVlanDevice),
					resource.TestCheckResourceAttr("opnsense_interfaces_assignment.test", "description", "UpdatedAssignmentTest"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

const (
	// testAccAssignmentVlanDevice is set explicitly on the VLAN resource so the
	// assignment can reference it directly. Left blank, OPNsense generates the
	// device name and it would have to be looked up via
	// opnsense_interfaces_overview_all.
	testAccAssignmentVlanDevice = "vlan0.147"
	testAccAssignmentVlanParent = "vtnet0"
	testAccAssignmentVlanTag    = 147
)

func testAccAssignmentResourceConfig(description string) string {
	return fmt.Sprintf(`
resource "opnsense_interfaces_vlan" "test" {
  description = "assignment acc test"
  tag         = %[1]d
  priority    = 0
  parent      = %[2]q
  device      = %[3]q
}

resource "opnsense_interfaces_assignment" "test" {
  device      = opnsense_interfaces_vlan.test.device
  description = %[4]q
}
`, testAccAssignmentVlanTag, testAccAssignmentVlanParent, testAccAssignmentVlanDevice, description)
}
