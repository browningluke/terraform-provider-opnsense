package dnsmasq

import (
	"context"
	"fmt"
	"regexp"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/browningluke/terraform-provider-opnsense/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type rangeResourceModel struct {
	Interface     types.String `tfsdk:"interface"`
	Tag           types.String `tfsdk:"tag"`
	StartAddress  types.String `tfsdk:"start_address"`
	EndAddress    types.String `tfsdk:"end_address"`
	SubnetMask    types.String `tfsdk:"subnet_mask"`
	Constructor   types.String `tfsdk:"constructor"`
	Mode          types.Set    `tfsdk:"mode"`
	LeaseTime     types.Int64  `tfsdk:"lease_time"`
	DomainType    types.String `tfsdk:"domain_type"`
	Domain        types.String `tfsdk:"domain"`
	DisableHASync types.Bool   `tfsdk:"nosync"`
	Description   types.String `tfsdk:"description"`

	// IPv6 Settings
	PrefixLength     types.Int64  `tfsdk:"prefix_length"`
	RaMode           types.Set    `tfsdk:"ra_mode"`
	RaPriority       types.String `tfsdk:"ra_priority"`
	RaMTU            types.Int64  `tfsdk:"ra_mtu"`
	RaInterval       types.Int64  `tfsdk:"ra_interval"`
	RaRouterLifetime types.Int64  `tfsdk:"ra_router_lifetime"`

	Id types.String `tfsdk:"id"`
}

func rangeResourceSchema() schema.Schema {
	return schema.Schema{
		Version:             1,
		MarkdownDescription: "Configure DHCP range configuration for dnsmasq.",

		Attributes: map[string]schema.Attribute{
			"interface": schema.StringAttribute{
				MarkdownDescription: "Interface to serve this range",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "Optional tag to set for requests matching this range which can be used to selectively match DHCP options",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`),
						"must be UUID",
					),
				},
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"start_address": schema.StringAttribute{
				MarkdownDescription: "Start of the range, e.g. 192.168.1.100 for DHCPv4, 2000::1 for DHCPv6 or when a constructor is using a suffix like ::1. To reveal IPv6 related options, enter a IPv6 address. When using router advertisements, it is possible to use a constructor with :: as the start address and no end address.",
				Required:            true,
				Validators: []validator.String{
					validators.IP(),
				},
			},
			"end_address": schema.StringAttribute{
				MarkdownDescription: "End of the range.",
				Optional:            true,
				Validators: []validator.String{
					validators.IP(),
				},
			},
			"subnet_mask": schema.StringAttribute{
				MarkdownDescription: "Leave empty to auto-calculate the subnet mask from the interface or the network class of the start address. If a DHCP relay forwards IPv4 DHCP Discovers to Dnsmasq, setting a subnet mask is required in most cases.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"constructor": schema.StringAttribute{
				MarkdownDescription: "Interface to use to calculate a DHCPv6 or RA range. Start address can then be specified as a suffix (e.g. ::, ::1 or ::400).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"mode": schema.SetAttribute{
				MarkdownDescription: "Mode flags to set for this range, 'static' means no addresses will be automatically assigned.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf(
						"static",
					)),
				},
			},
			"lease_time": schema.Int64Attribute{
				MarkdownDescription: "Defines how long the addresses (leases) given out by the server are valid (in seconds). Set 0 for infinite.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(86400),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"domain_type": schema.StringAttribute{
				MarkdownDescription: "Choose if the domain will only match clients in this range, or all clients in any subnets on the selected interface. If you create both IPv4 and IPv6 ranges, setting this to \"Interface\" on both ranges is recommended.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("range", "interface"),
				},
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Offer this domain to DHCP clients.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"nosync": schema.BoolAttribute{
				MarkdownDescription: "Ignore this range from being transfered or updated by ha sync.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "You may enter a description here for your reference (not parsed).",
				Optional:            true,
			},

			// -----------------
			// IPv6 / RA
			// -----------------
			"prefix_length": schema.Int64Attribute{
				MarkdownDescription: "Prefix length offered to the client. Custom values in this field will be ignored if Router Advertisements are enabled, as SLAAC will only work with a prefix length of 64.",
				Optional:            true,
				Computed:            false,
				Validators: []validator.Int64{
					int64validator.Between(1, 64),
				},
			},
			"ra_mode": schema.SetAttribute{
				MarkdownDescription: "Control how IPv6 clients receive their addresses. Enabling Router Advertisements in general settings will enable it for all configured DHCPv6 ranges with the managed address bits set, and the use SLAAC bit reset. To change this default, select a combination of the possible options here. \"slaac\", \"ra-stateless\" and \"ra-names\" can be freely combined, all other options shall remain single selections.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            false,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf("ra-only", "slaac", "ra-names", "ra-stateless", "ra-advrouter", "off-link")),
				},
			},
			"ra_priority": schema.StringAttribute{
				MarkdownDescription: "Priority of the RA announcements.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("high", "", "low"),
				},
				Default: stringdefault.StaticString(""),
			},
			"ra_mtu": schema.Int64Attribute{
				MarkdownDescription: "Optional MTU to send to clients via Router Advertisements. If unsure leave empty.",
				Optional:            true,
				Computed:            false,
			},
			"ra_interval": schema.Int64Attribute{
				MarkdownDescription: "Time (seconds) between Router Advertisements.",
				Optional:            true,
				Computed:            false,
			},
			"ra_router_lifetime": schema.Int64Attribute{
				MarkdownDescription: "The lifetime of the route may be changed or set to zero, which allows a router to advertise prefixes but not a route via itself. When using HA, setting a short timespan here is adviced for faster IPv6 failover. A good combination could be 10 seconds RA interval and 30 seconds RA router lifetime. Going lower than that can pose issues in busy networks.",
				Optional:            true,
				Computed:            false,
			},

			// -----------------
			// Meta
			// -----------------
			"id": schema.StringAttribute{
				MarkdownDescription: "Range UUID.",
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func rangeDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Read DHCP range configuration from dnsmasq.",

		Attributes: map[string]dschema.Attribute{
			"interface": dschema.StringAttribute{
				MarkdownDescription: "Interface to serve this range",
				Computed:            true,
			},
			"tag": dschema.StringAttribute{
				MarkdownDescription: "Optional tag to set for requests matching this range which can be used to selectively match DHCP options",
				Computed:            true,
			},
			"start_address": dschema.StringAttribute{
				MarkdownDescription: "Start of the range, e.g. 192.168.1.100 for DHCPv4, 2000::1 for DHCPv6 or when a constructor is using a suffix like ::1. To reveal IPv6 related options, enter a IPv6 address. When using router advertisements, it is possible to use a constructor with :: as the start address and no end address.",
				Computed:            true,
			},
			"end_address": dschema.StringAttribute{
				MarkdownDescription: "End of the range.",
				Computed:            true,
			},
			"subnet_mask": dschema.StringAttribute{
				MarkdownDescription: "Leave empty to auto-calculate the subnet mask from the interface or the network class of the start address. If a DHCP relay forwards IPv4 DHCP Discovers to Dnsmasq, setting a subnet mask is required in most cases.",
				Computed:            true,
			},
			"constructor": dschema.StringAttribute{
				MarkdownDescription: "Interface to use to calculate a DHCPv6 or RA range. Start address can then be specified as a suffix (e.g. ::, ::1 or ::400).",
				Computed:            true,
			},
			"mode": dschema.SetAttribute{
				MarkdownDescription: "Mode flags to set for this range, 'static' means no addresses will be automatically assigned.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"lease_time": dschema.Int64Attribute{
				MarkdownDescription: "Defines how long the addresses (leases) given out by the server are valid (in seconds). Set 0 for infinite.",
				Computed:            true,
			},
			"domain_type": dschema.StringAttribute{
				MarkdownDescription: "Choose if the domain will only match clients in this range, or all clients in any subnets on the selected interface. If you create both IPv4 and IPv6 ranges, setting this to \"Interface\" on both ranges is recommended.",
				Computed:            true,
			},
			"domain": dschema.StringAttribute{
				MarkdownDescription: "Offer this domain to DHCP clients.",
				Computed:            true,
			},
			"nosync": dschema.BoolAttribute{
				MarkdownDescription: "Ignore this range from being transfered or updated by ha sync.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "You may enter a description here for your reference (not parsed).",
				Computed:            true,
			},

			// IPv6 / RA
			"prefix_length": dschema.Int64Attribute{
				MarkdownDescription: "Prefix length offered to the client. Custom values in this field will be ignored if Router Advertisements are enabled, as SLAAC will only work with a prefix length of 64.",
				Computed:            true,
			},
			"ra_mode": dschema.SetAttribute{
				MarkdownDescription: "Control how IPv6 clients receive their addresses. Enabling Router Advertisements in general settings will enable it for all configured DHCPv6 ranges with the managed address bits set, and the use SLAAC bit reset. To change this default, select a combination of the possible options here. \"slaac\", \"ra-stateless\" and \"ra-names\" can be freely combined, all other options shall remain single selections.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ra_priority": dschema.StringAttribute{
				MarkdownDescription: "Priority of the RA announcements.",
				Computed:            true,
			},
			"ra_mtu": dschema.Int64Attribute{
				MarkdownDescription: "Optional MTU to send to clients via Router Advertisements. If unsure leave empty.",
				Computed:            true,
			},
			"ra_interval": dschema.Int64Attribute{
				MarkdownDescription: "Time (seconds) between Router Advertisements.",
				Computed:            true,
			},
			"ra_router_lifetime": dschema.Int64Attribute{
				MarkdownDescription: "The lifetime of the route may be changed or set to zero, which allows a router to advertise prefixes but not a route via itself. When using HA, setting a short timespan here is adviced for faster IPv6 failover. A good combination could be 10 seconds RA interval and 30 seconds RA router lifetime. Going lower than that can pose issues in busy networks.",
				Computed:            true,
			},
			"id": dschema.StringAttribute{
				MarkdownDescription: "Range UUID.",
				Required:            true,
			},
		},
	}
}

func mapValues(theSet types.Set, errorMessage string) ([]string, error) {
	var values []string

	diags := theSet.ElementsAs(
		context.Background(),
		&values,
		false,
	)

	if diags.HasError() {
		return nil, fmt.Errorf("%s", errorMessage)
	}

	return values, nil
}

func convertRangeSchemaToStruct(d *rangeResourceModel) (*dnsmasq.Range, error) {
	resultRange := &dnsmasq.Range{
		Interface:     api.SelectedMap(d.Interface.ValueString()),
		Tag:           api.SelectedMap(d.Tag.ValueString()),
		StartAddress:  d.StartAddress.ValueString(),
		EndAddress:    d.EndAddress.ValueString(),
		SubnetMask:    d.SubnetMask.ValueString(),
		Constructor:   api.SelectedMap(d.Constructor.ValueString()),
		LeaseTime:     d.LeaseTime.String(),
		DomainType:    api.SelectedMap(d.DomainType.ValueString()),
		Domain:        d.Domain.ValueString(),
		DisableHASync: tools.BoolToString(d.DisableHASync.ValueBool()),
		Description:   d.Description.ValueString(),

		// IPv6 / RA
		RaPriority: api.SelectedMap(d.RaPriority.ValueString()),
	}

	if !d.RaMode.IsNull() && !d.RaMode.IsUnknown() {
		values, errorMessage := mapValues(d.RaMode, "failed to parse Ra Modes")

		if errorMessage != nil {
			return nil, errorMessage
		}
		resultRange.RaMode = api.SelectedMapList(values)
	}

	if !d.Mode.IsNull() && !d.Mode.IsUnknown() {
		values, errorMessage := mapValues(d.Mode, "failed to parse Modes")

		if errorMessage != nil {
			return nil, errorMessage
		}
		resultRange.Mode = api.SelectedMapList(values)
	}

	// The OPNsense API represents unset numeric IPv6 values as empty strings
	// (and prefix_length as -1). Do not turn those sentinels into Terraform
	// values, otherwise a configuration that omits the attributes never
	// reaches a stable state after refresh.
	if !d.PrefixLength.IsNull() && !d.PrefixLength.IsUnknown() {
		resultRange.PrefixLength = d.PrefixLength.String()
	}
	if !d.RaInterval.IsNull() && !d.RaInterval.IsUnknown() {
		resultRange.RaInterval = d.RaInterval.String()
	}
	if !d.RaMTU.IsNull() && !d.RaMTU.IsUnknown() {
		resultRange.RaMTU = d.RaMTU.String()
	}
	if !d.RaRouterLifetime.IsNull() && !d.RaRouterLifetime.IsUnknown() {
		resultRange.RaRouterLifetime = d.RaRouterLifetime.String()
	}

	return resultRange, nil
}

func convertRangeStructToSchema(d *dnsmasq.Range) (*rangeResourceModel, error) {
	raModeSet, diags := types.SetValueFrom(
		context.Background(),
		types.StringType,
		d.RaMode,
	)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert RaMode")
	}
	modeSet, diags := types.SetValueFrom(
		context.Background(),
		types.StringType,
		d.Mode,
	)
	if diags.HasError() {
		return nil, fmt.Errorf("Mode")
	}

	model := &rangeResourceModel{
		Interface:     types.StringValue(d.Interface.String()),
		Tag:           types.StringValue(d.Tag.String()),
		StartAddress:  types.StringValue(d.StartAddress),
		EndAddress:    types.StringValue(d.EndAddress),
		SubnetMask:    types.StringValue(d.SubnetMask),
		Constructor:   types.StringValue(d.Constructor.String()),
		LeaseTime:     types.Int64Value(tools.StringToInt64(d.LeaseTime)),
		DomainType:    types.StringValue(d.DomainType.String()),
		Domain:        types.StringValue(d.Domain),
		DisableHASync: types.BoolValue(tools.StringToBool(d.DisableHASync)),
		Description:   types.StringValue(d.Description),

		// IPv6 / RA
		RaPriority: types.StringValue(d.RaPriority.String()),
	}
	model.RaMode = raModeSet
	model.Mode = modeSet
	if d.PrefixLength != "" && d.PrefixLength != "-1" {
		model.PrefixLength = types.Int64Value(tools.StringToInt64(d.PrefixLength))
	} else {
		model.PrefixLength = types.Int64Null()
	}

	if d.RaInterval != "" {
		model.RaInterval = types.Int64Value(tools.StringToInt64(d.RaInterval))
	}
	if d.RaMTU != "" {
		model.RaMTU = types.Int64Value(tools.StringToInt64(d.RaMTU))
	}
	if d.RaRouterLifetime != "" {
		model.RaRouterLifetime = types.Int64Value(tools.StringToInt64(d.RaRouterLifetime))
	}
	tflog.Debug(context.Background(), "range state", map[string]interface{}{
		"disable_ha_sync": model.DisableHASync,
	})
	return model, nil
}
