package service

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/ipsec"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IpsecPskResourceModel describes the resource data model.
type IpsecPskResourceModel struct {
	IdentityLocal  types.String `tfsdk:"identity_local"`
	IdentityRemote types.String `tfsdk:"identity_remote"`
	PreSharedKey   types.String `tfsdk:"pre_shared_key"`
	Type           types.String `tfsdk:"type"`
	Description    types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func IpsecPskResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "IPsec Pre-Shared Keys (PSKs) are used for authenticating IPsec VPN connections.",

		Attributes: map[string]schema.Attribute{
			"identity_local": schema.StringAttribute{
				MarkdownDescription: "Local identity for the PSK.",
				Required:            true,
			},
			"identity_remote": schema.StringAttribute{
				MarkdownDescription: "Remote identity for the PSK.",
				Required:            true,
			},
			"pre_shared_key": schema.StringAttribute{
				MarkdownDescription: "The pre-shared key used for authentication.",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the pre-shared key. Valid values are 'PSK' (traditional pre-shared key) or 'EAP' (for EAP-MSCHAPv2 authentication).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("PSK"),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for the PSK.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Description: "UUID of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func IpsecPskDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "IPsec Pre-Shared Keys (PSKs) are used for authenticating IPsec VPN connections.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"identity_local": dschema.StringAttribute{
				MarkdownDescription: "Local identity for the PSK.",
				Computed:            true,
			},
			"identity_remote": dschema.StringAttribute{
				MarkdownDescription: "Remote identity for the PSK.",
				Computed:            true,
			},
			"pre_shared_key": dschema.StringAttribute{
				MarkdownDescription: "The pre-shared key used for authentication.",
				Computed:            true,
			},
			"type": dschema.StringAttribute{
				MarkdownDescription: "Type of the pre-shared key, e.g., 'psk'.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description for the PSK.",
				Computed:            true,
			},
		},
	}
}

func convertIpsecPskSchemaToStruct(d *IpsecPskResourceModel) (*ipsec.IPsecPSK, error) {
	return &ipsec.IPsecPSK{
		IdentityLocal:  d.IdentityLocal.ValueString(),
		IdentityRemote: d.IdentityRemote.ValueString(),
		PreSharedKey:   d.PreSharedKey.ValueString(),
		Type:           api.SelectedMap(d.Type.ValueString()),
		Description:    d.Description.ValueString(),
	}, nil
}

func convertIpsecPskStructToSchema(d *ipsec.IPsecPSK) (*IpsecPskResourceModel, error) {
	return &IpsecPskResourceModel{
		IdentityLocal:  types.StringValue(d.IdentityLocal),
		IdentityRemote: types.StringValue(d.IdentityRemote),
		PreSharedKey:   types.StringValue(d.PreSharedKey),
		Type:           types.StringValue(d.Type.String()),
		Description:    types.StringValue(d.Description),
	}, nil
}
