package haproxy

import (
	"context"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/haproxy"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAclValueFollowsExpression(t *testing.T) {
	for _, expr := range aclExpressions {
		t.Run(expr, func(t *testing.T) {
			acl, err := convertAclSchemaToStruct(&aclResourceModel{
				Name:       types.StringValue("acl"),
				Expression: types.StringValue(expr),
				Value:      types.StringValue("v-" + expr),
			})
			require.NoError(t, err)
			assert.Equal(t, "v-"+expr, *aclValueField(acl, expr))

			back, err := convertAclStructToSchema(acl)
			require.NoError(t, err)
			assert.Equal(t, expr, back.Expression.ValueString())
			assert.Equal(t, "v-"+expr, back.Value.ValueString())
		})
	}
}

func TestAclHostAndPath(t *testing.T) {
	acl, err := convertAclSchemaToStruct(&aclResourceModel{
		Expression: types.StringValue("hdr_end"),
		Value:      types.StringValue(".example.com"),
	})
	require.NoError(t, err)
	assert.Equal(t, ".example.com", acl.HdrEnd)
	assert.Empty(t, acl.PathBeg, "only the field of the chosen expression is set")
}

func TestAclUnsupportedExpressionOnRead(t *testing.T) {
	_, err := convertAclStructToSchema(&haproxy.Acl{Expression: api.SelectedMap("http_method")})
	assert.ErrorContains(t, err, "http_method")
}

func TestFrontendActionsKeepOrder(t *testing.T) {
	ctx := context.Background()
	order := []string{"6b1a0c6e-0000-4000-8000-00000000000b", "6b1a0c6e-0000-4000-8000-00000000000a"}

	list, diags := types.ListValueFrom(ctx, types.StringType, order)
	require.False(t, diags.HasError())
	fe, diags := convertFrontendActionsSchemaToStruct(ctx, &frontendActionsResourceModel{ActionIds: list})
	require.False(t, diags.HasError())
	assert.Equal(t, api.SelectedMapListOrdered(order), fe.LinkedActions)

	model, diags := convertFrontendActionsStructToSchema(ctx, "fe-id", fe)
	require.False(t, diags.HasError())
	var got []string
	require.False(t, model.ActionIds.ElementsAs(ctx, &got, false).HasError())
	assert.Equal(t, order, got)
	assert.Equal(t, "fe-id", model.FrontendId.ValueString())
}

func TestServerPortSentinel(t *testing.T) {
	s, err := convertServerSchemaToStruct(&serverResourceModel{Port: types.Int64Value(-1)})
	require.NoError(t, err)
	assert.Equal(t, "", s.Port)

	m, err := convertServerStructToSchema(&haproxy.Server{Port: ""})
	require.NoError(t, err)
	assert.Equal(t, int64(-1), m.Port.ValueInt64())
}
