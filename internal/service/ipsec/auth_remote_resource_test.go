package ipsec_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIpsecAuthRemoteResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAuthRemoteResourceConfig(
					"1",                      // enabled
					"0",                      // round
					"psk",                    // authentication
					"remote@example.com",     // auth_id
					"",                       // eap_id (empty)
					[]string{},               // certificates
					[]string{},               // public_keys
					"Test IPsec Auth Remote", // description
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "enabled", "1"),
					resource.TestCheckResourceAttrPair("opnsense_ipsec_auth_remote.test", "ipsec_connection", "opnsense_ipsec_connection.parent", "id"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "round", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "authentication", "psk"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "auth_id", "remote@example.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "eap_id", ""),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "certificates.#", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "public_keys.#", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "description", "Test IPsec Auth Remote"),
					resource.TestCheckResourceAttrSet("opnsense_ipsec_auth_remote.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_ipsec_auth_remote.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAuthRemoteResourceConfig(
					"1",                              // enabled
					"0",                              // round - updated
					"psk",                            // authentication - updated
					"updated-remote@example.com",     // auth_id - updated
					"",                               // eap_id - updated
					[]string{},                       // certificates - updated
					[]string{},                       // public_keys - updated
					"Updated Test IPsec Auth Remote", // description - updated
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "enabled", "1"),
					resource.TestCheckResourceAttrPair("opnsense_ipsec_auth_remote.test", "ipsec_connection", "opnsense_ipsec_connection.parent", "id"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "round", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "authentication", "psk"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "auth_id", "updated-remote@example.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "eap_id", ""),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "certificates.#", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "public_keys.#", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "description", "Updated Test IPsec Auth Remote"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAuthRemoteResourceConfig(
	enabled string,
	round string,
	authentication string,
	authID string,
	eapID string,
	certificates []string,
	publicKeys []string,
	description string,
) string {
	certificatesLine := ""
	if len(certificates) > 0 {
		certificatesLine = fmt.Sprintf("  certificates = [\"%s\"]\n", strings.Join(certificates, `", "`))
	}
	publicKeysLine := ""
	if len(publicKeys) > 0 {
		publicKeysLine = fmt.Sprintf("  public_keys  = [\"%s\"]\n", strings.Join(publicKeys, `", "`))
	}

	return fmt.Sprintf(`
resource "opnsense_ipsec_connection" "parent" {
  enabled                  = "1"
  proposals                = ["aes128-sha256-modp2048"]
  unique                   = "no"
  aggressive               = "0"
  version                  = "2"
  mobike                   = "1"
  local_addresses          = ["192.168.1.1"]
  remote_addresses         = ["10.0.0.1"]
  local_port               = ""
  remote_port              = ""
  udp_encapsulation        = "0"
  reauthentication_time    = "3600"
  rekey_time               = "1800"
  ike_lifetime             = "3600"
  dpd_delay                = "10"
  dpd_timeout              = "60"
  send_certificate_request = "1"
  send_certificate         = "ifasked"
  keying_tries             = "1"
  description              = "Test IPsec Connection for Child"
}

resource "opnsense_ipsec_auth_remote" "test" {
  enabled          = %[1]q
  ipsec_connection = opnsense_ipsec_connection.parent.id
  round            = %[2]q
  authentication   = %[3]q
  auth_id          = %[4]q
  eap_id           = %[5]q
%[6]s%[7]s  description    = %[8]q
}
`, enabled, round, authentication, authID, eapID,
		certificatesLine, publicKeysLine, description)
}
