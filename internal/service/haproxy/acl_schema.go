package haproxy

import (
	"fmt"
	"strings"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/haproxy"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type aclResourceModel struct {
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Expression    types.String `tfsdk:"expression"`
	Value         types.String `tfsdk:"value"`
	Negate        types.Bool   `tfsdk:"negate"`
	CaseSensitive types.Bool   `tfsdk:"case_sensitive"`

	Id types.String `tfsdk:"id"`
}

// aclValueField returns the API field that holds the value for an expression.
// OPNsense keeps one field per expression, named like the expression.
func aclValueField(a *haproxy.Acl, expression string) *string {
	fields := map[string]*string{
		"hdr":        &a.Hdr,
		"hdr_beg":    &a.HdrBeg,
		"hdr_end":    &a.HdrEnd,
		"hdr_reg":    &a.HdrReg,
		"hdr_sub":    &a.HdrSub,
		"path":       &a.Path,
		"path_beg":   &a.PathBeg,
		"path_dir":   &a.PathDir,
		"path_end":   &a.PathEnd,
		"path_reg":   &a.PathReg,
		"path_sub":   &a.PathSub,
		"src":        &a.Src,
		"ssl_fc_sni": &a.SSLFcSNI,
		"ssl_sni":    &a.SSLSNI,
		"custom_acl": &a.CustomAcl,
	}
	return fields[expression]
}

// aclExpressions lists the expressions aclValueField knows.
var aclExpressions = []string{"custom_acl", "hdr", "hdr_beg", "hdr_end", "hdr_reg", "hdr_sub", "path", "path_beg",
	"path_dir", "path_end", "path_reg", "path_sub", "src", "ssl_fc_sni", "ssl_sni"}

func aclResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "A HAProxy condition (ACL). Rules (`opnsense_haproxy_action`) use conditions to decide when they apply. Requires the `os-haproxy` plugin.",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name to identify this condition. Must not contain `.`, `,`, `;`, brackets, braces or tabs.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for this condition. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"expression": schema.StringAttribute{
				MarkdownDescription: "HAProxy fetch the condition tests, e.g. `hdr_end` (Host header ends with), `path_beg` (path starts with), `src` (source IP) or `custom_acl` (raw ACL expression). One of `" + strings.Join(aclExpressions, "`, `") + "`.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf(aclExpressions...)},
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Value matched by `expression`, e.g. `.example.com` for `hdr_end` or `/.well-known/acme-challenge/` for `path_beg`.",
				Required:            true,
			},
			"negate": schema.BoolAttribute{
				MarkdownDescription: "When enabled, the result of the condition is inverted. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"case_sensitive": schema.BoolAttribute{
				MarkdownDescription: "When enabled, the match is case-sensitive. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the condition.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func convertAclSchemaToStruct(d *aclResourceModel) (*haproxy.Acl, error) {
	acl := &haproxy.Acl{
		Name:          d.Name.ValueString(),
		Description:   d.Description.ValueString(),
		Expression:    api.SelectedMap(d.Expression.ValueString()),
		Negate:        tools.BoolToString(d.Negate.ValueBool()),
		CaseSensitive: tools.BoolToString(d.CaseSensitive.ValueBool()),
	}
	field := aclValueField(acl, d.Expression.ValueString())
	if field == nil {
		return nil, fmt.Errorf("expression %q is not supported", d.Expression.ValueString())
	}
	*field = d.Value.ValueString()
	return acl, nil
}

func convertAclStructToSchema(d *haproxy.Acl) (*aclResourceModel, error) {
	field := aclValueField(d, d.Expression.String())
	if field == nil {
		return nil, fmt.Errorf("expression %q is not supported by this resource, supported: %v", d.Expression.String(), aclExpressions)
	}
	return &aclResourceModel{
		Name:          types.StringValue(d.Name),
		Description:   types.StringValue(d.Description),
		Expression:    types.StringValue(d.Expression.String()),
		Value:         types.StringValue(*field),
		Negate:        types.BoolValue(tools.StringToBool(d.Negate)),
		CaseSensitive: types.BoolValue(tools.StringToBool(d.CaseSensitive)),
	}, nil
}
