package interfaces

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newBridgeResource,
		newVipResource,
		newVlanResource,
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newBridgeDataSource,
		newVipDataSource,
		newVlanDataSource,
		newOverviewInterfaceDataSource,
		newOverviewAllDataSource,
	}
}
