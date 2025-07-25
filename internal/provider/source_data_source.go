package provider

import (
	"context"
	"fmt"

	"github.com/aws-gopher/terraform-provider-unstructured/internal/datasource_source"
	"github.com/aws-gopher/terraform-provider-unstructured/internal/resource_source"
	"github.com/aws-gopher/unstructured-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = (*sourceDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*sourceDataSource)(nil)

func NewSourceDataSource() datasource.DataSource {
	return &sourceDataSource{}
}

type sourceDataSource struct {
	client *unstructured.Client
}

func (d *sourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (d *sourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_source.SourceDataSourceSchema(ctx)
}

func (d *sourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *sourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_source.SourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the source by ID
	source, err := d.client.GetSource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting source", err.Error())
		return
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, resource_source.SourceToModel(ctx, source, resp.Diagnostics))...)
	if resp.Diagnostics.HasError() {
		return
	}
}
