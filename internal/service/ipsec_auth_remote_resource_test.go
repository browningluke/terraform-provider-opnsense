package service

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIpsecAuthRemoteResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIpsecAuthRemoteResourceConfig(
					"1",                              // enabled
					"connection-uuid-123",            // connection
					"1",                              // round
					"pubkey",                         // authentication
					"remote@example.com",             // auth_id
					"",                               // eap_id (empty)
					[]string{"cert-uuid-1"},          // certificates
					[]string{"pubkey-uuid-1"},        // public_keys
					"Test IPsec Auth Remote",         // description
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "connection", "connection-uuid-123"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "round", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "authentication", "pubkey"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "auth_id", "remote@example.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "eap_id", ""),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "certificates.#", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "certificates.0", "cert-uuid-1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "public_keys.#", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "public_keys.0", "pubkey-uuid-1"),
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
				Config: testAccIpsecAuthRemoteResourceConfig(
					"1",                                                // enabled
					"connection-uuid-123",                              // connection
					"2",                                                // round - updated
					"psk",                                              // authentication - updated
					"updated-remote@example.com",                       // auth_id - updated
					"eap-remote@example.com",                           // eap_id - updated
					[]string{"cert-uuid-1", "cert-uuid-2"},            // certificates - updated
					[]string{"pubkey-uuid-1", "pubkey-uuid-2"},        // public_keys - updated
					"Updated Test IPsec Auth Remote",                   // description - updated
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "round", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "authentication", "psk"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "auth_id", "updated-remote@example.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "eap_id", "eap-remote@example.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "certificates.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "public_keys.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "description", "Updated Test IPsec Auth Remote"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIpsecAuthRemoteResource_MinimalConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIpsecAuthRemoteResourceConfigMinimal(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "round", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "authentication", "psk"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "auth_id", ""),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "eap_id", ""),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "certificates.#", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "public_keys.#", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "description", ""),
					resource.TestCheckResourceAttrSet("opnsense_ipsec_auth_remote.test", "id"),
				),
			},
		},
	})
}

func TestAccIpsecAuthRemoteResource_CertificateAuth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIpsecAuthRemoteResourceConfigCertificate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "authentication", "pubkey"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "auth_id", "CN=remote.example.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "certificates.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_remote.test", "description", "Certificate Authentication Test"),
				),
			},
		},
	})
}

func testAccIpsecAuthRemoteResourceConfig(
	enabled string,
	connection string,
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
resource "opnsense_ipsec_auth_remote" "test" {
  enabled        = %[1]q
  connection     = %[2]q
  round          = %[3]q
  authentication = %[4]q
  auth_id        = %[5]q
  eap_id         = %[6]q
%[7]s%[8]s  description    = %[9]q
}
`, enabled, connection, round, authentication, authID, eapID, 
   certificatesLine, publicKeysLine, description)
}

func testAccIpsecAuthRemoteResourceConfigMinimal() string {
	return `
resource "opnsense_ipsec_auth_remote" "test" {
  enabled        = "1"
  connection     = "connection-uuid-minimal"
  authentication = "psk"
}
`
}

func testAccIpsecAuthRemoteResourceConfigCertificate() string {
	return `
resource "opnsense_ipsec_auth_remote" "test" {
  enabled        = "1"
  connection     = "connection-uuid-cert"
  round          = "1"
  authentication = "pubkey"
  auth_id        = "CN=remote.example.com"
  certificates   = ["cert-uuid-ca", "cert-uuid-remote"]
  description    = "Certificate Authentication Test"
}
`
}