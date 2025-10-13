package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func ProtoV6ProviderServerFactory(ctx context.Context) (func() tfprotov6.ProviderServer, provider.Provider, error) {
	opnsense, err := NewProvider(ctx)
	if err != nil {
		return nil, nil, err
	}

	server := providerserver.NewProtocol6(opnsense)

	return server, opnsense, nil
}
