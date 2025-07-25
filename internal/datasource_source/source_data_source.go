package datasource_source

import (
	"context"
	"time"

	"github.com/aws-gopher/unstructured-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SourceToModel(ctx context.Context, source *unstructured.Source, diagnostics diag.Diagnostics) *SourceModel {
	if source == nil {
		return nil
	}

	model := &SourceModel{
		Id:        types.StringValue(source.ID),
		Name:      types.StringValue(source.Name),
		CreatedAt: types.StringValue(source.CreatedAt.Format(time.RFC3339)),
		UpdatedAt: types.StringValue(source.UpdatedAt.Format(time.RFC3339)),
	}

	// Set the appropriate nested config block based on the source type
	switch source.Type {
	case unstructured.ConnectorTypeAzure:
		if config, ok := source.Config.(*unstructured.AzureSourceConnectorConfig); ok {
			model.Azure = AzureValue{
				RemoteUrl: types.StringValue(config.RemoteURL),
				Recursive: types.BoolValue(config.Recursive),
			}
			if config.AccountName != nil {
				model.Azure.AccountName = types.StringValue(*config.AccountName)
			}
			if config.AccountKey != nil {
				model.Azure.AccountKey = types.StringValue(*config.AccountKey)
			}
			if config.ConnectionString != nil {
				model.Azure.ConnectionString = types.StringValue(*config.ConnectionString)
			}
			if config.SASToken != nil {
				model.Azure.SasToken = types.StringValue(*config.SASToken)
			}
		}

	case unstructured.ConnectorTypeBox:
		if config, ok := source.Config.(*unstructured.BoxSourceConnectorConfig); ok {
			model.Box = BoxValue{
				BoxAppConfig: types.StringValue(config.BoxAppConfig),
				Recursive:    types.BoolValue(config.Recursive),
			}
		}

	case unstructured.ConnectorTypeConfluence:
		if config, ok := source.Config.(*unstructured.ConfluenceSourceConnectorConfig); ok {
			// Convert []string to types.List
			spacesList, diags := types.ListValueFrom(ctx, types.StringType, config.Spaces)
			if diags.HasError() {
				diagnostics.Append(diags...)
			}
			model.Confluence = ConfluenceValue{
				Url:                       types.StringValue(config.URL),
				Username:                  types.StringValue(config.Username),
				Cloud:                     types.BoolValue(config.Cloud),
				MaxNumOfSpaces:            types.Int64Value(int64(config.MaxNumOfSpaces)),
				MaxNumOfDocsFromEachSpace: types.Int64Value(int64(config.MaxNumOfDocsFromEachSpace)),
				Spaces:                    spacesList,
			}
			if config.Password != nil {
				model.Confluence.Password = types.StringValue(*config.Password)
			}
			if config.APIToken != nil {
				model.Confluence.ApiToken = types.StringValue(*config.APIToken)
			}
			if config.Token != nil {
				model.Confluence.Token = types.StringValue(*config.Token)
			}
			if config.ExtractImages != nil {
				model.Confluence.ExtractImages = types.BoolValue(*config.ExtractImages)
			}
			if config.ExtractFiles != nil {
				model.Confluence.ExtractFiles = types.BoolValue(*config.ExtractFiles)
			}
		}

	case unstructured.ConnectorTypeCouchbase:
		if config, ok := source.Config.(*unstructured.CouchbaseSourceConnectorConfig); ok {
			model.Couchbase = CouchbaseValue{
				Bucket:           types.StringValue(config.Bucket),
				ConnectionString: types.StringValue(config.ConnectionString),
				BatchSize:        types.Int64Value(int64(config.BatchSize)),
				Username:         types.StringValue(config.Username),
				Password:         types.StringValue(config.Password),
				CollectionId:     types.StringValue(config.CollectionID),
			}
			if config.Scope != nil {
				model.Couchbase.Scope = types.StringValue(*config.Scope)
			}
			if config.Collection != nil {
				model.Couchbase.Collection = types.StringValue(*config.Collection)
			}
		}

	case unstructured.ConnectorTypeDatabricksVolumes:
		if config, ok := source.Config.(*unstructured.DatabricksVolumesConnectorConfig); ok {
			model.DatabricksVolumes = DatabricksVolumesValue{
				Host:         types.StringValue(config.Host),
				Catalog:      types.StringValue(config.Catalog),
				Volume:       types.StringValue(config.Volume),
				VolumePath:   types.StringValue(config.VolumePath),
				ClientSecret: types.StringValue(config.ClientSecret),
				ClientId:     types.StringValue(config.ClientID),
			}
			if config.Schema != nil {
				model.DatabricksVolumes.Schema = types.StringValue(*config.Schema)
			}
		}

	case unstructured.ConnectorTypeDropbox:
		if config, ok := source.Config.(*unstructured.DropboxSourceConnectorConfig); ok {
			model.Dropbox = DropboxValue{
				Token:     types.StringValue(config.Token),
				RemoteUrl: types.StringValue(config.RemoteURL),
				Recursive: types.BoolValue(config.Recursive),
			}
		}

	case unstructured.ConnectorTypeElasticsearch:
		if config, ok := source.Config.(*unstructured.ElasticsearchConnectorConfig); ok {
			// Convert []string to types.List
			hostsList, diags := types.ListValueFrom(ctx, types.StringType, config.Hosts)
			if diags.HasError() {
				diagnostics.Append(diags...)
			}
			model.Elasticsearch = ElasticsearchValue{
				Hosts:     hostsList,
				IndexName: types.StringValue(config.IndexName),
				EsApiKey:  types.StringValue(config.ESAPIKey),
			}
		}

	case unstructured.ConnectorTypeGCS:
		if config, ok := source.Config.(*unstructured.GCSSourceConnectorConfig); ok {
			model.Gcs = GcsValue{
				RemoteUrl:         types.StringValue(config.RemoteURL),
				ServiceAccountKey: types.StringValue(config.ServiceAccountKey),
				Recursive:         types.BoolValue(config.Recursive),
			}
		}

	case unstructured.ConnectorTypeGoogleDrive:
		if config, ok := source.Config.(*unstructured.GoogleDriveSourceConnectorConfig); ok {
			// Convert []string to types.List
			extensionsList, diags := types.ListValueFrom(ctx, types.StringType, config.Extensions)
			if diags.HasError() {
				diagnostics.Append(diags...)
			}
			model.GoogleDrive = GoogleDriveValue{
				DriveId:           types.StringValue(config.DriveID),
				ServiceAccountKey: types.StringValue(config.ServiceAccountKey),
				Recursive:         types.BoolValue(config.Recursive),
				Extensions:        extensionsList,
			}
		}

	case unstructured.ConnectorTypeJira:
		if config, ok := source.Config.(*unstructured.JiraSourceConnectorConfig); ok {
			// Convert []string to types.List
			projectsList, diags := types.ListValueFrom(ctx, types.StringType, config.Projects)
			if diags.HasError() {
				diagnostics.Append(diags...)
			}
			boardsList, diags := types.ListValueFrom(ctx, types.StringType, config.Boards)
			if diags.HasError() {
				diagnostics.Append(diags...)
			}
			issuesList, diags := types.ListValueFrom(ctx, types.StringType, config.Issues)
			if diags.HasError() {
				diagnostics.Append(diags...)
			}
			statusFiltersList, diags := types.ListValueFrom(ctx, types.StringType, config.StatusFilters)
			if diags.HasError() {
				diagnostics.Append(diags...)
			}
			model.Jira = JiraValue{
				Url:           types.StringValue(config.URL),
				Username:      types.StringValue(config.Username),
				Projects:      projectsList,
				Boards:        boardsList,
				Issues:        issuesList,
				StatusFilters: statusFiltersList,
			}
			if config.Password != nil {
				model.Jira.Password = types.StringValue(*config.Password)
			}
			if config.Token != nil {
				model.Jira.Token = types.StringValue(*config.Token)
			}
			if config.Cloud != nil {
				model.Jira.Cloud = types.BoolValue(*config.Cloud)
			}
			if config.DownloadAttachments != nil {
				model.Jira.DownloadAttachments = types.BoolValue(*config.DownloadAttachments)
			}
		}

	case unstructured.ConnectorTypeKafkaCloud:
		if config, ok := source.Config.(*unstructured.KafkaCloudSourceConnectorConfig); ok {
			model.KafkaCloud = KafkaCloudValue{
				BootstrapServers:     types.StringValue(config.BootstrapServers),
				Port:                 types.Int64Value(int64(config.Port)),
				Topic:                types.StringValue(config.Topic),
				KafkaApiKey:          types.StringValue(config.KafkaAPIKey),
				Secret:               types.StringValue(config.Secret),
				NumMessagesToConsume: types.Int64Value(int64(config.NumMessagesToConsume)),
			}
			if config.GroupID != nil {
				model.KafkaCloud.GroupId = types.StringValue(*config.GroupID)
			}
		}

	case unstructured.ConnectorTypeMongoDB:
		if config, ok := source.Config.(*unstructured.MongoDBConnectorConfig); ok {
			model.Mongodb = MongodbValue{
				Database:   types.StringValue(config.Database),
				Collection: types.StringValue(config.Collection),
				Uri:        types.StringValue(config.URI),
			}
		}

	case unstructured.ConnectorTypeOneDrive:
		if config, ok := source.Config.(*unstructured.OneDriveSourceConnectorConfig); ok {
			model.Onedrive = OnedriveValue{
				ClientId:     types.StringValue(config.ClientID),
				UserPname:    types.StringValue(config.UserPName),
				Tenant:       types.StringValue(config.Tenant),
				AuthorityUrl: types.StringValue(config.AuthorityURL),
				ClientCred:   types.StringValue(config.ClientCred),
				Recursive:    types.BoolValue(config.Recursive),
				Path:         types.StringValue(config.Path),
			}
		}

	case unstructured.ConnectorTypeOutlook:
		if config, ok := source.Config.(*unstructured.OutlookSourceConnectorConfig); ok {
			// Convert []string to types.List
			outlookFoldersList, diags := types.ListValueFrom(ctx, types.StringType, config.OutlookFolders)
			if diags.HasError() {
				diagnostics.Append(diags...)
			}
			model.Outlook = OutlookValue{
				ClientId:       types.StringValue(config.ClientID),
				ClientCred:     types.StringValue(config.ClientCred),
				Recursive:      types.BoolValue(config.Recursive),
				UserEmail:      types.StringValue(config.UserEmail),
				OutlookFolders: outlookFoldersList,
			}
			if config.AuthorityURL != nil {
				model.Outlook.AuthorityUrl = types.StringValue(*config.AuthorityURL)
			}
			if config.Tenant != nil {
				model.Outlook.Tenant = types.StringValue(*config.Tenant)
			}
		}

	case unstructured.ConnectorTypePostgres:
		if config, ok := source.Config.(*unstructured.PostgresSourceConnectorConfig); ok {
			// Convert []string to types.List
			fieldsList, diags := types.ListValueFrom(ctx, types.StringType, config.Fields)
			if diags.HasError() {
				diagnostics.Append(diags...)
			}
			model.Postgres = PostgresValue{
				Host:      types.StringValue(config.Host),
				Database:  types.StringValue(config.Database),
				Port:      types.Int64Value(int64(config.Port)),
				Username:  types.StringValue(config.Username),
				Password:  types.StringValue(config.Password),
				TableName: types.StringValue(config.TableName),
				BatchSize: types.Int64Value(int64(config.BatchSize)),
				IdColumn:  types.StringValue(config.IDColumn),
				Fields:    fieldsList,
			}
		}

	case unstructured.ConnectorTypeS3:
		if config, ok := source.Config.(*unstructured.S3SourceConnectorConfig); ok {
			model.S3 = S3Value{
				RemoteUrl: types.StringValue(config.RemoteURL),
				Anonymous: types.BoolValue(config.Anonymous),
				Recursive: types.BoolValue(config.Recursive),
			}
			if config.Key != nil {
				model.S3.Key = types.StringValue(*config.Key)
			}
			if config.Secret != nil {
				model.S3.Secret = types.StringValue(*config.Secret)
			}
			if config.Token != nil {
				model.S3.Token = types.StringValue(*config.Token)
			}
			if config.EndpointURL != nil {
				model.S3.EndpointUrl = types.StringValue(*config.EndpointURL)
			}
		}

	case unstructured.ConnectorTypeSalesforce:
		if config, ok := source.Config.(*unstructured.SalesforceSourceConnectorConfig); ok {
			// Convert []string to types.List
			categoriesList, diags := types.ListValueFrom(ctx, types.StringType, config.Categories)
			if diags.HasError() {
				diagnostics.Append(diags...)
			}
			model.Salesforce = SalesforceValue{
				Username:    types.StringValue(config.Username),
				ConsumerKey: types.StringValue(config.ConsumerKey),
				PrivateKey:  types.StringValue(config.PrivateKey),
				Categories:  categoriesList,
			}
		}

	case unstructured.ConnectorTypeSharePoint:
		if config, ok := source.Config.(*unstructured.SharePointSourceConnectorConfig); ok {
			model.Sharepoint = SharepointValue{
				Site:       types.StringValue(config.Site),
				Tenant:     types.StringValue(config.Tenant),
				UserPname:  types.StringValue(config.UserPName),
				ClientId:   types.StringValue(config.ClientID),
				ClientCred: types.StringValue(config.ClientCred),
				Recursive:  types.BoolValue(config.Recursive),
			}
			if config.AuthorityURL != nil {
				model.Sharepoint.AuthorityUrl = types.StringValue(*config.AuthorityURL)
			}
			if config.Path != nil {
				model.Sharepoint.Path = types.StringValue(*config.Path)
			}
		}

	case unstructured.ConnectorTypeSnowflake:
		if config, ok := source.Config.(*unstructured.SnowflakeSourceConnectorConfig); ok {
			model.Snowflake = SnowflakeValue{
				Account:  types.StringValue(config.Account),
				Role:     types.StringValue(config.Role),
				User:     types.StringValue(config.User),
				Password: types.StringValue(config.Password),
				Host:     types.StringValue(config.Host),
				Database: types.StringValue(config.Database),
			}
			if config.Port != nil {
				model.Snowflake.Port = types.Int64Value(int64(*config.Port))
			}
			if config.Schema != nil {
				model.Snowflake.Schema = types.StringValue(*config.Schema)
			}
			if config.TableName != nil {
				model.Snowflake.TableName = types.StringValue(*config.TableName)
			}
			if config.BatchSize != nil {
				model.Snowflake.BatchSize = types.Int64Value(int64(*config.BatchSize))
			}
			if config.IDColumn != nil {
				model.Snowflake.RecordIdKey = types.StringValue(*config.IDColumn)
			}
		}

	case unstructured.ConnectorTypeZendesk:
		if config, ok := source.Config.(*unstructured.ZendeskSourceConnectorConfig); ok {
			model.Zendesk = ZendeskValue{
				Subdomain: types.StringValue(config.Subdomain),
				Email:     types.StringValue(config.Email),
				ApiToken:  types.StringValue(config.APIToken),
			}
			if config.ItemType != nil {
				model.Zendesk.ItemType = types.StringValue(*config.ItemType)
			}
			if config.BatchSize != nil {
				model.Zendesk.BatchSize = types.Int64Value(int64(*config.BatchSize))
			}
		}

	default:
		diagnostics.AddError(
			"Unsupported source type",
			"Source type '"+source.Type+"' is not supported",
		)
	}

	return model
}
