package firewall

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNatDescriptionValidator(t *testing.T) {
	tests := map[string]bool{
		"WG_ProtonDE NAT":               true,
		"Redirect ZT to ZT_DEBIAN":      true,
		"tardis: audit (infra repo) #1": true,
		"Disable ZT NAT":                true,
		"first line\nsecond line":       false,
		"carriage\rreturn":              false,
		"nul\x00byte":                   false,
		"form\ffeed":                    false,
	}

	for description, valid := range tests {
		request := validator.StringRequest{
			Path:        path.Root("description"),
			ConfigValue: types.StringValue(description),
		}
		response := &validator.StringResponse{}
		natDescriptionValidator.ValidateString(context.Background(), request, response)

		if response.Diagnostics.HasError() == valid {
			t.Errorf("description %q: valid = %t, got errors: %v", description, valid, response.Diagnostics)
		}
	}
}
