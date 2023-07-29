package service

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/interfaces"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/tools"
)

// InterfacesVlanResourceModel describes the resource data model.
type InterfacesVlanResourceModel struct {
	Description types.String `tfsdk:"description"`
	Tag         types.Int64  `tfsdk:"tag"`
	Priority    types.Int64  `tfsdk:"priority"`
	Parent      types.String `tfsdk:"parent"`
	Device      types.String `tfsdk:"device"`

	Id types.String `tfsdk:"id"`
}

func InterfacesVlanResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "VLANs (Virtual LANs) can be used to segment a single physical network into multiple virtual networks.",

		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Optional:            true,
			},
			"tag": schema.Int64Attribute{
				MarkdownDescription: "802.1Q VLAN tag.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 4094),
				},
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "802.1Q VLAN PCP (priority code point). Defaults to `0`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 7),
				},
			},
			"parent": schema.StringAttribute{
				MarkdownDescription: "VLAN capable interface to attach the VLAN to, e.g. `vtnet0`.",
				Required:            true,
			},
			"device": schema.StringAttribute{
				MarkdownDescription: "Custom VLAN name. Custom names are possible, but only if the start of the name matches the required prefix and contains numeric characters or dots, e.g. `vlan0.1.2` or `qinq0.3.4`. Set to `\"\"` to generate a device name. Defaults to `\"\"`",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the VLAN.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func InterfacesVlanDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "VLANs (Virtual LANs) can be used to segment a single physical network into multiple virtual networks.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Computed:            true,
			},
			"tag": dschema.Int64Attribute{
				MarkdownDescription: "802.1Q VLAN tag.",
				Computed:            true,
			},
			"priority": dschema.Int64Attribute{
				MarkdownDescription: "802.1Q VLAN PCP (priority code point).",
				Computed:            true,
			},
			"parent": dschema.StringAttribute{
				MarkdownDescription: "VLAN capable interface to attach the VLAN to, e.g. `vtnet0`.",
				Computed:            true,
			},
			"device": dschema.StringAttribute{
				MarkdownDescription: "Custom VLAN name. Custom names are possible, but only if the start of the name matches the required prefix and contains numeric characters or dots, e.g. `vlan0.1.2` or `qinq0.3.4`.",
				Computed:            true,
			},
		},
	}
}

func convertInterfacesVlanSchemaToStruct(d *InterfacesVlanResourceModel) (*interfaces.Vlan, error) {
	return &interfaces.Vlan{
		Description: d.Description.ValueString(),
		Tag:         tools.Int64ToString(d.Tag.ValueInt64()),
		Priority:    api.SelectedMap(tools.Int64ToString(d.Priority.ValueInt64())),
		Parent:      api.SelectedMap(d.Parent.ValueString()),
		Device:      d.Device.ValueString(),
	}, nil
}

func convertInterfacesVlanStructToSchema(d *interfaces.Vlan) (*InterfacesVlanResourceModel, error) {
	return &InterfacesVlanResourceModel{
		Description: tools.StringOrNull(d.Description),
		Tag:         tools.StringToInt64Null(d.Tag),
		Priority:    tools.StringToInt64Null(d.Priority.String()),
		Parent:      types.StringValue(d.Parent.String()),
		Device:      types.StringValue(d.Device),
	}, nil
}
