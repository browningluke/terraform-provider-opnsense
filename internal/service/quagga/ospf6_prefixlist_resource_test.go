package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaOSPF6PrefixListResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaOSPF6PrefixListResourceConfig("TEST_OSPF6_PL", "permit", "10", "2001:db8::/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_prefix_list.test", "name", "TEST_OSPF6_PL"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_prefix_list.test", "action", "permit"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_prefix_list.test", "network", "2001:db8::/32"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_prefix_list.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_ospf6_prefix_list.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_ospf6_prefix_list.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaOSPF6PrefixListResourceConfig("TEST_OSPF6_PL", "deny", "10", "2001:db8:1::/48"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_prefix_list.test", "action", "deny"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf6_prefix_list.test", "network", "2001:db8:1::/48"),
				),
			},
		},
	})
}

func testAccQuaggaOSPF6PrefixListResourceConfig(name, action, seqNumber, network string) string {
	return `
resource "opnsense_quagga_ospf6_prefix_list" "test" {
  name            = "` + name + `"
  action          = "` + action + `"
  sequence_number = "` + seqNumber + `"
  network         = "` + network + `"
}
`
}
