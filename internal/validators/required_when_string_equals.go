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

// requiredWhenStringEqualsOneOfValidator is a resource-level validator that ensures
// when a condition attribute equals one of the specified values, a dependent attribute must be set.
type requiredWhenStringEqualsOneOfValidator struct {
	dependentPath path.Expression
	conditionPath path.Expression
	triggerValues []string
}

// Description returns a plain text description of the validator's behavior.
func (v requiredWhenStringEqualsOneOfValidator) Description(_ context.Context) string {
	return fmt.Sprintf("When %q equals one of: %s, %q must be set",
		v.conditionPath.String(), strings.Join(v.triggerValues, ", "), v.dependentPath.String())
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior.
func (v requiredWhenStringEqualsOneOfValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateResource performs the validation.
func (v requiredWhenStringEqualsOneOfValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	// Read the condition value
	conditionValue, diags := v.getStringValue(ctx, req, v.conditionPath)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Condition unknown/missing — nothing to enforce here
	if conditionValue == nil {
		return
	}

	// Check whether the condition value triggers the requirement
	triggered := false
	for _, t := range v.triggerValues {
		if *conditionValue == t {
			triggered = true
			break
		}
	}
	if !triggered {
		return
	}

	// Condition is triggered — dependent must be set
	isSet, diags := v.isDependentFieldSet(ctx, req, v.dependentPath)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	if isSet {
		return
	}

	resp.Diagnostics.AddError(
		"Missing Required Attribute",
		fmt.Sprintf("%s must be set when %s equals one of: %s",
			v.dependentPath.String(),
			v.conditionPath.String(),
			strings.Join(v.triggerValues, ", "),
		),
	)
}

// isDependentFieldSet checks if the dependent field is set (non-default/non-empty).
// Supports all Terraform types: Int64, String, Bool, Float64, Number, Set, List, Map, Object.
// Returns false if the field is null, unknown, or considered "not set" based on type-specific logic.
func (v requiredWhenStringEqualsOneOfValidator) isDependentFieldSet(ctx context.Context, req resource.ValidateConfigRequest, expression path.Expression) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	matchedPaths, pathDiags := req.Config.PathMatches(ctx, expression)
	diags.Append(pathDiags...)
	if pathDiags.HasError() {
		return false, diags
	}
	if len(matchedPaths) == 0 {
		return false, diags
	}

	var matchedPathValue attr.Value
	getDiags := req.Config.GetAttribute(ctx, matchedPaths[0], &matchedPathValue)
	diags.Append(getDiags...)
	if getDiags.HasError() {
		return false, diags
	}

	if matchedPathValue.IsNull() || matchedPathValue.IsUnknown() {
		return false, diags
	}

	switch value := matchedPathValue.(type) {
	case types.Int64:
		return value.ValueInt64() != -1, diags
	case types.String:
		return value.ValueString() != "", diags
	case types.Bool:
		return true, diags
	case types.Float64:
		return true, diags
	case types.Number:
		return true, diags
	case types.Set:
		return len(value.Elements()) > 0, diags
	case types.List:
		return len(value.Elements()) > 0, diags
	case types.Map:
		return len(value.Elements()) > 0, diags
	case types.Object:
		return true, diags
	default:
		return true, diags
	}
}

// getStringValue extracts a String value from the given path expression.
// Returns nil if the value is null or unknown.
func (v requiredWhenStringEqualsOneOfValidator) getStringValue(ctx context.Context, req resource.ValidateConfigRequest, expression path.Expression) (*string, diag.Diagnostics) {
	var diags diag.Diagnostics

	matchedPaths, pathDiags := req.Config.PathMatches(ctx, expression)
	diags.Append(pathDiags...)
	if pathDiags.HasError() {
		return nil, diags
	}
	if len(matchedPaths) == 0 {
		return nil, diags
	}

	var matchedPathValue attr.Value
	getDiags := req.Config.GetAttribute(ctx, matchedPaths[0], &matchedPathValue)
	diags.Append(getDiags...)
	if getDiags.HasError() {
		return nil, diags
	}

	if matchedPathValue.IsNull() || matchedPathValue.IsUnknown() {
		return nil, diags
	}

	var value types.String
	valueDiags := tfsdk.ValueAs(ctx, matchedPathValue, &value)
	diags.Append(valueDiags...)
	if valueDiags.HasError() {
		return nil, diags
	}

	stringValue := value.ValueString()
	return &stringValue, diags
}

// RequiredWhenStringEqualsOneOf returns a validator that ensures when the condition
// field equals one of the specified trigger values, the dependent field must be set.
//
// Supports all Terraform types for the dependent field; see RequiresStringEqualsOneOf
// for the "is set" semantics applied per type. The condition field must be a String.
//
// Example: Ensure path_expression is set when type is "urljson":
//
//	validators.RequiredWhenStringEqualsOneOf(
//	    path.MatchRoot("path_expression"),
//	    path.MatchRoot("type"),
//	    []string{"urljson"},
//	)
func RequiredWhenStringEqualsOneOf(dependentPath, conditionPath path.Expression, triggerValues []string) resource.ConfigValidator {
	return requiredWhenStringEqualsOneOfValidator{
		dependentPath: dependentPath,
		conditionPath: conditionPath,
		triggerValues: triggerValues,
	}
}
