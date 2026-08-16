package diagnostics

import (
	"context"
	"sort"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/diagnostics"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// netflowResourceModel describes the resource data model.
// This is a SINGLETON resource - it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type netflowResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	Interfaces           types.Set    `tfsdk:"interfaces"`
	EgressOnlyInterfaces types.Set    `tfsdk:"egress_only_interfaces"`
	Version              types.String `tfsdk:"version"`
	Targets              types.Set    `tfsdk:"targets"`
	ActiveTimeout        types.Int64  `tfsdk:"active_timeout"`
	InactiveTimeout      types.Int64  `tfsdk:"inactive_timeout"`
}

func netflowResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages Netflow diagnostics configuration. This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_diagnostics_netflow.netflow netflow_settings\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `netflow_settings`. Use this value when importing: `terraform import opnsense_diagnostics_netflow.netflow netflow_settings`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable Netflow data collection. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"interfaces": schema.SetAttribute{
				MarkdownDescription: "Interfaces to capture Netflow data on. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"egress_only_interfaces": schema.SetAttribute{
				MarkdownDescription: "Interfaces to capture egress-only Netflow data on. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Netflow protocol version to export. Defaults to `\"v9\"`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("v9"),
				Validators: []validator.String{
					stringvalidator.OneOf("v5", "v9"),
				},
			},
			"targets": schema.SetAttribute{
				MarkdownDescription: "Netflow collector targets, as `host:port` pairs. Defaults to `[]`.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(tools.EmptySetValue(types.StringType)),
			},
			"active_timeout": schema.Int64Attribute{
				MarkdownDescription: "Active flow timeout, in seconds. Defaults to `1800`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1800),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"inactive_timeout": schema.Int64Attribute{
				MarkdownDescription: "Inactive flow timeout, in seconds. Defaults to `15`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(15),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
		},
	}
}

// convertNetflowSchemaToStruct converts TF schema to upstream API struct
func convertNetflowSchemaToStruct(d *netflowResourceModel) (*diagnostics.Netflow, error) {
	result := &diagnostics.Netflow{}

	result.Collect.Enable = tools.BoolToString(d.Enabled.ValueBool())
	result.Capture.Version = api.SelectedMap(d.Version.ValueString())
	result.ActiveTimeout = tools.Int64ToString(d.ActiveTimeout.ValueInt64())
	result.InactiveTimeout = tools.Int64ToString(d.InactiveTimeout.ValueInt64())

	if !d.Interfaces.IsNull() && !d.Interfaces.IsUnknown() {
		var interfaces []string
		d.Interfaces.ElementsAs(context.Background(), &interfaces, false)
		sort.Strings(interfaces)
		result.Capture.Interfaces = api.SelectedMapList(interfaces)
	}

	if !d.EgressOnlyInterfaces.IsNull() && !d.EgressOnlyInterfaces.IsUnknown() {
		var egressOnly []string
		d.EgressOnlyInterfaces.ElementsAs(context.Background(), &egressOnly, false)
		sort.Strings(egressOnly)
		result.Capture.EgressOnly = api.SelectedMapList(egressOnly)
	}

	if !d.Targets.IsNull() && !d.Targets.IsUnknown() {
		var targets []string
		d.Targets.ElementsAs(context.Background(), &targets, false)
		sort.Strings(targets)
		result.Capture.Targets = api.SelectedMapList(targets)
	}

	return result, nil
}

// convertNetflowStructToSchema converts upstream API struct to TF schema
func convertNetflowStructToSchema(netflow *diagnostics.NetflowMonad) (*netflowResourceModel, error) {
	model := &netflowResourceModel{
		Enabled:         types.BoolValue(tools.StringToBool(netflow.Netflow.Collect.Enable)),
		Version:         types.StringValue(netflow.Netflow.Capture.Version.String()),
		ActiveTimeout:   types.Int64Value(tools.StringToInt64(netflow.Netflow.ActiveTimeout)),
		InactiveTimeout: types.Int64Value(tools.StringToInt64(netflow.Netflow.InactiveTimeout)),
		Interfaces:      tools.StringSliceToSet(netflow.Netflow.Capture.Interfaces),
		EgressOnlyInterfaces: tools.StringSliceToSet(
			netflow.Netflow.Capture.EgressOnly,
		),
		Targets: tools.StringSliceToSet(netflow.Netflow.Capture.Targets),
	}

	return model, nil
}
