package quagga

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &bgpASPathResource{} },
		func() resource.Resource { return &bgpCommunityListResource{} },
		func() resource.Resource { return &bgpNeighborResource{} },
		func() resource.Resource { return &bgpPrefixListResource{} },
		func() resource.Resource { return &bgpRouteMapResource{} },
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return &bgpASPathDataSource{} },
		func() datasource.DataSource { return &bgpCommunityListDataSource{} },
		func() datasource.DataSource { return &bgpNeighborDataSource{} },
		func() datasource.DataSource { return &bgpPrefixListDataSource{} },
		func() datasource.DataSource { return &bgpRouteMapDataSource{} },
	}
}
