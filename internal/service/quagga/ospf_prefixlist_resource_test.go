package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaOSPFPrefixListResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaOSPFPrefixListResourceConfig("TEST_OSPF_PL", "permit", "10", "10.0.0.0/8"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_prefix_list.test", "name", "TEST_OSPF_PL"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_prefix_list.test", "action", "permit"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_prefix_list.test", "network", "10.0.0.0/8"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_prefix_list.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_ospf_prefix_list.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_ospf_prefix_list.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaOSPFPrefixListResourceConfig("TEST_OSPF_PL", "deny", "10", "172.16.0.0/12"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_prefix_list.test", "action", "deny"),
					resource.TestCheckResourceAttr("opnsense_quagga_ospf_prefix_list.test", "network", "172.16.0.0/12"),
				),
			},
		},
	})
}

func testAccQuaggaOSPFPrefixListResourceConfig(name, action, seqNumber, network string) string {
	return `
resource "opnsense_quagga_ospf_prefix_list" "test" {
  name            = "` + name + `"
  action          = "` + action + `"
  sequence_number = "` + seqNumber + `"
  network         = "` + network + `"
}
`
}
