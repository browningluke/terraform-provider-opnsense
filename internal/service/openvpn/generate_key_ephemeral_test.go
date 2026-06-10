package openvpn_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// echoFactories pairs the opnsense provider with the upstream "echo" provider
// from terraform-plugin-testing. The echo data source surfaces the otherwise-
// state-less ephemeral value to TestCheckResourceAttr.
var echoFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"echo": echoprovider.NewProviderServer(),
}

func TestAccOpenvpnGenerateKeyEphemeral_default(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_10_0)},
		ProtoV6ProviderFactories: mergeFactories(acctest.ProtoV6ProviderFactories, echoFactories),
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "opnsense_openvpn_generate_key" "test" {}

provider "echo" {
  data = ephemeral.opnsense_openvpn_generate_key.test
}

resource "echo" "test" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// key_type is unset (null) when the user omits it; echo
					// reports null attributes as absent, so we just verify
					// the generated key shape.
					resource.TestMatchResourceAttr("echo.test", "data.key", regexp.MustCompile(`(?s)BEGIN OpenVPN Static key V1.*END OpenVPN Static key V1`)),
				),
			},
		},
	})
}

func TestAccOpenvpnGenerateKeyEphemeral_tlsAuth(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_10_0)},
		ProtoV6ProviderFactories: mergeFactories(acctest.ProtoV6ProviderFactories, echoFactories),
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "opnsense_openvpn_generate_key" "test" {
  key_type = "tls-auth"
}

provider "echo" {
  data = ephemeral.opnsense_openvpn_generate_key.test
}

resource "echo" "test" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test", "data.key_type", "tls-auth"),
					resource.TestMatchResourceAttr("echo.test", "data.key", regexp.MustCompile(`(?s)BEGIN OpenVPN Static key V1.*END OpenVPN Static key V1`)),
				),
			},
		},
	})
}

func TestAccOpenvpnGenerateKeyEphemeral_invalidKeyType(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_10_0)},
		ProtoV6ProviderFactories: mergeFactories(acctest.ProtoV6ProviderFactories, echoFactories),
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "opnsense_openvpn_generate_key" "test" {
  key_type = "totally-not-a-real-type"
}

provider "echo" {
  data = ephemeral.opnsense_openvpn_generate_key.test
}

resource "echo" "test" {}
`,
				ExpectError: regexp.MustCompile(`(?s)key_type.*value must be one of`),
			},
		},
	})
}

func mergeFactories(m ...map[string]func() (tfprotov6.ProviderServer, error)) map[string]func() (tfprotov6.ProviderServer, error) {
	out := map[string]func() (tfprotov6.ProviderServer, error){}
	for _, src := range m {
		for k, v := range src {
			out[k] = v
		}
	}
	return out
}
