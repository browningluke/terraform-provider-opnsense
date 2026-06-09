package openvpn

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newInstanceResource,
		newStaticKeyResource,
		newClientOverwriteResource,
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newInstanceDataSource,
		newStaticKeyDataSource,
		newClientOverwriteDataSource,
	}
}
