package openvpn

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ ephemeral.EphemeralResource = &generateKeyEphemeral{}
var _ ephemeral.EphemeralResourceWithConfigure = &generateKeyEphemeral{}

func newGenerateKeyEphemeral() ephemeral.EphemeralResource {
	return &generateKeyEphemeral{}
}

type generateKeyEphemeral struct {
	client opnsense.Client
}

type generateKeyEphemeralModel struct {
	KeyType types.String `tfsdk:"key_type"`
	Key     types.String `tfsdk:"key"`
}

func (r *generateKeyEphemeral) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_openvpn_generate_key"
}

func (r *generateKeyEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a fresh OpenVPN static key by calling `/openvpn/instances/gen_key` on each apply. The generated `key` is exposed only in-memory during the run — it is never written to Terraform state. Ephemeral values can be consumed by provider configurations, by other ephemeral blocks, and by `ephemeral = true` outputs; they cannot be passed to a regular resource argument (use a write-only attribute or a secrets-manager provider instead).",

		Attributes: map[string]schema.Attribute{
			"key_type": schema.StringAttribute{
				MarkdownDescription: "Which key flavour to generate. One of `secret` (default, plain shared secret), `tls-auth`, `tls-crypt`, `tls-crypt-v2-server`, or `tls-crypt-v2-client`. Defaults to `secret` when unset.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("secret", "tls-auth", "tls-crypt", "tls-crypt-v2-server", "tls-crypt-v2-client"),
				},
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "The freshly generated key material in OpenVPN static-key format.",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *generateKeyEphemeral) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	apiClient, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected EphemeralResource Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = opnsense.NewClient(apiClient)
}

func (r *generateKeyEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data generateKeyEphemeralModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var keyTypePtr *string
	if !data.KeyType.IsNull() && !data.KeyType.IsUnknown() {
		v := data.KeyType.ValueString()
		keyTypePtr = &v
	}

	res, err := r.client.Openvpn().ServiceGenKey(ctx, keyTypePtr)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to generate openvpn key, got error: %s", err))
		return
	}
	if res.Result != "ok" {
		kt := "secret"
		if keyTypePtr != nil {
			kt = *keyTypePtr
		}
		resp.Diagnostics.AddError(
			"OpenVPN gen_key returned a non-ok result",
			fmt.Sprintf("gen_key returned %q for type %q. Check that the OPNsense version supports this key type.", res.Result, kt),
		)
		return
	}

	data.Key = types.StringValue(res.Key)
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
