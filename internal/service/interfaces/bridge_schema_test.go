package interfaces

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/interfaces"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConvertBridgeSchemaToStruct(t *testing.T) {
	model := &bridgeResourceModel{
		Members:        mustSet(t, "opt1", "opt2"),
		LinkLocal:      types.BoolValue(true),
		EnableStp:      types.BoolValue(true),
		StpProto:       types.StringValue("rstp"),
		StpInterfaces:  mustSet(t, "opt1"),
		MaxAge:         types.Int64Value(20),
		ForwardDelay:   types.Int64Null(),
		HoldCount:      types.Int64Null(),
		MaxAddresses:   types.Int64Null(),
		AddressTimeout: types.Int64Null(),
		SpanInterface:  types.StringValue(""),
		EdgePorts:      mustSet(t),
		AutoEdgePorts:  mustSet(t),
		PtpPorts:       mustSet(t),
		AutoPtpPorts:   mustSet(t),
		StaticPorts:    mustSet(t),
		PrivatePorts:   mustSet(t),
		Description:    types.StringValue("test bridge"),
		Device:         types.StringValue(""),
	}

	bridge, err := convertBridgeSchemaToStruct(model)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if got := len(bridge.Members); got != 2 {
		t.Errorf("expected 2 members, got %d", got)
	}
	if bridge.LinkLocal != "1" {
		t.Errorf("expected linklocal \"1\", got %q", bridge.LinkLocal)
	}
	if bridge.EnableStp != "1" {
		t.Errorf("expected enablestp \"1\", got %q", bridge.EnableStp)
	}
	if bridge.StpProto.String() != "rstp" {
		t.Errorf("expected proto \"rstp\", got %q", bridge.StpProto.String())
	}
	if bridge.MaxAge != "20" {
		t.Errorf("expected maxage \"20\", got %q", bridge.MaxAge)
	}
	if bridge.ForwardDelay != "" {
		t.Errorf("expected empty fwdelay for null value, got %q", bridge.ForwardDelay)
	}
	if bridge.Description != "test bridge" {
		t.Errorf("expected description \"test bridge\", got %q", bridge.Description)
	}
}

func TestConvertBridgeStructToSchema(t *testing.T) {
	bridge := &interfaces.Bridge{
		Device:         "bridge0",
		Members:        []string{"opt1", "opt2"},
		LinkLocal:      "0",
		EnableStp:      "1",
		StpProto:       "stp",
		StpInterfaces:  []string{"opt1"},
		MaxAge:         "20",
		ForwardDelay:   "",
		HoldCount:      "",
		MaxAddresses:   "",
		AddressTimeout: "",
		SpanInterface:  "",
		EdgePorts:      []string{},
		AutoEdgePorts:  []string{},
		PtpPorts:       []string{},
		AutoPtpPorts:   []string{},
		StaticPorts:    []string{},
		PrivatePorts:   []string{},
		Description:    "test bridge",
	}

	model, err := convertBridgeStructToSchema(bridge)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if model.Device.ValueString() != "bridge0" {
		t.Errorf("expected device \"bridge0\", got %q", model.Device.ValueString())
	}
	if got := len(model.Members.Elements()); got != 2 {
		t.Errorf("expected 2 members, got %d", got)
	}
	if model.LinkLocal.ValueBool() {
		t.Error("expected link_local false")
	}
	if !model.EnableStp.ValueBool() {
		t.Error("expected enable_stp true")
	}
	if model.StpProto.ValueString() != "stp" {
		t.Errorf("expected stp_proto \"stp\", got %q", model.StpProto.ValueString())
	}
	if model.MaxAge.ValueInt64() != 20 {
		t.Errorf("expected max_age 20, got %d", model.MaxAge.ValueInt64())
	}
	if !model.ForwardDelay.IsNull() {
		t.Error("expected forward_delay to be null for empty API value")
	}
	if model.Description.ValueString() != "test bridge" {
		t.Errorf("expected description \"test bridge\", got %q", model.Description.ValueString())
	}
}

func mustSet(t *testing.T, elems ...string) types.Set {
	t.Helper()

	var values []types.String
	for _, e := range elems {
		values = append(values, types.StringValue(e))
	}
	set, diags := types.SetValueFrom(t.Context(), types.StringType, values)
	if diags.HasError() {
		t.Fatalf("failed to build set: %v", diags)
	}
	return set
}
