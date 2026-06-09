package openvpn

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/openvpn"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// staticKeyResourceModel describes the resource data model.
type staticKeyResourceModel struct {
	Mode        types.String `tfsdk:"mode"`
	Key         types.String `tfsdk:"key"`
	Description types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func staticKeyResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "OpenVPN static keys are pre-shared TLS or static keys used by OpenVPN instances for `tls-auth`, `tls-crypt`, `tls-crypt-v2` or peer-to-peer secret authentication.",

		Attributes: map[string]schema.Attribute{
			"mode": schema.StringAttribute{
				MarkdownDescription: "The static-key mode. One of `auth`, `crypt`, or `crypt-v2`. Defaults to `crypt`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("crypt"),
				Validators: []validator.String{
					stringvalidator.OneOf("auth", "crypt", "crypt-v2"),
				},
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "The static key payload. Use the OpenVPN-formatted key (e.g. output of `openvpn --genkey secret`).",
				Required:            true,
				Sensitive:           true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for this static key.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the static key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func staticKeyDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "OpenVPN static keys are pre-shared TLS or static keys used by OpenVPN instances.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"mode": dschema.StringAttribute{
				MarkdownDescription: "The static-key mode.",
				Computed:            true,
			},
			"key": dschema.StringAttribute{
				MarkdownDescription: "The static key payload.",
				Computed:            true,
				Sensitive:           true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Description for this static key.",
				Computed:            true,
			},
		},
	}
}

func convertStaticKeySchemaToStruct(d *staticKeyResourceModel) (*openvpn.StaticKey, error) {
	return &openvpn.StaticKey{
		Mode:        api.SelectedMap(d.Mode.ValueString()),
		Key:         d.Key.ValueString(),
		Description: d.Description.ValueString(),
	}, nil
}

func convertStaticKeyStructToSchema(d *openvpn.StaticKey) (*staticKeyResourceModel, error) {
	return &staticKeyResourceModel{
		Mode:        types.StringValue(d.Mode.String()),
		Key:         types.StringValue(d.Key),
		Description: types.StringValue(d.Description),
	}, nil
}
