package diagnostics_test

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDiagnosticsNetflowResource tests the singleton netflow resource.
//
// Because this resource blocks creation (terraform import must be used instead),
// the test begins with an import step rather than an apply step. Only
// active_timeout/inactive_timeout are toggled, since those don't start actual
// packet capture on the test VM's interfaces.
func TestAccDiagnosticsNetflowResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import the singleton resource. Create is blocked, so we must import first.
			{
				Config:             testAccNetflowResourceConfig(1800, 15),
				ResourceName:       "opnsense_diagnostics_netflow.netflow",
				ImportState:        true,
				ImportStateId:      "netflow_settings",
				ImportStatePersist: true,
			},
			// Apply the baseline config and verify key attributes round-trip correctly.
			{
				Config: testAccNetflowResourceConfig(1800, 15),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_diagnostics_netflow.netflow", "id", "netflow_settings"),
					resource.TestCheckResourceAttr("opnsense_diagnostics_netflow.netflow", "active_timeout", "1800"),
					resource.TestCheckResourceAttr("opnsense_diagnostics_netflow.netflow", "inactive_timeout", "15"),
				),
			},
			// Update timeouts and verify changes are applied upstream.
			{
				Config: testAccNetflowResourceConfig(1801, 16),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_diagnostics_netflow.netflow", "active_timeout", "1801"),
					resource.TestCheckResourceAttr("opnsense_diagnostics_netflow.netflow", "inactive_timeout", "16"),
				),
			},
			// Restore timeouts and verify.
			{
				Config: testAccNetflowResourceConfig(1800, 15),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_diagnostics_netflow.netflow", "active_timeout", "1800"),
					resource.TestCheckResourceAttr("opnsense_diagnostics_netflow.netflow", "inactive_timeout", "15"),
				),
			},
			// Delete testing: automatically removes from state only (no upstream change).
		},
	})
}

// TestAccDiagnosticsNetflowResource_CreateBlocked verifies that attempting to create
// this singleton resource without importing it first returns a clear error.
func TestAccDiagnosticsNetflowResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `resource "opnsense_diagnostics_netflow" "netflow" {}`,
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

func testAccNetflowResourceConfig(activeTimeout, inactiveTimeout int) string {
	return `
resource "opnsense_diagnostics_netflow" "netflow" {
  version          = "v9"
  active_timeout   = ` + strconv.Itoa(activeTimeout) + `
  inactive_timeout = ` + strconv.Itoa(inactiveTimeout) + `
}
`
}
