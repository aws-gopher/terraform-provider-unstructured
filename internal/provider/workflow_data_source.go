package provider

import (
	"context"
	"fmt"

	"github.com/aws-gopher/terraform-provider-unstructured/internal/datasource_workflow"
	"github.com/aws-gopher/unstructured-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = (*workflowDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*workflowDataSource)(nil)

func NewWorkflowDataSource() datasource.DataSource {
	return &workflowDataSource{}
}

type workflowDataSource struct {
	client *unstructured.Client
}

func (d *workflowDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

func (d *workflowDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_workflow.WorkflowDataSourceSchema(ctx)
}

func (d *workflowDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *workflowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_workflow.WorkflowModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	workflow, err := d.client.GetWorkflow(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting workflow", err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, datasource_workflow.WorkflowToModel(ctx, workflow, resp.Diagnostics))...)
}
