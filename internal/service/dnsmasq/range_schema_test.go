package dnsmasq

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	upstream "github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertRangeSchemaToStruct(t *testing.T) {
	model := &rangeResourceModel{
		Interface:     types.StringValue("lan"),
		StartAddress:  types.StringValue("192.168.10.100"),
		EndAddress:    types.StringValue("192.168.10.200"),
		LeaseTime:     types.Int64Value(3600),
		DomainType:    types.StringValue("range"),
		DisableHASync: types.BoolValue(true),
		Mode:          types.SetValueMust(types.StringType, []attr.Value{types.StringValue("static")}),
	}
	rangeObj, err := convertRangeSchemaToStruct(model)
	require.NoError(t, err)
	assert.Equal(t, "lan", rangeObj.Interface.String())
	assert.Equal(t, "3600", rangeObj.LeaseTime)
	assert.Equal(t, "range", rangeObj.DomainType.String())
	assert.Equal(t, []string{"static"}, []string(rangeObj.Mode))
	assert.Equal(t, "1", rangeObj.DisableHASync)
}

func TestConvertRangeStructToSchema(t *testing.T) {
	model, err := convertRangeStructToSchema(&upstream.Range{
		Interface:        api.SelectedMap("lan"),
		StartAddress:     "2001:db8::10",
		LeaseTime:        "7200",
		PrefixLength:     "64",
		RaPriority:       api.SelectedMap("high"),
		RaMTU:            "1500",
		RaInterval:       "30",
		RaRouterLifetime: "90",
	})
	require.NoError(t, err)
	assert.Equal(t, types.StringValue("lan"), model.Interface)
	assert.Equal(t, types.Int64Value(7200), model.LeaseTime)
	assert.Equal(t, types.Int64Value(64), model.PrefixLength)
	assert.Equal(t, types.StringValue("high"), model.RaPriority)
	assert.Equal(t, types.Int64Value(1500), model.RaMTU)
	assert.Equal(t, types.Int64Value(30), model.RaInterval)
	assert.Equal(t, types.Int64Value(90), model.RaRouterLifetime)
}

func TestConvertRangeStructToSchemaTreatsUnsetPrefixLengthAsNull(t *testing.T) {
	model, err := convertRangeStructToSchema(&upstream.Range{
		StartAddress: "192.0.2.100",
		PrefixLength: "-1",
	})
	require.NoError(t, err)
	assert.True(t, model.PrefixLength.IsNull())
}
