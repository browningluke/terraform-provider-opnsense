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

// QuaggaBGPCommunityListResourceModel describes the resource data model.
type QuaggaBGPCommunityListResourceModel struct {
	Enabled        types.Bool   `tfsdk:"enabled"`
	Description    types.String `tfsdk:"description"`
	Number         types.Int64  `tfsdk:"number"`
	SequenceNumber types.Int64  `tfsdk:"seq_number"`
	Action         types.String `tfsdk:"action"`
	Community      types.String `tfsdk:"community"`

	Id types.String `tfsdk:"id"`
}

func quaggaBGPCommunityListResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure community lists for BGP.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this community list. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "An optional description for this prefix list. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"number": schema.Int64Attribute{
				MarkdownDescription: "Set the number of your Community-List. 1-99 are standard lists while 100-500 are expanded lists.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 500),
				},
			},
			"seq_number": schema.Int64Attribute{
				MarkdownDescription: "The ACL sequence number (10-99).",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.Between(10, 99),
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
			"community": schema.StringAttribute{
				MarkdownDescription: "The community you want to match. You can also regex and it is not validated so please be careful.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the community list.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func QuaggaBGPCommunityListDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure community lists for BGP.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this community list.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "An optional description for this prefix list.",
				Computed:            true,
			},
			"number": schema.Int64Attribute{
				MarkdownDescription: "Set the number of your Community-List. 1-99 are standard lists while 100-500 are expanded lists.",
				Computed:            true,
			},
			"seq_number": schema.Int64Attribute{
				MarkdownDescription: "The ACL sequence number (10-99).",
				Computed:            true,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Set permit for match or deny to negate the rule.",
				Computed:            true,
			},
			"community": schema.StringAttribute{
				MarkdownDescription: "The community you want to match. You can also regex and it is not validated so please be careful.",
				Computed:            true,
			},
		},
	}
}

func convertQuaggaBGPCommunityListSchemaToStruct(d *QuaggaBGPCommunityListResourceModel) (*quagga.BGPCommunityList, error) {
	return &quagga.BGPCommunityList{
		Enabled:        tools.BoolToString(d.Enabled.ValueBool()),
		Description:    d.Description.ValueString(),
		Number:         tools.Int64ToString(d.Number.ValueInt64()),
		SequenceNumber: tools.Int64ToString(d.SequenceNumber.ValueInt64()),
		Action:         api.SelectedMap(d.Action.ValueString()),
		Community:      d.Community.ValueString(),
	}, nil
}

func convertQuaggaBGPCommunityListStructToSchema(d *quagga.BGPCommunityList) (*QuaggaBGPCommunityListResourceModel, error) {
	return &QuaggaBGPCommunityListResourceModel{
		Enabled:        types.BoolValue(tools.StringToBool(d.Enabled)),
		Description:    types.StringValue(d.Description),
		Number:         types.Int64Value(tools.StringToInt64(d.Number)),
		SequenceNumber: types.Int64Value(tools.StringToInt64(d.SequenceNumber)),
		Action:         types.StringValue(d.Action.String()),
		Community:      types.StringValue(d.Community),
	}, nil
}
