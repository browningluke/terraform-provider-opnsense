package openvpn_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOpenvpnClientOverwriteResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClientOverwriteResourceConfig("acctest-client", "acctest-cso", "10.50.0.0/24"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_openvpn_client_overwrite.test", "common_name", "acctest-client"),
					resource.TestCheckResourceAttr("opnsense_openvpn_client_overwrite.test", "description", "acctest-cso"),
					resource.TestCheckResourceAttr("opnsense_openvpn_client_overwrite.test", "tunnel_network", "10.50.0.0/24"),
					resource.TestCheckResourceAttr("opnsense_openvpn_client_overwrite.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_openvpn_client_overwrite.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_openvpn_client_overwrite.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccClientOverwriteResourceConfig("acctest-client", "acctest-cso-upd", "10.51.0.0/24"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_openvpn_client_overwrite.test", "description", "acctest-cso-upd"),
					resource.TestCheckResourceAttr("opnsense_openvpn_client_overwrite.test", "tunnel_network", "10.51.0.0/24"),
				),
			},
		},
	})
}

func testAccClientOverwriteResourceConfig(commonName, description, tunnelNetwork string) string {
	return fmt.Sprintf(`
resource "opnsense_openvpn_client_overwrite" "test" {
  common_name    = %[1]q
  description    = %[2]q
  tunnel_network = %[3]q
}
`, commonName, description, tunnelNetwork)
}
