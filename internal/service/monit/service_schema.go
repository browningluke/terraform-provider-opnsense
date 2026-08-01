package monit

import (
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/monit"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type serviceResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Type         types.String `tfsdk:"type"`
	Pidfile      types.String `tfsdk:"pidfile"`
	Match        types.String `tfsdk:"match"`
	Path         types.String `tfsdk:"path"`
	Timeout      types.String `tfsdk:"timeout"`
	StartTimeout types.String `tfsdk:"start_timeout"`
	Address      types.String `tfsdk:"address"`
	Interface    types.String `tfsdk:"interface"`
	Start        types.String `tfsdk:"start"`
	Stop         types.String `tfsdk:"stop"`
	Tests        types.Set    `tfsdk:"tests"`
	Depends      types.Set    `tfsdk:"depends"`
	Polltime     types.String `tfsdk:"polltime"`
}

func serviceResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a Monit monitored service (process, file, filesystem, host, or custom check).",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the service.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable monitoring of this service. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Unique name for this monitored service.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for this service.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of service check: `process`, `file`, `fifo`, `filesystem`, `directory`, `host`, `system`, `custom`, `network`.",
				Required:            true,
			},
			"pidfile": schema.StringAttribute{
				MarkdownDescription: "Absolute path to the PID file. Used when `type` is `process`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"match": schema.StringAttribute{
				MarkdownDescription: "Pattern used to match a process by command line when `type` is `process` without a `pidfile`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Absolute path to the file, directory, or FIFO being monitored.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"timeout": schema.StringAttribute{
				MarkdownDescription: "Number of cycles to wait for a start/stop action to complete before timing out. Defaults to `300`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("300"),
			},
			"start_timeout": schema.StringAttribute{
				MarkdownDescription: "Grace period, in seconds, given to the service after it starts before checks begin failing it. Defaults to `0`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("0"),
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "Hostname or IP address to check. Used when `type` is `host` or `network`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"interface": schema.StringAttribute{
				MarkdownDescription: "Network interface to monitor. Used when `type` is `network`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"start": schema.StringAttribute{
				MarkdownDescription: "Command (with arguments) used to start the service.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"stop": schema.StringAttribute{
				MarkdownDescription: "Command (with arguments) used to stop the service.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"tests": schema.SetAttribute{
				MarkdownDescription: "Set of Monit test UUIDs to apply to this service. Test objects are managed outside this provider (via the OPNsense UI or API directly).",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"depends": schema.SetAttribute{
				MarkdownDescription: "Set of other Monit service UUIDs this service depends on.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"polltime": schema.StringAttribute{
				MarkdownDescription: "Custom cron-like poll schedule for this service. When empty, the global Monit polling interval is used.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
		},
	}
}

func serviceDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Retrieves a Monit monitored service by UUID.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "UUID of the service.",
			},
			"enabled": dschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether monitoring of this service is enabled.",
			},
			"name": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the monitored service.",
			},
			"description": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description of the service.",
			},
			"type": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Type of service check.",
			},
			"pidfile": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Absolute path to the PID file.",
			},
			"match": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Process match pattern.",
			},
			"path": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Absolute path to the file, directory, or FIFO being monitored.",
			},
			"timeout": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Number of cycles to wait for a start/stop action to complete.",
			},
			"start_timeout": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Grace period, in seconds, given to the service after it starts.",
			},
			"address": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Hostname or IP address being checked.",
			},
			"interface": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Network interface being monitored.",
			},
			"start": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Command used to start the service.",
			},
			"stop": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Command used to stop the service.",
			},
			"tests": dschema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "Set of Monit test UUIDs applied to this service.",
				ElementType:         types.StringType,
			},
			"depends": dschema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "Set of other Monit service UUIDs this service depends on.",
				ElementType:         types.StringType,
			},
			"polltime": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Custom cron-like poll schedule for this service.",
			},
		},
	}
}

func convertServiceSchemaToStruct(d *serviceResourceModel) (*monit.Service, error) {
	return &monit.Service{
		Enabled:      tools.BoolToString(d.Enabled.ValueBool()),
		Name:         d.Name.ValueString(),
		Description:  d.Description.ValueString(),
		Type:         api.SelectedMap(d.Type.ValueString()),
		Pidfile:      d.Pidfile.ValueString(),
		Match:        d.Match.ValueString(),
		Path:         d.Path.ValueString(),
		Timeout:      d.Timeout.ValueString(),
		Starttimeout: d.StartTimeout.ValueString(),
		Address:      d.Address.ValueString(),
		Interface:    api.SelectedMap(d.Interface.ValueString()),
		Start:        d.Start.ValueString(),
		Stop:         d.Stop.ValueString(),
		Tests:        api.SelectedMapList(tools.SetToStringSlice(d.Tests)),
		Depends:      api.SelectedMapList(tools.SetToStringSlice(d.Depends)),
		Polltime:     d.Polltime.ValueString(),
	}, nil
}

func convertServiceStructToSchema(d *monit.Service) (*serviceResourceModel, error) {
	return &serviceResourceModel{
		Enabled:      types.BoolValue(tools.StringToBool(d.Enabled)),
		Name:         types.StringValue(d.Name),
		Description:  types.StringValue(d.Description),
		Type:         types.StringValue(d.Type.String()),
		Pidfile:      types.StringValue(d.Pidfile),
		Match:        types.StringValue(d.Match),
		Path:         types.StringValue(d.Path),
		Timeout:      types.StringValue(d.Timeout),
		StartTimeout: types.StringValue(d.Starttimeout),
		Address:      types.StringValue(d.Address),
		Interface:    types.StringValue(d.Interface.String()),
		Start:        types.StringValue(d.Start),
		Stop:         types.StringValue(d.Stop),
		Tests:        tools.StringSliceToSet(d.Tests),
		Depends:      tools.StringSliceToSet(d.Depends),
		Polltime:     types.StringValue(d.Polltime),
	}, nil
}
