package interfaces

import (
	"regexp"

	"github.com/browningluke/opnsense-go/pkg/interfaces"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// assignmentResourceModel describes the resource data model.
type assignmentResourceModel struct {
	Device      types.String `tfsdk:"device"`
	Description types.String `tfsdk:"description"`
	Lock        types.Bool   `tfsdk:"lock"`
	Id          types.String `tfsdk:"id"`
}

func assignmentResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Assigns a physical (or virtual) device to a logical interface, such as `wan`, `lan`, or an optional `optN` interface. " +
			"Requires OPNsense 26.7 or later; the underlying API endpoint doesn't exist on earlier versions. " +
			"OPNsense's REST API does not currently expose the interface's general configuration (IPv4/IPv6 addressing mode, enable, MTU, etc.) " +
			"— only this device assignment, description, and lock flag can be managed here.\n\n" +
			"To manage the existing `wan` or `lan` interface, import it first: `terraform import opnsense_interfaces_assignment.wan wan`. " +
			"Creating a new assignment adds an additional optional interface (`opt1`, `opt2`, ...); `terraform destroy` unassigns the interface.",

		Attributes: map[string]schema.Attribute{
			"device": schema.StringAttribute{
				MarkdownDescription: "The physical (or virtual) device to assign to this interface, e.g. `vtnet0`.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed). " +
					"OPNsense only accepts letters, digits and underscores here — spaces and other " +
					"punctuation are rejected by the API.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9_]{0,255}$`),
						"must contain only letters, digits and underscores, and be at most 255 characters",
					),
				},
			},
			"lock": schema.BoolAttribute{
				// OPNsense's own model defaults lock to 1, but a locked interface cannot be
				// removed ("Interface locked, unset lock first before removal"), which would
				// leave `terraform destroy` unable to unassign what it created. Default to
				// false so the resource lifecycle is complete; set it explicitly to opt in.
				MarkdownDescription: "Prevent removal of the interface. Defaults to `false`. " +
					"Note this differs from the OPNsense web UI, which defaults to locked; a locked " +
					"interface cannot be unassigned, so `terraform destroy` fails until `lock` is set back to `false`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Logical interface identifier (e.g. `wan`, `lan`, `opt1`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func assignmentDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Assigns a physical (or virtual) device to a logical interface, such as `wan`, `lan`, or an optional `optN` interface.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "Logical interface identifier (e.g. `wan`, `lan`, `opt1`).",
				Required:            true,
			},
			"device": dschema.StringAttribute{
				MarkdownDescription: "The physical (or virtual) device assigned to this interface.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Computed:            true,
			},
			"lock": dschema.BoolAttribute{
				MarkdownDescription: "Whether removal of the interface is prevented.",
				Computed:            true,
			},
		},
	}
}

func convertAssignmentSchemaToStruct(d *assignmentResourceModel) (*interfaces.Assignment, error) {
	return &interfaces.Assignment{
		Device:      api.SelectedMap(d.Device.ValueString()),
		Description: d.Description.ValueString(),
		Lock:        tools.BoolToString(d.Lock.ValueBool()),
	}, nil
}

func convertAssignmentStructToSchema(d *interfaces.Assignment) (*assignmentResourceModel, error) {
	return &assignmentResourceModel{
		Device:      types.StringValue(d.Device.String()),
		Description: tools.StringOrNull(d.Description),
		Lock:        types.BoolValue(tools.StringToBool(d.Lock)),
	}, nil
}
