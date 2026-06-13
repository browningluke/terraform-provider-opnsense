package quagga_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuaggaBGPCommunityListResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccQuaggaBGPCommunityListResourceConfig("permit", "65000:100"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_communitylist.test", "number", "10"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_communitylist.test", "action", "permit"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_communitylist.test", "community", "65000:100"),
					resource.TestCheckResourceAttrSet("opnsense_quagga_bgp_communitylist.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_quagga_bgp_communitylist.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuaggaBGPCommunityListResourceConfig("deny", "65000:200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_communitylist.test", "action", "deny"),
					resource.TestCheckResourceAttr("opnsense_quagga_bgp_communitylist.test", "community", "65000:200"),
				),
			},
		},
	})
}

func testAccQuaggaBGPCommunityListResourceConfig(action, community string) string {
	return `
resource "opnsense_quagga_bgp_communitylist" "test" {
  number     = 10
  seq_number = 10
  action     = "` + action + `"
  community  = "` + community + `"
}
`
}
