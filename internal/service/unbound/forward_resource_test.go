package unbound_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUnboundForwardResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundForwardResourceConfig("", "1.1.1.1", 53),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_forward.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_unbound_forward.test", "domain", ""),
					resource.TestCheckResourceAttr("opnsense_unbound_forward.test", "server_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("opnsense_unbound_forward.test", "server_port", "53"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_forward.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_forward.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccUnboundForwardResourceConfig("", "1.0.0.1", 53),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_forward.test", "server_ip", "1.0.0.1"),
				),
			},
		},
	})
}

func TestAccUnboundForwardResource_DomainScoped(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundForwardResourceConfig("example.com", "192.168.1.53", 53),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_forward.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("opnsense_unbound_forward.test", "server_ip", "192.168.1.53"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_forward.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_forward.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccUnboundForwardResource_DoT(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundForwardResourceConfigDoT("", "1.1.1.1", 853, "cloudflare-dns.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_forward.test", "server_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("opnsense_unbound_forward.test", "server_port", "853"),
					resource.TestCheckResourceAttr("opnsense_unbound_forward.test", "verify_cn", "cloudflare-dns.com"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_forward.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_forward.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccUnboundForwardResource_Disabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnboundForwardResourceConfigDisabled("", "1.1.1.1", 53),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_unbound_forward.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet("opnsense_unbound_forward.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_unbound_forward.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccUnboundForwardResourceConfig(domain, serverIP string, serverPort int) string {
	return fmt.Sprintf(`
resource "opnsense_unbound_forward" "test" {
  domain      = %[1]q
  server_ip   = %[2]q
  server_port = %[3]d
}
`, domain, serverIP, serverPort)
}

func testAccUnboundForwardResourceConfigDoT(domain, serverIP string, serverPort int, verifyCN string) string {
	return fmt.Sprintf(`
resource "opnsense_unbound_forward" "test" {
  domain      = %[1]q
  server_ip   = %[2]q
  server_port = %[3]d
  verify_cn   = %[4]q
}
`, domain, serverIP, serverPort, verifyCN)
}

func testAccUnboundForwardResourceConfigDisabled(domain, serverIP string, serverPort int) string {
	return fmt.Sprintf(`
resource "opnsense_unbound_forward" "test" {
  enabled     = false
  domain      = %[1]q
  server_ip   = %[2]q
  server_port = %[3]d
}
`, domain, serverIP, serverPort)
}
