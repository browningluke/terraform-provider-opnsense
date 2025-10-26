package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// numericGreaterThanOrEqualValidator is a resource-level validator that ensures
// one numeric attribute value is greater than or equal to another numeric attribute value.
// Supports Int64, Float64, and Number types.
type numericGreaterThanOrEqualValidator struct {
	greaterOrEqualPath path.Expression
	thanPath           path.Expression
}

// Description returns a plain text description of the validator's behavior.
func (v numericGreaterThanOrEqualValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Ensures the value at %q is greater than or equal to the value at %q",
		v.greaterOrEqualPath.String(), v.thanPath.String())
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior.
func (v numericGreaterThanOrEqualValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateResource performs the validation.
func (v numericGreaterThanOrEqualValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	// Get the "greater or equal" value
	greaterOrEqualValue, diags := v.getNumericValue(ctx, req, v.greaterOrEqualPath)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Skip validation if greater or equal value is not set or is default
	if greaterOrEqualValue == nil || *greaterOrEqualValue == -1 {
		return
	}

	// Get the "than" value
	thanValue, diags := v.getNumericValue(ctx, req, v.thanPath)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Skip validation if than value is not set or is default
	if thanValue == nil || *thanValue == -1 {
		return
	}

	// Perform the comparison
	if *greaterOrEqualValue < *thanValue {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			fmt.Sprintf("The %s value (%v) must be greater than or equal to the %s value (%v).",
				v.greaterOrEqualPath.String(), *greaterOrEqualValue,
				v.thanPath.String(), *thanValue,
			),
		)
	}
}

// getNumericValue extracts a numeric value (Int64, Float64, or Number) from the given path expression.
// Returns the value as a float64 for comparison purposes.
// Returns nil if the value is null or unknown.
func (v numericGreaterThanOrEqualValidator) getNumericValue(ctx context.Context, req resource.ValidateConfigRequest, expression path.Expression) (*float64, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Find paths matching the expression
	matchedPaths, pathDiags := req.Config.PathMatches(ctx, expression)
	diags.Append(pathDiags...)
	if pathDiags.HasError() {
		return nil, diags
	}

	// We expect exactly one match for our use case
	if len(matchedPaths) == 0 {
		// Path doesn't exist, which is okay - it might be optional
		return nil, diags
	}

	// Use the first matched path
	matchedPath := matchedPaths[0]

	// Get the generic attr.Value at the matched path
	var matchedPathValue attr.Value
	getDiags := req.Config.GetAttribute(ctx, matchedPath, &matchedPathValue)
	diags.Append(getDiags...)
	if getDiags.HasError() {
		return nil, diags
	}

	// If the value is null or unknown, we cannot compare
	if matchedPathValue.IsNull() || matchedPathValue.IsUnknown() {
		return nil, diags
	}

	// Try to extract as different numeric types
	var result float64

	// Try Int64
	if int64Value, ok := matchedPathValue.(basetypes.Int64Valuable); ok {
		var value types.Int64
		valueDiags := tfsdk.ValueAs(ctx, int64Value, &value)
		diags.Append(valueDiags...)
		if valueDiags.HasError() {
			return nil, diags
		}
		result = float64(value.ValueInt64())
		return &result, diags
	}

	// Try Float64
	if float64Value, ok := matchedPathValue.(basetypes.Float64Valuable); ok {
		var value types.Float64
		valueDiags := tfsdk.ValueAs(ctx, float64Value, &value)
		diags.Append(valueDiags...)
		if valueDiags.HasError() {
			return nil, diags
		}
		result = value.ValueFloat64()
		return &result, diags
	}

	// Try Number (bigfloat)
	if numberValue, ok := matchedPathValue.(basetypes.NumberValuable); ok {
		var value types.Number
		valueDiags := tfsdk.ValueAs(ctx, numberValue, &value)
		diags.Append(valueDiags...)
		if valueDiags.HasError() {
			return nil, diags
		}
		bigFloat := value.ValueBigFloat()
		floatVal, _ := bigFloat.Float64()
		result = floatVal
		return &result, diags
	}

	// If we get here, the type is not numeric
	diags.AddError(
		"Invalid Attribute Type",
		fmt.Sprintf("The attribute at %s is not a numeric type (Int64, Float64, or Number).",
			expression.String()),
	)
	return nil, diags
}

// NumericGreaterThanOrEqual returns a validator that ensures the value at greaterOrEqualPath
// is greater than or equal to the value at thanPath. Both paths must point to numeric attributes
// (Int64, Float64, or Number). Values of -1 are treated as "not set" and will skip validation.
//
// Example usage:
//
//	validators.NumericGreaterThanOrEqual(
//	    path.MatchRoot("stateful_firewall").AtName("adaptive_timeouts").AtName("end"),
//	    path.MatchRoot("stateful_firewall").AtName("max").AtName("states"),
//	)
func NumericGreaterThanOrEqual(greaterOrEqualPath, thanPath path.Expression) resource.ConfigValidator {
	return numericGreaterThanOrEqualValidator{
		greaterOrEqualPath: greaterOrEqualPath,
		thanPath:           thanPath,
	}
}
