package routes

import (
	"testing"

	"github.com/browningluke/opnsense-go/pkg/routes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConvertRouteSchemaToStructSendsBothEnabledKeys(t *testing.T) {
	route, err := convertRouteSchemaToStruct(&routeResourceModel{
		Enabled: types.BoolValue(false),
		Gateway: types.StringValue("WAN_DHCP"),
		Network: types.StringValue("192.0.2.0/24"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if route.Enabled != "0" || route.Disabled != "1" {
		t.Fatalf("expected enabled=0 disabled=1, got enabled=%q disabled=%q", route.Enabled, route.Disabled)
	}
}

func TestConvertRouteStructToSchemaEnabledField(t *testing.T) {
	tests := []struct {
		name     string
		route    routes.Route
		expected bool
	}{
		{"26.1.10+ enabled", routes.Route{Enabled: "1"}, true},
		{"26.1.10+ disabled", routes.Route{Enabled: "0"}, false},
		{"pre-26.1.10 enabled", routes.Route{Disabled: "0"}, true},
		{"pre-26.1.10 disabled", routes.Route{Disabled: "1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := convertRouteStructToSchema(&tt.route)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if model.Enabled.ValueBool() != tt.expected {
				t.Fatalf("expected enabled=%v, got %v", tt.expected, model.Enabled.ValueBool())
			}
		})
	}
}
