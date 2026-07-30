package firewall

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	opnfirewall "github.com/browningluke/opnsense-go/pkg/firewall"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestConvertAliasSchemaToStructIgnoresIPProtocolForHostAlias(t *testing.T) {
	data := &aliasResourceModel{
		Enabled:    types.BoolValue(true),
		Name:       types.StringValue("testalias"),
		Type:       types.StringValue("host"),
		IPProtocol: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("IPv4")}),
		Interface:  types.StringValue(""),
		Content: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("192.168.1.100"),
			types.StringValue("192.168.1.101"),
		}),
		Categories:  types.SetValueMust(types.StringType, []attr.Value{}),
		UpdateFreq:  types.Float64Value(-1),
		Statistics:  types.BoolValue(false),
		Description: types.StringNull(),
	}

	result, err := convertAliasSchemaToStruct(data)

	require.NoError(t, err)
	require.Equal(t, "host", result.Type.String())
	require.Empty(t, result.IPProtocol)
}

func TestConvertAliasSchemaToStructDefaultsIPProtocolForApplicableAlias(t *testing.T) {
	data := &aliasResourceModel{
		Enabled:     types.BoolValue(true),
		Name:        types.StringValue("externalalias"),
		Type:        types.StringValue("external"),
		IPProtocol:  types.SetNull(types.StringType),
		Interface:   types.StringValue(""),
		Content:     types.SetValueMust(types.StringType, []attr.Value{}),
		Categories:  types.SetValueMust(types.StringType, []attr.Value{}),
		UpdateFreq:  types.Float64Value(-1),
		Statistics:  types.BoolValue(false),
		Description: types.StringNull(),
	}

	result, err := convertAliasSchemaToStruct(data)

	require.NoError(t, err)
	require.Equal(t, "external", result.Type.String())
	require.Equal(t, api.SelectedMapList{"IPv4"}, result.IPProtocol)
}

func TestConvertAliasStructToSchemaIgnoresIPProtocolForHostAlias(t *testing.T) {
	result, err := convertAliasStructToSchema(&opnfirewall.Alias{
		Enabled:    "1",
		Name:       "testalias",
		Type:       api.SelectedMap("host"),
		IPProtocol: api.SelectedMapList{"IPv4"},
		Interface:  api.SelectedMap(""),
		Content: api.SelectedMapListNL{
			"192.168.1.100",
			"192.168.1.101",
		},
		Categories: []string{},
		UpdateFreq: "-1",
		Statistics: "0",
	})

	require.NoError(t, err)
	require.Equal(t, "host", result.Type.ValueString())
	require.Empty(t, result.IPProtocol.Elements())
}

func TestConvertAliasStructToSchemaDefaultsIPProtocolForApplicableAlias(t *testing.T) {
	result, err := convertAliasStructToSchema(&opnfirewall.Alias{
		Enabled:    "1",
		Name:       "externalalias",
		Type:       api.SelectedMap("external"),
		IPProtocol: api.SelectedMapList{},
		Interface:  api.SelectedMap(""),
		Content:    api.SelectedMapListNL{},
		Categories: []string{},
		UpdateFreq: "-1",
		Statistics: "0",
	})

	require.NoError(t, err)
	require.Equal(t, "external", result.Type.ValueString())
	require.Len(t, result.IPProtocol.Elements(), 1)
	require.Contains(t, result.IPProtocol.Elements(), types.StringValue("IPv4"))
}
