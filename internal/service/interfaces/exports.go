package interfaces

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &vipResource{} },
		func() resource.Resource { return &vlanResource{} },
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return &vipDataSource{} },
		func() datasource.DataSource { return &vlanDataSource{} },
		func() datasource.DataSource { return &overviewInterfaceDataSource{} },
		func() datasource.DataSource { return &overviewAllDataSource{} },
	}
}
