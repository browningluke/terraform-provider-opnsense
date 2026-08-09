package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newAliasResource,
		newCategoryResource,
		newFilterResource,
		newGeoIPSettingsResource,
		newNATResource,
		newNATOneToOneResource,
		newNATPortForwardResource,
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newAliasDataSource,
		newCategoryDataSource,
		newFilterDataSource,
		newGeoIPSettingsDataSource,
		newNATDataSource,
		newNATOneToOneDataSource,
		newNATPortForwardDataSource,
	}
}
