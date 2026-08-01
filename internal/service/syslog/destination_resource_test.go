package syslog_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSyslogDestinationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSyslogDestinationResourceConfig("192.0.2.1", "514"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_syslog_destination.test", "hostname", "192.0.2.1"),
					resource.TestCheckResourceAttr("opnsense_syslog_destination.test", "port", "514"),
					resource.TestCheckResourceAttr("opnsense_syslog_destination.test", "transport", "udp4"),
					resource.TestCheckResourceAttr("opnsense_syslog_destination.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_syslog_destination.test", "description", "acc test destination"),
					resource.TestCheckResourceAttrSet("opnsense_syslog_destination.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_syslog_destination.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSyslogDestinationResourceConfig("192.0.2.1", "10514"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_syslog_destination.test", "port", "10514"),
				),
			},
		},
	})
}

func TestAccSyslogDestinationResource_Disabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSyslogDestinationResourceConfigDisabled("192.0.2.1", "514"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_syslog_destination.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet("opnsense_syslog_destination.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_syslog_destination.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSyslogDestinationResource_WithFacility(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSyslogDestinationResourceConfigWithFacility("192.0.2.1", "514"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_syslog_destination.test", "facility.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_syslog_destination.test", "facility.*", "local0"),
					resource.TestCheckResourceAttrSet("opnsense_syslog_destination.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_syslog_destination.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSyslogDestinationResourceConfig(hostname, port string) string {
	return fmt.Sprintf(`
resource "opnsense_syslog_destination" "test" {
  hostname    = %[1]q
  port        = %[2]q
  description = "acc test destination"
}
`, hostname, port)
}

func testAccSyslogDestinationResourceConfigDisabled(hostname, port string) string {
	return fmt.Sprintf(`
resource "opnsense_syslog_destination" "test" {
  enabled     = false
  hostname    = %[1]q
  port        = %[2]q
  description = "acc test destination disabled"
}
`, hostname, port)
}

func testAccSyslogDestinationResourceConfigWithFacility(hostname, port string) string {
	return fmt.Sprintf(`
resource "opnsense_syslog_destination" "test" {
  hostname    = %[1]q
  port        = %[2]q
  facility    = ["local0"]
  description = "acc test destination with facility"
}
`, hostname, port)
}
