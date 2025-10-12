package provider_test

import (
	"context"
	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/browningluke/terraform-provider-opnsense/internal/provider"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestProvider(t *testing.T) {
	opnsense, err := provider.NewProvider(context.Background())
	if opnsense == nil || err != nil {
		t.Fatal("provider.New() returned nil")
	}
}

func TestProvider_Configure(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
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
