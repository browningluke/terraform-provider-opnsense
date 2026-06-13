package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaBGPPrefixListResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaBGPPrefixListResourceConfig("TEST_PL", "IPv4", "permit", "10.0.0.0/8"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_prefixlist.test", "name", "TEST_PL"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_prefixlist.test", "ip_version", "IPv4"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_prefixlist.test", "action", "permit"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_prefixlist.test", "network", "10.0.0.0/8"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_bgp_prefixlist.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_bgp_prefixlist.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaBGPPrefixListResourceConfig("TEST_PL", "IPv4", "deny", "192.168.0.0/16"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_prefixlist.test", "action", "deny"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_prefixlist.test", "network", "192.168.0.0/16"),
				),
			},
		},
	})
}

func testAccQuaggaBGPPrefixListResourceConfig(name, ipVersion, action, network string) string {
	return `
resource "opnsense_quagga_bgp_prefixlist" "test" {
  name       = "` + name + `"
  ip_version = "` + ipVersion + `"
  number     = 10
  action     = "` + action + `"
  network    = "` + network + `"
}
`
}
