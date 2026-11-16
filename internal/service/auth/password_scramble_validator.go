package auth

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// passwordScrambleConflictValidator rejects configs that set both a password
// and scramble_password = true (default), which would discard the password.
type passwordScrambleConflictValidator struct{}

func (v passwordScrambleConflictValidator) Description(ctx context.Context) string {
	return "password cannot be set when scramble_password is true"
}

func (v passwordScrambleConflictValidator) MarkdownDescription(ctx context.Context) string {
	return "password cannot be set when scramble_password is true"
}

func (v passwordScrambleConflictValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var password types.String
	var scramble types.Bool

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password"), &password)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("scramble_password"), &scramble)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if password.IsNull() || password.IsUnknown() || password.ValueString() == "" {
		return
	}

	// Default scramble_password is true when the attribute is omitted.
	scrambleOn := scramble.IsNull() || scramble.IsUnknown() || scramble.ValueBool()
	if !scrambleOn {
		return
	}

	resp.Diagnostics.AddAttributeError(
		path.Root("password"),
		"Conflicting password and scramble_password",
		"Set scramble_password = false when providing password, or omit password to keep the default scrambled password.",
	)
}
