package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/ipsec"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IpsecAuthLocalResourceModel describes the resource data model.
type IpsecAuthLocalResourceModel struct {
	Enabled        types.String `tfsdk:"enabled"`
	Connection     types.String `tfsdk:"connection"`
	Round          types.String `tfsdk:"round"`
	Authentication types.String `tfsdk:"authentication"`
	AuthId         types.String `tfsdk:"auth_id"`
	EAPId          types.String `tfsdk:"eap_id"`
	Certificates   types.Set    `tfsdk:"certificates"`
	PublicKeys     types.Set    `tfsdk:"public_keys"`
	Description    types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func IpsecAuthLocalResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "IPsec AuthLocal Resources are used for phase 1 authentication of IPsec VPN connections.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.StringAttribute{
				MarkdownDescription: "Enable or disable the AuthLocal Resource.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("1"),
			},
			"connection": schema.StringAttribute{
				MarkdownDescription: "The parent connection UUID.",
				Required:            true,
			},
			"round": schema.StringAttribute{
				MarkdownDescription: "Authentication round for the AuthLocal Resource.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("1"),
			},
			"authentication": schema.StringAttribute{
				MarkdownDescription: "Authentication method for the AuthLocal Resource.",
				Required:            true,
			},
			"auth_id": schema.StringAttribute{
				MarkdownDescription: "Authentication ID for the AuthLocal Resource.",
				Optional:            true,
			},
			"eap_id": schema.StringAttribute{
				MarkdownDescription: "EAP ID for the AuthLocal Resource.",
				Optional:            true,
			},
			"certificates": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of certificates for the AuthLocal Resource.",
				Optional:            true,
			},
			"public_keys": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of public keys for the AuthLocal Resource.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for the AuthLocal Resource.",
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

func IpsecAuthLocalDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "IPsec AuthLocal Resources are used for phase 1 authentication of IPsec VPN connections.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.StringAttribute{
				MarkdownDescription: "Enable or disable the AuthLocal Resource.",
				Computed:            true,
			},
			"connection": dschema.StringAttribute{
				MarkdownDescription: "Connection ID for the AuthLocal Resource.",
				Computed:            true,
			},
			"round": dschema.StringAttribute{
				MarkdownDescription: "Authentication round for the AuthLocal Resource.",
				Computed:            true,
			},
			"authentication": dschema.StringAttribute{
				MarkdownDescription: "Authentication method for the AuthLocal Resource.",
				Computed:            true,
			},
			"auth_id": dschema.StringAttribute{
				MarkdownDescription: "Authentication ID for the AuthLocal Resource.",
				Computed:            true,
			},
			"eap_id": dschema.StringAttribute{
				MarkdownDescription: "EAP ID for the AuthLocal Resource.",
				Computed:            true,
			},
			"certificates": dschema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of certificates for the AuthLocal Resource.",
				Computed:            true,
			},
			"public_keys": dschema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of public keys for the AuthLocal Resource.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description for the AuthLocal Resource.",
				Computed:            true,
			},
		},
	}
}

func convertIpsecAuthLocalSchemaToStruct(d *IpsecAuthLocalResourceModel) (*ipsec.IPsecAuthLocal, error) {
	var certificatesList []string
	d.Certificates.ElementsAs(context.Background(), &certificatesList, false)
	sort.Strings(certificatesList)

	var publicKeysList []string
	d.PublicKeys.ElementsAs(context.Background(), &publicKeysList, false)
	sort.Strings(publicKeysList)

	return &ipsec.IPsecAuthLocal{
		Enabled:        d.Enabled.ValueString(),
		Connection:     api.SelectedMap(d.Connection.ValueString()),
		Round:          d.Round.ValueString(),
		Authentication: api.SelectedMap(d.Authentication.ValueString()),
		Id:             d.AuthId.ValueString(),
		EAPId:          d.EAPId.ValueString(),
		Certificates:   api.SelectedMapList(certificatesList),
		PublicKeys:     api.SelectedMapList(publicKeysList),
		Description:    d.Description.ValueString(),
	}, nil
}

func convertIpsecAuthLocalStructToSchema(d *ipsec.IPsecAuthLocal) (*IpsecAuthLocalResourceModel, error) {
	// Convert Set fields
	certificates, diag := types.SetValueFrom(context.TODO(), types.StringType, d.Certificates)
	if diag.HasError() {
		return nil, fmt.Errorf("error converting certificates: %v", diag)
	}
	publicKeys, diag := types.SetValueFrom(context.TODO(), types.StringType, d.PublicKeys)
	if diag.HasError() {
		return nil, fmt.Errorf("error converting public keys: %v", diag)
	}

	return &IpsecAuthLocalResourceModel{
		Enabled:        types.StringValue(d.Enabled),
		Connection:     types.StringValue(d.Connection.String()),
		Round:          types.StringValue(d.Round),
		Authentication: types.StringValue(d.Authentication.String()),
		AuthId:         types.StringValue(d.Id),
		EAPId:          types.StringValue(d.EAPId),
		Certificates:   certificates,
		PublicKeys:     publicKeys,
		Description:    types.StringValue(d.Description),
		Id:             types.StringValue(""), // ID will be set after creation
	}, nil
}
