package service

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/tools"
)

// QuaggaBGPASPathResourceModel describes the resource data model.
type QuaggaBGPASPathResourceModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Description types.String `tfsdk:"description"`
	Number      types.Int64  `tfsdk:"number"`
	Action      types.String `tfsdk:"action"`
	AS          types.String `tfsdk:"as"`

	Id types.String `tfsdk:"id"`
}

func quaggaBGPASPathResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure AS Path lists for BGP.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this AS path. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "An optional description for this AS path. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"number": schema.Int64Attribute{
				MarkdownDescription: "The ACL rule number (0-4294967294); keep in mind that there are no sequence numbers with AS-Path lists. When you want to add a new line between you have to completely remove the ACL!",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967294),
				},
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Set permit for match or deny to negate the rule. Defaults to `\"permit\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("permit"),
				Validators: []validator.String{
					stringvalidator.OneOf("permit", "deny"),
				},
			},
			"as": schema.StringAttribute{
				MarkdownDescription: "The AS pattern you want to match, regexp allowed (e.g. `.$` or `_1$`). It's not validated so please be careful!",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the AS path.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func QuaggaBGPASPathDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure AS Path lists for BGP.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this AS path.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "An optional description for this AS path.",
				Computed:            true,
			},
			"number": dschema.Int64Attribute{
				MarkdownDescription: "The ACL rule number (0-4294967294); keep in mind that there are no sequence numbers with AS-Path lists. When you want to add a new line between you have to completely remove the ACL!",
				Computed:            true,
			},
			"action": dschema.StringAttribute{
				MarkdownDescription: "Set permit for match or deny to negate the rule.",
				Computed:            true,
			},
			"as": dschema.StringAttribute{
				MarkdownDescription: "The AS pattern you want to match, regexp allowed (e.g. `.$` or `_1$`). It's not validated so please be careful!",
				Computed:            true,
			},
		},
	}
}

func convertQuaggaBGPASPathSchemaToStruct(d *QuaggaBGPASPathResourceModel) (*quagga.BGPASPath, error) {
	return &quagga.BGPASPath{
		Enabled:     tools.BoolToString(d.Enabled.ValueBool()),
		Description: d.Description.ValueString(),
		Number:      tools.Int64ToString(d.Number.ValueInt64()),
		Action:      api.SelectedMap(d.Action.ValueString()),
		AS:          d.AS.ValueString(),
	}, nil
}

func convertQuaggaBGPASPathStructToSchema(d *quagga.BGPASPath) (*QuaggaBGPASPathResourceModel, error) {
	return &QuaggaBGPASPathResourceModel{
		Enabled:     types.BoolValue(tools.StringToBool(d.Enabled)),
		Description: types.StringValue(d.Description),
		Number:      types.Int64Value(tools.StringToInt64(d.Number)),
		Action:      types.StringValue(d.Action.String()),
		AS:          types.StringValue(d.AS),
	}, nil
}
