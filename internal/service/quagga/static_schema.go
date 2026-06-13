package quagga

import (
	"github.com/browningluke/opnsense-go/pkg/quagga"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// staticResourceModel describes the resource data model.
// This is a SINGLETON resource — it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type staticResourceModel struct {
	Id      types.String `tfsdk:"id"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func staticResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages the Quagga staticd (static routing daemon) settings. This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_quagga_static.static quagga_static\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `quagga_static`. Use this value when importing: `terraform import opnsense_quagga_static.static quagga_static`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable the staticd daemon. staticd is required to manage static routes via the FRRouting stack. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func convertStaticSchemaToStruct(d *staticResourceModel) (*quagga.QuaggaStatic, error) {
	return &quagga.QuaggaStatic{
		Enabled: tools.BoolToString(d.Enabled.ValueBool()),
	}, nil
}

func convertStaticStructToSchema(d *quagga.QuaggaStatic) (*staticResourceModel, error) {
	return &staticResourceModel{
		Enabled: types.BoolValue(tools.StringToBool(d.Enabled)),
	}, nil
}
