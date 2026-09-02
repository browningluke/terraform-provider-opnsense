package dnsmasq

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type optionResourceModel struct {
	Type         types.String `tfsdk:"type"`
	OptionV4     types.Int64  `tfsdk:"option"`
	OptionV6     types.Int64  `tfsdk:"option6"`
	Interface    types.String `tfsdk:"interface"`
	TypeSetTags  types.Set    `tfsdk:"tag"`
	TypeMatchTag types.String `tfsdk:"set_tag"`
	Value        types.String `tfsdk:"value"`
	Force        types.Bool   `tfsdk:"force"`
	Description  types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func optionResourceSchema() schema.Schema {
	return schema.Schema{
		Version:             1,
		MarkdownDescription: "Configure DHCP options for dnsmasq.",
		Attributes: map[string]schema.Attribute{
			"interface": schema.StringAttribute{
				MarkdownDescription: "This adds a single interface as tag so this DHCP option can match the interface of a DHCP range.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "\"Set\" option to send it to a client in a DHCP offer or \"Match\" option to dynamically tag clients that send it in the initial DHCP request.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("set", "match"),
				},
			},
			"option": schema.Int64Attribute{
				MarkdownDescription: "DHCPv4 option to offer to the client.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"option6": schema.Int64Attribute{
				MarkdownDescription: "DHCPv6 option to offer to the client.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			// used for type = set
			"tag": schema.SetAttribute{
				MarkdownDescription: "If the optional tags are given then this option is only sent when all the tags are matched. Can be optionally combined with an interface tag.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			// used for type = match (the name is wrong, this is a bug in OPNsense)
			"set_tag": schema.StringAttribute{
				MarkdownDescription: "Tag to set for requests matching this range which can be used to selectively match DHCP options",
				Optional:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Value (or values) to send to the client. The special address 0.0.0.0 or [::] is taken to mean \"the address of the machine running Dnsmasq\". When using \"Match\", leave empty to match on the option only.",
				Optional:            true,
			},
			"force": schema.BoolAttribute{
				MarkdownDescription: "Always send the option, also when the client does not ask for it in the parameter request list.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "You may enter a description here for your reference (not parsed).",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID of the option.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func optionDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure DHCP options for dnsmasq.",
		Attributes: map[string]dschema.Attribute{
			"interface": dschema.StringAttribute{
				MarkdownDescription: "This adds a single interface as tag so this DHCP option can match the interface of a DHCP range.",
				Computed:            true,
			},
			"type": dschema.StringAttribute{
				MarkdownDescription: "\"Set\" option to send it to a client in a DHCP offer or \"Match\" option to dynamically tag clients that send it in the initial DHCP request.",
				Computed:            true,
			},
			"option": dschema.Int64Attribute{
				MarkdownDescription: "DHCPv4 option to offer to the client.",
				Computed:            true,
			},
			"option6": dschema.Int64Attribute{
				MarkdownDescription: "DHCPv6 option to offer to the client.",
				Computed:            true,
			},
			"tag": dschema.SetAttribute{
				MarkdownDescription: "If the optional tags are given then this option is only sent when all the tags are matched. Can be optionally combined with an interface tag.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"set_tag": dschema.StringAttribute{
				MarkdownDescription: "Tag to set for requests matching this range which can be used to selectively match DHCP options",
				Computed:            true,
			},
			"value": dschema.StringAttribute{
				MarkdownDescription: "Value (or values) to send to the client. The special address 0.0.0.0 or [::] is taken to mean \"the address of the machine running Dnsmasq\". When using \"Match\", leave empty to match on the option only.",
				Computed:            true,
			},
			"force": dschema.BoolAttribute{
				MarkdownDescription: "Always send the option, also when the client does not ask for it in the parameter request list.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "You may enter a description here for your reference (not parsed).",
				Computed:            true,
			},
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the option.",
				Required:            true,
			},
		},
	}
}

func convertOptionSchemaToStruct(d *optionResourceModel) (*dnsmasq.Option, error) {
	resultOption := &dnsmasq.Option{}

	// Value
	if !d.Value.IsNull() && !d.Value.IsUnknown() {
		resultOption.Value = d.Value.ValueString()
	}

	// Force
	if !d.Force.IsNull() && !d.Force.IsUnknown() {
		resultOption.Force = tools.BoolToString(d.Force.ValueBool())
	}

	// Type
	if !d.Type.IsNull() && !d.Type.IsUnknown() {
		resultOption.Type = api.SelectedMap(d.Type.ValueString())
	}

	// Interface
	if !d.Interface.IsNull() && !d.Interface.IsUnknown() {
		resultOption.Interface = api.SelectedMap(d.Interface.ValueString())
	}

	// TypeMatchTag
	if !d.TypeMatchTag.IsNull() && !d.TypeMatchTag.IsUnknown() {
		resultOption.TypeMatchTag = api.SelectedMap(d.TypeMatchTag.ValueString())
	}

	// Description
	if !d.Description.IsNull() && !d.Description.IsUnknown() {
		resultOption.Description = d.Description.ValueString()
	}

	// OptionV4
	if !d.OptionV4.IsNull() && !d.OptionV4.IsUnknown() {
		resultOption.OptionV4 =
			api.SelectedMap(tools.Int64ToString(d.OptionV4.ValueInt64()))
	}

	// OptionV6
	if !d.OptionV6.IsNull() && !d.OptionV6.IsUnknown() {
		resultOption.OptionV6 =
			api.SelectedMap(tools.Int64ToString(d.OptionV6.ValueInt64()))
	}

	// Tags
	if !d.TypeSetTags.IsNull() && !d.TypeSetTags.IsUnknown() {
		var tags []string

		diags := d.TypeSetTags.ElementsAs(
			context.Background(),
			&tags,
			false,
		)

		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse tags")
		}

		resultOption.TypeSetTags = api.SelectedMapList(tags)
	}

	return resultOption, nil
}

func convertOptionStructToSchema(d *dnsmasq.Option) (*optionResourceModel, error) {

	tagSet := types.SetNull(types.StringType)
	if len(d.TypeSetTags) > 0 {
		convertedTagSet, diags := types.SetValueFrom(
			context.Background(),
			types.StringType,
			d.TypeSetTags,
		)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert tags")
		}
		tagSet = convertedTagSet
	}

	option := &optionResourceModel{}

	// Type
	if d.Type.String() == "" {
		option.Type = types.StringNull()
	} else {
		option.Type = types.StringValue(d.Type.String())
	}

	// Interface
	if d.Interface.String() == "" {
		option.Interface = types.StringNull()
	} else {
		option.Interface = types.StringValue(d.Interface.String())
	}

	// TypeMatchTag
	if d.TypeMatchTag.String() == "" {
		option.TypeMatchTag = types.StringNull()
	} else {
		option.TypeMatchTag = types.StringValue(d.TypeMatchTag.String())
	}

	// Description
	if d.Description == "" {
		option.Description = types.StringValue("")
	} else {
		option.Description = types.StringValue(d.Description)
	}

	// Value
	if d.Value == "" {
		option.Value = types.StringNull()
	} else {
		option.Value = types.StringValue(d.Value)
	}

	// Force
	if d.Force == "" {
		option.Force = types.BoolNull()
	} else {
		option.Force = types.BoolValue(tools.StringToBool(d.Force))
	}

	// OptionV4
	if d.OptionV4.String() == "" || d.OptionV4.String() == "0" {
		option.OptionV4 = types.Int64Null()
	} else {
		option.OptionV4 =
			types.Int64Value(tools.StringToInt64(d.OptionV4.String()))
	}

	// OptionV6
	if d.OptionV6.String() == "" || d.OptionV6.String() == "0" {
		option.OptionV6 = types.Int64Null()
	} else {
		option.OptionV6 =
			types.Int64Value(tools.StringToInt64(d.OptionV6.String()))
	}

	// Tags
	option.TypeSetTags = tagSet

	return option, nil
}
