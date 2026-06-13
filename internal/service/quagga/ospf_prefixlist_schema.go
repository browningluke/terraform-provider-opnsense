package quagga

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ospfPrefixListResourceModel describes the resource data model.
type ospfPrefixListResourceModel struct {
	Enabled        types.Bool   `tfsdk:"enabled"`
	Name           types.String `tfsdk:"name"`
	SequenceNumber types.String `tfsdk:"sequence_number"`
	Action         types.String `tfsdk:"action"`
	Network        types.String `tfsdk:"network"`

	Id types.String `tfsdk:"id"`
}

func ospfPrefixListResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure OSPF prefix lists.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this prefix list. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of this prefix list.",
				Required:            true,
			},
			"sequence_number": schema.StringAttribute{
				MarkdownDescription: "The sequence number. Defaults to `\"\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Set permit for match or deny to negate the rule. One of `\"permit\"`, `\"deny\"`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("permit", "deny"),
				},
			},
			"network": schema.StringAttribute{
				MarkdownDescription: "The network to match.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the OSPF prefix list.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func ospfPrefixListDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure OSPF prefix lists.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this prefix list.",
				Computed:            true,
			},
			"name": dschema.StringAttribute{
				MarkdownDescription: "The name of this prefix list.",
				Computed:            true,
			},
			"sequence_number": dschema.StringAttribute{
				MarkdownDescription: "The sequence number.",
				Computed:            true,
			},
			"action": dschema.StringAttribute{
				MarkdownDescription: "Set permit for match or deny to negate the rule.",
				Computed:            true,
			},
			"network": dschema.StringAttribute{
				MarkdownDescription: "The network to match.",
				Computed:            true,
			},
		},
	}
}

func convertOSPFPrefixListSchemaToStruct(d *ospfPrefixListResourceModel) (*quagga.OSPFPrefixList, error) {
	return &quagga.OSPFPrefixList{
		Enabled:        tools.BoolToString(d.Enabled.ValueBool()),
		Name:           d.Name.ValueString(),
		SequenceNumber: d.SequenceNumber.ValueString(),
		Action:         api.SelectedMap(d.Action.ValueString()),
		Network:        d.Network.ValueString(),
	}, nil
}

func convertOSPFPrefixListStructToSchema(d *quagga.OSPFPrefixList) (*ospfPrefixListResourceModel, error) {
	return &ospfPrefixListResourceModel{
		Enabled:        types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:           types.StringValue(d.Name),
		SequenceNumber: types.StringValue(d.SequenceNumber),
		Action:         types.StringValue(d.Action.String()),
		Network:        types.StringValue(d.Network),
	}, nil
}
