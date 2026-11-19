package unbound

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/unbound"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConvertHostOverrideSchemaToStructTTL(t *testing.T) {
	tests := []struct {
		name     string
		ttl      int64
		expected string
	}{
		{"unset leaves the field empty", -1, ""},
		{"zero is a valid ttl", 0, "0"},
		{"explicit ttl", 300, "300"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			override, err := convertHostOverrideSchemaToStruct(&hostOverrideResourceModel{
				Hostname: types.StringValue("host"),
				Domain:   types.StringValue("example.com"),
				Server:   types.StringValue("192.0.2.1"),
				TTL:      types.Int64Value(tt.ttl),
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if override.TTL != tt.expected {
				t.Fatalf("expected ttl=%q, got %q", tt.expected, override.TTL)
			}
		})
	}
}

func TestConvertHostOverrideStructToSchemaTTL(t *testing.T) {
	tests := []struct {
		name     string
		ttl      string
		expected int64
	}{
		{"empty means unbound's default", "", -1},
		{"zero is a valid ttl", "0", 0},
		{"explicit ttl", "300", 300},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := convertHostOverrideStructToSchema(&unbound.HostOverride{TTL: tt.ttl})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if model.TTL.ValueInt64() != tt.expected {
				t.Fatalf("expected ttl=%d, got %d", tt.expected, model.TTL.ValueInt64())
			}
		})
	}
}
