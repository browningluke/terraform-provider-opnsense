package dnsmasq_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsmasqRangeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsmasqRangeResourceConfig("192.0.2.100", "192.0.2.120", 3600, "initial range"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_range.test", "start_address", "192.0.2.100"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_range.test", "end_address", "192.0.2.120"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_range.test", "lease_time", "3600"),
					resource.TestCheckResourceAttrSet("opnsense_dnsmasq_range.test", "id"),
				),
			},
			{ResourceName: "opnsense_dnsmasq_range.test", ImportState: true, ImportStateVerify: true},
			{
				Config: testAccDnsmasqRangeResourceConfig("192.0.2.130", "192.0.2.150", 7200, "updated range"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_range.test", "start_address", "192.0.2.130"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_range.test", "end_address", "192.0.2.150"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_range.test", "lease_time", "7200"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_range.test", "description", "updated range"),
				),
			},
		},
	})
}

func testAccDnsmasqRangeResourceConfig(startAddress, endAddress string, leaseTime int, description string) string {
	return fmt.Sprintf(`
resource "opnsense_dnsmasq_range" "test" {
  start_address = %q
  end_address   = %q
  lease_time    = %d
  domain_type   = "range"
  description   = %q
}
`, startAddress, endAddress, leaseTime, description)
}
