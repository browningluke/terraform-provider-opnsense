package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfacesBridgeResource(t *testing.T) {
	member := acctest.BridgeMemberPreCheck(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBridgeResourceConfig(member, "Test bridge", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("opnsense_interfaces_bridge.test", "id"),
					resource.TestCheckResourceAttr("opnsense_interfaces_bridge.test", "members.#", "1"),
					resource.TestCheckResourceAttr("opnsense_interfaces_bridge.test", "members.0", member),
					resource.TestCheckResourceAttr("opnsense_interfaces_bridge.test", "description", "Test bridge"),
					resource.TestCheckResourceAttr("opnsense_interfaces_bridge.test", "enable_stp", "false"),
					resource.TestCheckResourceAttr("opnsense_interfaces_bridge.test", "stp_proto", "rstp"),
					resource.TestCheckResourceAttr("opnsense_interfaces_bridge.test", "link_local", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_interfaces_bridge.test",
				ImportState:       true,
				ImportStateVerify: true,
				// device is "" in config (auto-generate) but holds the actual
				// device name after import.
				ImportStateVerifyIgnore: []string{"device"},
			},
			// Update and Read testing
			{
				Config: testAccBridgeResourceConfig(member, "Updated bridge", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_bridge.test", "description", "Updated bridge"),
					resource.TestCheckResourceAttr("opnsense_interfaces_bridge.test", "enable_stp", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccBridgeResourceConfig(member string, description string, enableStp bool) string {
	return fmt.Sprintf(`
resource "opnsense_interfaces_bridge" "test" {
  members     = [%[1]q]
  description = %[2]q
  enable_stp  = %[3]t
}
`, member, description, enableStp)
}
