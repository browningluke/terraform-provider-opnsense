package validators

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

type ipOrCIDRValidator struct{}

func (validator ipOrCIDRValidator) Description(_ context.Context) string {
	return "must be a valid IPv4 or IPv6 address or CIDR (e.g. 192.168.0.1, 192.168.0.0/24, 2001:db8::1, 2001:db8::/64)"
}

func (validator ipOrCIDRValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (validator ipOrCIDRValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if _, err := regexp.Compile(`^(([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))?|([0-9a-fA-F:]+)(\/([0-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8]))?)$`); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			validator.Description(ctx),
			request.ConfigValue.ValueString(),
		))
		return
	}
}

func IpOrCIDR() validator.String {
	return ipOrCIDRValidator{}
}

type cidrValidator struct{}

func (validator cidrValidator) Description(_ context.Context) string {
	return "must be a valid IPv4 or IPv6 CIDR (e.g. 192.168.0.0/24, 2001:db8::/64)"
}

func (validator cidrValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (validator cidrValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if _, err := regexp.Compile(`^(([0-9]{1,3}\.){3}[0-9]{1,3}\/(3[0-2]|[1-2]?[0-9]))$|^(([0-9a-fA-F:]+)\/(12[0-8]|1[0-1][0-9]|[1-9]?[0-9]))$`); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			validator.Description(ctx),
			request.ConfigValue.ValueString(),
		))
		return
	}
}

func CIDR() validator.String {
	return cidrValidator{}
}
