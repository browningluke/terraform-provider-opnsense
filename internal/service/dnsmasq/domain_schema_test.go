package dnsmasq

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestConvertDnsmasqDomainToStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    *domainResourceModel
		expected *dnsmasq.Domain
	}{
		{
			name: "full domain conversion",
			input: &domainResourceModel{
				Sequence:      types.Int64Value(10),
				Domain:        types.StringValue("example.com"),
				FirewallAlias: types.StringValue("example_alias"),
				SourceIp:      types.StringValue("192.168.1.10"),
				Port:          types.Int64Value(5353),
				Ip:            types.StringValue("192.168.1.53"),
				Description:   types.StringValue("Test domain override"),
			},
			expected: &dnsmasq.Domain{
				Sequence:      "10",
				Domain:        "example.com",
				FirewallAlias: api.SelectedMap("example_alias"),
				SourceIp:      "192.168.1.10",
				Port:          "5353",
				Ip:            "192.168.1.53",
				Description:   "Test domain override",
			},
		},
		{
			name: "domain conversion without optional port",
			input: &domainResourceModel{
				Sequence:      types.Int64Value(20),
				Domain:        types.StringValue("internal.example.com"),
				FirewallAlias: types.StringValue(""),
				SourceIp:      types.StringValue("192.168.1.1"),
				Port:          types.Int64Null(),
				Ip:            types.StringValue("10.0.0.53"),
				Description:   types.StringValue("Without port"),
			},
			expected: &dnsmasq.Domain{
				Sequence:      "20",
				Domain:        "internal.example.com",
				FirewallAlias: api.SelectedMap(""),
				SourceIp:      "192.168.1.1",
				Ip:            "10.0.0.53",
				Description:   "Without port",
			},
		},
		{
			name: "domain conversion with empty optional strings",
			input: &domainResourceModel{
				Sequence:      types.Int64Value(1),
				Domain:        types.StringValue("example.org"),
				FirewallAlias: types.StringValue(""),
				SourceIp:      types.StringValue(""),
				Port:          types.Int64Null(),
				Ip:            types.StringValue(""),
				Description:   types.StringValue(""),
			},
			expected: &dnsmasq.Domain{
				Sequence:      "1",
				Domain:        "example.org",
				FirewallAlias: api.SelectedMap(""),
				SourceIp:      "",
				Ip:            "",
				Description:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertDomainSchemaToStruct(tt.input)

			assert.NoError(t, err)
			assert.NotNil(t, result)

			assert.Equal(t, tt.expected.Sequence, result.Sequence)
			assert.Equal(t, tt.expected.Domain, result.Domain)
			assert.Equal(t, tt.expected.SourceIp, result.SourceIp)
			assert.Equal(t, tt.expected.Port, result.Port)
			assert.Equal(t, tt.expected.Ip, result.Ip)
			assert.Equal(t, tt.expected.Description, result.Description)

			assert.Equal(
				t,
				tt.expected.FirewallAlias.String(),
				result.FirewallAlias.String(),
			)
		})
	}
}

func TestConvertDomainStructToSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    *dnsmasq.Domain
		expected *domainResourceModel
	}{
		{
			name: "full domain conversion",
			input: &dnsmasq.Domain{
				Sequence:      "10",
				Domain:        "example.com",
				FirewallAlias: api.SelectedMap("example_alias"),
				SourceIp:      "192.168.1.10",
				Port:          "5353",
				Ip:            "192.168.1.53",
				Description:   "Test domain override",
			},
			expected: &domainResourceModel{
				Sequence:      types.Int64Value(10),
				Domain:        types.StringValue("example.com"),
				FirewallAlias: types.StringValue("example_alias"),
				SourceIp:      types.StringValue("192.168.1.10"),
				Port:          types.Int64Value(5353),
				Ip:            types.StringValue("192.168.1.53"),
				Description:   types.StringValue("Test domain override"),
			},
		},
		{
			name: "domain conversion without port",
			input: &dnsmasq.Domain{
				Sequence:      "20",
				Domain:        "internal.example.com",
				FirewallAlias: api.SelectedMap(""),
				SourceIp:      "192.168.1.1",
				Ip:            "10.0.0.53",
				Description:   "Without port",
			},
			expected: &domainResourceModel{
				Sequence:      types.Int64Value(20),
				Domain:        types.StringValue("internal.example.com"),
				FirewallAlias: types.StringValue(""),
				SourceIp:      types.StringValue("192.168.1.1"),
				Port:          types.Int64Null(),
				Ip:            types.StringValue("10.0.0.53"),
				Description:   types.StringValue("Without port"),
			},
		},
		{
			name: "domain conversion without sequence and port",
			input: &dnsmasq.Domain{
				Domain:        "example.org",
				FirewallAlias: api.SelectedMap(""),
				SourceIp:      "",
				Ip:            "",
				Description:   "",
			},
			expected: &domainResourceModel{
				Sequence:      types.Int64Null(),
				Domain:        types.StringValue("example.org"),
				FirewallAlias: types.StringValue(""),
				SourceIp:      types.StringValue(""),
				Port:          types.Int64Null(),
				Ip:            types.StringValue(""),
				Description:   types.StringValue(""),
			},
		},
		{
			name: "zero-like empty sequence and port handling",
			input: &dnsmasq.Domain{
				Sequence:      "",
				Domain:        "example.net",
				FirewallAlias: api.SelectedMap("dns_alias"),
				SourceIp:      "10.0.0.1",
				Port:          "",
				Ip:            "10.0.0.53",
				Description:   "Empty numeric fields",
			},
			expected: &domainResourceModel{
				Sequence:      types.Int64Null(),
				Domain:        types.StringValue("example.net"),
				FirewallAlias: types.StringValue("dns_alias"),
				SourceIp:      types.StringValue("10.0.0.1"),
				Port:          types.Int64Null(),
				Ip:            types.StringValue("10.0.0.53"),
				Description:   types.StringValue("Empty numeric fields"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertDomainStructToSchema(tt.input)

			assert.NoError(t, err)
			assert.NotNil(t, result)

			assert.Equal(t, tt.expected.Sequence, result.Sequence)
			assert.Equal(t, tt.expected.Domain, result.Domain)
			assert.Equal(t, tt.expected.FirewallAlias, result.FirewallAlias)
			assert.Equal(t, tt.expected.SourceIp, result.SourceIp)
			assert.Equal(t, tt.expected.Port, result.Port)
			assert.Equal(t, tt.expected.Ip, result.Ip)
			assert.Equal(t, tt.expected.Description, result.Description)
		})
	}
}
