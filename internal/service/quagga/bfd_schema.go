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

// bfdResourceModel describes the resource data model.
// This is a SINGLETON resource — it manages existing upstream configuration
// that cannot be created or destroyed via Terraform.
type bfdResourceModel struct {
	Id      types.String `tfsdk:"id"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func bfdResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages Quagga BFD (Bidirectional Forwarding Detection) settings. This is a singleton resource that manages existing upstream configuration.\n\n" +
			"**Important:** This resource must be imported before it can be managed:\n" +
			"```bash\n" +
			"terraform import opnsense_quagga_bfd.bfd quagga_bfd\n" +
			"```\n\n" +
			"After importing, you can manage the configuration with `terraform apply`. " +
			"Running `terraform destroy` will remove the resource from state but will NOT modify the upstream configuration.",

		Version: 1,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Always set to `quagga_bfd`. Use this value when importing: `terraform import opnsense_quagga_bfd.bfd quagga_bfd`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable the BFD daemon. BFD provides fast failure detection for routing protocols such as BGP and OSPF. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func convertBFDSchemaToStruct(d *bfdResourceModel) (*quagga.QuaggaBFD, error) {
	return &quagga.QuaggaBFD{
		Enabled: tools.BoolToString(d.Enabled.ValueBool()),
	}, nil
}

func convertBFDStructToSchema(d *quagga.QuaggaBFD) (*bfdResourceModel, error) {
	return &bfdResourceModel{
		Enabled: types.BoolValue(tools.StringToBool(d.Enabled)),
	}, nil
}
