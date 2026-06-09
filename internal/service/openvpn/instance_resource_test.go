package openvpn_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOpenvpnInstanceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceResourceConfig(7001, "acctest-server", "10.99.97.0/24", 11951),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_openvpn_instance.test", "description", "acctest-server"),
					resource.TestCheckResourceAttr("opnsense_openvpn_instance.test", "role", "server"),
					resource.TestCheckResourceAttr("opnsense_openvpn_instance.test", "server", "10.99.97.0/24"),
					resource.TestCheckResourceAttr("opnsense_openvpn_instance.test", "port", "11951"),
					resource.TestCheckResourceAttr("opnsense_openvpn_instance.test", "vpn_id", "7001"),
					resource.TestCheckResourceAttrSet("opnsense_openvpn_instance.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_openvpn_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInstanceResourceConfig(7001, "acctest-server-upd", "10.99.97.0/24", 11952),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_openvpn_instance.test", "description", "acctest-server-upd"),
					resource.TestCheckResourceAttr("opnsense_openvpn_instance.test", "port", "11952"),
				),
			},
		},
	})
}

func testAccInstanceResourceConfig(vpnID int, description, server string, port int) string {
	return fmt.Sprintf(`
resource "opnsense_openvpn_instance" "test" {
  vpn_id                  = %[1]d
  description             = %[2]q
  role                    = "server"
  dev_type                = "tun"
  protocol                = "udp"
  topology                = "subnet"
  server                  = %[3]q
  port                    = %[4]d
  verify_client_cert      = "none"
  username_as_common_name = true
  auth_mode               = ["Local Database"]
}
`, vpnID, description, server, port)
}
