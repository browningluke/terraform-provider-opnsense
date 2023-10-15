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

// QuaggaBGPPrefixListResourceModel describes the resource data model.
type QuaggaBGPPrefixListResourceModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
	IPVersion   types.String `tfsdk:"ip_version"`
	Number      types.Int64  `tfsdk:"number"`
	Action      types.String `tfsdk:"action"`
	Network     types.String `tfsdk:"network"`

	Id types.String `tfsdk:"id"`
}

func quaggaBGPPrefixListResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure prefix lists for BGP.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this prefix list. Defaults to `true`.",
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
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of this prefix list.",
				Required:            true,
			},
			"ip_version": schema.StringAttribute{
				MarkdownDescription: "Set the IP version to use. Defaults to `\"IPv4\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("IPv4"),
				Validators: []validator.String{
					stringvalidator.OneOf("IPv4", "IPv6"),
				},
			},
			"number": schema.Int64Attribute{
				MarkdownDescription: "The ACL sequence number (1-4294967294).",
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
			"network": schema.StringAttribute{
				MarkdownDescription: "The network pattern you want to match. You can also add \"ge\" or \"le\" additions after the network statement. It's not validated so please be careful!",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the prefix list.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func QuaggaBGPPrefixListDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure prefix lists for BGP.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this prefix list.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "An optional description for this prefix list.",
				Computed:            true,
			},
			"name": dschema.StringAttribute{
				MarkdownDescription: "The name of this prefix list.",
				Computed:            true,
			},
			"ip_version": dschema.StringAttribute{
				MarkdownDescription: "Set the IP version to use.",
				Computed:            true,
			},
			"number": dschema.Int64Attribute{
				MarkdownDescription: "The ACL sequence number (1-4294967294).",
				Computed:            true,
			},
			"action": dschema.StringAttribute{
				MarkdownDescription: "Set permit for match or deny to negate the rule.",
				Computed:            true,
			},
			"network": dschema.StringAttribute{
				MarkdownDescription: "The network pattern you want to match. You can also add \"ge\" or \"le\" additions after the network statement. It's not validated so please be careful!",
				Computed:            true,
			},
		},
	}
}

func convertQuaggaBGPPrefixListSchemaToStruct(d *QuaggaBGPPrefixListResourceModel) (*quagga.BGPPrefixList, error) {
	return &quagga.BGPPrefixList{
		Enabled:        tools.BoolToString(d.Enabled.ValueBool()),
		Description:    d.Description.ValueString(),
		Name:           d.Name.ValueString(),
		IPVersion:      api.SelectedMap(d.IPVersion.ValueString()),
		SequenceNumber: tools.Int64ToString(d.Number.ValueInt64()),
		Action:         api.SelectedMap(d.Action.ValueString()),
		Network:        d.Network.ValueString(),
	}, nil
}

func convertQuaggaBGPPrefixListStructToSchema(d *quagga.BGPPrefixList) (*QuaggaBGPPrefixListResourceModel, error) {
	return &QuaggaBGPPrefixListResourceModel{
		Enabled:     types.BoolValue(tools.StringToBool(d.Enabled)),
		Description: types.StringValue(d.Description),
		Name:        types.StringValue(d.Name),
		IPVersion:   types.StringValue(d.IPVersion.String()),
		Number:      types.Int64Value(tools.StringToInt64(d.SequenceNumber)),
		Action:      types.StringValue(d.Action.String()),
		Network:     types.StringValue(d.Network),
	}, nil
}
