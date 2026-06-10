package openvpn_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const exampleStaticKey = `#
# 2048 bit OpenVPN static key
#
-----BEGIN OpenVPN Static key V1-----
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
1112233344455566778899aabbccddeeff
-----END OpenVPN Static key V1-----`

func TestAccOpenvpnStaticKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticKeyResourceConfig("acctest-key", "auth"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_openvpn_static_key.test", "description", "acctest-key"),
					resource.TestCheckResourceAttr("opnsense_openvpn_static_key.test", "mode", "auth"),
					resource.TestCheckResourceAttrSet("opnsense_openvpn_static_key.test", "id"),
				),
			},
			{
				ResourceName:            "opnsense_openvpn_static_key.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key"},
			},
			{
				Config: testAccStaticKeyResourceConfig("acctest-key-upd", "crypt"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_openvpn_static_key.test", "description", "acctest-key-upd"),
					resource.TestCheckResourceAttr("opnsense_openvpn_static_key.test", "mode", "crypt"),
				),
			},
		},
	})
}

func testAccStaticKeyResourceConfig(description, mode string) string {
	return fmt.Sprintf(`
resource "opnsense_openvpn_static_key" "test" {
  description = %[1]q
  mode        = %[2]q
  key         = <<-EOT
%[3]s
EOT
}
`, description, mode, exampleStaticKey)
}
