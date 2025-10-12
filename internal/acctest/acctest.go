package acctest

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	ProviderName = "opnsense"
)

var (
	ProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error) = protoV6ProviderFactoriesInit(context.Background(), ProviderName)
)

func protoV6ProviderFactoriesInit(ctx context.Context, providerNames ...string) map[string]func() (tfprotov6.ProviderServer, error) {
	factories := make(map[string]func() (tfprotov6.ProviderServer, error))

	for _, name := range providerNames {
		if name == ProviderName {
			serverFactory, _, err := provider.ProtoV6ProviderServerFactory(ctx)
			if err != nil {
				log.Fatal(err)
			}
			factories[name] = func() (tfprotov6.ProviderServer, error) {
				return serverFactory(), nil
			}
		}
	}

	return factories
}

func AccPreCheck(t *testing.T) {
	if v := os.Getenv("OPNSENSE_URI"); v == "" {
		t.Fatal("OPNSENSE_URI must be set for acceptance tests")
	}
	if v := os.Getenv("OPNSENSE_API_KEY"); v == "" {
		t.Fatal("OPNSENSE_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("OPNSENSE_API_SECRET"); v == "" {
		t.Fatal("OPNSENSE_API_SECRET must be set for acceptance tests")
	}
}
