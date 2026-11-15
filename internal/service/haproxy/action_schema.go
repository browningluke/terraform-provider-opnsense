package haproxy

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/haproxy"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/browningluke/terraform-provider-opnsense/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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

	Id types.String `tfsdk:"id"`
}

func actionResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "A HAProxy rule (action): what HAProxy does when its conditions match. A rule takes effect once it is linked to a frontend with `opnsense_haproxy_frontend_actions`. Only the `use_backend` type is supported so far. Requires the `os-haproxy` plugin.",
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
				MarkdownDescription: "UUIDs of the conditions (`opnsense_haproxy_acl`) this rule tests.",
				ElementType:         types.StringType,
				Required:            true,
				Validators:          []validator.Set{setvalidator.SizeAtLeast(1)},
			},
			"operator": schema.StringAttribute{
				MarkdownDescription: "How multiple conditions are combined, `and` or `or`. Defaults to `and`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("and"),
				Validators:          []validator.String{stringvalidator.OneOf("and", "or")},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Rule type. Only `use_backend` is supported so far. Defaults to `use_backend`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("use_backend"),
				Validators:          []validator.String{stringvalidator.OneOf("use_backend")},
			},
			"use_backend": schema.StringAttribute{
				MarkdownDescription: "UUID of the backend pool (`opnsense_haproxy_backend`) that receives matching requests.",
				Required:            true,
				Validators:          []validator.String{validators.IsUUIDv4()},
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
	}, nil
}
