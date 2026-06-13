package trust_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrustCertResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTrustCertResourceConfig("Test Cert", "server.example.com", "server.example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_trust_cert.test", "description", "Test Cert"),
					resource.TestCheckResourceAttr("opnsense_trust_cert.test", "action", "internal"),
					resource.TestCheckResourceAttr("opnsense_trust_cert.test", "common_name", "server.example.com"),
					resource.TestCheckResourceAttr("opnsense_trust_cert.test", "cert_type", "server_cert"),
					resource.TestCheckResourceAttr("opnsense_trust_cert.test", "altnames_dns", "server.example.com"),
					resource.TestCheckResourceAttrSet("opnsense_trust_cert.test", "id"),
					resource.TestCheckResourceAttrSet("opnsense_trust_cert.test", "ref_id"),
					resource.TestCheckResourceAttrSet("opnsense_trust_cert.test", "crt"),
					resource.TestCheckResourceAttrSet("opnsense_trust_cert.test", "crt_payload"),
					resource.TestCheckResourceAttrSet("opnsense_trust_cert.test", "valid_to"),
				),
			},
			{
				ResourceName:            "opnsense_trust_cert.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"action", "prv", "prv_payload"},
			},
			{
				Config: testAccTrustCertResourceConfig("Test Cert Updated", "server.example.com", "server.example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_trust_cert.test", "description", "Test Cert Updated"),
				),
			},
		},
	})
}

func testAccTrustCertResourceConfig(description, commonName, altnamesDns string) string {
	return fmt.Sprintf(`
resource "opnsense_trust_ca" "test_ca" {
  description = "Test CA for Cert"
  action      = "internal"
  key_type    = "2048"
  digest      = "sha256"
  lifetime    = "3650"
  country     = "US"
  common_name = "test-ca.example.com"
}

resource "opnsense_trust_cert" "test" {
  description   = %[1]q
  action        = "internal"
  caref         = opnsense_trust_ca.test_ca.ref_id
  key_type      = "2048"
  digest        = "sha256"
  cert_type     = "server_cert"
  lifetime      = "397"
  country       = "US"
  common_name   = %[2]q
  altnames_dns  = %[3]q
}
`, description, commonName, altnamesDns)
}
