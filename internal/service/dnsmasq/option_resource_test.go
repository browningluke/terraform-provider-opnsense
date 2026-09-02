package dnsmasq_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsmasqOptionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccPreCheck(t) }, ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDnsmasqOptionResourceConfig(6, "8.8.8.8", "initial option"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_option.test", "type", "set"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_option.test", "option", "6"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_option.test", "value", "8.8.8.8"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_option.test", "description", "initial option"),
					resource.TestCheckResourceAttrSet("opnsense_dnsmasq_option.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_dnsmasq_option.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDnsmasqOptionResourceConfig(42, "example.test", "updated option"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_option.test", "option", "42"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_option.test", "value", "example.test"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_option.test", "description", "updated option"),
				),
			},
		},
	})
}

func testAccDnsmasqOptionResourceConfig(option int, value, description string) string {
	return fmt.Sprintf(`
resource "opnsense_dnsmasq_option" "test" {
  type        = "set"
  option      = %d
  value       = %q
  description = %q
}
`, option, value, description)
}
