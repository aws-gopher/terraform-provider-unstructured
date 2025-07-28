package provider

import (
	"context"
	"fmt"

	"github.com/aws-gopher/terraform-provider-unstructured/internal/resource_source"
	"github.com/aws-gopher/unstructured-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = (*sourceResource)(nil)
var _ resource.ResourceWithConfigure = (*sourceResource)(nil)
var _ resource.ResourceWithImportState = (*sourceResource)(nil)

func NewSourceResource() resource.Resource {
	return &sourceResource{}
}

type sourceResource struct {
	client *unstructured.Client
}

func (r *sourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (r *sourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_source.SourceResourceSchema(ctx)
}

func (r *sourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*unstructured.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *unstructured.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// validateSourceConfig ensures only one nested block is provided.
func (r *sourceResource) validateSourceConfig(data *resource_source.SourceModel) error {
	configCount := 0
	if !data.S3.IsNull() {
		configCount++
	}
	if !data.Postgres.IsNull() {
		configCount++
	}
	if !data.Azure.IsNull() {
		configCount++
	}
	if !data.GoogleDrive.IsNull() {
		configCount++
	}
	if !data.Salesforce.IsNull() {
		configCount++
	}

	if configCount == 0 {
		return fmt.Errorf("exactly one source configuration block must be provided (s3, postgres, azure, google_drive, or salesforce)")
	}
	if configCount > 1 {
		return fmt.Errorf("only one source configuration block can be provided (s3, postgres, azure, google_drive, or salesforce)")
	}
	return nil
}

// getSourceConfig converts the Terraform model to the appropriate API config.
func (r *sourceResource) getSourceConfig(data *resource_source.SourceModel) (unstructured.SourceConfigInput, error) {
	if !data.S3.IsNull() {
		config := &unstructured.S3SourceConnectorConfigInput{
			RemoteURL: data.S3.RemoteUrl.ValueString(),
		}

		if !data.S3.Anonymous.IsNull() && !data.S3.Anonymous.IsUnknown() {
			anonymous := data.S3.Anonymous.ValueBool()
			config.Anonymous = &anonymous
		}
		if !data.S3.Key.IsNull() && !data.S3.Key.IsUnknown() {
			key := data.S3.Key.ValueString()
			config.Key = &key
		}
		if !data.S3.Secret.IsNull() && !data.S3.Secret.IsUnknown() {
			secret := data.S3.Secret.ValueString()
			config.Secret = &secret
		}
		if !data.S3.Token.IsNull() && !data.S3.Token.IsUnknown() {
			token := data.S3.Token.ValueString()
			config.Token = &token
		}
		if !data.S3.EndpointUrl.IsNull() && !data.S3.EndpointUrl.IsUnknown() {
			endpointUrl := data.S3.EndpointUrl.ValueString()
			config.EndpointURL = &endpointUrl
		}
		if !data.S3.Recursive.IsNull() && !data.S3.Recursive.IsUnknown() {
			recursive := data.S3.Recursive.ValueBool()
			config.Recursive = &recursive
		}

		return config, nil
	}

	if !data.Postgres.IsNull() {
		config := &unstructured.PostgresSourceConnectorConfigInput{
			Host:      data.Postgres.Host.ValueString(),
			Database:  data.Postgres.Database.ValueString(),
			Port:      int(data.Postgres.Port.ValueInt64()),
			Username:  data.Postgres.Username.ValueString(),
			Password:  data.Postgres.Password.ValueString(),
			TableName: data.Postgres.TableName.ValueString(),
			BatchSize: int(data.Postgres.BatchSize.ValueInt64()),
		}
		if !data.Postgres.IdColumn.IsNull() && !data.Postgres.IdColumn.IsUnknown() {
			idColumn := data.Postgres.IdColumn.ValueString()
			config.IDColumn = &idColumn
		}

		// Convert fields list
		if !data.Postgres.Fields.IsNull() && !data.Postgres.Fields.IsUnknown() {
			var fields []string
			data.Postgres.Fields.ElementsAs(context.Background(), &fields, false)
			config.Fields = fields
		}

		return config, nil
	}

	if !data.Azure.IsNull() {
		config := &unstructured.AzureSourceConnectorConfigInput{
			RemoteURL: data.Azure.RemoteUrl.ValueString(),
		}
		if !data.Azure.ConnectionString.IsNull() && !data.Azure.ConnectionString.IsUnknown() {
			connectionString := data.Azure.ConnectionString.ValueString()
			config.ConnectionString = &connectionString
		}
		return config, nil
	}

	if !data.GoogleDrive.IsNull() {
		config := &unstructured.GoogleDriveSourceConnectorConfigInput{
			DriveID: data.GoogleDrive.DriveId.ValueString(),
		}
		if !data.GoogleDrive.ServiceAccountKey.IsNull() && !data.GoogleDrive.ServiceAccountKey.IsUnknown() {
			serviceAccountKey := data.GoogleDrive.ServiceAccountKey.ValueString()
			config.ServiceAccountKey = &serviceAccountKey
		}

		if !data.GoogleDrive.Extensions.IsNull() && !data.GoogleDrive.Extensions.IsUnknown() {
			var extensions []string
			data.GoogleDrive.Extensions.ElementsAs(context.Background(), &extensions, false)
			config.Extensions = extensions
		}
		if !data.GoogleDrive.Recursive.IsNull() && !data.GoogleDrive.Recursive.IsUnknown() {
			recursive := data.GoogleDrive.Recursive.ValueBool()
			config.Recursive = &recursive
		}

		return config, nil
	}

	if !data.Salesforce.IsNull() {
		config := &unstructured.SalesforceSourceConnectorConfigInput{
			Username:    data.Salesforce.Username.ValueString(),
			ConsumerKey: data.Salesforce.ConsumerKey.ValueString(),
			PrivateKey:  data.Salesforce.PrivateKey.ValueString(),
		}

		if !data.Salesforce.Categories.IsNull() && !data.Salesforce.Categories.IsUnknown() {
			var categories []string
			data.Salesforce.Categories.ElementsAs(context.Background(), &categories, false)
			config.Categories = categories
		}

		return config, nil
	}

	return nil, fmt.Errorf("no valid source configuration found")
}

func (r *sourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_source.SourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that exactly one source configuration is provided
	if err := r.validateSourceConfig(&data); err != nil {
		resp.Diagnostics.AddError("Invalid source configuration", err.Error())
		return
	}

	// Get the source configuration
	config, err := r.getSourceConfig(&data)
	if err != nil {
		resp.Diagnostics.AddError("Error creating source configuration", err.Error())
		return
	}

	// Create the source
	source, err := r.client.CreateSource(ctx, unstructured.CreateSourceRequest{
		Name:   data.Name.ValueString(),
		Config: config,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating source", err.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, resource_source.SourceToModel(ctx, source, resp.Diagnostics))...)
}

func (r *sourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_source.SourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the source by ID
	source, err := r.client.GetSource(ctx, data.Id.ValueString())
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

func (r *sourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_source.SourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state resource_source.SourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that exactly one source configuration is provided
	if err := r.validateSourceConfig(&data); err != nil {
		resp.Diagnostics.AddError("Invalid source configuration", err.Error())
		return
	}

	// Get the source configuration
	config, err := r.getSourceConfig(&data)
	if err != nil {
		resp.Diagnostics.AddError("Error creating source configuration", err.Error())
		return
	}

	// Update the source
	source, err := r.client.UpdateSource(ctx, unstructured.UpdateSourceRequest{
		Config: config,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating source", err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, resource_source.SourceToModel(ctx, source, resp.Diagnostics))...)
}

func (r *sourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_source.SourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the source
	err := r.client.DeleteSource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting source", err.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *sourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID
	source, err := r.client.GetSource(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing source", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, resource_source.SourceToModel(ctx, source, resp.Diagnostics))...)
}
