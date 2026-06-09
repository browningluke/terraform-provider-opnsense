package unbound_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUnboundDomainOverrideResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.DomainOverridePreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundDomainOverrideResourceConfig("internal.example.com", "192.168.1.53"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_domain_override.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_unbound_domain_override.test", "domain", "internal.example.com"),
					resource.TestCheckResourceAttr("opnsense_unbound_domain_override.test", "server", "192.168.1.53"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_domain_override.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_domain_override.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccUnboundDomainOverrideResourceConfig("internal.example.com", "192.168.1.54"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_domain_override.test", "server", "192.168.1.54"),
				),
			},
		},
	})
}

func TestAccUnboundDomainOverrideResource_Disabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.DomainOverridePreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundDomainOverrideResourceConfigDisabled("internal.example.com", "192.168.1.53"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_domain_override.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_domain_override.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_domain_override.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccUnboundDomainOverrideResource_WithPortSuffix(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.DomainOverridePreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundDomainOverrideResourceConfig("internal.example.com", "192.168.1.53@5353"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_domain_override.test", "server", "192.168.1.53@5353"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_domain_override.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_domain_override.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccUnboundDomainOverrideResource_WithDescription(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.DomainOverridePreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundDomainOverrideResourceConfigWithDescription("internal.example.com", "192.168.1.53", "Internal DNS zone"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_domain_override.test", "description", "Internal DNS zone"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_domain_override.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_domain_override.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccUnboundDomainOverrideResourceConfig(domain, server string) string {
	return fmt.Sprintf(`
resource "opnsense_unbound_domain_override" "test" {
  domain = %[1]q
  server = %[2]q
}
`, domain, server)
}

func testAccUnboundDomainOverrideResourceConfigDisabled(domain, server string) string {
	return fmt.Sprintf(`
resource "opnsense_unbound_domain_override" "test" {
  enabled = false
  domain  = %[1]q
  server  = %[2]q
}
`, domain, server)
}

func testAccUnboundDomainOverrideResourceConfigWithDescription(domain, server, description string) string {
	return fmt.Sprintf(`
resource "opnsense_unbound_domain_override" "test" {
  domain      = %[1]q
  server      = %[2]q
  description = %[3]q
}
`, domain, server, description)
}
