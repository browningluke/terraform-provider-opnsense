package auth

import (
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type userPasswordEphemeralModel struct {
	// Id       types.String `tfsdk:"id"`
	Password types.String `tfsdk:"password"`
}

func userPasswordEphemeralSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "User Password ephemeral schema description",

		Attributes: map[string]schema.Attribute{
			// "id": schema.StringAttribute{
			// 	Computed:            true,
			// 	MarkdownDescription: "Id of the resource",
			// },
			"password": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Password of the user",
				Sensitive:           true,
			},
		},
	}
}

// func convertPasswordSchemaToStruct(scheme *userPasswordEphemeralModel) (*auth.User, error) {
// 	return &auth.User{
// 		// UserId:   scheme.Id.ValueString(),
// 		Password: scheme.Password.ValueString(),
// 	}, nil
// }

// func convertPasswordStructToSchema(strct *auth.User) (*userPasswordEphemeralModel, error) {
// 	return &userPasswordEphemeralModel{
// 		// Id:       types.StringValue(strct.UserId),
// 		Password: types.StringValue(strct.Password),
// 	}, nil
// }
