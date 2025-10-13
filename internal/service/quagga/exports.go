package quagga

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newBGPASPathResource,
		newBGPCommunityListResource,
		newBGPNeighborResource,
		newBGPPrefixListResource,
		newBGPRouteMapResource,
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newBGPASPathDataSource,
		newBGPCommunityListDataSource,
		newBGPNeighborDataSource,
		newBGPPrefixListDataSource,
		newBGPRouteMapDataSource,
	}
}
