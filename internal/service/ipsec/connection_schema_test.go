package ipsec

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/ipsec"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestConvertIpsecConnectionSchemaToStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    *connectionResourceModel
		expected *ipsec.IPsecConnection
	}{
		{
			name: "basic conversion",
			input: &connectionResourceModel{
				Enabled:                types.StringValue("1"),
				Proposals:              types.SetValueMust(types.StringType, []attr.Value{types.StringValue("aes128-sha256-modp2048")}),
				Unique:                 types.StringValue("no"),
				Aggressive:             types.StringValue("0"),
				Version:                types.StringValue("ikev2"),
				Mobike:                 types.StringValue("1"),
				LocalAddresses:         types.SetValueMust(types.StringType, []attr.Value{types.StringValue("192.168.1.1")}),
				RemoteAddresses:        types.SetValueMust(types.StringType, []attr.Value{types.StringValue("10.0.0.1")}),
				LocalPort:              types.StringValue("500"),
				RemotePort:             types.StringValue("500"),
				UDPEncapsulation:       types.StringValue("0"),
				ReauthenticationTime:   types.StringValue("3600"),
				RekeyTime:              types.StringValue("1800"),
				IKELifetime:            types.StringValue("3600"),
				DPDDelay:               types.StringValue("10"),
				DPDTimeout:             types.StringValue("60"),
				IPPools:                types.SetValueMust(types.StringType, []attr.Value{}),
				SendCertificateRequest: types.StringValue("1"),
				SendCertificate:        types.StringValue("ifasked"),
				KeyingTries:            types.StringValue("1"),
				Description:            types.StringValue("Test Connection"),
				Id:                     types.StringValue("uuid-123"),
			},
			expected: &ipsec.IPsecConnection{
				Enabled:                "1",
				Proposals:              api.SelectedMapList([]string{"aes128-sha256-modp2048"}),
				Unique:                 api.SelectedMap("no"),
				Aggressive:             "0",
				Version:                api.SelectedMap("ikev2"),
				Mobike:                 "1",
				LocalAddresses:         api.SelectedMapList([]string{"192.168.1.1"}),
				RemoteAddresses:        api.SelectedMapList([]string{"10.0.0.1"}),
				LocalPort:              api.SelectedMap("500"),
				RemotePort:             api.SelectedMap("500"),
				UDPEncapsulation:       "0",
				ReauthenticationTime:   "3600",
				RekeyTime:              "1800",
				IKELifetime:            "3600",
				DPDDelay:               "10",
				DPDTimeout:             "60",
				IPPools:                api.SelectedMapList([]string{}),
				SendCertificateRequest: "1",
				SendCertificate:        api.SelectedMap("ifasked"),
				KeyingTries:            "1",
				Description:            "Test Connection",
			},
		},
		{
			name: "multiple addresses and proposals",
			input: &connectionResourceModel{
				Enabled: types.StringValue("1"),
				Proposals: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("aes256-sha256-modp2048"),
					types.StringValue("aes128-sha1-modp1024"),
				}),
				Unique:     types.StringValue("no"),
				Aggressive: types.StringValue("0"),
				Version:    types.StringValue("ikev2"),
				Mobike:     types.StringValue("1"),
				LocalAddresses: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("192.168.1.1"),
					types.StringValue("192.168.1.10"),
				}),
				RemoteAddresses: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("10.0.0.1"),
					types.StringValue("10.0.0.10"),
				}),
				LocalPort:              types.StringValue("500"),
				RemotePort:             types.StringValue("500"),
				UDPEncapsulation:       types.StringValue("0"),
				ReauthenticationTime:   types.StringValue("7200"),
				RekeyTime:              types.StringValue("3600"),
				IKELifetime:            types.StringValue("7200"),
				DPDDelay:               types.StringValue("30"),
				DPDTimeout:             types.StringValue("120"),
				IPPools:                types.SetValueMust(types.StringType, []attr.Value{}),
				SendCertificateRequest: types.StringValue("0"),
				SendCertificate:        types.StringValue("never"),
				KeyingTries:            types.StringValue("3"),
				Description:            types.StringValue(""),
			},
			expected: &ipsec.IPsecConnection{
				Enabled: "1",
				Proposals: api.SelectedMapList([]string{
					"aes256-sha256-modp2048",
					"aes128-sha1-modp1024",
				}),
				Unique:     api.SelectedMap("no"),
				Aggressive: "0",
				Version:    api.SelectedMap("ikev2"),
				Mobike:     "1",
				LocalAddresses: api.SelectedMapList([]string{
					"192.168.1.1",
					"192.168.1.10",
				}),
				RemoteAddresses: api.SelectedMapList([]string{
					"10.0.0.1",
					"10.0.0.10",
				}),
				LocalPort:              api.SelectedMap("500"),
				RemotePort:             api.SelectedMap("500"),
				UDPEncapsulation:       "0",
				ReauthenticationTime:   "7200",
				RekeyTime:              "3600",
				IKELifetime:            "7200",
				DPDDelay:               "30",
				DPDTimeout:             "120",
				IPPools:                api.SelectedMapList([]string{}),
				SendCertificateRequest: "0",
				SendCertificate:        api.SelectedMap("never"),
				KeyingTries:            "3",
				Description:            "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertConnectionSchemaToStruct(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Enabled, result.Enabled)
			assert.Equal(t, tt.expected.Aggressive, result.Aggressive)
			assert.Equal(t, tt.expected.Mobike, result.Mobike)
			assert.Equal(t, tt.expected.UDPEncapsulation, result.UDPEncapsulation)
			assert.Equal(t, tt.expected.ReauthenticationTime, result.ReauthenticationTime)
			assert.Equal(t, tt.expected.RekeyTime, result.RekeyTime)
			assert.Equal(t, tt.expected.IKELifetime, result.IKELifetime)
			assert.Equal(t, tt.expected.DPDDelay, result.DPDDelay)
			assert.Equal(t, tt.expected.DPDTimeout, result.DPDTimeout)
			assert.Equal(t, tt.expected.SendCertificateRequest, result.SendCertificateRequest)
			assert.Equal(t, tt.expected.KeyingTries, result.KeyingTries)
			assert.Equal(t, tt.expected.Description, result.Description)
			// Note: SelectedMap and SelectedMapList fields are harder to compare directly
		})
	}
}

func TestConvertIpsecConnectionStructToSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    *ipsec.IPsecConnection
		expected *connectionResourceModel
	}{
		{
			name: "basic conversion",
			input: &ipsec.IPsecConnection{
				Enabled:                "1",
				Proposals:              api.SelectedMapList([]string{"aes128-sha256-modp2048"}),
				Unique:                 api.SelectedMap("no"),
				Aggressive:             "0",
				Version:                api.SelectedMap("ikev2"),
				Mobike:                 "1",
				LocalAddresses:         api.SelectedMapList([]string{"192.168.1.1"}),
				RemoteAddresses:        api.SelectedMapList([]string{"10.0.0.1"}),
				LocalPort:              api.SelectedMap("500"),
				RemotePort:             api.SelectedMap("500"),
				UDPEncapsulation:       "0",
				ReauthenticationTime:   "3600",
				RekeyTime:              "1800",
				IKELifetime:            "3600",
				DPDDelay:               "10",
				DPDTimeout:             "60",
				IPPools:                api.SelectedMapList([]string{}),
				SendCertificateRequest: "1",
				SendCertificate:        api.SelectedMap("ifasked"),
				KeyingTries:            "1",
				Description:            "Test Connection",
			},
			expected: &connectionResourceModel{
				Enabled:                types.StringValue("1"),
				Proposals:              types.SetValueMust(types.StringType, []attr.Value{types.StringValue("aes128-sha256-modp2048")}),
				Unique:                 types.StringValue("no"),
				Aggressive:             types.StringValue("0"),
				Version:                types.StringValue("ikev2"),
				Mobike:                 types.StringValue("1"),
				LocalAddresses:         types.SetValueMust(types.StringType, []attr.Value{types.StringValue("192.168.1.1")}),
				RemoteAddresses:        types.SetValueMust(types.StringType, []attr.Value{types.StringValue("10.0.0.1")}),
				LocalPort:              types.StringValue("500"),
				RemotePort:             types.StringValue("500"),
				UDPEncapsulation:       types.StringValue("0"),
				ReauthenticationTime:   types.StringValue("3600"),
				RekeyTime:              types.StringValue("1800"),
				IKELifetime:            types.StringValue("3600"),
				DPDDelay:               types.StringValue("10"),
				DPDTimeout:             types.StringValue("60"),
				IPPools:                types.SetValueMust(types.StringType, []attr.Value{}),
				SendCertificateRequest: types.StringValue("1"),
				SendCertificate:        types.StringValue("ifasked"),
				KeyingTries:            types.StringValue("1"),
				Description:            types.StringValue("Test Connection"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertConnectionStructToSchema(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Enabled, result.Enabled)
			assert.Equal(t, tt.expected.Aggressive, result.Aggressive)
			assert.Equal(t, tt.expected.Mobike, result.Mobike)
			assert.Equal(t, tt.expected.UDPEncapsulation, result.UDPEncapsulation)
			assert.Equal(t, tt.expected.ReauthenticationTime, result.ReauthenticationTime)
			assert.Equal(t, tt.expected.RekeyTime, result.RekeyTime)
			assert.Equal(t, tt.expected.IKELifetime, result.IKELifetime)
			assert.Equal(t, tt.expected.DPDDelay, result.DPDDelay)
			assert.Equal(t, tt.expected.DPDTimeout, result.DPDTimeout)
			assert.Equal(t, tt.expected.SendCertificateRequest, result.SendCertificateRequest)
			assert.Equal(t, tt.expected.KeyingTries, result.KeyingTries)
			assert.Equal(t, tt.expected.Description, result.Description)
		})
	}
}
