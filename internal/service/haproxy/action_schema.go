package haproxy

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/haproxy"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/browningluke/terraform-provider-opnsense/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type actionResourceModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	TestType    types.String `tfsdk:"test_type"`
	LinkedAcls  types.Set    `tfsdk:"linked_acls"`
	Operator    types.String `tfsdk:"operator"`
	Type        types.String `tfsdk:"type"`
	UseBackend  types.String `tfsdk:"use_backend"`

	HttpRequestAction types.String `tfsdk:"http_request_action"`
	HttpRequestOption types.String `tfsdk:"http_request_option"`

	Id types.String `tfsdk:"id"`
}

func actionResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "A HAProxy rule (action): what HAProxy does when its conditions match. A rule takes effect once it is linked to a frontend with `opnsense_haproxy_frontend_actions`. Supported types are `use_backend` and `http-request` with the `redirect` keyword. HAProxy runs all `http-request` rules of a frontend before any `use_backend` rule, whatever their order in the list. Requires the `os-haproxy` plugin.",
		Version:             1,

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this rule. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name to identify this rule. Must not contain `.`, `,`, `;`, brackets, braces or tabs.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for this rule. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"test_type": schema.StringAttribute{
				MarkdownDescription: "`if` runs the rule when the conditions match, `unless` when they do not. Defaults to `if`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("if"),
				Validators:          []validator.String{stringvalidator.OneOf("if", "unless")},
			},
			"linked_acls": schema.SetAttribute{
				MarkdownDescription: "UUIDs of the conditions (`opnsense_haproxy_acl`) this rule tests. Empty means the rule applies to every request. Defaults to `[]`.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"operator": schema.StringAttribute{
				MarkdownDescription: "How multiple conditions are combined, `and` or `or`. Defaults to `and`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("and"),
				Validators:          []validator.String{stringvalidator.OneOf("and", "or")},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Rule type: `use_backend` sends matching requests to `use_backend`, `http-request` runs `http_request_action`.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("use_backend", "http-request")},
			},
			"use_backend": schema.StringAttribute{
				MarkdownDescription: "UUID of the backend pool (`opnsense_haproxy_backend`) that receives matching requests. Required for type `use_backend`. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators:          []validator.String{validators.IsUUIDv4()},
			},
			"http_request_action": schema.StringAttribute{
				MarkdownDescription: "`http-request` keyword. Only `redirect` is supported so far. Required for type `http-request`. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators:          []validator.String{stringvalidator.OneOf("", "redirect")},
			},
			"http_request_option": schema.StringAttribute{
				MarkdownDescription: "Arguments of the `http-request` keyword, e.g. `scheme https code 301` for a permanent redirect to HTTPS with the same host and path. Required for type `http-request`. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func convertActionSchemaToStruct(d *actionResourceModel) (*haproxy.Action, error) {
	return &haproxy.Action{
		Enabled:     tools.BoolToString(d.Enabled.ValueBool()),
		Name:        d.Name.ValueString(),
		Description: d.Description.ValueString(),
		TestType:    api.SelectedMap(d.TestType.ValueString()),
		LinkedAcls:  tools.SetToStringSlice(d.LinkedAcls),
		Operator:    api.SelectedMap(d.Operator.ValueString()),
		Type:        api.SelectedMap(d.Type.ValueString()),
		UseBackend:  api.SelectedMap(d.UseBackend.ValueString()),

		HttpRequestAction: api.SelectedMap(d.HttpRequestAction.ValueString()),
		HttpRequestOption: d.HttpRequestOption.ValueString(),
	}, nil
}

func convertActionStructToSchema(d *haproxy.Action) (*actionResourceModel, error) {
	return &actionResourceModel{
		Enabled:     types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:        types.StringValue(d.Name),
		Description: types.StringValue(d.Description),
		TestType:    types.StringValue(d.TestType.String()),
		LinkedAcls:  tools.StringSliceToSet(d.LinkedAcls),
		Operator:    types.StringValue(d.Operator.String()),
		Type:        types.StringValue(d.Type.String()),
		UseBackend:  types.StringValue(d.UseBackend.String()),

		HttpRequestAction: types.StringValue(d.HttpRequestAction.String()),
		HttpRequestOption: types.StringValue(d.HttpRequestOption),
	}, nil
}
