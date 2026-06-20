package unbound

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &aclResource{} },
		func() resource.Resource { return &domainOverrideResource{} },
		func() resource.Resource { return &forwardResource{} },
		func() resource.Resource { return &hostAliasResource{} },
		func() resource.Resource { return &hostOverrideResource{} },
		func() resource.Resource { return &settingsResource{} },
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return &aclDataSource{} },
		func() datasource.DataSource { return &domainOverrideDataSource{} },
		func() datasource.DataSource { return &forwardDataSource{} },
		func() datasource.DataSource { return &hostAliasDataSource{} },
		func() datasource.DataSource { return &hostOverrideDataSource{} },
		func() datasource.DataSource { return &settingsDataSource{} },
	}
}
