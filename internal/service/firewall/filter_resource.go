package firewall

import (
	"context"
	"errors"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/errs"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/browningluke/terraform-provider-opnsense/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &filterResource{}
var _ resource.ResourceWithConfigure = &filterResource{}
var _ resource.ResourceWithImportState = &filterResource{}
var _ resource.ResourceWithConfigValidators = &filterResource{}
var _ resource.ResourceWithUpgradeState = &filterResource{}

func newFilterResource() resource.Resource {
	return &filterResource{}
}

// filterResource defines the resource implementation.
type filterResource struct {
	client opnsense.Client
}

func (r *filterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_filter"
}

func (r *filterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = filterResourceSchema()
}

func (r *filterResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		// Ensure adaptive end is greater than or equal to adaptive start
		validators.NumericGreaterThanOrEqual(
			path.MatchRoot("stateful_firewall").AtName("adaptive_timeouts").AtName("end"),
			path.MatchRoot("stateful_firewall").AtName("adaptive_timeouts").AtName("start"),
		),
		// Ensure adaptive end is greater than or equal to max states
		validators.NumericGreaterThanOrEqual(
			path.MatchRoot("stateful_firewall").AtName("adaptive_timeouts").AtName("end"),
			path.MatchRoot("stateful_firewall").AtName("max").AtName("states"),
		),
		// Ensure max.states is only set for TCP protocols
		validators.RequiresStringEqualsOneOf(
			path.MatchRoot("stateful_firewall").AtName("max").AtName("states"),
			path.MatchRoot("filter").AtName("protocol"),
			[]string{"TCP", "TCP/UDP"},
		),
		// Ensure icmp_type is only set for ICMP protocol
		validators.RequiresStringEqualsOneOf(
			path.MatchRoot("filter").AtName("icmp_type"),
			path.MatchRoot("filter").AtName("protocol"),
			[]string{"ICMP"},
		),
	}
}

func (r *filterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *opnsense.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = opnsense.NewClient(apiClient)
}

func (r *filterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *filterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	resourceStruct, err := convertFilterSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse firwall filter, got error: %s", err))
		return
	}

	// Add firewall filter to unbound
	id, err := r.client.Firewall().AddFilter(ctx, resourceStruct)
	if err != nil {
		if id != "" {
			// Tag new resource with ID from OPNsense
			data.Id = types.StringValue(id)

			// Save data into Terraform state
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create firewall filter, got error: %s", err))
		return
	}

	// Tag new resource with ID from OPNsense
	data.Id = types.StringValue(id)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *filterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *filterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get firewall filter from OPNsense unbound API
	resourceStruct, err := r.client.Firewall().GetFilter(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, fmt.Sprintf("firewall filter not present in remote, removing from state"))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read firewall filter, got error: %s", err))
		return
	}

	// Convert OPNsense struct to TF schema
	resourceModel, err := convertFilterStructToSchema(resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read firewall filter, got error: %s", err))
		return
	}

	// ID cannot be added by convert... func, have to add here
	resourceModel.Id = data.Id

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}

func (r *filterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *filterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert TF schema OPNsense struct
	resourceStruct, err := convertFilterSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse firewall filter, got error: %s", err))
		return
	}

	// Update firewall filter in unbound
	err = r.client.Firewall().UpdateFilter(ctx, data.Id.ValueString(), resourceStruct)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create firewall filter, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *filterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *filterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Firewall().DeleteFilter(ctx, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete firewall filter, got error: %s", err))
		return
	}
}

func (r *filterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *filterResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	schemaV0 := filterResourceSchemaV0()
	return map[int64]resource.StateUpgrader{
		// Upgrade from version 0 (old flat schema) to version 1 (nested schema)
		0: {
			PriorSchema:   &schemaV0,
			StateUpgrader: upgradeFilterStateV0toV1,
		},
	}
}

// upgradeFilterStateV0toV1 migrates state from schema version 0 to version 1.
// Schema v0 had a flat structure with top-level attributes.
// Schema v1 reorganizes into nested blocks: interface, filter, stateful_firewall, etc.
func upgradeFilterStateV0toV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	tflog.Info(ctx, "Upgrading filter resource state from v0 to v1")

	// Parse the old state from RawState JSON
	var oldState filterResourceModelV0

	// Read old state
	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to read old state during upgrade")
		return
	}

	// Extract source and destination objects
	var oldSource *firewallLocation
	var oldDest *firewallLocation

	if !oldState.Source.IsNull() {
		oldSource = &firewallLocation{}
		resp.Diagnostics.Append(oldState.Source.As(ctx, oldSource, basetypes.ObjectAsOptions{})...)
	}

	if !oldState.Destination.IsNull() {
		oldDest = &firewallLocation{}
		resp.Diagnostics.Append(oldState.Destination.As(ctx, oldDest, basetypes.ObjectAsOptions{})...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Build new v1 state structure
	newState := &filterResourceModel{
		// Top-level fields (preserved)
		Enabled:     oldState.Enabled,
		Sequence:    oldState.Sequence,
		Description: oldState.Description,
		Id:          oldState.Id,

		// New top-level fields (defaults)
		NoXMLRPCSync: types.BoolValue(false),
		Categories:   types.SetValueMust(types.StringType, []attr.Value{}),

		// Interface block (transform Set → nested with invert)
		Interface: &filterInterfaceBlock{
			Invert:    types.BoolValue(false),
			Interface: oldState.Interface,
		},

		// Filter block (consolidate many old top-level fields)
		Filter: &filterFilterBlock{
			Quick:         oldState.Quick,
			Action:        oldState.Action,
			AllowOptions:  types.BoolValue(false), // New field, default
			Direction:     oldState.Direction,
			IPProtocol:    oldState.IPProtocol,
			Protocol:      oldState.Protocol,
			ICMPType:      types.SetValueMust(types.StringType, []attr.Value{}), // New field, default
			Source:        oldSource,
			Destination:   oldDest,
			Log:           oldState.Log,
			TCPFlags:      types.SetValueMust(types.StringType, []attr.Value{}), // New field, default
			TCPFlagsOutOf: types.SetValueMust(types.StringType, []attr.Value{}), // New field, default
			Schedule:      types.StringValue(""),                                // New field, default
		},

		// Stateful Firewall block (new, all defaults)
		StatefulFirewall: &filterStatefulFirewallBlock{
			Type:    types.StringValue("keep"),
			Policy:  types.StringValue(""),
			Timeout: types.Int64Value(-1),
			AdaptiveTimeouts: &filterAdaptiveTimeouts{
				Start: types.Int64Value(-1),
				End:   types.Int64Value(-1),
			},
			Max: &filterMax{
				States:            types.Int64Value(-1),
				SourceNodes:       types.Int64Value(-1),
				SourceStates:      types.Int64Value(-1),
				SourceConnections: types.Int64Value(-1),
				NewConnections: &filterNewConnections{
					Count:   types.Int64Value(-1),
					Seconds: types.Int64Value(-1),
				},
			},
			OverloadTable: types.StringValue(""),
			NoPfsync:      types.BoolValue(false),
		},

		// Traffic Shaping block (new, all defaults)
		TrafficShaping: &filterTrafficShapingBlock{
			Shaper:        types.StringValue(""),
			ReverseShaper: types.StringValue(""),
		},

		// Source Routing block (gateway moved here)
		SourceRouting: &filterSourceRoutingBlock{
			Gateway:        oldState.Gateway,
			DisableReplyTo: types.BoolValue(false), // New field, default
			ReplyTo:        types.StringValue(""),  // New field, default
		},

		// Priority block (new, all defaults)
		Priority: &filterPriorityBlock{
			Match:       types.Int64Value(-1),
			Set:         types.Int64Value(-1),
			LowDelaySet: types.Int64Value(-1),
			MatchTOS:    types.StringValue(""),
		},

		// Internal Tagging block (new, all defaults)
		InternalTagging: &filterInternalTaggingBlock{
			SetLocal:   types.StringValue(""),
			MatchLocal: types.StringValue(""),
		},
	}

	// Set the upgraded state
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)

	if !resp.Diagnostics.HasError() {
		tflog.Info(ctx, "Successfully upgraded filter resource state from v0 to v1", map[string]any{
			"id": oldState.Id.ValueString(),
		})
	}
}
