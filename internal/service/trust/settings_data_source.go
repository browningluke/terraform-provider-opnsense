package trust

import (
	"github.com/browningluke/opnsense-go/pkg/trust"
	"github.com/browningluke/terraform-provider-opnsense/internal/provider/generic"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var genericSettingsDataSource = &generic.GenericDataSource[trust.TrustSettings, settingsResourceModel]{
	Definition: settingsDefinition,
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = genericSettingsDataSource
var _ datasource.DataSourceWithConfigure = genericSettingsDataSource

func newSettingsDataSource() datasource.DataSource {
	return genericSettingsDataSource
}
