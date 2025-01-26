package service

import (
	"context"
	"github.com/browningluke/opnsense-go/pkg/bind"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-opnsense/internal/tools"
)

// BindAclResourceModel describes the resource data model.
type BindAclResourceModel struct {
	Enabled  types.Bool   `tfsdk:"enabled"`
	Name     types.String `tfsdk:"name"`
	Networks types.Set    `tfsdk:"networks"`

	Id types.String `tfsdk:"id"`
}

func BindAclResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "ACLs are named lists of networks that can be used to configure who can access which resources in Bind.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this ACL. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Arbitrary name that is used to reference this ACL. Should be a string between 1 and 32 characters. Allowed characters are 0-9, a-z, A-Z, _ and -.",
				Required:            true,
				Validators: []validator.String{
					// TODO check maximum string length again
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"networks": schema.SetAttribute{
				MarkdownDescription: "The list of subnets to include in this ACL.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func BindAclDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "ACLs are named lists of networks that can be used to configure who can access which resources in Bind.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Enable this ACL.",
				Computed:            true,
			},
			"name": dschema.StringAttribute{
				MarkdownDescription: "Arbitrary name that is used to reference this ACL. Should be a string between 1 and 32 characters. Allowed characters are 0-9, a-z, A-Z, _ and -.",
				Computed:            true,
			},
			"networks": dschema.SetAttribute{
				MarkdownDescription: "The list of subnets to include in this ACL.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func convertBindAclSchemaToStruct(d *BindAclResourceModel) (*bind.Acl, error) {
	// Parse 'Networks'
	var networksList string
	d.Networks.ElementsAs(context.Background(), &networksList, false)

	return &bind.Acl{
		Enabled:  tools.BoolToString(d.Enabled.ValueBool()),
		Name:     d.Name.ValueString(),
		Networks: tools.SetToStringSlice(d.Networks),
	}, nil
}

func convertBindAclStructToSchema(d *bind.Acl) (*BindAclResourceModel, error) {
	model := &BindAclResourceModel{
		Enabled:  types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:     types.StringValue(d.Name),
		Networks: tools.StringSliceToSet(d.Networks),
	}

	return model, nil
}
