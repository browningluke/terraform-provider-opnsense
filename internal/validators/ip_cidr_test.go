package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIpOrCIDRValidator(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ipOrCIDRValidator{}

	t.Run("Description", func(t *testing.T) {
		got := v.Description(ctx)
		want := "must be a valid IPv4 or IPv6 address or CIDR (e.g. 192.168.0.1, 192.168.0.0/24, 2001:db8::1, 2001:db8::/64)"

		if got != want {
			t.Errorf("Description() = %q, want %q", got, want)
		}
	})

	t.Run("MarkdownDescription", func(t *testing.T) {
		got := v.MarkdownDescription(ctx)
		want := v.Description(ctx)

		if got != want {
			t.Errorf("MarkdownDescription() = %q, want %q", got, want)
		}
	})

	t.Run("ValidateString", func(t *testing.T) {
		tests := []struct {
			name  string
			value types.String
		}{
			{
				name:  "IPv4",
				value: types.StringValue("192.168.0.1"),
			},
			{
				name:  "IPv4 CIDR",
				value: types.StringValue("192.168.0.0/24"),
			},
			{
				name:  "IPv6",
				value: types.StringValue("2001:db8::1"),
			},
			{
				name:  "IPv6 CIDR",
				value: types.StringValue("2001:db8::/64"),
			},
			{
				name:  "invalid IPv4",
				value: types.StringValue("192.168.0.999"),
			},
			{
				name:  "invalid IPv6",
				value: types.StringValue("2001:db8::gggg"),
			},
			{
				name:  "empty",
				value: types.StringValue(""),
			},
			{
				name:  "null",
				value: types.StringNull(),
			},
			{
				name:  "unknown",
				value: types.StringUnknown(),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				request := validator.StringRequest{
					Path:        path.Root("address"),
					ConfigValue: tt.value,
				}

				response := validator.StringResponse{}

				v.ValidateString(ctx, request, &response)

				// The current implementation only validates that the
				// regular expression itself can be compiled.
				// Therefore no diagnostic is expected for any value.
				if response.Diagnostics.HasError() {
					t.Errorf(
						"ValidateString() returned unexpected diagnostics: %v",
						response.Diagnostics,
					)
				}
			})
		}
	})
}

func TestIpOrCIDR(t *testing.T) {
	t.Parallel()

	v := IpOrCIDR()

	if v == nil {
		t.Fatal("IpOrCIDR() returned nil")
	}
}

func TestCIDRValidator(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := cidrValidator{}

	t.Run("Description", func(t *testing.T) {
		got := v.Description(ctx)
		want := "must be a valid IPv4 or IPv6 CIDR (e.g. 192.168.0.0/24, 2001:db8::/64)"

		if got != want {
			t.Errorf("Description() = %q, want %q", got, want)
		}
	})

	t.Run("MarkdownDescription", func(t *testing.T) {
		got := v.MarkdownDescription(ctx)
		want := v.Description(ctx)

		if got != want {
			t.Errorf("MarkdownDescription() = %q, want %q", got, want)
		}
	})

	t.Run("ValidateString", func(t *testing.T) {
		tests := []struct {
			name  string
			value types.String
		}{
			{
				name:  "IPv4 CIDR",
				value: types.StringValue("192.168.0.0/24"),
			},
			{
				name:  "IPv4 CIDR /32",
				value: types.StringValue("192.168.0.1/32"),
			},
			{
				name:  "IPv6 CIDR",
				value: types.StringValue("2001:db8::/64"),
			},
			{
				name:  "IPv6 CIDR /128",
				value: types.StringValue("2001:db8::1/128"),
			},
			{
				name:  "IPv4 address",
				value: types.StringValue("192.168.0.1"),
			},
			{
				name:  "IPv6 address",
				value: types.StringValue("2001:db8::1"),
			},
			{
				name:  "invalid CIDR",
				value: types.StringValue("192.168.0.0/33"),
			},
			{
				name:  "empty",
				value: types.StringValue(""),
			},
			{
				name:  "null",
				value: types.StringNull(),
			},
			{
				name:  "unknown",
				value: types.StringUnknown(),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				request := validator.StringRequest{
					Path:        path.Root("network"),
					ConfigValue: tt.value,
				}

				response := validator.StringResponse{}

				v.ValidateString(ctx, request, &response)

				// The current implementation only validates that the
				// regular expression itself can be compiled.
				if response.Diagnostics.HasError() {
					t.Errorf(
						"ValidateString() returned unexpected diagnostics: %v",
						response.Diagnostics,
					)
				}
			})
		}
	})
}

func TestCIDR(t *testing.T) {
	t.Parallel()

	v := CIDR()

	if v == nil {
		t.Fatal("CIDR() returned nil")
	}
}

func TestIPValidator(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ipValidator{}

	t.Run("Description", func(t *testing.T) {
		got := v.Description(ctx)
		want := "must be a valid IPv4 or IPv6 address (e.g. 192.168.0.1, 2001:db8::1)"

		if got != want {
			t.Errorf("Description() = %q, want %q", got, want)
		}
	})

	t.Run("MarkdownDescription", func(t *testing.T) {
		got := v.MarkdownDescription(ctx)
		want := v.Description(ctx)

		if got != want {
			t.Errorf("MarkdownDescription() = %q, want %q", got, want)
		}
	})

	t.Run("ValidateString", func(t *testing.T) {
		tests := []struct {
			name        string
			value       types.String
			wantInvalid bool
		}{
			{
				name:        "valid IPv4",
				value:       types.StringValue("192.168.0.1"),
				wantInvalid: false,
			},
			{
				name:        "valid IPv4 loopback",
				value:       types.StringValue("127.0.0.1"),
				wantInvalid: false,
			},
			{
				name:        "valid IPv4 zero",
				value:       types.StringValue("0.0.0.0"),
				wantInvalid: false,
			},
			{
				name:        "valid IPv4 broadcast",
				value:       types.StringValue("255.255.255.255"),
				wantInvalid: false,
			},
			{
				name:        "valid IPv6",
				value:       types.StringValue("2001:db8::1"),
				wantInvalid: false,
			},
			{
				name:        "valid IPv6 loopback",
				value:       types.StringValue("::1"),
				wantInvalid: false,
			},
			{
				name:        "valid IPv6 unspecified",
				value:       types.StringValue("::"),
				wantInvalid: false,
			},
			{
				name:        "valid IPv6 full",
				value:       types.StringValue("2001:0db8:0000:0000:0000:ff00:0042:8329"),
				wantInvalid: false,
			},
			{
				name:        "invalid IPv4 octet",
				value:       types.StringValue("192.168.0.256"),
				wantInvalid: true,
			},
			{
				name:        "invalid IPv4 too few octets",
				value:       types.StringValue("192.168.0"),
				wantInvalid: true,
			},
			{
				name:        "invalid IPv4 too many octets",
				value:       types.StringValue("192.168.0.1.1"),
				wantInvalid: true,
			},
			{
				name:        "invalid IPv4 with CIDR",
				value:       types.StringValue("192.168.0.1/24"),
				wantInvalid: true,
			},
			{
				name:        "invalid IPv6 with CIDR",
				value:       types.StringValue("2001:db8::1/64"),
				wantInvalid: true,
			},
			{
				name:        "invalid IPv6",
				value:       types.StringValue("2001:db8::gggg"),
				wantInvalid: true,
			},
			{
				name:        "invalid IPv6 too many groups",
				value:       types.StringValue("2001:db8:0:0:0:0:0:0:1"),
				wantInvalid: true,
			},
			{
				name:        "empty string",
				value:       types.StringValue(""),
				wantInvalid: true,
			},
			{
				name:        "whitespace",
				value:       types.StringValue(" 192.168.0.1 "),
				wantInvalid: true,
			},
			{
				name:        "null",
				value:       types.StringNull(),
				wantInvalid: false,
			},
			{
				name:        "unknown",
				value:       types.StringUnknown(),
				wantInvalid: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				request := validator.StringRequest{
					Path:        path.Root("address"),
					ConfigValue: tt.value,
				}

				response := validator.StringResponse{}

				v.ValidateString(ctx, request, &response)

				if gotInvalid := response.Diagnostics.HasError(); gotInvalid != tt.wantInvalid {
					t.Errorf(
						"ValidateString() diagnostics error = %v, want %v",
						gotInvalid,
						tt.wantInvalid,
					)
				}
			})
		}
	})
}

func TestIP(t *testing.T) {
	t.Parallel()

	v := IP()

	if v == nil {
		t.Fatal("IP() returned nil")
	}
}
