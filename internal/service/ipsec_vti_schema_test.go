package service

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/ipsec"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestConvertIpsecVtiSchemaToStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    *IpsecVtiResourceModel
		expected *ipsec.IPsecVTI
	}{
		{
			name: "basic_conversion",
			input: &IpsecVtiResourceModel{
				Enabled:         types.StringValue("1"),
				RequestID:       types.StringValue("1234"),
				LocalIP:         types.StringValue("192.168.1.10"),
				RemoteIP:        types.StringValue("203.0.113.10"),
				TunnelLocalIP:   types.StringValue("10.0.1.1"),
				TunnelRemoteIP:  types.StringValue("10.0.1.2"),
				TunnelLocalIP2:  types.StringValue("10.0.2.1"),
				TunnelRemoteIP2: types.StringValue("10.0.2.2"),
				Description:     types.StringValue("Test VTI"),
				Id:              types.StringValue("test-id"),
			},
			expected: &ipsec.IPsecVTI{
				Enabled:         "1",
				RequestID:       "1234",
				LocalIP:         "192.168.1.10",
				RemoteIP:        "203.0.113.10",
				TunnelLocalIP:   "10.0.1.1",
				TunnelRemoteIP:  "10.0.1.2",
				TunnelLocalIP2:  "10.0.2.1",
				TunnelRemoteIP2: "10.0.2.2",
				Description:     "Test VTI",
			},
		},
		{
			name: "minimal_required_fields",
			input: &IpsecVtiResourceModel{
				Enabled:         types.StringValue("1"),
				RequestID:       types.StringValue("100"),
				LocalIP:         types.StringValue("172.16.1.1"),
				RemoteIP:        types.StringValue("172.16.2.1"),
				TunnelLocalIP:   types.StringValue("10.100.1.1"),
				TunnelRemoteIP:  types.StringValue("10.100.1.2"),
				TunnelLocalIP2:  types.StringValue(""),
				TunnelRemoteIP2: types.StringValue(""),
				Description:     types.StringValue(""),
				Id:              types.StringValue("minimal-id"),
			},
			expected: &ipsec.IPsecVTI{
				Enabled:         "1",
				RequestID:       "100",
				LocalIP:         "172.16.1.1",
				RemoteIP:        "172.16.2.1",
				TunnelLocalIP:   "10.100.1.1",
				TunnelRemoteIP:  "10.100.1.2",
				TunnelLocalIP2:  "",
				TunnelRemoteIP2: "",
				Description:     "",
			},
		},
		{
			name: "disabled_vti",
			input: &IpsecVtiResourceModel{
				Enabled:         types.StringValue("0"),
				RequestID:       types.StringValue("9999"),
				LocalIP:         types.StringValue("10.10.10.1"),
				RemoteIP:        types.StringValue("20.20.20.1"),
				TunnelLocalIP:   types.StringValue("172.16.10.1"),
				TunnelRemoteIP:  types.StringValue("172.16.10.2"),
				TunnelLocalIP2:  types.StringValue("172.16.20.1"),
				TunnelRemoteIP2: types.StringValue("172.16.20.2"),
				Description:     types.StringValue("Disabled VTI for testing"),
				Id:              types.StringValue("disabled-id"),
			},
			expected: &ipsec.IPsecVTI{
				Enabled:         "0",
				RequestID:       "9999",
				LocalIP:         "10.10.10.1",
				RemoteIP:        "20.20.20.1",
				TunnelLocalIP:   "172.16.10.1",
				TunnelRemoteIP:  "172.16.10.2",
				TunnelLocalIP2:  "172.16.20.1",
				TunnelRemoteIP2: "172.16.20.2",
				Description:     "Disabled VTI for testing",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertIpsecVtiSchemaToStruct(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Enabled, result.Enabled)
			assert.Equal(t, tt.expected.RequestID, result.RequestID)
			assert.Equal(t, tt.expected.LocalIP, result.LocalIP)
			assert.Equal(t, tt.expected.RemoteIP, result.RemoteIP)
			assert.Equal(t, tt.expected.TunnelLocalIP, result.TunnelLocalIP)
			assert.Equal(t, tt.expected.TunnelRemoteIP, result.TunnelRemoteIP)
			assert.Equal(t, tt.expected.TunnelLocalIP2, result.TunnelLocalIP2)
			assert.Equal(t, tt.expected.TunnelRemoteIP2, result.TunnelRemoteIP2)
			assert.Equal(t, tt.expected.Description, result.Description)
		})
	}
}

func TestConvertIpsecVtiStructToSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    *ipsec.IPsecVTI
		expected *IpsecVtiResourceModel
	}{
		{
			name: "basic_conversion",
			input: &ipsec.IPsecVTI{
				Enabled:         "1",
				RequestID:       "5678",
				LocalIP:         "192.168.50.1",
				RemoteIP:        "203.0.113.50",
				TunnelLocalIP:   "10.50.1.1",
				TunnelRemoteIP:  "10.50.1.2",
				TunnelLocalIP2:  "10.60.1.1",
				TunnelRemoteIP2: "10.60.1.2",
				Description:     "Converted VTI",
			},
			expected: &IpsecVtiResourceModel{
				Enabled:         types.StringValue("1"),
				RequestID:       types.StringValue("5678"),
				LocalIP:         types.StringValue("192.168.50.1"),
				RemoteIP:        types.StringValue("203.0.113.50"),
				TunnelLocalIP:   types.StringValue("10.50.1.1"),
				TunnelRemoteIP:  types.StringValue("10.50.1.2"),
				TunnelLocalIP2:  types.StringValue("10.60.1.1"),
				TunnelRemoteIP2: types.StringValue("10.60.1.2"),
				Description:     types.StringValue("Converted VTI"),
			},
		},
		{
			name: "empty_optional_fields",
			input: &ipsec.IPsecVTI{
				Enabled:         "0",
				RequestID:       "",
				LocalIP:         "172.16.100.1",
				RemoteIP:        "172.16.200.1",
				TunnelLocalIP:   "10.200.1.1",
				TunnelRemoteIP:  "10.200.1.2",
				TunnelLocalIP2:  "",
				TunnelRemoteIP2: "",
				Description:     "",
			},
			expected: &IpsecVtiResourceModel{
				Enabled:         types.StringValue("0"),
				RequestID:       types.StringValue(""),
				LocalIP:         types.StringValue("172.16.100.1"),
				RemoteIP:        types.StringValue("172.16.200.1"),
				TunnelLocalIP:   types.StringValue("10.200.1.1"),
				TunnelRemoteIP:  types.StringValue("10.200.1.2"),
				TunnelLocalIP2:  types.StringValue(""),
				TunnelRemoteIP2: types.StringValue(""),
				Description:     types.StringValue(""),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertIpsecVtiStructToSchema(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Enabled, result.Enabled)
			assert.Equal(t, tt.expected.RequestID, result.RequestID)
			assert.Equal(t, tt.expected.LocalIP, result.LocalIP)
			assert.Equal(t, tt.expected.RemoteIP, result.RemoteIP)
			assert.Equal(t, tt.expected.TunnelLocalIP, result.TunnelLocalIP)
			assert.Equal(t, tt.expected.TunnelRemoteIP, result.TunnelRemoteIP)
			assert.Equal(t, tt.expected.TunnelLocalIP2, result.TunnelLocalIP2)
			assert.Equal(t, tt.expected.TunnelRemoteIP2, result.TunnelRemoteIP2)
			assert.Equal(t, tt.expected.Description, result.Description)
		})
	}
}

func TestConvertIpsecVtiRoundTrip(t *testing.T) {
	// Test that schema -> struct -> schema conversion preserves data
	original := &IpsecVtiResourceModel{
		Enabled:         types.StringValue("1"),
		RequestID:       types.StringValue("12345"),
		LocalIP:         types.StringValue("192.168.1.100"),
		RemoteIP:        types.StringValue("203.0.113.100"),
		TunnelLocalIP:   types.StringValue("10.100.1.1"),
		TunnelRemoteIP:  types.StringValue("10.100.1.2"),
		TunnelLocalIP2:  types.StringValue("10.200.1.1"),
		TunnelRemoteIP2: types.StringValue("10.200.1.2"),
		Description:     types.StringValue("Round trip test VTI"),
		Id:              types.StringValue("round-trip-id"),
	}

	// Convert to struct
	vtiStruct, err := convertIpsecVtiSchemaToStruct(original)
	assert.NoError(t, err)

	// Convert back to schema
	result, err := convertIpsecVtiStructToSchema(vtiStruct)
	assert.NoError(t, err)

	// Verify data is preserved (excluding ID which is not part of the struct)
	assert.Equal(t, original.Enabled, result.Enabled)
	assert.Equal(t, original.RequestID, result.RequestID)
	assert.Equal(t, original.LocalIP, result.LocalIP)
	assert.Equal(t, original.RemoteIP, result.RemoteIP)
	assert.Equal(t, original.TunnelLocalIP, result.TunnelLocalIP)
	assert.Equal(t, original.TunnelRemoteIP, result.TunnelRemoteIP)
	assert.Equal(t, original.TunnelLocalIP2, result.TunnelLocalIP2)
	assert.Equal(t, original.TunnelRemoteIP2, result.TunnelRemoteIP2)
	assert.Equal(t, original.Description, result.Description)
}

func TestIpsecVtiSchemaValidation(t *testing.T) {
	schema := IpsecVtiResourceSchema()

	// Verify required fields
	requiredFields := []string{"local_ip", "remote_ip", "tunnel_local_ip", "tunnel_remote_ip", "request_id"}
	for _, field := range requiredFields {
		attr := schema.Attributes[field]
		assert.NotNil(t, attr, "Field %s should exist", field)
		assert.True(t, attr.IsRequired(), "Field %s should be required", field)
	}

	// Verify optional+computed fields
	optionalComputedFields := []string{"enabled"}
	for _, field := range optionalComputedFields {
		attr := schema.Attributes[field]
		assert.NotNil(t, attr, "Field %s should exist", field)
		assert.True(t, attr.IsOptional(), "Field %s should be optional", field)
		assert.True(t, attr.IsComputed(), "Field %s should be computed", field)
	}

	// Verify optional only fields (with defaults)
	optionalOnlyFields := []string{"tunnel_local_ip2", "tunnel_remote_ip2", "description"}
	for _, field := range optionalOnlyFields {
		attr := schema.Attributes[field]
		assert.NotNil(t, attr, "Field %s should exist", field)
		assert.True(t, attr.IsOptional(), "Field %s should be optional", field)
	}

	// Verify computed fields
	computedFields := []string{"id"}
	for _, field := range computedFields {
		attr := schema.Attributes[field]
		assert.NotNil(t, attr, "Field %s should exist", field)
		assert.True(t, attr.IsComputed(), "Field %s should be computed", field)
	}
}

func TestIpsecVtiDataSourceSchema(t *testing.T) {
	schema := IpsecVtiDataSourceSchema()

	// Verify all fields are computed except id
	allFields := []string{"enabled", "request_id", "local_ip", "remote_ip", "tunnel_local_ip", "tunnel_remote_ip", "tunnel_local_ip2", "tunnel_remote_ip2", "description"}
	for _, field := range allFields {
		attr := schema.Attributes[field]
		assert.NotNil(t, attr, "Field %s should exist", field)
		assert.True(t, attr.IsComputed(), "Field %s should be computed", field)
	}

	// Verify id is required
	idAttr := schema.Attributes["id"]
	assert.NotNil(t, idAttr, "id field should exist")
	assert.True(t, idAttr.IsRequired(), "id field should be required")
}
