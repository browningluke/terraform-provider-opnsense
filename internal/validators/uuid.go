package validators

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type uuidv4Validator struct{}

func (v uuidv4Validator) Description(_ context.Context) string {
	return "must be a valid UUIDv4 (e.g. 1ae521bb-05e2-43c1-8e1f-7e34b53dc015)"
}

func (v uuidv4Validator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v uuidv4Validator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	// Allow empty strings (for default values)
	if value == "" {
		return
	}

	// UUIDv4 format: xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
	// where y is one of [8, 9, a, b]
	matched, _ := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, value)
	if !matched {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

func IsUUIDv4() validator.String {
	return uuidv4Validator{}
}
