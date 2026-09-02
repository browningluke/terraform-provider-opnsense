package dnsmasq_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsmasqTagResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDnsmasqTagResourceConfig("myTagValue"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_tag.test", "tag", "myTagValue"),
					resource.TestCheckResourceAttrSet("opnsense_dnsmasq_tag.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_dnsmasq_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDnsmasqTagResourceConfig("changeTagValue"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_tag.test", "tag", "changeTagValue"),
				),
			},
			// Delete testing occurs automatically
		},
	})
}

func testAccDnsmasqTagResourceConfig(
	name string,
) string {
	return fmt.Sprintf(`
resource "opnsense_dnsmasq_tag" "test" {
	tag       = "%s"
}
`, name)
}
