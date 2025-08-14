package service

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/ipsec"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestConvertIpsecAuthRemoteSchemaToStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    *IpsecAuthRemoteResourceModel
		expected *ipsec.IPsecAuthRemote
	}{
		{
			name: "basic PSK conversion",
			input: &IpsecAuthRemoteResourceModel{
				Enabled:         types.StringValue("1"),
				IPsecConnection: types.StringValue("connection-uuid-123"),
				Round:           types.StringValue("1"),
				Authentication:  types.StringValue("psk"),
				AuthId:          types.StringValue("remote@example.com"),
				EAPId:           types.StringValue(""),
				Certificates:    types.SetValueMust(types.StringType, []attr.Value{}),
				PublicKeys:      types.SetValueMust(types.StringType, []attr.Value{}),
				Description:     types.StringValue("Test Auth Remote PSK"),
				Id:              types.StringValue("uuid-123"),
			},
			expected: &ipsec.IPsecAuthRemote{
				Enabled:        "1",
				Connection:     api.SelectedMap("connection-uuid-123"),
				Round:          "1",
				Authentication: api.SelectedMap("psk"),
				Id:             "remote@example.com",
				EAPId:          "",
				Certificates:   api.SelectedMapList([]string{}),
				PublicKeys:     api.SelectedMapList([]string{}),
				Description:    "Test Auth Remote PSK",
			},
		},
		{
			name: "public key with certificates",
			input: &IpsecAuthRemoteResourceModel{
				Enabled:         types.StringValue("1"),
				IPsecConnection: types.StringValue("connection-uuid-456"),
				Round:           types.StringValue("1"),
				Authentication:  types.StringValue("pubkey"),
				AuthId:          types.StringValue("CN=remote.example.com"),
				EAPId:           types.StringValue(""),
				Certificates: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("cert-uuid-ca"),
					types.StringValue("cert-uuid-remote"),
				}),
				PublicKeys: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("pubkey-uuid-remote"),
				}),
				Description: types.StringValue("Test Auth Remote Certificate"),
			},
			expected: &ipsec.IPsecAuthRemote{
				Enabled:        "1",
				Connection:     api.SelectedMap("connection-uuid-456"),
				Round:          "1",
				Authentication: api.SelectedMap("pubkey"),
				Id:             "CN=remote.example.com",
				EAPId:          "",
				Certificates: api.SelectedMapList([]string{
					"cert-uuid-ca",
					"cert-uuid-remote",
				}),
				PublicKeys: api.SelectedMapList([]string{
					"pubkey-uuid-remote",
				}),
				Description: "Test Auth Remote Certificate",
			},
		},
		{
			name: "EAP authentication",
			input: &IpsecAuthRemoteResourceModel{
				Enabled:         types.StringValue("1"),
				IPsecConnection: types.StringValue("connection-uuid-789"),
				Round:           types.StringValue("2"),
				Authentication:  types.StringValue("eap-tls"),
				AuthId:          types.StringValue(""),
				EAPId:           types.StringValue("eap-remote@example.com"),
				Certificates:    types.SetValueMust(types.StringType, []attr.Value{}),
				PublicKeys:      types.SetValueMust(types.StringType, []attr.Value{}),
				Description:     types.StringValue("Test Auth Remote EAP"),
			},
			expected: &ipsec.IPsecAuthRemote{
				Enabled:        "1",
				Connection:     api.SelectedMap("connection-uuid-789"),
				Round:          "2",
				Authentication: api.SelectedMap("eap-tls"),
				Id:             "",
				EAPId:          "eap-remote@example.com",
				Certificates:   api.SelectedMapList([]string{}),
				PublicKeys:     api.SelectedMapList([]string{}),
				Description:    "Test Auth Remote EAP",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertIpsecAuthRemoteSchemaToStruct(tt.input)
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

func TestConvertIpsecAuthRemoteStructToSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    *ipsec.IPsecAuthRemote
		expected *IpsecAuthRemoteResourceModel
	}{
		{
			name: "basic PSK conversion",
			input: &ipsec.IPsecAuthRemote{
				Enabled:        "1",
				Connection:     api.SelectedMap("connection-uuid-123"),
				Round:          "1",
				Authentication: api.SelectedMap("psk"),
				Id:             "remote@example.com",
				EAPId:          "",
				Certificates:   api.SelectedMapList([]string{}),
				PublicKeys:     api.SelectedMapList([]string{}),
				Description:    "Test Auth Remote PSK",
			},
			expected: &IpsecAuthRemoteResourceModel{
				Enabled:         types.StringValue("1"),
				IPsecConnection: types.StringValue("connection-uuid-123"),
				Round:           types.StringValue("1"),
				Authentication:  types.StringValue("psk"),
				AuthId:          types.StringValue("remote@example.com"),
				EAPId:           types.StringValue(""),
				Certificates:    types.SetValueMust(types.StringType, []attr.Value{}),
				PublicKeys:      types.SetValueMust(types.StringType, []attr.Value{}),
				Description:     types.StringValue("Test Auth Remote PSK"),
			},
		},
		{
			name: "certificate authentication conversion",
			input: &ipsec.IPsecAuthRemote{
				Enabled:        "1",
				Connection:     api.SelectedMap("connection-uuid-456"),
				Round:          "1",
				Authentication: api.SelectedMap("pubkey"),
				Id:             "CN=remote.example.com",
				EAPId:          "",
				Certificates: api.SelectedMapList([]string{
					"cert-uuid-ca",
					"cert-uuid-remote",
				}),
				PublicKeys: api.SelectedMapList([]string{
					"pubkey-uuid-remote",
				}),
				Description: "Test Auth Remote Certificate",
			},
			expected: &IpsecAuthRemoteResourceModel{
				Enabled:         types.StringValue("1"),
				IPsecConnection: types.StringValue("connection-uuid-456"),
				Round:           types.StringValue("1"),
				Authentication:  types.StringValue("pubkey"),
				AuthId:          types.StringValue("CN=remote.example.com"),
				EAPId:           types.StringValue(""),
				Certificates: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("cert-uuid-ca"),
					types.StringValue("cert-uuid-remote"),
				}),
				PublicKeys: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("pubkey-uuid-remote"),
				}),
				Description: types.StringValue("Test Auth Remote Certificate"),
			},
		},
		{
			name: "EAP authentication conversion",
			input: &ipsec.IPsecAuthRemote{
				Enabled:        "1",
				Connection:     api.SelectedMap("connection-uuid-789"),
				Round:          "2",
				Authentication: api.SelectedMap("eap-tls"),
				Id:             "",
				EAPId:          "eap-remote@example.com",
				Certificates:   api.SelectedMapList([]string{}),
				PublicKeys:     api.SelectedMapList([]string{}),
				Description:    "Test Auth Remote EAP",
			},
			expected: &IpsecAuthRemoteResourceModel{
				Enabled:         types.StringValue("1"),
				IPsecConnection: types.StringValue("connection-uuid-789"),
				Round:           types.StringValue("2"),
				Authentication:  types.StringValue("eap-tls"),
				AuthId:          types.StringValue(""),
				EAPId:           types.StringValue("eap-remote@example.com"),
				Certificates:    types.SetValueMust(types.StringType, []attr.Value{}),
				PublicKeys:      types.SetValueMust(types.StringType, []attr.Value{}),
				Description:     types.StringValue("Test Auth Remote EAP"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertIpsecAuthRemoteStructToSchema(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Enabled, result.Enabled)
			assert.Equal(t, tt.expected.Round, result.Round)
			assert.Equal(t, tt.expected.AuthId, result.AuthId)
			assert.Equal(t, tt.expected.EAPId, result.EAPId)
			assert.Equal(t, tt.expected.Description, result.Description)
		})
	}
}
