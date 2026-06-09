package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type requiredWhenBoolValidator struct {
	expression path.Expression
	value      bool
}

func (v requiredWhenBoolValidator) Description(_ context.Context) string {
	return fmt.Sprintf("must be set to a non-empty value when %s is %v", v.expression, v.value)
}

func (v requiredWhenBoolValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v requiredWhenBoolValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Pass: this attribute already has a non-empty value.
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() && req.ConfigValue.ValueString() != "" {
		return
	}

	matchedPaths, diags := req.Config.PathMatches(ctx, v.expression)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, p := range matchedPaths {
		var other types.Bool
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, p, &other)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Defer judgement until the other attribute is known.
		if other.IsUnknown() {
			return
		}

		// If the other attribute is unset, the schema default fires at plan
		// time; we treat "unset" as "not the trigger value" and let plan
		// re-validate when the value is concrete.
		if other.IsNull() {
			continue
		}

		if other.ValueBool() == v.value {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Missing required attribute",
				fmt.Sprintf("Attribute %q must be set to a non-empty value when %q is %v.", req.Path, p, v.value),
			)
		}
	}
}

// RequiredWhenBool returns a validator that errors if the attribute is
// null/empty while the bool attribute at the given path equals value.
func RequiredWhenBool(expression path.Expression, value bool) validator.String {
	return requiredWhenBoolValidator{expression: expression, value: value}
}
