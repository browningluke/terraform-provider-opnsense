package auth_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAuthUserResource(t *testing.T) {
	name := "tfacc-auth-user"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthUserResourceConfig(name, "acc user", "/bin/sh"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_auth_user.test", "name", name),
					resource.TestCheckResourceAttr("opnsense_auth_user.test", "full_name", "acc user"),
					resource.TestCheckResourceAttr("opnsense_auth_user.test", "shell", "/bin/sh"),
					resource.TestCheckResourceAttr("opnsense_auth_user.test", "scramble_password", "true"),
					resource.TestCheckResourceAttr("opnsense_auth_user.test", "groups.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_auth_user.test", "groups.*", "admins"),
					resource.TestCheckResourceAttrSet("opnsense_auth_user.test", "id"),
				),
			},
			{
				ResourceName:            "opnsense_auth_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "scramble_password"},
			},
			{
				Config: testAccAuthUserResourceConfig(name, "acc user updated", "/bin/sh"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_auth_user.test", "full_name", "acc user updated"),
					resource.TestCheckResourceAttr("opnsense_auth_user.test", "scramble_password", "true"),
				),
			},
		},
	})
}

func testAccAuthUserResourceConfig(name, fullName, shell string) string {
	return fmt.Sprintf(`
resource "opnsense_auth_user" "test" {
  name              = %q
  full_name         = %q
  shell             = %q
  scramble_password = true
  groups            = ["admins"]
}
`, name, fullName, shell)
}
