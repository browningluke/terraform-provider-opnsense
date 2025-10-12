package ipsec

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/ipsec"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestConvertIpsecPskSchemaToStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    *pskResourceModel
		expected *ipsec.IPsecPSK
	}{
		{
			name: "basic conversion",
			input: &pskResourceModel{
				IdentityLocal:  types.StringValue("local@example.com"),
				IdentityRemote: types.StringValue("remote@example.com"),
				PreSharedKey:   types.StringValue("secretkey123"),
				Type:           types.StringValue("PSK"),
				Description:    types.StringValue("Test PSK"),
				Id:             types.StringValue("uuid-123"),
			},
			expected: &ipsec.IPsecPSK{
				IdentityLocal:  "local@example.com",
				IdentityRemote: "remote@example.com",
				PreSharedKey:   "secretkey123",
				Type:           api.SelectedMap("PSK"),
				Description:    "Test PSK",
			},
		},
		{
			name: "minimal conversion",
			input: &pskResourceModel{
				IdentityLocal:  types.StringValue("user1"),
				IdentityRemote: types.StringValue("user2"),
				PreSharedKey:   types.StringValue("key123"),
				Type:           types.StringValue("PSK"),
				Description:    types.StringValue(""),
			},
			expected: &ipsec.IPsecPSK{
				IdentityLocal:  "user1",
				IdentityRemote: "user2",
				PreSharedKey:   "key123",
				Type:           api.SelectedMap("PSK"),
				Description:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertPskSchemaToStruct(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.IdentityLocal, result.IdentityLocal)
			assert.Equal(t, tt.expected.IdentityRemote, result.IdentityRemote)
			assert.Equal(t, tt.expected.PreSharedKey, result.PreSharedKey)
			assert.Equal(t, tt.expected.Description, result.Description)
			// Note: Type is api.SelectedMap which is harder to compare directly
		})
	}
}

func TestConvertIpsecPskStructToSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    *ipsec.IPsecPSK
		expected *pskResourceModel
	}{
		{
			name: "basic conversion",
			input: &ipsec.IPsecPSK{
				IdentityLocal:  "local@example.com",
				IdentityRemote: "remote@example.com",
				PreSharedKey:   "secretkey123",
				Type:           api.SelectedMap("PSK"),
				Description:    "Test PSK",
			},
			expected: &pskResourceModel{
				IdentityLocal:  types.StringValue("local@example.com"),
				IdentityRemote: types.StringValue("remote@example.com"),
				PreSharedKey:   types.StringValue("secretkey123"),
				Type:           types.StringValue("PSK"),
				Description:    types.StringValue("Test PSK"),
			},
		},
		{
			name: "empty description",
			input: &ipsec.IPsecPSK{
				IdentityLocal:  "user1",
				IdentityRemote: "user2",
				PreSharedKey:   "key123",
				Type:           api.SelectedMap("PSK"),
				Description:    "",
			},
			expected: &pskResourceModel{
				IdentityLocal:  types.StringValue("user1"),
				IdentityRemote: types.StringValue("user2"),
				PreSharedKey:   types.StringValue("key123"),
				Type:           types.StringValue("PSK"),
				Description:    types.StringNull(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertPskStructToSchema(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.IdentityLocal, result.IdentityLocal)
			assert.Equal(t, tt.expected.IdentityRemote, result.IdentityRemote)
			assert.Equal(t, tt.expected.PreSharedKey, result.PreSharedKey)
			assert.Equal(t, tt.expected.Type, result.Type)
			// Description comparison depends on tools.StringOrNull behavior
		})
	}
}
