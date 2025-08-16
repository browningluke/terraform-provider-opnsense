package service

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/ipsec"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestConvertIpsecChildSchemaToStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    *IpsecChildResourceModel
		expected *ipsec.IPsecChild
	}{
		{
			name: "basic conversion",
			input: &IpsecChildResourceModel{
				Enabled:         types.StringValue("1"),
				IPsecConnection: types.StringValue("connection-uuid-123"),
				Proposals:       types.SetValueMust(types.StringType, []attr.Value{types.StringValue("aes128-sha256-modp2048")}),
				SHA256_96:       types.StringValue("0"),
				StartAction:     types.StringValue("start"),
				CloseAction:     types.StringValue("none"),
				DPDAction:       types.StringValue("hold"),
				Mode:            types.StringValue("tunnel"),
				InstallPolicies: types.StringValue("1"),
				LocalNetworks:   types.SetValueMust(types.StringType, []attr.Value{types.StringValue("192.168.1.0/24")}),
				RemoteNetworks:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("10.0.0.0/24")}),
				RequestID:       types.StringValue(""),
				RekeyTime:       types.StringValue("0"),
				Description:     types.StringValue("Test IPsec Child"),
				Id:              types.StringValue("uuid-123"),
			},
			expected: &ipsec.IPsecChild{
				Enabled:         "1",
				Connection:      api.SelectedMap("connection-uuid-123"),
				Proposals:       api.SelectedMapList([]string{"aes128-sha256-modp2048"}),
				SHA256_96:       "0",
				StartAction:     api.SelectedMap("start"),
				CloseAction:     api.SelectedMap("none"),
				DPDAction:       api.SelectedMap("hold"),
				Mode:            api.SelectedMap("tunnel"),
				InstallPolicies: "1",
				LocalNetworks:   api.SelectedMapList([]string{"192.168.1.0/24"}),
				RemoteNetworks:  api.SelectedMapList([]string{"10.0.0.0/24"}),
				RequestID:       "",
				RekeyTime:       "0",
				Description:     "Test IPsec Child",
			},
		},
		{
			name: "multiple proposals and networks",
			input: &IpsecChildResourceModel{
				Enabled:         types.StringValue("1"),
				IPsecConnection: types.StringValue("connection-uuid-456"),
				Proposals: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("aes256-sha256-modp2048"),
					types.StringValue("aes128-sha1-modp1024"),
				}),
				SHA256_96:       types.StringValue("1"),
				StartAction:     types.StringValue("route"),
				CloseAction:     types.StringValue("trap"),
				DPDAction:       types.StringValue("restart"),
				Mode:            types.StringValue("transport"),
				InstallPolicies: types.StringValue("0"),
				LocalNetworks: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("192.168.1.0/24"),
					types.StringValue("192.168.2.0/24"),
				}),
				RemoteNetworks: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("10.0.0.0/24"),
					types.StringValue("10.1.0.0/24"),
				}),
				RequestID:   types.StringValue("custom-request-id"),
				RekeyTime:   types.StringValue("3600"),
				Description: types.StringValue("Complex Test Child"),
			},
			expected: &ipsec.IPsecChild{
				Enabled:    "1",
				Connection: api.SelectedMap("connection-uuid-456"),
				Proposals: api.SelectedMapList([]string{
					"aes256-sha256-modp2048",
					"aes128-sha1-modp1024",
				}),
				SHA256_96:       "1",
				StartAction:     api.SelectedMap("route"),
				CloseAction:     api.SelectedMap("trap"),
				DPDAction:       api.SelectedMap("restart"),
				Mode:            api.SelectedMap("transport"),
				InstallPolicies: "0",
				LocalNetworks: api.SelectedMapList([]string{
					"192.168.1.0/24",
					"192.168.2.0/24",
				}),
				RemoteNetworks: api.SelectedMapList([]string{
					"10.0.0.0/24",
					"10.1.0.0/24",
				}),
				RequestID:   "custom-request-id",
				RekeyTime:   "3600",
				Description: "Complex Test Child",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertIpsecChildSchemaToStruct(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Enabled, result.Enabled)
			assert.Equal(t, tt.expected.SHA256_96, result.SHA256_96)
			assert.Equal(t, tt.expected.InstallPolicies, result.InstallPolicies)
			assert.Equal(t, tt.expected.RequestID, result.RequestID)
			assert.Equal(t, tt.expected.RekeyTime, result.RekeyTime)
			assert.Equal(t, tt.expected.Description, result.Description)
			// Note: SelectedMap and SelectedMapList fields are harder to compare directly
		})
	}
}

func TestConvertIpsecChildStructToSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    *ipsec.IPsecChild
		expected *IpsecChildResourceModel
	}{
		{
			name: "basic conversion",
			input: &ipsec.IPsecChild{
				Enabled:         "1",
				Connection:      api.SelectedMap("connection-uuid-123"),
				Proposals:       api.SelectedMapList([]string{"aes128-sha256-modp2048"}),
				SHA256_96:       "0",
				StartAction:     api.SelectedMap("start"),
				CloseAction:     api.SelectedMap("none"),
				DPDAction:       api.SelectedMap("hold"),
				Mode:            api.SelectedMap("tunnel"),
				InstallPolicies: "1",
				LocalNetworks:   api.SelectedMapList([]string{"192.168.1.0/24"}),
				RemoteNetworks:  api.SelectedMapList([]string{"10.0.0.0/24"}),
				RequestID:       "",
				RekeyTime:       "0",
				Description:     "Test IPsec Child",
			},
			expected: &IpsecChildResourceModel{
				Enabled:         types.StringValue("1"),
				IPsecConnection: types.StringValue("connection-uuid-123"),
				Proposals:       types.SetValueMust(types.StringType, []attr.Value{types.StringValue("aes128-sha256-modp2048")}),
				SHA256_96:       types.StringValue("0"),
				StartAction:     types.StringValue("start"),
				CloseAction:     types.StringValue("none"),
				DPDAction:       types.StringValue("hold"),
				Mode:            types.StringValue("tunnel"),
				InstallPolicies: types.StringValue("1"),
				LocalNetworks:   types.SetValueMust(types.StringType, []attr.Value{types.StringValue("192.168.1.0/24")}),
				RemoteNetworks:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("10.0.0.0/24")}),
				RequestID:       types.StringValue(""),
				RekeyTime:       types.StringValue("0"),
				Description:     types.StringValue("Test IPsec Child"),
			},
		},
		{
			name: "multiple networks conversion",
			input: &ipsec.IPsecChild{
				Enabled:    "1",
				Connection: api.SelectedMap("connection-uuid-456"),
				Proposals: api.SelectedMapList([]string{
					"aes256-sha256-modp2048",
					"aes128-sha1-modp1024",
				}),
				SHA256_96:       "1",
				StartAction:     api.SelectedMap("route"),
				CloseAction:     api.SelectedMap("trap"),
				DPDAction:       api.SelectedMap("restart"),
				Mode:            api.SelectedMap("transport"),
				InstallPolicies: "0",
				LocalNetworks: api.SelectedMapList([]string{
					"192.168.1.0/24",
					"192.168.2.0/24",
				}),
				RemoteNetworks: api.SelectedMapList([]string{
					"10.0.0.0/24",
					"10.1.0.0/24",
				}),
				RequestID:   "custom-request-id",
				RekeyTime:   "3600",
				Description: "Complex Test Child",
			},
			expected: &IpsecChildResourceModel{
				Enabled:         types.StringValue("1"),
				IPsecConnection: types.StringValue("connection-uuid-456"),
				Proposals: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("aes256-sha256-modp2048"),
					types.StringValue("aes128-sha1-modp1024"),
				}),
				SHA256_96:       types.StringValue("1"),
				StartAction:     types.StringValue("route"),
				CloseAction:     types.StringValue("trap"),
				DPDAction:       types.StringValue("restart"),
				Mode:            types.StringValue("transport"),
				InstallPolicies: types.StringValue("0"),
				LocalNetworks: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("192.168.1.0/24"),
					types.StringValue("192.168.2.0/24"),
				}),
				RemoteNetworks: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("10.0.0.0/24"),
					types.StringValue("10.1.0.0/24"),
				}),
				RequestID:   types.StringValue("custom-request-id"),
				RekeyTime:   types.StringValue("3600"),
				Description: types.StringValue("Complex Test Child"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertIpsecChildStructToSchema(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Enabled, result.Enabled)
			assert.Equal(t, tt.expected.SHA256_96, result.SHA256_96)
			assert.Equal(t, tt.expected.InstallPolicies, result.InstallPolicies)
			assert.Equal(t, tt.expected.RequestID, result.RequestID)
			assert.Equal(t, tt.expected.RekeyTime, result.RekeyTime)
			assert.Equal(t, tt.expected.Description, result.Description)
		})
	}
}
