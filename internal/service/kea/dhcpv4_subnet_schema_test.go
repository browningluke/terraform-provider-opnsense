package kea

import (
	"testing"

	apiKea "github.com/browningluke/opnsense-go/pkg/kea"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestDhcpv4SubnetValidLifetimeSchemaToStruct(t *testing.T) {
	model, err := convertDhcpv4SubnetStructToSchema(&apiKea.SubnetV4{})
	require.NoError(t, err)
	model.ValidLifetime = types.Int64Value(86400)

	result, err := convertDhcpv4SubnetSchemaToStruct(model)
	require.NoError(t, err)
	require.Equal(t, "86400", result.ValidLifetime)
}

func TestDhcpv4SubnetValidLifetimeSchemaToStructUnset(t *testing.T) {
	model, err := convertDhcpv4SubnetStructToSchema(&apiKea.SubnetV4{})
	require.NoError(t, err)

	result, err := convertDhcpv4SubnetSchemaToStruct(model)
	require.NoError(t, err)
	require.Empty(t, result.ValidLifetime)
}

func TestDhcpv4SubnetValidLifetimeStructToSchema(t *testing.T) {
	model, err := convertDhcpv4SubnetStructToSchema(&apiKea.SubnetV4{ValidLifetime: "86400"})
	require.NoError(t, err)
	require.Equal(t, int64(86400), model.ValidLifetime.ValueInt64())

	model, err = convertDhcpv4SubnetStructToSchema(&apiKea.SubnetV4{})
	require.NoError(t, err)
	require.True(t, model.ValidLifetime.IsNull())
}
