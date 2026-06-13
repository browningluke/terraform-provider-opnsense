package wireguard_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWireguardClientResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWireguardClientResourceConfig(
					"test-client",
					"ng28gpVw6F9YYsXTdPtfglAsGU1ET+ePrTFDfOzGniw=",
					"10.0.1.2/32",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_wireguard_client.test", "name", "test-client"),
					resource.TestCheckResourceAttr("opnsense_wireguard_client.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_wireguard_client.test", "public_key", "ng28gpVw6F9YYsXTdPtfglAsGU1ET+ePrTFDfOzGniw="),
					resource.TestCheckResourceAttr("opnsense_wireguard_client.test", "tunnel_address.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_wireguard_client.test", "tunnel_address.*", "10.0.1.2/32"),
					resource.TestCheckResourceAttrSet("opnsense_wireguard_client.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_wireguard_client.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWireguardClientResourceConfig(
					"test-client-updated",
					"Ls+7gpPwsSkNUc7kJbPUbAVHOE079+ACZBzQLIe7s0k=",
					"10.0.1.3/32",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_wireguard_client.test", "name", "test-client-updated"),
					resource.TestCheckResourceAttr("opnsense_wireguard_client.test", "public_key", "Ls+7gpPwsSkNUc7kJbPUbAVHOE079+ACZBzQLIe7s0k="),
					resource.TestCheckTypeSetElemAttr("opnsense_wireguard_client.test", "tunnel_address.*", "10.0.1.3/32"),
				),
			},
		},
	})
}

func TestAccWireguardClientResource_WithPSK(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWireguardClientResourceWithPSKConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_wireguard_client.test", "name", "test-client-psk"),
					resource.TestCheckResourceAttr("opnsense_wireguard_client.test", "psk", "C+S3B796lTBat+4NnTUHHWk1yvA+0753a0im4U7iNaw="),
				),
			},
		},
	})
}

func testAccWireguardClientResourceConfig(name, pubkey, tunnelAddr string) string {
	return `
resource "opnsense_wireguard_client" "test" {
  name           = "` + name + `"
  public_key     = "` + pubkey + `"
  tunnel_address = ["` + tunnelAddr + `"]
}
`
}

func testAccWireguardClientResourceWithPSKConfig() string {
	return `
resource "opnsense_wireguard_client" "test" {
  name           = "test-client-psk"
  public_key     = "jJd2OYHYUlMZH5OizDKftCoam8dl9BAZhdXNDXs4M0c="
  psk            = "C+S3B796lTBat+4NnTUHHWk1yvA+0753a0im4U7iNaw="
  tunnel_address = ["10.0.1.4/32"]
}
`
}
