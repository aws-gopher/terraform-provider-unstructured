package provider

import (
	"context"
	"fmt"

	"github.com/aws-gopher/terraform-provider-unstructured/internal/datasource_destination"
	"github.com/aws-gopher/unstructured-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = (*destinationDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*destinationDataSource)(nil)

func NewDestinationDataSource() datasource.DataSource {
	return &destinationDataSource{}
}

type destinationDataSource struct {
	client *unstructured.Client
}

func (d *destinationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination"
}

func (d *destinationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_destination.DestinationDataSourceSchema(ctx)
}

func (d *destinationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*unstructured.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *unstructured.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *destinationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_destination.DestinationModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the destination by ID
	destination, err := d.client.GetDestination(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting destination", err.Error())
		return
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, datasource_destination.DestinationToModel(ctx, destination, resp.Diagnostics))...)
	if resp.Diagnostics.HasError() {
		return
	}
}
