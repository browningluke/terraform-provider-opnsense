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

// IpsecChildResourceModel describes the resource data model.
type IpsecChildResourceModel struct {
	Enabled         types.String `tfsdk:"enabled"`
	IPsecConnection types.String `tfsdk:"ipsec_connection"`
	Proposals       types.Set    `tfsdk:"proposals"`
	SHA256_96       types.String `tfsdk:"sha256_96"`
	StartAction     types.String `tfsdk:"start_action"`
	CloseAction     types.String `tfsdk:"close_action"`
	DPDAction       types.String `tfsdk:"dpd_action"`
	Mode            types.String `tfsdk:"mode"`
	InstallPolicies types.String `tfsdk:"install_policies"`
	LocalNetworks   types.Set    `tfsdk:"local_networks"`
	RemoteNetworks  types.Set    `tfsdk:"remote_networks"`
	RequestID       types.String `tfsdk:"request_id"`
	RekeyTime       types.String `tfsdk:"rekey_time"`
	Description     types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func IpsecChildResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "IPsec Child Resources are used for phase 2 of IPsec VPN connections.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.StringAttribute{
				MarkdownDescription: "Enable or disable the Child Resource.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("1"),
			},
			"ipsec_connection": schema.StringAttribute{
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

func IpsecChildDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "IPsec Child Resources are used for phase 2 of IPsec VPN connections.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.StringAttribute{
				MarkdownDescription: "Enable or disable the Child Resource.",
				Computed:            true,
			},
			"ipsec_connection": dschema.StringAttribute{
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

func convertIpsecChildSchemaToStruct(d *IpsecChildResourceModel) (*ipsec.IPsecChild, error) {
	// Convert lists to string slices
	var proposalsList []string
	d.Proposals.ElementsAs(context.Background(), &proposalsList, false)
	sort.Strings(proposalsList)

	var localNetworksList []string
	d.LocalNetworks.ElementsAs(context.Background(), &localNetworksList, false)
	sort.Strings(localNetworksList)

	var remoteNetworksList []string
	d.RemoteNetworks.ElementsAs(context.Background(), &remoteNetworksList, false)
	sort.Strings(remoteNetworksList)

	return &ipsec.IPsecChild{
		Enabled:         d.Enabled.ValueString(),
		Connection:      api.SelectedMap(d.IPsecConnection.ValueString()),
		Proposals:       api.SelectedMapList(proposalsList),
		SHA256_96:       d.SHA256_96.ValueString(),
		StartAction:     api.SelectedMap(d.StartAction.ValueString()),
		CloseAction:     api.SelectedMap(d.CloseAction.ValueString()),
		DPDAction:       api.SelectedMap(d.DPDAction.ValueString()),
		Mode:            api.SelectedMap(d.Mode.ValueString()),
		InstallPolicies: d.InstallPolicies.ValueString(),
		LocalNetworks:   api.SelectedMapList(localNetworksList),
		RemoteNetworks:  api.SelectedMapList(remoteNetworksList),
		RequestID:       d.RequestID.ValueString(),
		RekeyTime:       d.RekeyTime.ValueString(),
		Description:     d.Description.ValueString(),
	}, nil
}

func convertIpsecChildStructToSchema(d *ipsec.IPsecChild) (*IpsecChildResourceModel, error) {
	// Convert List fields
	sort.Strings(d.Proposals)
	proposals, diag := types.SetValueFrom(context.TODO(), types.StringType, d.Proposals)
	if diag.HasError() {
		return nil, fmt.Errorf("error converting proposals: %v", diag)
	}
	sort.Strings(d.LocalNetworks)
	localNetworks, diag := types.SetValueFrom(context.TODO(), types.StringType, d.LocalNetworks)
	if diag.HasError() {
		return nil, fmt.Errorf("error converting local networks: %v", diag)
	}
	sort.Strings(d.RemoteNetworks)
	remoteNetworks, diag := types.SetValueFrom(context.TODO(), types.StringType, d.RemoteNetworks)
	if diag.HasError() {
		return nil, fmt.Errorf("error converting remote networks: %v", diag)
	}

	return &IpsecChildResourceModel{
		Enabled:         types.StringValue(d.Enabled),
		IPsecConnection: types.StringValue(d.Connection.String()),
		Proposals:       proposals,
		SHA256_96:       types.StringValue(d.SHA256_96),
		StartAction:     types.StringValue(d.StartAction.String()),
		CloseAction:     types.StringValue(d.CloseAction.String()),
		DPDAction:       types.StringValue(d.DPDAction.String()),
		Mode:            types.StringValue(d.Mode.String()),
		InstallPolicies: types.StringValue(d.InstallPolicies),
		LocalNetworks:   localNetworks,
		RemoteNetworks:  remoteNetworks,
		RequestID:       types.StringValue(d.RequestID),
		RekeyTime:       types.StringValue(d.RekeyTime),
		Description:     types.StringValue(d.Description),
		Id:              types.StringValue(""), // ID will be set after creation
	}, nil
}
