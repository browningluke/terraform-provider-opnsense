package service

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/ipsec"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestConvertIpsecAuthLocalSchemaToStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    *IpsecAuthLocalResourceModel
		expected *ipsec.IPsecAuthLocal
	}{
		{
			name: "basic PSK conversion",
			input: &IpsecAuthLocalResourceModel{
				Enabled:        types.StringValue("1"),
				Connection:     types.StringValue("connection-uuid-123"),
				Round:          types.StringValue("1"),
				Authentication: types.StringValue("psk"),
				AuthId:         types.StringValue("local@example.com"),
				EAPId:          types.StringValue(""),
				Certificates:   types.SetValueMust(types.StringType, []attr.Value{}),
				PublicKeys:     types.SetValueMust(types.StringType, []attr.Value{}),
				Description:    types.StringValue("Test Auth Local PSK"),
				Id:             types.StringValue("uuid-123"),
			},
			expected: &ipsec.IPsecAuthLocal{
				Enabled:        "1",
				Connection:     api.SelectedMap("connection-uuid-123"),
				Round:          "1",
				Authentication: api.SelectedMap("psk"),
				Id:             "local@example.com",
				EAPId:          "",
				Certificates:   api.SelectedMapList([]string{}),
				PublicKeys:     api.SelectedMapList([]string{}),
				Description:    "Test Auth Local PSK",
			},
		},
		{
			name: "public key with certificates",
			input: &IpsecAuthLocalResourceModel{
				Enabled:        types.StringValue("1"),
				Connection:     types.StringValue("connection-uuid-456"),
				Round:          types.StringValue("1"),
				Authentication: types.StringValue("pubkey"),
				AuthId:         types.StringValue("CN=local.example.com"),
				EAPId:          types.StringValue(""),
				Certificates: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("cert-uuid-1"),
					types.StringValue("cert-uuid-2"),
				}),
				PublicKeys: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("pubkey-uuid-1"),
				}),
				Description: types.StringValue("Test Auth Local Certificate"),
			},
			expected: &ipsec.IPsecAuthLocal{
				Enabled:        "1",
				Connection:     api.SelectedMap("connection-uuid-456"),
				Round:          "1",
				Authentication: api.SelectedMap("pubkey"),
				Id:             "CN=local.example.com",
				EAPId:          "",
				Certificates: api.SelectedMapList([]string{
					"cert-uuid-1",
					"cert-uuid-2",
				}),
				PublicKeys: api.SelectedMapList([]string{
					"pubkey-uuid-1",
				}),
				Description: "Test Auth Local Certificate",
			},
		},
		{
			name: "EAP authentication",
			input: &IpsecAuthLocalResourceModel{
				Enabled:        types.StringValue("1"),
				Connection:     types.StringValue("connection-uuid-789"),
				Round:          types.StringValue("2"),
				Authentication: types.StringValue("eap-radius"),
				AuthId:         types.StringValue(""),
				EAPId:          types.StringValue("eap-user@example.com"),
				Certificates:   types.SetValueMust(types.StringType, []attr.Value{}),
				PublicKeys:     types.SetValueMust(types.StringType, []attr.Value{}),
				Description:    types.StringValue("Test Auth Local EAP"),
			},
			expected: &ipsec.IPsecAuthLocal{
				Enabled:        "1",
				Connection:     api.SelectedMap("connection-uuid-789"),
				Round:          "2",
				Authentication: api.SelectedMap("eap-radius"),
				Id:             "",
				EAPId:          "eap-user@example.com",
				Certificates:   api.SelectedMapList([]string{}),
				PublicKeys:     api.SelectedMapList([]string{}),
				Description:    "Test Auth Local EAP",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertIpsecAuthLocalSchemaToStruct(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Enabled, result.Enabled)
			assert.Equal(t, tt.expected.Round, result.Round)
			assert.Equal(t, tt.expected.Id, result.Id)
			assert.Equal(t, tt.expected.EAPId, result.EAPId)
			assert.Equal(t, tt.expected.Description, result.Description)
			// Note: SelectedMap and SelectedMapList fields are harder to compare directly
		})
	}
}

func TestConvertIpsecAuthLocalStructToSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    *ipsec.IPsecAuthLocal
		expected *IpsecAuthLocalResourceModel
	}{
		{
			name: "basic PSK conversion",
			input: &ipsec.IPsecAuthLocal{
				Enabled:        "1",
				Connection:     api.SelectedMap("connection-uuid-123"),
				Round:          "1",
				Authentication: api.SelectedMap("psk"),
				Id:             "local@example.com",
				EAPId:          "",
				Certificates:   api.SelectedMapList([]string{}),
				PublicKeys:     api.SelectedMapList([]string{}),
				Description:    "Test Auth Local PSK",
			},
			expected: &IpsecAuthLocalResourceModel{
				Enabled:        types.StringValue("1"),
				Connection:     types.StringValue("connection-uuid-123"),
				Round:          types.StringValue("1"),
				Authentication: types.StringValue("psk"),
				AuthId:         types.StringValue("local@example.com"),
				EAPId:          types.StringValue(""),
				Certificates:   types.SetValueMust(types.StringType, []attr.Value{}),
				PublicKeys:     types.SetValueMust(types.StringType, []attr.Value{}),
				Description:    types.StringValue("Test Auth Local PSK"),
			},
		},
		{
			name: "certificate authentication conversion",
			input: &ipsec.IPsecAuthLocal{
				Enabled:        "1",
				Connection:     api.SelectedMap("connection-uuid-456"),
				Round:          "1",
				Authentication: api.SelectedMap("pubkey"),
				Id:             "CN=local.example.com",
				EAPId:          "",
				Certificates: api.SelectedMapList([]string{
					"cert-uuid-1",
					"cert-uuid-2",
				}),
				PublicKeys: api.SelectedMapList([]string{
					"pubkey-uuid-1",
				}),
				Description: "Test Auth Local Certificate",
			},
			expected: &IpsecAuthLocalResourceModel{
				Enabled:        types.StringValue("1"),
				Connection:     types.StringValue("connection-uuid-456"),
				Round:          types.StringValue("1"),
				Authentication: types.StringValue("pubkey"),
				AuthId:         types.StringValue("CN=local.example.com"),
				EAPId:          types.StringValue(""),
				Certificates: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("cert-uuid-1"),
					types.StringValue("cert-uuid-2"),
				}),
				PublicKeys: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("pubkey-uuid-1"),
				}),
				Description: types.StringValue("Test Auth Local Certificate"),
			},
		},
		{
			name: "EAP authentication conversion",
			input: &ipsec.IPsecAuthLocal{
				Enabled:        "1",
				Connection:     api.SelectedMap("connection-uuid-789"),
				Round:          "2",
				Authentication: api.SelectedMap("eap-radius"),
				Id:             "",
				EAPId:          "eap-user@example.com",
				Certificates:   api.SelectedMapList([]string{}),
				PublicKeys:     api.SelectedMapList([]string{}),
				Description:    "Test Auth Local EAP",
			},
			expected: &IpsecAuthLocalResourceModel{
				Enabled:        types.StringValue("1"),
				Connection:     types.StringValue("connection-uuid-789"),
				Round:          types.StringValue("2"),
				Authentication: types.StringValue("eap-radius"),
				AuthId:         types.StringValue(""),
				EAPId:          types.StringValue("eap-user@example.com"),
				Certificates:   types.SetValueMust(types.StringType, []attr.Value{}),
				PublicKeys:     types.SetValueMust(types.StringType, []attr.Value{}),
				Description:    types.StringValue("Test Auth Local EAP"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertIpsecAuthLocalStructToSchema(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Enabled, result.Enabled)
			assert.Equal(t, tt.expected.Round, result.Round)
			assert.Equal(t, tt.expected.AuthId, result.AuthId)
			assert.Equal(t, tt.expected.EAPId, result.EAPId)
			assert.Equal(t, tt.expected.Description, result.Description)
		})
	}
}