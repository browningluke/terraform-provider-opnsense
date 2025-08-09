package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"opnsense": providerserver.NewProtocol6WithError(New("test")()),
}

func TestProvider(t *testing.T) {
	provider := New("test")()
	if provider == nil {
		t.Fatal("provider.New() returned nil")
	}
}

func TestProvider_Configure(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						return nil
					},
				),
			},
		},
	})
}

const testAccProviderConfig = `
provider "opnsense" {
  # Configuration will be loaded from environment variables
  # OPNSENSE_URI, OPNSENSE_API_KEY, OPNSENSE_API_SECRET
}
`
