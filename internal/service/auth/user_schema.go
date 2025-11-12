package auth

import (
	"github.com/browningluke/opnsense-go/pkg/auth"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type userResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Disabled types.Bool   `tfsdk:"disabled"`
	Name     types.String `tfsdk:"name"`
	Password types.String `tfsdk:"password"`
}

func userResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "User schema description",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id of the resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disabled": schema.BoolAttribute{
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If user is disabled",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the user",
			},
			"password": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Password of the user",
				Sensitive:           true,
			},
		},
	}
}

func userDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "User data schema description",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Id of the resource",
			},
		},
	}
}

func convertUserSchemaToStruct(scheme *userResourceModel) (*auth.User, error) {
	return &auth.User{
		UserId:   scheme.Id.ValueString(),
		Disabled: tools.BoolToString(scheme.Disabled.ValueBool()),
		Name:     scheme.Name.ValueString(),
		Password: scheme.Password.ValueString(),
	}, nil
}

func convertUserStructToSchema(strct *auth.User) (*userResourceModel, error) {
	return &userResourceModel{
		Id:       types.StringValue(strct.UserId),
		Disabled: types.BoolValue(tools.StringToBool(strct.Disabled)),
		Name:     types.StringValue(strct.Name),
		Password: types.StringValue(strct.Password),
	}, nil
}
