package dnsmasq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettingsResourceSchemaContainsAllConfigurationBlocks(t *testing.T) {
	s := settingsResourceSchema()
	assert.Equal(t, int64(1), s.Version)
	for _, name := range []string{"id", "enabled", "interface", "strict_interface_binding", "dns", "dns_query_forwarding", "dhcp", "legacy"} {
		assert.Contains(t, s.Attributes, name)
	}
	for _, name := range []string{"dns", "dns_query_forwarding", "dhcp", "legacy"} {
		attribute, ok := s.Attributes[name].(schema.SingleNestedAttribute)
		require.Truef(t, ok, "%s must be a nested configuration block", name)
		assert.True(t, attribute.Optional)
		assert.True(t, attribute.Computed)
	}
}

func TestConvertSettingsSchemaRoundTrip(t *testing.T) {
	model := &settingsResourceModel{
		Enabled:                types.BoolValue(true),
		StrictInterfaceBinding: types.BoolValue(true),
		Interface:              types.SetValueMust(types.StringType, []attr.Value{types.StringValue("wan"), types.StringValue("lan")}),
		DNS: &settingsDNSBlock{
			Port:                 types.Int64Value(5353),
			DNSSEC:               types.BoolValue(true),
			NoHostLookup:         types.BoolValue(true),
			LogQueries:           types.BoolValue(true),
			MaxConcurrentQueries: types.Int64Value(150),
			CacheSize:            types.Int64Value(1000),
			LocalEntryTTL:        types.Int64Value(60),
			NoIdent:              types.BoolValue(true),
		},
		DNSQueryForwarding: &settingsDNSQueryForwardingBlock{
			QuerySequentially:          types.BoolValue(true),
			RequireDomain:              types.BoolValue(true),
			DoNotForwardSystemDNS:      types.BoolValue(true),
			DoNotForwardPrivateReverse: types.BoolValue(true),
			AddMAC:                     types.StringValue("standard"),
			AddSubnet:                  types.BoolValue(true),
			StripSubnet:                types.BoolValue(true),
		},
		DHCP: &settingsDHCPBlock{
			NoDHCPInterfaces:      types.SetValueMust(types.StringType, []attr.Value{types.StringValue("opt1")}),
			FQDN:                  types.BoolValue(true),
			DefaultDomain:         types.StringValue("example.test"),
			LocalDomain:           types.BoolValue(true),
			MaxLeases:             types.Int64Value(200),
			Authoritative:         types.BoolValue(true),
			ReplyDelay:            types.Int64Value(5),
			RegisterFirewallRules: types.BoolValue(true),
			RouterAdvertisements:  types.BoolValue(true),
			HostPing:              types.BoolValue(true),
			DisableHASync:         types.BoolValue(true),
			LogDhcp:               types.BoolValue(true),
			LogQuiet:              types.BoolValue(true),
		},
		Legacy: &settingsLegacyBlock{
			RegisterISCDHCP4Leases:    types.BoolValue(true),
			DHCPDomainOverride:        types.StringValue("legacy.test"),
			RegisterDHCPStaticMapping: types.BoolValue(true),
			PreferDHCP:                types.BoolValue(true),
		},
	}
	settings, err := convertSettingsSchemaToStruct(model)
	require.NoError(t, err)
	assert.Equal(t, "1", settings.Dnsmasq.IsEnabled)
	assert.Equal(t, "1", settings.Dnsmasq.StrictInterfaceBinding)
	assert.Equal(t, []string{"lan", "wan"}, []string(settings.Dnsmasq.Interface))
	assert.Equal(t, "5353", settings.Dnsmasq.DNS_Port)
	assert.Equal(t, "standard", settings.Dnsmasq.DNS_QF_AddMac.String())
	assert.Equal(t, "200", settings.Dnsmasq.DHCPSettings.MaxLeases)
	assert.Equal(t, "legacy.test", settings.Dnsmasq.Legacy_DhcpDomainOverride)
	roundTripped, err := convertSettingsStructToSchema(settings)
	require.NoError(t, err)
	assert.Equal(t, types.StringValue("dnsmasq_settings"), roundTripped.Id)
	assert.Equal(t, model.Enabled, roundTripped.Enabled)
	assert.Equal(t, types.SetValueMust(types.StringType, []attr.Value{
		types.StringValue("lan"),
		types.StringValue("wan"),
	}), roundTripped.Interface)
	assert.Equal(t, model.DNS, roundTripped.DNS)
	assert.Equal(t, model.DNSQueryForwarding, roundTripped.DNSQueryForwarding)
	assert.Equal(t, model.DHCP, roundTripped.DHCP)
	assert.Equal(t, model.Legacy, roundTripped.Legacy)
}
