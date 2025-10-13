package ipsec_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIpsecPskResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPskResourceConfig("user1@example.com", "peer1@example.com", "supersecretkey123", "PSK", "Test PSK for VPN"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "identity_local", "user1@example.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "identity_remote", "peer1@example.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "pre_shared_key", "supersecretkey123"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "type", "PSK"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "description", "Test PSK for VPN"),
					resource.TestCheckResourceAttrSet("opnsense_ipsec_psk.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_ipsec_psk.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccPskResourceConfig("user1@example.com", "peer1@example.com", "newsupersecretkey456", "PSK", "Updated PSK for VPN"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "identity_local", "user1@example.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "identity_remote", "peer1@example.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "pre_shared_key", "newsupersecretkey456"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "type", "PSK"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "description", "Updated PSK for VPN"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIpsecPskResource_MinimalConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPskResourceConfigMinimal("testuser", "testpeer", "minimalsecretkey"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "identity_local", "testuser"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "identity_remote", "testpeer"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "pre_shared_key", "minimalsecretkey"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "type", "PSK"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "description", ""),
					resource.TestCheckResourceAttrSet("opnsense_ipsec_psk.test", "id"),
				),
			},
		},
	})
}

func TestAccIpsecPskResource_IPAddresses(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPskResourceConfig("192.168.1.10", "10.0.0.5", "ipbasedpskkey", "PSK", "IP-based PSK"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "identity_local", "192.168.1.10"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "identity_remote", "10.0.0.5"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "pre_shared_key", "ipbasedpskkey"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "type", "PSK"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "description", "IP-based PSK"),
				),
			},
		},
	})
}

func TestAccIpsecPskResource_LongKey(t *testing.T) {
	longKey := "verylongsecretkeywiththisislongerthanusualbutshouldbefinetohandle123456789"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPskResourceConfig("local@domain.com", "remote@domain.com", longKey, "PSK", "Long key test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "identity_local", "local@domain.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "identity_remote", "remote@domain.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "pre_shared_key", longKey),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "type", "PSK"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "description", "Long key test"),
				),
			},
		},
	})
}

func TestAccIpsecPskResource_SpecialCharacters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPskResourceConfig("useralphanumeric123", "peeralphanumeric456", "key!@#$%^&*()_+-=", "PSK", "Special chars test - & more"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "identity_local", "useralphanumeric123"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "identity_remote", "peeralphanumeric456"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "pre_shared_key", "key!@#$%^&*()_+-="),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "type", "PSK"),
					resource.TestCheckResourceAttr("opnsense_ipsec_psk.test", "description", "Special chars test - & more"),
				),
			},
		},
	})
}

func testAccPskResourceConfig(identityLocal, identityRemote, preSharedKey, pskType, description string) string {
	return fmt.Sprintf(`
resource "opnsense_ipsec_psk" "test" {
  identity_local  = %[1]q
  identity_remote = %[2]q
  pre_shared_key  = %[3]q
  type           = %[4]q
  description    = %[5]q
}
`, identityLocal, identityRemote, preSharedKey, pskType, description)
}

func testAccPskResourceConfigMinimal(identityLocal, identityRemote, preSharedKey string) string {
	return fmt.Sprintf(`
resource "opnsense_ipsec_psk" "test" {
  identity_local  = %[1]q
  identity_remote = %[2]q
  pre_shared_key  = %[3]q
}
`, identityLocal, identityRemote, preSharedKey)
}
