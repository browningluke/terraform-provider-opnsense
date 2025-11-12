package auth

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ ephemeral.EphemeralResource = &userPasswordEphemeral{}
var _ ephemeral.EphemeralResourceWithConfigure = &userPasswordEphemeral{}

type userPasswordEphemeral struct {
	client opnsense.Client
}

func newUserPasswordEphemeral() ephemeral.EphemeralResource {
	return &userPasswordEphemeral{}
}

// Metadata implements ephemeral.EphemeralResource.
func (u *userPasswordEphemeral) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_user_password"
}

// Schema implements ephemeral.EphemeralResource.
func (u *userPasswordEphemeral) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = userPasswordEphemeralSchema()
}

// Configure implements ephemeral.EphemeralResourceWithConfigure.
func (u *userPasswordEphemeral) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *opnsense.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	u.client = opnsense.NewClient(apiClient)
}

// Open implements ephemeral.EphemeralResource.
func (u *userPasswordEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data userPasswordEphemeralModel

	// Read Terraform config data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("\n%+v\n", data))

	// Typically ephemeral resources will make external calls, however this example
	// hardcodes setting the token attribute to a specific value for brevity.
	// data.Password = data.Password

	// Save data into ephemeral result data
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
