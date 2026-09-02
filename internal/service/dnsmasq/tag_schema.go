package dnsmasq

import (
	"regexp"

	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type tagResourceModel struct {
	Tag types.String `tfsdk:"tag"`

	Id types.String `tfsdk:"id"`
}

func tagResourceSchema() schema.Schema {
	return schema.Schema{
		Version:             1,
		MarkdownDescription: "Configure tags for dnsmasq.",
		Attributes: map[string]schema.Attribute{

			"tag": schema.StringAttribute{
				MarkdownDescription: "An alphanumeric label which marks a network so that DHCP options may be specified on a per-network basis.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9]+$`),
						"must contain only alphanumeric characters",
					),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID of the tag.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func tagDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Configure tag for dnsmasq.",
		Attributes: map[string]dschema.Attribute{

			"tag": schema.StringAttribute{
				MarkdownDescription: "An alphanumeric label which marks a network so that DHCP options may be specified on a per-network basis.",
				Computed:            true,
			},
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the tag.",
				Required:            true,
			},
		},
	}
}

func convertTagSchemaToStruct(d *tagResourceModel) (*dnsmasq.Tag, error) {
	return &dnsmasq.Tag{
		Tag: d.Tag.ValueString(),
	}, nil
}

func convertTagStructToSchema(d *dnsmasq.Tag) (*tagResourceModel, error) {
	return &tagResourceModel{
		Tag: types.StringValue(d.Tag),
	}, nil
}
