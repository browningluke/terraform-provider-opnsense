package trust_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrustCaResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTrustCaResourceConfig("Test Internal CA", "internal", "US", "test.example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_trust_ca.test", "description", "Test Internal CA"),
					resource.TestCheckResourceAttr("opnsense_trust_ca.test", "action", "internal"),
					resource.TestCheckResourceAttr("opnsense_trust_ca.test", "common_name", "test.example.com"),
					resource.TestCheckResourceAttr("opnsense_trust_ca.test", "key_type", "2048"),
					resource.TestCheckResourceAttr("opnsense_trust_ca.test", "digest", "sha256"),
					resource.TestCheckResourceAttrSet("opnsense_trust_ca.test", "id"),
					resource.TestCheckResourceAttrSet("opnsense_trust_ca.test", "ref_id"),
					resource.TestCheckResourceAttrSet("opnsense_trust_ca.test", "crt"),
					resource.TestCheckResourceAttrSet("opnsense_trust_ca.test", "crt_payload"),
					resource.TestCheckResourceAttrSet("opnsense_trust_ca.test", "valid_to"),
				),
			},
			{
				ResourceName:            "opnsense_trust_ca.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"action", "lifetime", "prv", "prv_payload"},
			},
			{
				Config: testAccTrustCaResourceConfig("Test Internal CA Updated", "internal", "US", "test.example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_trust_ca.test", "description", "Test Internal CA Updated"),
				),
			},
		},
	})
}

func testAccTrustCaResourceConfig(description, action, country, commonName string) string {
	return fmt.Sprintf(`
resource "opnsense_trust_ca" "test" {
  description = %[1]q
  action      = %[2]q
  key_type    = "2048"
  digest      = "sha256"
  lifetime    = "3650"
  country     = %[3]q
  common_name = %[4]q
}
`, description, action, country, commonName)
}
