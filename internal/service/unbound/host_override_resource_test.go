package unbound_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUnboundHostOverrideResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundHostOverrideResourceConfig("testhost", "example.com", "192.168.1.10"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "hostname", "testhost"),
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "type", "A"),
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "server", "192.168.1.10"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_host_override.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_host_override.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccUnboundHostOverrideResourceConfig("testhost", "example.com", "192.168.1.11"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "server", "192.168.1.11"),
				),
			},
		},
	})
}

func TestAccUnboundHostOverrideResource_Disabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundHostOverrideResourceConfigDisabled("testhost", "example.com", "192.168.1.10"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_host_override.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_host_override.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccUnboundHostOverrideResource_AAAA(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundHostOverrideResourceConfigAAAA("testhost", "example.com", "fd00::1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "type", "AAAA"),
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "server", "fd00::1"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_host_override.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_host_override.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccUnboundHostOverrideResource_MX(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundHostOverrideResourceConfigMX("example.com", "mail.example.com", 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "type", "MX"),
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "mx_host", "mail.example.com"),
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "mx_priority", "10"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_host_override.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_host_override.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccUnboundHostOverrideResourceConfigMX("example.com", "mail.example.com", 20),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "mx_priority", "20"),
				),
			},
		},
	})
}

func TestAccUnboundHostOverrideResource_Wildcard(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundHostOverrideResourceConfig("*", "example.com", "192.168.1.10"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "hostname", "*"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_host_override.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_host_override.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccUnboundHostOverrideResource_WithDescription(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundHostOverrideResourceConfigWithDescription("testhost", "example.com", "192.168.1.10", "Test host override"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_host_override.test", "description", "Test host override"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_host_override.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_host_override.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccUnboundHostOverrideResourceConfig(hostname, domain, server string) string {
	return fmt.Sprintf(`
resource "opnsense_unbound_host_override" "test" {
  hostname = %[1]q
  domain   = %[2]q
  server   = %[3]q
}
`, hostname, domain, server)
}

func testAccUnboundHostOverrideResourceConfigDisabled(hostname, domain, server string) string {
	return fmt.Sprintf(`
resource "opnsense_unbound_host_override" "test" {
  enabled  = false
  hostname = %[1]q
  domain   = %[2]q
  server   = %[3]q
}
`, hostname, domain, server)
}

func testAccUnboundHostOverrideResourceConfigAAAA(hostname, domain, server string) string {
	return fmt.Sprintf(`
resource "opnsense_unbound_host_override" "test" {
  hostname = %[1]q
  domain   = %[2]q
  type     = "AAAA"
  server   = %[3]q
}
`, hostname, domain, server)
}

func testAccUnboundHostOverrideResourceConfigMX(domain, mxHost string, mxPriority int) string {
	return fmt.Sprintf(`
resource "opnsense_unbound_host_override" "test" {
  hostname    = "example.com"
  domain      = %[1]q
  type        = "MX"
  mx_host     = %[2]q
  mx_priority = %[3]d
}
`, domain, mxHost, mxPriority)
}

func testAccUnboundHostOverrideResourceConfigWithDescription(hostname, domain, server, description string) string {
	return fmt.Sprintf(`
resource "opnsense_unbound_host_override" "test" {
  hostname    = %[1]q
  domain      = %[2]q
  server      = %[3]q
  description = %[4]q
}
`, hostname, domain, server, description)
}
