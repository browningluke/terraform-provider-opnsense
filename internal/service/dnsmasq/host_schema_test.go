package dnsmasq

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestConvertDnsmasqHostToStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    *hostResourceModel
		expected *dnsmasq.Host
	}{
		{
			name: "full host conversion",
			input: &hostResourceModel{
				Hostname:        types.StringValue("server"),
				Domain:          types.StringValue("example.com"),
				IsLocalDomain:   types.BoolValue(true),
				IpAddresses:     tools.StringSliceToSet([]string{"192.168.1.10", "fd00::10"}),
				AliasRecords:    tools.StringSliceToSet([]string{"server.example.net", "server.local"}),
				CnameRecords:    tools.StringSliceToSet([]string{"web"}),
				ClientID:        types.StringValue("01:02:03:04"),
				HarwareAdresses: tools.StringSliceToSet([]string{"00:11:22:33:44:55", "AA:BB:CC:DD:EE"}),
				Tag:             types.StringValue("server-tag"),
				IsIgnored:       types.BoolValue(true),
				Description:     types.StringValue("Test host"),
				Comment:         types.StringValue("Test comment"),
			},
			expected: &dnsmasq.Host{
				Hostname:          "server",
				Domain:            "example.com",
				IsLocalDomain:     "1",
				IpAddresses:       []string{"192.168.1.10", "fd00::10"},
				AliasRecords:      []string{"server.example.net", "server.local"},
				CnameRecords:      []string{"web"},
				ClientId:          "01:02:03:04",
				HardwareAddresses: []string{"00:11:22:33:44:55", "AA:BB:CC:DD:EE"},
				Tag:               api.SelectedMap("server-tag"),
				IsIgnored:         "1",
				Description:       "Test host",
				Comments:          "Test comment",
			},
		},
		{
			name: "host conversion with empty optional values",
			input: &hostResourceModel{
				Hostname:        types.StringValue("server"),
				Domain:          types.StringValue(""),
				IsLocalDomain:   types.BoolValue(false),
				IpAddresses:     tools.StringSliceToSet([]string{"192.168.1.10"}),
				AliasRecords:    tools.StringSliceToSet([]string{}),
				CnameRecords:    tools.StringSliceToSet([]string{}),
				ClientID:        types.StringValue(""),
				HarwareAdresses: tools.StringSliceToSet([]string{}),
				Tag:             types.StringValue(""),
				IsIgnored:       types.BoolValue(false),
				Description:     types.StringValue(""),
				Comment:         types.StringValue(""),
			},
			expected: &dnsmasq.Host{
				Hostname:          "server",
				Domain:            "",
				IsLocalDomain:     "0",
				IpAddresses:       []string{"192.168.1.10"},
				AliasRecords:      []string{},
				CnameRecords:      []string{},
				ClientId:          "",
				HardwareAddresses: []string{},
				Tag:               api.SelectedMap(""),
				IsIgnored:         "0",
				Description:       "",
				Comments:          "",
			},
		},
		{
			name: "host conversion with multiple addresses and records",
			input: &hostResourceModel{
				Hostname:      types.StringValue("router"),
				Domain:        types.StringValue("home.arpa"),
				IsLocalDomain: types.BoolValue(true),
				IpAddresses: tools.StringSliceToSet([]string{
					"192.168.1.1",
					"192.168.2.1",
					"fd00::1",
				}),
				AliasRecords: tools.StringSliceToSet([]string{
					"gw.home.arpa",
					"gateway.home.arpa",
				}),
				CnameRecords: tools.StringSliceToSet([]string{
					"gateway",
					"router-gw",
				}),
				ClientID:        types.StringValue(""),
				HarwareAdresses: tools.StringSliceToSet([]string{"00:01:02:03:04:05"}),
				Tag:             types.StringValue("router"),
				IsIgnored:       types.BoolValue(false),
				Description:     types.StringValue("Router"),
				Comment:         types.StringValue("Gateway host"),
			},
			expected: &dnsmasq.Host{
				Hostname:      "router",
				Domain:        "home.arpa",
				IsLocalDomain: "1",
				IpAddresses: []string{
					"192.168.1.1",
					"192.168.2.1",
					"fd00::1",
				},
				AliasRecords: []string{
					"gw.home.arpa",
					"gateway.home.arpa",
				},
				CnameRecords: []string{
					"gateway",
					"router-gw",
				},
				ClientId: "",
				HardwareAddresses: []string{
					"00:01:02:03:04:05",
				},
				Tag:         api.SelectedMap("router"),
				IsIgnored:   "0",
				Description: "Router",
				Comments:    "Gateway host",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertHostSchemaToStruct(tt.input)

			assert.NoError(t, err)
			assert.NotNil(t, result)

			assert.Equal(t, tt.expected.Hostname, result.Hostname)
			assert.Equal(t, tt.expected.Domain, result.Domain)
			assert.Equal(t, tt.expected.IsLocalDomain, result.IsLocalDomain)
			assert.Equal(t, tt.expected.IpAddresses, result.IpAddresses)
			assert.Equal(t, tt.expected.AliasRecords, result.AliasRecords)
			assert.Equal(t, tt.expected.CnameRecords, result.CnameRecords)
			assert.Equal(t, tt.expected.ClientId, result.ClientId)
			assert.Equal(t, tt.expected.HardwareAddresses, result.HardwareAddresses)
			assert.Equal(t, tt.expected.IsIgnored, result.IsIgnored)
			assert.Equal(t, tt.expected.Description, result.Description)
			assert.Equal(t, tt.expected.Comments, result.Comments)

			assert.Equal(
				t,
				tt.expected.Tag.String(),
				result.Tag.String(),
			)
		})
	}

}

func TestConvertDnsmasqHostStructToSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    *dnsmasq.Host
		expected *hostResourceModel
	}{
		{
			name: "full host conversion",
			input: &dnsmasq.Host{
				Hostname:          "server",
				Domain:            "example.com",
				IsLocalDomain:     "1",
				IpAddresses:       []string{"192.168.1.10", "fd00::10"},
				AliasRecords:      []string{"server.example.net", "server.local"},
				CnameRecords:      []string{"web"},
				ClientId:          "01:02:03:04",
				HardwareAddresses: []string{"00:11:22:33:44:55", "AA:BB:CC:DD:EE"},
				Tag:               api.SelectedMap("server-tag"),
				IsIgnored:         "1",
				Description:       "Test host",
				Comments:          "Test comment",
			},
			expected: &hostResourceModel{
				Hostname:        types.StringValue("server"),
				Domain:          types.StringValue("example.com"),
				IsLocalDomain:   types.BoolValue(true),
				IpAddresses:     tools.StringSliceToSet([]string{"192.168.1.10", "fd00::10"}),
				AliasRecords:    tools.StringSliceToSet([]string{"server.example.net", "server.local"}),
				CnameRecords:    tools.StringSliceToSet([]string{"web"}),
				ClientID:        types.StringValue("01:02:03:04"),
				HarwareAdresses: tools.StringSliceToSet([]string{"00:11:22:33:44:55", "AA:BB:CC:DD:EE"}),
				Tag:             types.StringValue("server-tag"),
				IsIgnored:       types.BoolValue(true),
				Description:     types.StringValue("Test host"),
				Comment:         types.StringValue("Test comment"),
			},
		},
		{
			name: "host conversion with empty optional values",
			input: &dnsmasq.Host{
				Hostname:          "server",
				Domain:            "",
				IsLocalDomain:     "0",
				IpAddresses:       []string{"192.168.1.10"},
				AliasRecords:      []string{},
				CnameRecords:      []string{},
				ClientId:          "",
				HardwareAddresses: []string{},
				Tag:               api.SelectedMap(""),
				IsIgnored:         "0",
				Description:       "",
				Comments:          "",
			},
			expected: &hostResourceModel{
				Hostname:        types.StringValue("server"),
				Domain:          types.StringValue(""),
				IsLocalDomain:   types.BoolValue(false),
				IpAddresses:     tools.StringSliceToSet([]string{"192.168.1.10"}),
				AliasRecords:    tools.StringSliceToSet([]string{}),
				CnameRecords:    tools.StringSliceToSet([]string{}),
				ClientID:        types.StringValue(""),
				HarwareAdresses: tools.StringSliceToSet([]string{}),
				Tag:             types.StringValue(""),
				IsIgnored:       types.BoolValue(false),
				Description:     types.StringValue(""),
				Comment:         types.StringValue(""),
			},
		},
		{
			name: "host conversion with false ignore and local domain",
			input: &dnsmasq.Host{
				Hostname:          "client",
				Domain:            "home.arpa",
				IsLocalDomain:     "1",
				IpAddresses:       []string{"10.0.0.20"},
				AliasRecords:      []string{"client.home.arpa"},
				CnameRecords:      []string{},
				ClientId:          "",
				HardwareAddresses: []string{"12:34:56:78:9A"},
				Tag:               api.SelectedMap("client"),
				IsIgnored:         "0",
				Description:       "Client host",
				Comments:          "Client comment",
			},
			expected: &hostResourceModel{
				Hostname:        types.StringValue("client"),
				Domain:          types.StringValue("home.arpa"),
				IsLocalDomain:   types.BoolValue(true),
				IpAddresses:     tools.StringSliceToSet([]string{"10.0.0.20"}),
				AliasRecords:    tools.StringSliceToSet([]string{"client.home.arpa"}),
				CnameRecords:    tools.StringSliceToSet([]string{}),
				ClientID:        types.StringValue(""),
				HarwareAdresses: tools.StringSliceToSet([]string{"12:34:56:78:9A"}),
				Tag:             types.StringValue("client"),
				IsIgnored:       types.BoolValue(false),
				Description:     types.StringValue("Client host"),
				Comment:         types.StringValue("Client comment"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertHostStructToSchema(tt.input)

			assert.NoError(t, err)
			assert.NotNil(t, result)

			assert.Equal(t, tt.expected.Hostname, result.Hostname)
			assert.Equal(t, tt.expected.Domain, result.Domain)
			assert.Equal(t, tt.expected.IsLocalDomain, result.IsLocalDomain)
			assert.Equal(t, tt.expected.IpAddresses, result.IpAddresses)
			assert.Equal(t, tt.expected.AliasRecords, result.AliasRecords)
			assert.Equal(t, tt.expected.CnameRecords, result.CnameRecords)
			assert.Equal(t, tt.expected.ClientID, result.ClientID)
			assert.Equal(t, tt.expected.HarwareAdresses, result.HarwareAdresses)
			assert.Equal(t, tt.expected.Tag, result.Tag)
			assert.Equal(t, tt.expected.IsIgnored, result.IsIgnored)
			assert.Equal(t, tt.expected.Description, result.Description)
			assert.Equal(t, tt.expected.Comment, result.Comment)
		})
	}

}
