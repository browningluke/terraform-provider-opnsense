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

/*
type IPsecAuthRemote struct {
	Enabled        string              `json:"enabled"`
	Connection     api.SelectedMap     `json:"connection"`
	Round          string              `json:"round"`
	Authentication api.SelectedMap     `json:"auth"`
	Id             string              `json:"id"`
	EAPId          string              `json:"eap_id"`
	Certificates   api.SelectedMapList `json:"certs"`
	PublicKeys     api.SelectedMapList `json:"public_keys"`
	Description    string              `json:"description"`
}
*/

// IpsecAuthRemoteResourceModel describes the resource data model.
type IpsecAuthRemoteResourceModel struct {
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

func IpsecAuthRemoteResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "IPsec AuthRemote Resources are used for phase 1 authentication of IPsec VPN connections.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.StringAttribute{
				MarkdownDescription: "Enable or disable the Child Resource.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("1"),
			},
			"connection": schema.StringAttribute{
				MarkdownDescription: "The parent connection UUID.",
				Required:            true,
			},
			"proposals": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of proposals for the Child Resource.",
				Required:            true,
			},
			"sha256_96": schema.StringAttribute{
				MarkdownDescription: "Enable or disable SHA256_96.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("0"),
			},
			"start_action": schema.StringAttribute{
				MarkdownDescription: "Start action for the Child Resource.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("start"),
			},
			"close_action": schema.StringAttribute{
				MarkdownDescription: "Close action for the Child Resource.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("none"),
			},
			"dpd_action": schema.StringAttribute{
				MarkdownDescription: "DPD action for the Child Resource.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("hold"),
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "Mode for the Child Resource.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("tunnel"),
			},
			"install_policies": schema.StringAttribute{
				MarkdownDescription: "Install policies for the Child Resource.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("1"),
			},
			"local_networks": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of local networks for the Child Resource.",
				Required:            true,
			},
			"remote_networks": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of remote networks for the Child Resource.",
				Required:            true,
			},
			"request_id": schema.StringAttribute{
				MarkdownDescription: "Request ID for the Child Resource.",
				Optional:            true,
			},
			"rekey_time": schema.StringAttribute{
				MarkdownDescription: "Rekey time for the Child Resource in seconds.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("0"),
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

func IpsecAuthRemoteDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "IPsec AuthRemote Resources are used for phase 1 authentication of IPsec VPN connections.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.StringAttribute{
				MarkdownDescription: "Enable or disable the Child Resource.",
				Computed:            true,
			},
			"connection": dschema.StringAttribute{
				MarkdownDescription: "Connection ID for the Child Resource.",
				Computed:            true,
			},
			"proposals": dschema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of proposals for the Child Resource.",
				Computed:            true,
			},
			"sha256_96": dschema.StringAttribute{
				MarkdownDescription: "Enable or disable SHA256_96.",
				Computed:            true,
			},
			"start_action": dschema.StringAttribute{
				MarkdownDescription: "Start action for the Child Resource.",
				Computed:            true,
			},
			"close_action": dschema.StringAttribute{
				MarkdownDescription: "Close action for the Child Resource.",
				Computed:            true,
			},
			"dpd_action": dschema.StringAttribute{
				MarkdownDescription: "DPD action for the Child Resource.",
				Computed:            true,
			},
			"mode": dschema.StringAttribute{
				MarkdownDescription: "Mode for the Child Resource.",
				Computed:            true,
			},
			"install_policies": dschema.StringAttribute{
				MarkdownDescription: "Install policies for the Child Resource.",
				Computed:            true,
			},
			"local_networks": dschema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of local networks for the Child Resource.",
				Computed:            true,
			},
			"remote_networks": dschema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of remote networks for the Child Resource.",
				Computed:            true,
			},
			"request_id": dschema.StringAttribute{
				MarkdownDescription: "Request ID for the Child Resource.",
				Computed:            true,
			},
			"rekey_time": dschema.StringAttribute{
				MarkdownDescription: "Rekey time for the Child Resource in seconds.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description for the PSK.",
				Computed:            true,
			},
		},
	}
}

func convertIpsecAuthRemoteSchemaToStruct(d *IpsecAuthRemoteResourceModel) (*ipsec.IPsecAuthRemote, error) {
	var certificatesList []string
	d.Certificates.ElementsAs(context.Background(), &certificatesList, false)
	sort.Strings(certificatesList)

	var publicKeysList []string
	d.PublicKeys.ElementsAs(context.Background(), &publicKeysList, false)
	sort.Strings(publicKeysList)

	return &ipsec.IPsecAuthRemote{
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

func convertIpsecAuthRemoteStructToSchema(d *ipsec.IPsecAuthRemote) (*IpsecAuthRemoteResourceModel, error) {
	// Convert Set fields
	certificates, diag := types.SetValueFrom(context.TODO(), types.StringType, d.Certificates)
	if diag.HasError() {
		return nil, fmt.Errorf("error converting certificates: %v", diag)
	}
	publicKeys, diag := types.SetValueFrom(context.TODO(), types.StringType, d.PublicKeys)
	if diag.HasError() {
		return nil, fmt.Errorf("error converting public keys: %v", diag)
	}

	return &IpsecAuthRemoteResourceModel{
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
