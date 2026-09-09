package dnsmasq

import (
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/provider/generic"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var genericHostDataSource = &generic.GenericDataSource[dnsmasq.Host, hostResourceModel]{
	Definition: hostDefinition,
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = genericHostDataSource
var _ datasource.DataSourceWithConfigure = genericHostDataSource

func newClientDataSource() datasource.DataSource {
	return genericHostDataSource
}
