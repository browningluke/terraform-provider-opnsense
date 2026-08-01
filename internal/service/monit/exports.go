package monit

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newAlertResource,
		newServiceResource,
		newSettingsResource,
	}
}

func DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newAlertDataSource,
		newServiceDataSource,
		newSettingsDataSource,
	}
}
