package monit_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMonitServiceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitServiceResourceConfig("test-service", "Test service"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_monit_service.test", "name", "test-service"),
					resource.TestCheckResourceAttr("opnsense_monit_service.test", "description", "Test service"),
					resource.TestCheckResourceAttr("opnsense_monit_service.test", "type", "process"),
					resource.TestCheckResourceAttrSet("opnsense_monit_service.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_monit_service.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitServiceResourceConfig("test-service", "Test service updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_monit_service.test", "description", "Test service updated"),
				),
			},
		},
	})
}

func testAccMonitServiceResourceConfig(name, description string) string {
	return fmt.Sprintf(`
resource "opnsense_monit_service" "test" {
  name        = %[1]q
  description = %[2]q
  type        = "process"
  pidfile     = "/var/run/test.pid"
  timeout     = "300"
}
`, name, description)
}
