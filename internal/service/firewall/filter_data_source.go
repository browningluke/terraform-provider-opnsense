package firewall

import (
	"github.com/browningluke/opnsense-go/pkg/firewall"
	"github.com/browningluke/terraform-provider-opnsense/internal/provider/generic"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var genericFilterDataSource = &generic.GenericDataSource[firewall.Filter, filterResourceModel]{
	Definition: filterDefinition,
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = genericFilterDataSource
var _ datasource.DataSourceWithConfigure = genericFilterDataSource

func newFilterDataSource() datasource.DataSource {
	return genericFilterDataSource
}
