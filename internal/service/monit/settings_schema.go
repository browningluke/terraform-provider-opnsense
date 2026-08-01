package monit

import (
	"strings"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/monit"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type settingsResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	Interval                  types.String `tfsdk:"interval"`
	Startdelay                types.String `tfsdk:"startdelay"`
	Mailserver                types.String `tfsdk:"mailserver"`
	Port                      types.String `tfsdk:"port"`
	Username                  types.String `tfsdk:"username"`
	Password                  types.String `tfsdk:"password"`
	Ssl                       types.Bool   `tfsdk:"ssl"`
	SslVersion                types.String `tfsdk:"ssl_version"`
	SslVerify                 types.Bool   `tfsdk:"ssl_verify"`
	Logfile                   types.String `tfsdk:"logfile"`
	Statefile                 types.String `tfsdk:"statefile"`
	EventqueuePath            types.String `tfsdk:"eventqueue_path"`
	EventqueueSlots           types.String `tfsdk:"eventqueue_slots"`
	HttpdEnabled              types.Bool   `tfsdk:"httpd_enabled"`
	HttpdUsername             types.String `tfsdk:"httpd_username"`
	HttpdPassword             types.String `tfsdk:"httpd_password"`
	HttpdPort                 types.String `tfsdk:"httpd_port"`
	HttpdAllow                types.String `tfsdk:"httpd_allow"`
	MmonitUrl                 types.String `tfsdk:"mmonit_url"`
	MmonitTimeout             types.String `tfsdk:"mmonit_timeout"`
	MmonitRegisterCredentials types.Bool   `tfsdk:"mmonit_register_credentials"`
}

func settingsResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages Monit general settings (polling interval, mail server, web GUI). This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_monit_settings.settings monit_settings\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `monit_settings`. Use this value when importing: `terraform import opnsense_monit_settings.settings monit_settings`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable the Monit service monitoring daemon.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"interval": schema.StringAttribute{
				MarkdownDescription: "Polling interval, in seconds, between two checks of the services.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"startdelay": schema.StringAttribute{
				MarkdownDescription: "Time to wait, in seconds, before the first check after Monit startup (allows daemons time to start).",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"mailserver": schema.StringAttribute{
				MarkdownDescription: "SMTP server used to send alert notifications.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"port": schema.StringAttribute{
				MarkdownDescription: "SMTP server port.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username used to authenticate against the SMTP server.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password used to authenticate against the SMTP server.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"ssl": schema.BoolAttribute{
				MarkdownDescription: "When enabled, use SSL/TLS to connect to the SMTP server.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"ssl_version": schema.StringAttribute{
				MarkdownDescription: "SSL/TLS version to use for the SMTP connection: `auto`, `tlsv1`, `tlsv11`, `tlsv12`, `tlsv13`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"ssl_verify": schema.BoolAttribute{
				MarkdownDescription: "When enabled, verify the SMTP server's SSL certificate.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"logfile": schema.StringAttribute{
				MarkdownDescription: "Absolute path to the Monit log file, or `syslog` to log via syslog.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"statefile": schema.StringAttribute{
				MarkdownDescription: "Absolute path to the Monit state file.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"eventqueue_path": schema.StringAttribute{
				MarkdownDescription: "Directory used to queue alert events when they cannot be delivered immediately.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"eventqueue_slots": schema.StringAttribute{
				MarkdownDescription: "Maximum number of events that can be queued.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"httpd_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable the Monit HTTP(S) web interface.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"httpd_username": schema.StringAttribute{
				MarkdownDescription: "Username for the Monit web interface.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"httpd_password": schema.StringAttribute{
				MarkdownDescription: "Password for the Monit web interface. Required when `httpd_enabled` is `true`.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"httpd_port": schema.StringAttribute{
				MarkdownDescription: "Port the Monit web interface listens on.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"httpd_allow": schema.StringAttribute{
				MarkdownDescription: "Comma-separated list of hosts/networks allowed to access the Monit web interface.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"mmonit_url": schema.StringAttribute{
				MarkdownDescription: "URL of an M/Monit server to register with.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"mmonit_timeout": schema.StringAttribute{
				MarkdownDescription: "Connection timeout, in seconds, for the M/Monit server.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"mmonit_register_credentials": schema.BoolAttribute{
				MarkdownDescription: "When enabled, register this system's credentials with the configured M/Monit server.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func settingsDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Reads Monit general settings from the upstream system.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `monit_settings`.",
			},
			"enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether the Monit service monitoring daemon is enabled.",
				Computed:            true,
			},
			"interval": dschema.StringAttribute{
				MarkdownDescription: "Polling interval, in seconds, between two checks of the services.",
				Computed:            true,
			},
			"startdelay": dschema.StringAttribute{
				MarkdownDescription: "Time to wait, in seconds, before the first check after Monit startup.",
				Computed:            true,
			},
			"mailserver": dschema.StringAttribute{
				MarkdownDescription: "SMTP server used to send alert notifications.",
				Computed:            true,
			},
			"port": dschema.StringAttribute{
				MarkdownDescription: "SMTP server port.",
				Computed:            true,
			},
			"username": dschema.StringAttribute{
				MarkdownDescription: "Username used to authenticate against the SMTP server.",
				Computed:            true,
			},
			"password": dschema.StringAttribute{
				MarkdownDescription: "Password used to authenticate against the SMTP server.",
				Computed:            true,
				Sensitive:           true,
			},
			"ssl": dschema.BoolAttribute{
				MarkdownDescription: "Whether SSL/TLS is used to connect to the SMTP server.",
				Computed:            true,
			},
			"ssl_version": dschema.StringAttribute{
				MarkdownDescription: "SSL/TLS version used for the SMTP connection.",
				Computed:            true,
			},
			"ssl_verify": dschema.BoolAttribute{
				MarkdownDescription: "Whether the SMTP server's SSL certificate is verified.",
				Computed:            true,
			},
			"logfile": dschema.StringAttribute{
				MarkdownDescription: "Absolute path to the Monit log file, or `syslog`.",
				Computed:            true,
			},
			"statefile": dschema.StringAttribute{
				MarkdownDescription: "Absolute path to the Monit state file.",
				Computed:            true,
			},
			"eventqueue_path": dschema.StringAttribute{
				MarkdownDescription: "Directory used to queue alert events.",
				Computed:            true,
			},
			"eventqueue_slots": dschema.StringAttribute{
				MarkdownDescription: "Maximum number of events that can be queued.",
				Computed:            true,
			},
			"httpd_enabled": dschema.BoolAttribute{
				MarkdownDescription: "Whether the Monit HTTP(S) web interface is enabled.",
				Computed:            true,
			},
			"httpd_username": dschema.StringAttribute{
				MarkdownDescription: "Username for the Monit web interface.",
				Computed:            true,
			},
			"httpd_password": dschema.StringAttribute{
				MarkdownDescription: "Password for the Monit web interface.",
				Computed:            true,
				Sensitive:           true,
			},
			"httpd_port": dschema.StringAttribute{
				MarkdownDescription: "Port the Monit web interface listens on.",
				Computed:            true,
			},
			"httpd_allow": dschema.StringAttribute{
				MarkdownDescription: "Comma-separated list of hosts/networks allowed to access the Monit web interface.",
				Computed:            true,
			},
			"mmonit_url": dschema.StringAttribute{
				MarkdownDescription: "URL of an M/Monit server to register with.",
				Computed:            true,
			},
			"mmonit_timeout": dschema.StringAttribute{
				MarkdownDescription: "Connection timeout, in seconds, for the M/Monit server.",
				Computed:            true,
			},
			"mmonit_register_credentials": dschema.BoolAttribute{
				MarkdownDescription: "Whether this system's credentials are registered with the configured M/Monit server.",
				Computed:            true,
			},
		},
	}
}

// commaListToSlice splits a comma-separated option list into its keys;
// an empty string maps to nil so nothing is marked selected upstream.
func commaListToSlice(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func convertSettingsSchemaToStruct(d *settingsResourceModel) (*monit.Settings, error) {
	return &monit.Settings{
		Enabled:                   tools.BoolToString(d.Enabled.ValueBool()),
		Interval:                  d.Interval.ValueString(),
		Startdelay:                d.Startdelay.ValueString(),
		Mailserver:                api.SelectedMapList(commaListToSlice(d.Mailserver.ValueString())),
		Port:                      d.Port.ValueString(),
		Username:                  d.Username.ValueString(),
		Password:                  d.Password.ValueString(),
		Ssl:                       tools.BoolToString(d.Ssl.ValueBool()),
		SslVersion:                api.SelectedMap(d.SslVersion.ValueString()),
		SslVerify:                 tools.BoolToString(d.SslVerify.ValueBool()),
		Logfile:                   d.Logfile.ValueString(),
		Statefile:                 d.Statefile.ValueString(),
		EventqueuePath:            d.EventqueuePath.ValueString(),
		EventqueueSlots:           d.EventqueueSlots.ValueString(),
		HttpdEnabled:              tools.BoolToString(d.HttpdEnabled.ValueBool()),
		HttpdUsername:             d.HttpdUsername.ValueString(),
		HttpdPassword:             d.HttpdPassword.ValueString(),
		HttpdPort:                 d.HttpdPort.ValueString(),
		HttpdAllow:                api.SelectedMapList(commaListToSlice(d.HttpdAllow.ValueString())),
		MmonitUrl:                 d.MmonitUrl.ValueString(),
		MmonitTimeout:             d.MmonitTimeout.ValueString(),
		MmonitRegisterCredentials: tools.BoolToString(d.MmonitRegisterCredentials.ValueBool()),
	}, nil
}

func convertSettingsStructToSchema(d *monit.Settings) (*settingsResourceModel, error) {
	return &settingsResourceModel{
		Enabled:                   types.BoolValue(tools.StringToBool(d.Enabled)),
		Interval:                  types.StringValue(d.Interval),
		Startdelay:                types.StringValue(d.Startdelay),
		Mailserver:                types.StringValue(d.Mailserver.String()),
		Port:                      types.StringValue(d.Port),
		Username:                  types.StringValue(d.Username),
		Password:                  types.StringValue(d.Password),
		Ssl:                       types.BoolValue(tools.StringToBool(d.Ssl)),
		SslVersion:                types.StringValue(d.SslVersion.String()),
		SslVerify:                 types.BoolValue(tools.StringToBool(d.SslVerify)),
		Logfile:                   types.StringValue(d.Logfile),
		Statefile:                 types.StringValue(d.Statefile),
		EventqueuePath:            types.StringValue(d.EventqueuePath),
		EventqueueSlots:           types.StringValue(d.EventqueueSlots),
		HttpdEnabled:              types.BoolValue(tools.StringToBool(d.HttpdEnabled)),
		HttpdUsername:             types.StringValue(d.HttpdUsername),
		HttpdPassword:             types.StringValue(d.HttpdPassword),
		HttpdPort:                 types.StringValue(d.HttpdPort),
		HttpdAllow:                types.StringValue(d.HttpdAllow.String()),
		MmonitUrl:                 types.StringValue(d.MmonitUrl),
		MmonitTimeout:             types.StringValue(d.MmonitTimeout),
		MmonitRegisterCredentials: types.BoolValue(tools.StringToBool(d.MmonitRegisterCredentials)),
	}, nil
}
