package service

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIpsecAuthLocalResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIpsecAuthLocalResourceConfig(
					"1",                              // enabled
					"connection-uuid-123",            // connection
					"1",                              // round
					"pubkey",                         // authentication
					"local@example.com",              // auth_id
					"",                               // eap_id (empty)
					[]string{"cert-uuid-1"},          // certificates
					[]string{"pubkey-uuid-1"},        // public_keys
					"Test IPsec Auth Local",          // description
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "connection", "connection-uuid-123"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "round", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "authentication", "pubkey"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "auth_id", "local@example.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "eap_id", ""),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "certificates.#", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "certificates.0", "cert-uuid-1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "public_keys.#", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "public_keys.0", "pubkey-uuid-1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "description", "Test IPsec Auth Local"),
					resource.TestCheckResourceAttrSet("opnsense_ipsec_auth_local.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_ipsec_auth_local.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIpsecAuthLocalResourceConfig(
					"1",                                                // enabled
					"connection-uuid-123",                              // connection
					"2",                                                // round - updated
					"psk",                                              // authentication - updated
					"updated-local@example.com",                        // auth_id - updated
					"eap-user@example.com",                             // eap_id - updated
					[]string{"cert-uuid-1", "cert-uuid-2"},            // certificates - updated
					[]string{"pubkey-uuid-1", "pubkey-uuid-2"},        // public_keys - updated
					"Updated Test IPsec Auth Local",                    // description - updated
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "round", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "authentication", "psk"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "auth_id", "updated-local@example.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "eap_id", "eap-user@example.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "certificates.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "public_keys.#", "2"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "description", "Updated Test IPsec Auth Local"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIpsecAuthLocalResource_MinimalConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIpsecAuthLocalResourceConfigMinimal(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "enabled", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "round", "1"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "authentication", "psk"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "auth_id", ""),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "eap_id", ""),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "certificates.#", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "public_keys.#", "0"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "description", ""),
					resource.TestCheckResourceAttrSet("opnsense_ipsec_auth_local.test", "id"),
				),
			},
		},
	})
}

func TestAccIpsecAuthLocalResource_EAPConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIpsecAuthLocalResourceConfigEAP(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "authentication", "eap-radius"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "eap_id", "eap-test@example.com"),
					resource.TestCheckResourceAttr("opnsense_ipsec_auth_local.test", "description", "EAP Authentication Test"),
				),
			},
		},
	})
}

func testAccIpsecAuthLocalResourceConfig(
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
resource "opnsense_ipsec_auth_local" "test" {
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

func testAccIpsecAuthLocalResourceConfigMinimal() string {
	return `
resource "opnsense_ipsec_auth_local" "test" {
  enabled        = "1"
  connection     = "connection-uuid-minimal"
  authentication = "psk"
}
`
}

func testAccIpsecAuthLocalResourceConfigEAP() string {
	return `
resource "opnsense_ipsec_auth_local" "test" {
  enabled        = "1"
  connection     = "connection-uuid-eap"
  round          = "1"
  authentication = "eap-radius"
  eap_id         = "eap-test@example.com"
  description    = "EAP Authentication Test"
}
`
}