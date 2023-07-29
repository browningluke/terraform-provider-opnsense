package service

import (
	"context"
	"fmt"
	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &UnboundHostOverrideDataSource{}

func NewUnboundHostOverrideDataSource() datasource.DataSource {
	return &UnboundHostOverrideDataSource{}
}

// UnboundHostOverrideDataSource defines the data source implementation.
type UnboundHostOverrideDataSource struct {
	client opnsense.Client
}

func (d *UnboundHostOverrideDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_unbound_host_override"
}

func (d *UnboundHostOverrideDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = UnboundHostOverrideDataSourceSchema()
}

func (d *UnboundHostOverrideDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = opnsense.NewClient(apiClient)
}

func (d *UnboundHostOverrideDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *UnboundHostOverrideResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get resource from OPNsense API
	resource, err := d.client.Unbound().GetHostOverride(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read host_override, got error: %s", err))
		return
	}

	// Convert OPNsense struct to TF schema
	resourceModel, err := convertUnboundHostOverrideStructToSchema(resource)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read host_override, got error: %s", err))
		return
	}

	// ID cannot be added by convert... func, have to add here
	resourceModel.Id = data.Id

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &resourceModel)...)
}
