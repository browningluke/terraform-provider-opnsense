package dnsmasq

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewHostResource,
		newDomainResource,
		newOptionResource,
		newRangeResource,
		newTagResource,
		newSettingsResource,
		newBootResource,
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newHostDataSource,
		newDomainDataSource,
		newOptionDataSource,
		newRangeDataSource,
		newTagDataSource,
		newSettingsDataSource,
		newBootDataSource,
	}
}
