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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ospfAreaResourceModel describes the resource data model.
type ospfAreaResourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled"`
	AreaID  types.String `tfsdk:"area_id"`
	Type    types.String `tfsdk:"type"`

	Id types.String `tfsdk:"id"`
}

func ospfAreaResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Configure OSPF areas.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this OSPF area. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"area_id": schema.StringAttribute{
				MarkdownDescription: "The area ID in IPv4 dotted notation (e.g. `0.0.0.1`).",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The area type. One of `stub`, `stub no-summary`, `nssa`, `nssa no-summary`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("stub", "stub no-summary", "nssa", "nssa no-summary"),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the OSPF area.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func ospfAreaDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure OSPF areas.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this OSPF area.",
				Computed:            true,
			},
			"area_id": dschema.StringAttribute{
				MarkdownDescription: "The area ID in IPv4 dotted notation.",
				Computed:            true,
			},
			"type": dschema.StringAttribute{
				MarkdownDescription: "The area type.",
				Computed:            true,
			},
		},
	}
}

func convertOSPFAreaSchemaToStruct(d *ospfAreaResourceModel) (*quagga.OSPFArea, error) {
	return &quagga.OSPFArea{
		Enabled: tools.BoolToString(d.Enabled.ValueBool()),
		AreaID:  d.AreaID.ValueString(),
		Type:    api.SelectedMap(d.Type.ValueString()),
	}, nil
}

func convertOSPFAreaStructToSchema(d *quagga.OSPFArea) (*ospfAreaResourceModel, error) {
	return &ospfAreaResourceModel{
		Enabled: types.BoolValue(tools.StringToBool(d.Enabled)),
		AreaID:  types.StringValue(d.AreaID),
		Type:    types.StringValue(d.Type.String()),
	}, nil
}
