package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &aliasResource{} },
		func() resource.Resource { return &categoryResource{} },
		func() resource.Resource { return &filterResource{} },
		func() resource.Resource { return &natResource{} },
		func() resource.Resource { return &natOneToOneResource{} },
		func() resource.Resource { return &natPortForwardResource{} },
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return &aliasDataSource{} },
		func() datasource.DataSource { return &categoryDataSource{} },
		func() datasource.DataSource { return &filterDataSource{} },
		func() datasource.DataSource { return &natDataSource{} },
		func() datasource.DataSource { return &natOneToOneDataSource{} },
		func() datasource.DataSource { return &natPortForwardDataSource{} },
	}
}
