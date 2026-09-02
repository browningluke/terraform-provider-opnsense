package dnsmasq_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsmasqDomainResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDnsmasqDomainResourceConfig(1, "test.domain.tld", "192.168.1.1", 5353, "Test description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_domain.test", "sequence", "1"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_domain.test", "domain", "test.domain.tld"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_domain.test", "ip", "192.168.1.1"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_domain.test", "port", "5353"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_domain.test", "description", "Test description"),
					resource.TestCheckResourceAttrSet("opnsense_dnsmasq_domain.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_dnsmasq_domain.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDnsmasqDomainResourceConfig(2, "updated.domain.tld", "192.168.2.2", 5354, "Updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_domain.test", "sequence", "2"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_domain.test", "domain", "updated.domain.tld"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_domain.test", "ip", "192.168.2.2"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_domain.test", "port", "5354"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_domain.test", "description", "Updated description"),
				),
			},
			// Delete testing occurs automatically
		},
	})
}

func testAccDnsmasqDomainResourceConfig(
	sequence int64,
	domain string,
	ip string,
	port int64,
	description string,
) string {
	return fmt.Sprintf(`
resource "opnsense_dnsmasq_domain" "test" {
	sequence       = %d
	domain         = "%s"
	ip             = "%s"
	port           = %d
	description    = "%s"
}
`, sequence, domain, ip, port, description)
}
