package trust

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newCaResource,
		newCertResource,
		newSettingsResource,
	}
}

func DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newCaDataSource,
		newCertDataSource,
		newSettingsDataSource,
	}
}
