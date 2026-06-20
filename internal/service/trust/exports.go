package trust

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &caResource{} },
		func() resource.Resource { return &certResource{} },
		func() resource.Resource { return &settingsResource{} },
	}
}

func DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return &caDataSource{} },
		func() datasource.DataSource { return &certDataSource{} },
		func() datasource.DataSource { return &settingsDataSource{} },
	}
}
