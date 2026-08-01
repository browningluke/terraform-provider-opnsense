package monit_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMonitAlertResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitAlertResourceConfig("test@example.com", "Test alert"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_monit_alert.test", "recipient", "test@example.com"),
					resource.TestCheckResourceAttr("opnsense_monit_alert.test", "description", "Test alert"),
					resource.TestCheckResourceAttrSet("opnsense_monit_alert.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_monit_alert.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitAlertResourceConfig("test@example.com", "Test alert updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_monit_alert.test", "description", "Test alert updated"),
				),
			},
		},
	})
}

func testAccMonitAlertResourceConfig(recipient, description string) string {
	return fmt.Sprintf(`
resource "opnsense_monit_alert" "test" {
  recipient   = %[1]q
  description = %[2]q
  events      = ["action", "timeout"]
}
`, recipient, description)
}
