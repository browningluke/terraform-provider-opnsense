package firewall

import (
	"testing"

	opnfirewall "github.com/browningluke/opnsense-go/pkg/firewall"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestConvertNATPortForwardSchemaToStruct(t *testing.T) {
	data := &natPortForwardResourceModel{
		Enabled:    types.BoolValue(true),
		Sequence:   types.Int64Value(100),
		Interface:  types.StringValue("wan"),
		IPProtocol: types.StringValue("inet"),
		Protocol:   types.StringValue("tcp"),
		Source: &firewallLocation{
			Net:    types.StringValue("any"),
			Port:   types.StringValue(""),
			Invert: types.BoolValue(false),
		},
		Destination: &firewallLocation{
			Net:    types.StringValue("wanip"),
			Port:   types.StringValue("443"),
			Invert: types.BoolValue(false),
		},
		Target: &firewallTarget{
			IP:   types.StringValue("10.1.1.20"),
			Port: types.StringValue("443"),
		},
		Log:           types.BoolValue(true),
		NatReflection: types.StringValue("enable"),
		Description:   types.StringValue("WAN HTTPS to k3s Traefik ingress VIP"),
	}

	result, err := convertNATPortForwardSchemaToStruct(data)

	require.NoError(t, err)
	require.Equal(t, "0", result.Disabled)
	require.Equal(t, "100", result.Sequence)
	require.Equal(t, "wan", result.Interface.String())
	require.Equal(t, "inet", result.IPProtocol.String())
	require.Equal(t, "tcp", result.Protocol.String())
	require.Equal(t, "any", result.Source.Network)
	require.Equal(t, "", result.Source.Port)
	require.Equal(t, "0", result.Source.Invert)
	require.Equal(t, "wanip", result.Destination.Network)
	require.Equal(t, "443", result.Destination.Port)
	require.Equal(t, "0", result.Destination.Invert)
	require.Equal(t, "10.1.1.20", result.Target)
	require.Equal(t, "443", result.TargetPort)
	require.Equal(t, "1", result.Log)
	require.Equal(t, "purenat", result.NatReflection.String())
	require.Equal(t, "WAN HTTPS to k3s Traefik ingress VIP", result.Description)
}

func TestConvertNATPortForwardStructToSchema(t *testing.T) {
	result, err := convertNATPortForwardStructToSchema(&opnfirewall.NatPortForward{
		Disabled:   "1",
		Sequence:   "200",
		Interface:  "wan",
		IPProtocol: "inet",
		Protocol:   "tcp",
		Source: opnfirewall.NatPortForwardLocation{
			Network: "",
			Port:    "",
			Invert:  "0",
		},
		Destination: opnfirewall.NatPortForwardLocation{
			Network: "wanip",
			Port:    "8443",
			Invert:  "1",
		},
		Target:        "10.1.1.30",
		TargetPort:    "443",
		Log:           "0",
		NatReflection: "",
		Description:   "Updated WAN HTTPS",
	})

	require.NoError(t, err)
	require.False(t, result.Enabled.ValueBool())
	require.Equal(t, int64(200), result.Sequence.ValueInt64())
	require.Equal(t, "wan", result.Interface.ValueString())
	require.Equal(t, "inet", result.IPProtocol.ValueString())
	require.Equal(t, "tcp", result.Protocol.ValueString())
	require.Equal(t, "any", result.Source.Net.ValueString())
	require.Equal(t, "", result.Source.Port.ValueString())
	require.False(t, result.Source.Invert.ValueBool())
	require.Equal(t, "wanip", result.Destination.Net.ValueString())
	require.Equal(t, "8443", result.Destination.Port.ValueString())
	require.True(t, result.Destination.Invert.ValueBool())
	require.Equal(t, "10.1.1.30", result.Target.IP.ValueString())
	require.Equal(t, "443", result.Target.Port.ValueString())
	require.False(t, result.Log.ValueBool())
	require.Equal(t, "default", result.NatReflection.ValueString())
	require.Equal(t, "Updated WAN HTTPS", result.Description.ValueString())
}
