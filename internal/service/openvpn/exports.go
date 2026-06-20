package openvpn

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &instanceResource{} },
		func() resource.Resource { return &staticKeyResource{} },
		func() resource.Resource { return &clientOverwriteResource{} },
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return &instanceDataSource{} },
		func() datasource.DataSource { return &staticKeyDataSource{} },
		func() datasource.DataSource { return &clientOverwriteDataSource{} },
	}
}

func EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		func() ephemeral.EphemeralResource { return &generateKeyEphemeral{} },
	}
}
