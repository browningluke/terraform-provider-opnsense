package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// requiresStringEqualsOneOfValidator is a resource-level validator that ensures
// when a dependent attribute is set, a condition attribute must equal one of the specified values.
type requiresStringEqualsOneOfValidator struct {
	dependentPath path.Expression
	conditionPath path.Expression
	validValues   []string
}

// Description returns a plain text description of the validator's behavior.
func (v requiresStringEqualsOneOfValidator) Description(_ context.Context) string {
	return fmt.Sprintf("When %q is set, %q must equal one of: %s",
		v.dependentPath.String(), v.conditionPath.String(), strings.Join(v.validValues, ", "))
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior.
func (v requiresStringEqualsOneOfValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateResource performs the validation.
func (v requiresStringEqualsOneOfValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	// Check if the dependent field is set
	isSet, diags := v.isDependentFieldSet(ctx, req, v.dependentPath)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// If dependent field is not set, no validation needed
	if !isSet {
		return
	}

	// Dependent field is set, now check if condition field has a valid value
	conditionValue, diags := v.getStringValue(ctx, req, v.conditionPath)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// If condition value is not set, that's an error
	if conditionValue == nil {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			fmt.Sprintf("When %s is set, %s must be specified and equal one of: %s",
				v.dependentPath.String(),
				v.conditionPath.String(),
				strings.Join(v.validValues, ", "),
			),
		)
		return
	}

	// Check if the condition value is one of the valid values
	for _, validValue := range v.validValues {
		if *conditionValue == validValue {
			return // Valid!
		}
	}

	// If we get here, the condition value doesn't match any valid values
	resp.Diagnostics.AddError(
		"Invalid Attribute Combination",
		fmt.Sprintf("When %s is set, %s must equal one of: %s. Got: %q",
			v.dependentPath.String(),
			v.conditionPath.String(),
			strings.Join(v.validValues, ", "),
			*conditionValue,
		),
	)
}

// isDependentFieldSet checks if the dependent field is set (non-default/non-empty).
// Supports all Terraform types: Int64, String, Bool, Float64, Number, Set, List, Map, Object.
// Returns false if the field is null, unknown, or considered "not set" based on type-specific logic.
func (v requiresStringEqualsOneOfValidator) isDependentFieldSet(ctx context.Context, req resource.ValidateConfigRequest, expression path.Expression) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Find paths matching the expression
	matchedPaths, pathDiags := req.Config.PathMatches(ctx, expression)
	diags.Append(pathDiags...)
	if pathDiags.HasError() {
		return false, diags
	}

	// We expect exactly one match for our use case
	if len(matchedPaths) == 0 {
		// Path doesn't exist, which is okay - it might be optional
		return false, diags
	}

	// Use the first matched path
	matchedPath := matchedPaths[0]

	// Get the generic attr.Value at the matched path
	var matchedPathValue attr.Value
	getDiags := req.Config.GetAttribute(ctx, matchedPath, &matchedPathValue)
	diags.Append(getDiags...)
	if getDiags.HasError() {
		return false, diags
	}

	// If the value is null or unknown, it's not set
	if matchedPathValue.IsNull() || matchedPathValue.IsUnknown() {
		return false, diags
	}

	// Type-specific logic to determine if the field is "set"
	switch value := matchedPathValue.(type) {
	case types.Int64:
		// Int64: Not set if value is -1 (common convention for "unset" in this codebase)
		return value.ValueInt64() != -1, diags

	case types.String:
		// String: Not set if empty string
		return value.ValueString() != "", diags

	case types.Bool:
		// Bool: Always considered set if not null/unknown (false is valid)
		return true, diags

	case types.Float64:
		// Float64: Always considered set if not null/unknown
		return true, diags

	case types.Number:
		// Number: Always considered set if not null/unknown
		return true, diags

	case types.Set:
		// Set: Not set if empty
		return len(value.Elements()) > 0, diags

	case types.List:
		// List: Not set if empty
		return len(value.Elements()) > 0, diags

	case types.Map:
		// Map: Not set if empty
		return len(value.Elements()) > 0, diags

	case types.Object:
		// Object: Always considered set if not null/unknown
		return true, diags

	default:
		// Unknown type - treat as set to be safe and let other validators handle it
		return true, diags
	}
}

// getStringValue extracts a String value from the given path expression.
// Returns nil if the value is null or unknown.
func (v requiresStringEqualsOneOfValidator) getStringValue(ctx context.Context, req resource.ValidateConfigRequest, expression path.Expression) (*string, diag.Diagnostics) {
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

	// If the value is null or unknown, we cannot validate
	if matchedPathValue.IsNull() || matchedPathValue.IsUnknown() {
		return nil, diags
	}

	// Convert to types.String
	var value types.String
	valueDiags := tfsdk.ValueAs(ctx, matchedPathValue, &value)
	diags.Append(valueDiags...)
	if valueDiags.HasError() {
		return nil, diags
	}

	// Extract the string value
	stringValue := value.ValueString()
	return &stringValue, diags
}

// RequiresStringEqualsOneOf returns a validator that ensures when the dependent
// field is set, the condition field must equal one of the specified valid values.
//
// Supports all Terraform types for the dependent field:
//   - Int64: Not set if null/unknown or -1
//   - String: Not set if null/unknown or empty
//   - Bool: Not set if null/unknown (false is valid)
//   - Float64, Number: Not set if null/unknown
//   - Set, List: Not set if null/unknown or empty
//   - Map: Not set if null/unknown or empty
//   - Object: Not set if null/unknown
//
// The condition field must always be a String type.
//
// Example 1: Ensure max.states (Int64) can only be set when filter.protocol is TCP:
//
//	validators.RequiresStringEqualsOneOf(
//	    path.MatchRoot("stateful_firewall").AtName("max").AtName("states"),
//	    path.MatchRoot("filter").AtName("protocol"),
//	    []string{"TCP", "TCP/UDP"},
//	)
//
// Example 2: Ensure icmp_type (Set) can only be set when filter.protocol is ICMP:
//
//	validators.RequiresStringEqualsOneOf(
//	    path.MatchRoot("filter").AtName("icmp_type"),
//	    path.MatchRoot("filter").AtName("protocol"),
//	    []string{"ICMP"},
//	)
func RequiresStringEqualsOneOf(dependentPath, conditionPath path.Expression, validValues []string) resource.ConfigValidator {
	return requiresStringEqualsOneOfValidator{
		dependentPath: dependentPath,
		conditionPath: conditionPath,
		validValues:   validValues,
	}
}
