package wireguard

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &clientResource{} },
		func() resource.Resource { return &serverResource{} },
		func() resource.Resource { return &settingsResource{} },
	}
}

func DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return &clientDataSource{} },
		func() datasource.DataSource { return &serverDataSource{} },
		func() datasource.DataSource { return &settingsDataSource{} },
	}
}
