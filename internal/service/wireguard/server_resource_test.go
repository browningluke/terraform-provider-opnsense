package wireguard_test

import (
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWireguardServerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWireguardServerResourceConfig(
					"test-server",
					"iEn5o1vV+ijYQiWXjV9qprTCeOcx2g2VDmMwLcmp5Vk=",
					"jJd2OYHYUlMZH5OizDKftCoam8dl9BAZhdXNDXs4M0c=",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_wireguard_server.test", "name", "test-server"),
					resource.TestCheckResourceAttr("opnsense_wireguard_server.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_wireguard_server.test", "public_key", "jJd2OYHYUlMZH5OizDKftCoam8dl9BAZhdXNDXs4M0c="),
					resource.TestCheckResourceAttr("opnsense_wireguard_server.test", "private_key", "iEn5o1vV+ijYQiWXjV9qprTCeOcx2g2VDmMwLcmp5Vk="),
					resource.TestCheckResourceAttrSet("opnsense_wireguard_server.test", "id"),
					resource.TestCheckResourceAttrSet("opnsense_wireguard_server.test", "instance"),
				),
			},
			{
				ResourceName:            "opnsense_wireguard_server.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
			{
				Config: testAccWireguardServerResourceConfig(
					"test-server-updated",
					"2JUWD6TiFXO8lfecWgBC5jBOO2saNUWkkUxtRi41dUo=",
					"Ls+7gpPwsSkNUc7kJbPUbAVHOE079+ACZBzQLIe7s0k=",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_wireguard_server.test", "name", "test-server-updated"),
					resource.TestCheckResourceAttr("opnsense_wireguard_server.test", "public_key", "Ls+7gpPwsSkNUc7kJbPUbAVHOE079+ACZBzQLIe7s0k="),
				),
			},
		},
	})
}

func TestAccWireguardServerResource_WithTunnelAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWireguardServerResourceWithTunnelConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_wireguard_server.test", "name", "test-tunnel-server"),
					resource.TestCheckResourceAttr("opnsense_wireguard_server.test", "tunnel_address.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_wireguard_server.test", "tunnel_address.*", "10.0.0.1/24"),
					resource.TestCheckResourceAttr("opnsense_wireguard_server.test", "port", "51820"),
				),
			},
		},
	})
}

func testAccWireguardServerResourceConfig(name, privkey, pubkey string) string {
	return `
resource "opnsense_wireguard_server" "test" {
  name        = "` + name + `"
  private_key = "` + privkey + `"
  public_key  = "` + pubkey + `"
}
`
}

func testAccWireguardServerResourceWithTunnelConfig() string {
	return `
resource "opnsense_wireguard_server" "test" {
  name           = "test-tunnel-server"
  private_key    = "iD8W5IupKNnbt18tFt9PxpTLlBeBu+yxEUseLQT3OVQ="
  public_key     = "N6z94WfQPp2uLMs2uPA2CASRZBdPXUCJcZyrFufTkVE="
  tunnel_address = ["10.0.0.1/24"]
  port           = 51820
}
`
}
