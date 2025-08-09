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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IpsecConnectionResourceModel describes the resource data model.
type IpsecConnectionResourceModel struct {
	Enabled                types.String `tfsdk:"enabled"`
	Proposals              types.Set    `tfsdk:"proposals"`
	Unique                 types.String `tfsdk:"unique"`
	Aggressive             types.String `tfsdk:"aggressive"`
	Version                types.String `tfsdk:"version"`
	Mobike                 types.String `tfsdk:"mobike"`
	LocalAddresses         types.Set    `tfsdk:"local_addresses"`
	RemoteAddresses        types.Set    `tfsdk:"remote_addresses"`
	LocalPort              types.String `tfsdk:"local_port"`
	RemotePort             types.String `tfsdk:"remote_port"`
	UDPEncapsulation       types.String `tfsdk:"udp_encapsulation"`
	ReauthenticationTime   types.String `tfsdk:"reauthentication_time"`
	RekeyTime              types.String `tfsdk:"rekey_time"`
	IKELifetime            types.String `tfsdk:"ike_lifetime"`
	DPDDelay               types.String `tfsdk:"dpd_delay"`
	DPDTimeout             types.String `tfsdk:"dpd_timeout"`
	IPPools                types.Set    `tfsdk:"ip_pools"`
	SendCertificateRequest types.String `tfsdk:"send_certificate_request"`
	SendCertificate        types.String `tfsdk:"send_certificate"`
	KeyingTries            types.String `tfsdk:"keying_tries"`
	Description            types.String `tfsdk:"description"`

	Id types.String `tfsdk:"id"`
}

func IpsecConnectionResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "IPsec Connections are used for establishing secure communication channels.",

		Attributes: map[string]schema.Attribute{
			"enabled": schema.StringAttribute{
				MarkdownDescription: "Enable or disable the IPsec connection.",
				Required:            true,
			},
			"proposals": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of encryption proposals for the connection.",
				Required:            true,
			},
			"unique": schema.StringAttribute{
				MarkdownDescription: "Whether the connection should use unique IDs.",
				Required:            true,
			},
			"aggressive": schema.StringAttribute{
				MarkdownDescription: "Enable or disable aggressive mode.",
				Required:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "IKE version to use (e.g., '1', '2').",
				Required:            true,
			},
			"mobike": schema.StringAttribute{
				MarkdownDescription: "Enable or disable MOBIKE support.",
				Required:            true,
			},
			"local_addresses": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of local addresses for the connection.",
				Required:            true,
			},
			"remote_addresses": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of remote addresses for the connection.",
				Required:            true,
			},
			"local_port": schema.StringAttribute{
				MarkdownDescription: "Local port for the connection.",
				Required:            true,
			},
			"remote_port": schema.StringAttribute{
				MarkdownDescription: "Remote port for the connection.",
				Required:            true,
			},
			"udp_encapsulation": schema.StringAttribute{
				MarkdownDescription: "Enable or disable UDP encapsulation.",
				Required:            true,
			},
			"reauthentication_time": schema.StringAttribute{
				MarkdownDescription: "Time interval for reauthentication.",
				Required:            true,
			},
			"rekey_time": schema.StringAttribute{
				MarkdownDescription: "Time interval for rekeying.",
				Required:            true,
			},
			"ike_lifetime": schema.StringAttribute{
				MarkdownDescription: "IKE lifetime duration.",
				Required:            true,
			},
			"dpd_delay": schema.StringAttribute{
				MarkdownDescription: "Dead Peer Detection (DPD) delay.",
				Required:            true,
			},
			"dpd_timeout": schema.StringAttribute{
				MarkdownDescription: "Dead Peer Detection (DPD) timeout.",
				Required:            true,
			},
			"ip_pools": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of IP pools for the connection.",
				Optional:            true,
			},
			"send_certificate_request": schema.StringAttribute{
				MarkdownDescription: "Whether to send a certificate request.",
				Required:            true,
			},
			"send_certificate": schema.StringAttribute{
				MarkdownDescription: "Whether to send a certificate.",
				Required:            true,
			},
			"keying_tries": schema.StringAttribute{
				MarkdownDescription: "Number of keying tries.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for the IPsec connection.",
				Required:            true,
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

func IpsecConnectionDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "IPsec connections are used for establishing secure VPN tunnels.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "UUID of the resource.",
				Required:            true,
			},
			"enabled": dschema.StringAttribute{
				MarkdownDescription: "Enable or disable the IPsec connection.",
				Computed:            true,
			},
			"proposals": dschema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of encryption proposals for the connection.",
				Computed:            true,
			},
			"unique": dschema.StringAttribute{
				MarkdownDescription: "Whether the connection should use unique IDs.",
				Computed:            true,
			},
			"aggressive": dschema.StringAttribute{
				MarkdownDescription: "Enable or disable aggressive mode.",
				Computed:            true,
			},
			"version": dschema.StringAttribute{
				MarkdownDescription: "IKE version to use (e.g., '1', '2').",
				Computed:            true,
			},
			"mobike": dschema.StringAttribute{
				MarkdownDescription: "Enable or disable MOBIKE support.",
				Computed:            true,
			},
			"local_addresses": dschema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of local addresses for the connection.",
				Computed:            true,
			},
			"remote_addresses": dschema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of remote addresses for the connection.",
				Computed:            true,
			},
			"local_port": dschema.StringAttribute{
				MarkdownDescription: "Local port for the connection.",
				Computed:            true,
			},
			"remote_port": dschema.StringAttribute{
				MarkdownDescription: "Remote port for the connection.",
				Computed:            true,
			},
			"udp_encapsulation": dschema.StringAttribute{
				MarkdownDescription: "Enable or disable UDP encapsulation.",
				Computed:            true,
			},
			"reauthentication_time": dschema.StringAttribute{
				MarkdownDescription: "Time interval for reauthentication.",
				Computed:            true,
			},
			"rekey_time": dschema.StringAttribute{
				MarkdownDescription: "Time interval for rekeying.",
				Computed:            true,
			},
			"ike_lifetime": dschema.StringAttribute{
				MarkdownDescription: "IKE lifetime duration.",
				Computed:            true,
			},
			"dpd_delay": dschema.StringAttribute{
				MarkdownDescription: "Dead Peer Detection (DPD) delay.",
				Computed:            true,
			},
			"dpd_timeout": dschema.StringAttribute{
				MarkdownDescription: "Dead Peer Detection (DPD) timeout.",
				Computed:            true,
			},
			"ip_pools": dschema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of IP pools for the connection.",
				Computed:            true,
			},
			"send_certificate_request": dschema.StringAttribute{
				MarkdownDescription: "Whether to send a certificate request.",
				Computed:            true,
			},
			"send_certificate": dschema.StringAttribute{
				MarkdownDescription: "Whether to send a certificate.",
				Computed:            true,
			},
			"keying_tries": dschema.StringAttribute{
				MarkdownDescription: "Number of keying tries.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description for the IPsec connection.",
				Computed:            true,
			},
		},
	}
}

func convertIpsecConnectionSchemaToStruct(d *IpsecConnectionResourceModel) (*ipsec.IPsecConnection, error) {
	// Convert lists to string slices
	var proposalsList []string
	d.Proposals.ElementsAs(context.Background(), &proposalsList, false)
	sort.Strings(proposalsList)

	var localAddressesList []string
	d.LocalAddresses.ElementsAs(context.Background(), &localAddressesList, false)
	sort.Strings(localAddressesList)

	var remoteAddressesList []string
	d.RemoteAddresses.ElementsAs(context.Background(), &remoteAddressesList, false)
	sort.Strings(remoteAddressesList)

	var ipPoolsList []string
	d.IPPools.ElementsAs(context.Background(), &ipPoolsList, false)
	sort.Strings(ipPoolsList)

	// Handle optional port fields - only set if not empty
	var localPort, remotePort api.SelectedMap
	if d.LocalPort.ValueString() != "" {
		localPort = api.SelectedMap(d.LocalPort.ValueString())
	}
	if d.RemotePort.ValueString() != "" {
		remotePort = api.SelectedMap(d.RemotePort.ValueString())
	}

	return &ipsec.IPsecConnection{
		Enabled:                d.Enabled.ValueString(),
		Proposals:              api.SelectedMapList(proposalsList),
		Unique:                 api.SelectedMap(d.Unique.ValueString()),
		Aggressive:             d.Aggressive.ValueString(),
		Version:                api.SelectedMap(d.Version.ValueString()),
		Mobike:                 d.Mobike.ValueString(),
		LocalAddresses:         api.SelectedMapList(localAddressesList),
		RemoteAddresses:        api.SelectedMapList(remoteAddressesList),
		LocalPort:              localPort,
		RemotePort:             remotePort,
		UDPEncapsulation:       d.UDPEncapsulation.ValueString(),
		ReauthenticationTime:   d.ReauthenticationTime.ValueString(),
		RekeyTime:              d.RekeyTime.ValueString(),
		IKELifetime:            d.IKELifetime.ValueString(),
		DPDDelay:               d.DPDDelay.ValueString(),
		DPDTimeout:             d.DPDTimeout.ValueString(),
		IPPools:                api.SelectedMapList(ipPoolsList),
		SendCertificateRequest: d.SendCertificateRequest.ValueString(),
		SendCertificate:        api.SelectedMap(d.SendCertificate.ValueString()),
		KeyingTries:            d.KeyingTries.ValueString(),
		Description:            d.Description.ValueString(),
	}, nil
}

func convertIpsecConnectionStructToSchema(d *ipsec.IPsecConnection) (*IpsecConnectionResourceModel, error) {
	// Convert list fields
	sort.Strings(d.Proposals)
	proposals, diag := types.SetValueFrom(context.TODO(), types.StringType, d.Proposals)
	if diag.HasError() {
		return nil, fmt.Errorf("error converting proposals: %v", diag)
	}
	sort.Strings(d.LocalAddresses)
	localAddresses, diag := types.SetValueFrom(context.TODO(), types.StringType, d.LocalAddresses)
	if diag.HasError() {
		return nil, fmt.Errorf("error converting local addresses: %v", diag)
	}
	sort.Strings(d.RemoteAddresses)
	remoteAddresses, diag := types.SetValueFrom(context.TODO(), types.StringType, d.RemoteAddresses)
	if diag.HasError() {
		return nil, fmt.Errorf("error converting remote addresses: %v", diag)
	}
	sort.Strings(d.IPPools)
	ipPools, diag := types.SetValueFrom(context.TODO(), types.StringType, d.IPPools)
	if diag.HasError() {
		return nil, fmt.Errorf("error converting IP pools: %v", diag)
	}
	return &IpsecConnectionResourceModel{
		Enabled:                types.StringValue(d.Enabled),
		Proposals:              proposals,
		Unique:                 types.StringValue(d.Unique.String()),
		Aggressive:             types.StringValue(d.Aggressive),
		Version:                types.StringValue(d.Version.String()),
		Mobike:                 types.StringValue(d.Mobike),
		LocalAddresses:         localAddresses,
		RemoteAddresses:        remoteAddresses,
		LocalPort:              types.StringValue(d.LocalPort.String()),
		RemotePort:             types.StringValue(d.RemotePort.String()),
		UDPEncapsulation:       types.StringValue(d.UDPEncapsulation),
		ReauthenticationTime:   types.StringValue(d.ReauthenticationTime),
		RekeyTime:              types.StringValue(d.RekeyTime),
		IKELifetime:            types.StringValue(d.IKELifetime),
		DPDDelay:               types.StringValue(d.DPDDelay),
		DPDTimeout:             types.StringValue(d.DPDTimeout),
		IPPools:                ipPools,
		SendCertificateRequest: types.StringValue(d.SendCertificateRequest),
		SendCertificate:        types.StringValue(d.SendCertificate.String()),
		KeyingTries:            types.StringValue(d.KeyingTries),
		Description:            types.StringValue(d.Description),
	}, nil
}
