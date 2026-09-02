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

func TestConvertOptionSchemaToStruct(t *testing.T) {
	tags := types.SetValueMust(types.StringType, []attr.Value{types.StringValue("iot"), types.StringValue("lan")})
	option, err := convertOptionSchemaToStruct(&optionResourceModel{
		Type:        types.StringValue("set"),
		OptionV4:    types.Int64Value(6),
		Interface:   types.StringValue("lan"),
		TypeSetTags: tags,
		Value:       types.StringValue("8.8.8.8"),
		Force:       types.BoolValue(true),
		Description: types.StringValue("DNS server"),
	})
	require.NoError(t, err)
	assert.Equal(t, "set", option.Type.String())
	assert.Equal(t, "6", option.OptionV4.String())
	assert.Equal(t, "", option.OptionV6.String())
	assert.Equal(t, "lan", option.Interface.String())
	assert.Equal(t, []string{"iot", "lan"}, []string(option.TypeSetTags))
	assert.Equal(t, "8.8.8.8", option.Value)
	assert.Equal(t, "1", option.Force)
	assert.Equal(t, "DNS server", option.Description)
}

func TestConvertOptionStructToSchema(t *testing.T) {
	model, err := convertOptionStructToSchema(&upstream.Option{
		Type:         api.SelectedMap("match"),
		OptionV6:     api.SelectedMap("23"),
		TypeMatchTag: api.SelectedMap("pxe"),
		Value:        "",
		Description:  "classify PXE",
	})
	require.NoError(t, err)
	assert.Equal(t, types.StringValue("match"), model.Type)
	assert.True(t, model.OptionV4.IsNull())
	assert.Equal(t, types.Int64Value(23), model.OptionV6)
	assert.True(t, model.Interface.IsNull())
	assert.Equal(t, types.StringValue("pxe"), model.TypeMatchTag)
	assert.True(t, model.Value.IsNull())
	assert.True(t, model.Force.IsNull())
	assert.Equal(t, types.StringValue("classify PXE"), model.Description)
	assert.True(t, model.TypeSetTags.IsNull())
}
