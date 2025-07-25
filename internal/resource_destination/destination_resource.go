package resource_destination

import (
	"context"
	"time"

	"github.com/aws-gopher/unstructured-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func DestinationToModel(ctx context.Context, destination *unstructured.Destination, diagnostics diag.Diagnostics) *DestinationModel {
	if destination == nil {
		return nil
	}

	model := &DestinationModel{
		Id:        types.StringValue(destination.ID),
		Name:      types.StringValue(destination.Name),
		CreatedAt: types.StringValue(destination.CreatedAt.Format(time.RFC3339)),
		UpdatedAt: types.StringValue(destination.UpdatedAt.Format(time.RFC3339)),
	}

	// Set the appropriate nested config block based on the destination type
	switch destination.Type {
	case unstructured.ConnectorTypeAstraDB:
		if config, ok := destination.Config.(*unstructured.AstraDBConnectorConfig); ok {
			model.Astradb = AstradbValue{
				ApiEndpoint:    types.StringValue(config.APIEndpoint),
				BatchSize:      types.Int64Value(int64(config.BatchSize)),
				CollectionName: types.StringValue(config.CollectionName),
				Token:          types.StringValue(config.Token),
			}
			if config.Keyspace != nil {
				model.Astradb.Keyspace = types.StringValue(*config.Keyspace)
			}
		}

	case unstructured.ConnectorTypeAzureAISearch:
		if config, ok := destination.Config.(*unstructured.AzureAISearchConnectorConfig); ok {
			model.AzureAiSearch = AzureAiSearchValue{
				Endpoint: types.StringValue(config.Endpoint),
				Index:    types.StringValue(config.Index),
				Key:      types.StringValue(config.Key),
			}
		}

	case unstructured.ConnectorTypeCouchbase:
		if config, ok := destination.Config.(*unstructured.CouchbaseDestinationConnectorConfig); ok {
			model.Couchbase = CouchbaseValue{
				BatchSize:        types.Int64Value(int64(config.BatchSize)),
				Bucket:           types.StringValue(config.Bucket),
				ConnectionString: types.StringValue(config.ConnectionString),
				Password:         types.StringValue(config.Password),
				Username:         types.StringValue(config.Username),
			}
			if config.Scope != nil {
				model.Couchbase.Scope = types.StringValue(*config.Scope)
			}
			if config.Collection != nil {
				model.Couchbase.Collection = types.StringValue(*config.Collection)
				model.Couchbase.CollectionId = types.StringValue(*config.Collection)
			}
		}

	case unstructured.ConnectorTypeDatabricksVolumeDeltaTable:
		if config, ok := destination.Config.(*unstructured.DatabricksVDTDestinationConnectorConfig); ok {
			model.DatabricksVolumeDeltaTables = DatabricksVolumeDeltaTablesValue{
				Catalog:        types.StringValue(config.Catalog),
				ServerHostname: types.StringValue(config.ServerHostname),
				HttpPath:       types.StringValue(config.HTTPPath),
				Volume:         types.StringValue(config.Volume),
			}
			if config.Token != nil {
				model.DatabricksVolumeDeltaTables.Token = types.StringValue(*config.Token)
			}
			if config.ClientID != nil {
				model.DatabricksVolumeDeltaTables.ClientId = types.StringValue(*config.ClientID)
			}
			if config.ClientSecret != nil {
				model.DatabricksVolumeDeltaTables.ClientSecret = types.StringValue(*config.ClientSecret)
			}
			if config.Database != nil {
				model.DatabricksVolumeDeltaTables.Database = types.StringValue(*config.Database)
			}
			if config.TableName != nil {
				model.DatabricksVolumeDeltaTables.TableName = types.StringValue(*config.TableName)
			}
			if config.Schema != nil {
				model.DatabricksVolumeDeltaTables.Schema = types.StringValue(*config.Schema)
			}
			if config.VolumePath != nil {
				model.DatabricksVolumeDeltaTables.VolumePath = types.StringValue(*config.VolumePath)
			}
		}

	case unstructured.ConnectorTypeDatabricksVolumes:
		if config, ok := destination.Config.(*unstructured.DatabricksVolumesConnectorConfig); ok {
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

	case unstructured.ConnectorTypeDeltaTable:
		if config, ok := destination.Config.(*unstructured.DeltaTableConnectorConfig); ok {
			model.DeltaTable = DeltaTableValue{
				AwsAccessKeyId:     types.StringValue(config.AwsAccessKeyID),
				AwsSecretAccessKey: types.StringValue(config.AwsSecretAccessKey),
				AwsRegion:          types.StringValue(config.AwsRegion),
				TableUri:           types.StringValue(config.TableURI),
			}
		}

	case unstructured.ConnectorTypeElasticsearch:
		if config, ok := destination.Config.(*unstructured.ElasticsearchConnectorConfig); ok {
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
		if config, ok := destination.Config.(*unstructured.GCSDestinationConnectorConfig); ok {
			model.Gcs = GcsValue{
				RemoteUrl:         types.StringValue(config.RemoteURL),
				ServiceAccountKey: types.StringValue(config.ServiceAccountKey),
			}
		}

	case unstructured.ConnectorTypeIBMWatsonxS3:
		if config, ok := destination.Config.(*unstructured.IBMWatsonxS3DestinationConnectorConfig); ok {
			model.IbmWatsonxS3 = IbmWatsonxS3Value{
				IamApiKey:             types.StringValue(config.IAMApiKey),
				AccessKeyId:           types.StringValue(config.AccessKeyID),
				SecretAccessKey:       types.StringValue(config.SecretAccessKey),
				IcebergEndpoint:       types.StringValue(config.IcebergEndpoint),
				ObjectStorageEndpoint: types.StringValue(config.ObjectStorageEndpoint),
				ObjectStorageRegion:   types.StringValue(config.ObjectStorageRegion),
				Catalog:               types.StringValue(config.Catalog),
				Namespace:             types.StringValue(config.Namespace),
				Table:                 types.StringValue(config.Table),
			}
			if config.MaxRetriesConnection != nil {
				model.IbmWatsonxS3.MaxRetriesConnection = types.Int64Value(int64(*config.MaxRetriesConnection))
			}
			if config.MaxRetries != nil {
				model.IbmWatsonxS3.MaxRetries = types.Int64Value(int64(*config.MaxRetries))
			}
			if config.RecordIDKey != nil {
				model.IbmWatsonxS3.RecordIdKey = types.StringValue(*config.RecordIDKey)
			}
		}

	case unstructured.ConnectorTypeKafkaCloud:
		if config, ok := destination.Config.(*unstructured.KafkaCloudDestinationConnectorConfig); ok {
			model.KafkaCloud = KafkaCloudValue{
				BootstrapServers: types.StringValue(config.BootstrapServers),
				Topic:            types.StringValue(config.Topic),
				KafkaApiKey:      types.StringValue(config.KafkaAPIKey),
				Secret:           types.StringValue(config.Secret),
			}
			if config.Port != nil {
				model.KafkaCloud.Port = types.Int64Value(int64(*config.Port))
			}
			if config.GroupID != nil {
				model.KafkaCloud.GroupId = types.StringValue(*config.GroupID)
			}
			if config.BatchSize != nil {
				model.KafkaCloud.BatchSize = types.Int64Value(int64(*config.BatchSize))
			}
		}

	case unstructured.ConnectorTypeMilvus:
		if config, ok := destination.Config.(*unstructured.MilvusDestinationConnectorConfig); ok {
			model.Milvus = MilvusValue{
				Uri:            types.StringValue(config.URI),
				CollectionName: types.StringValue(config.CollectionName),
				RecordIdKey:    types.StringValue(config.RecordIDKey),
			}
			if config.User != nil {
				model.Milvus.User = types.StringValue(*config.User)
			}
			if config.Token != nil {
				model.Milvus.Token = types.StringValue(*config.Token)
			}
			if config.Password != nil {
				model.Milvus.Password = types.StringValue(*config.Password)
			}
			if config.DBName != nil {
				model.Milvus.DbName = types.StringValue(*config.DBName)
			}
		}

	case unstructured.ConnectorTypeMongoDB:
		if config, ok := destination.Config.(*unstructured.MongoDBConnectorConfig); ok {
			model.Mongodb = MongodbValue{
				Database:   types.StringValue(config.Database),
				Collection: types.StringValue(config.Collection),
				Uri:        types.StringValue(config.URI),
			}
		}

	case unstructured.ConnectorTypeMotherDuck:
		if config, ok := destination.Config.(*unstructured.MotherduckDestinationConnectorConfig); ok {
			model.Motherduck = MotherduckValue{
				Account:  types.StringValue(config.Account),
				Role:     types.StringValue(config.Role),
				User:     types.StringValue(config.User),
				Password: types.StringValue(config.Password),
				Host:     types.StringValue(config.Host),
				Database: types.StringValue(config.Database),
			}
			if config.Port != nil {
				model.Motherduck.Port = types.Int64Value(int64(*config.Port))
			}
			if config.Schema != nil {
				model.Motherduck.Schema = types.StringValue(*config.Schema)
			}
			if config.TableName != nil {
				model.Motherduck.TableName = types.StringValue(*config.TableName)
			}
			if config.BatchSize != nil {
				model.Motherduck.BatchSize = types.Int64Value(int64(*config.BatchSize))
			}
			if config.RecordIDKey != nil {
				model.Motherduck.RecordIdKey = types.StringValue(*config.RecordIDKey)
			}
		}

	case unstructured.ConnectorTypeNeo4j:
		if config, ok := destination.Config.(*unstructured.Neo4jDestinationConnectorConfig); ok {
			model.Neo4j = Neo4jValue{
				Uri:      types.StringValue(config.URI),
				Database: types.StringValue(config.Database),
				Username: types.StringValue(config.Username),
				Password: types.StringValue(config.Password),
			}
			if config.BatchSize != nil {
				model.Neo4j.BatchSize = types.Int64Value(int64(*config.BatchSize))
			}
		}

	case unstructured.ConnectorTypeOneDrive:
		if config, ok := destination.Config.(*unstructured.OneDriveDestinationConnectorConfig); ok {
			model.Onedrive = OnedriveValue{
				ClientId:     types.StringValue(config.ClientID),
				UserPname:    types.StringValue(config.UserPName),
				Tenant:       types.StringValue(config.Tenant),
				AuthorityUrl: types.StringValue(config.AuthorityURL),
				ClientCred:   types.StringValue(config.ClientCred),
				Path:         types.StringValue(config.RemoteURL), // Note: API uses RemoteURL but TF uses Path
			}
		}

	case unstructured.ConnectorTypePinecone:
		if config, ok := destination.Config.(*unstructured.PineconeDestinationConnectorConfig); ok {
			model.Pinecone = PineconeValue{
				IndexName: types.StringValue(config.IndexName),
				ApiKey:    types.StringValue(config.APIKey),
				Namespace: types.StringValue(config.Namespace),
			}
			if config.BatchSize != nil {
				model.Pinecone.BatchSize = types.Int64Value(int64(*config.BatchSize))
			}
		}

	case unstructured.ConnectorTypePostgres:
		if config, ok := destination.Config.(*unstructured.PostgresDestinationConnectorConfig); ok {
			model.Postgres = PostgresValue{
				Host:      types.StringValue(config.Host),
				Database:  types.StringValue(config.Database),
				Port:      types.Int64Value(int64(config.Port)),
				Username:  types.StringValue(config.Username),
				Password:  types.StringValue(config.Password),
				TableName: types.StringValue(config.TableName),
				BatchSize: types.Int64Value(int64(config.BatchSize)),
			}
		}

	case unstructured.ConnectorTypeQdrantCloud:
		if config, ok := destination.Config.(*unstructured.QdrantCloudDestinationConnectorConfig); ok {
			model.QdrantCloud = QdrantCloudValue{
				Url:            types.StringValue(config.URL),
				ApiKey:         types.StringValue(config.APIKey),
				CollectionName: types.StringValue(config.CollectionName),
			}
			if config.BatchSize != nil {
				model.QdrantCloud.BatchSize = types.Int64Value(int64(*config.BatchSize))
			}
		}

	case unstructured.ConnectorTypeRedis:
		if config, ok := destination.Config.(*unstructured.RedisDestinationConnectorConfig); ok {
			model.Redis = RedisValue{
				Host: types.StringValue(config.Host),
			}
			if config.Port != nil {
				model.Redis.Port = types.Int64Value(int64(*config.Port))
			}
			if config.Username != nil {
				model.Redis.Username = types.StringValue(*config.Username)
			}
			if config.Password != nil {
				model.Redis.Password = types.StringValue(*config.Password)
			}
			if config.URI != nil {
				model.Redis.Uri = types.StringValue(*config.URI)
			}
			if config.Database != nil {
				model.Redis.Database = types.Int64Value(int64(*config.Database))
			}
			if config.SSL != nil {
				model.Redis.Ssl = types.BoolValue(*config.SSL)
			}
			if config.BatchSize != nil {
				model.Redis.BatchSize = types.Int64Value(int64(*config.BatchSize))
			}
		}

	case unstructured.ConnectorTypeS3:
		if config, ok := destination.Config.(*unstructured.S3DestinationConnectorConfig); ok {
			model.S3 = S3Value{
				RemoteUrl: types.StringValue(config.RemoteURL),
				Anonymous: types.BoolValue(config.Anonymous),
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

	case unstructured.ConnectorTypeSnowflake:
		if config, ok := destination.Config.(*unstructured.SnowflakeDestinationConnectorConfig); ok {
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
			if config.RecordIDKey != nil {
				model.Snowflake.RecordIdKey = types.StringValue(*config.RecordIDKey)
			}
		}

	case unstructured.ConnectorTypeWeaviateCloud:
		if config, ok := destination.Config.(*unstructured.WeaviateDestinationConnectorConfig); ok {
			model.WeaviateCloud = WeaviateCloudValue{
				ClusterUrl: types.StringValue(config.ClusterURL),
				ApiKey:     types.StringValue(config.APIKey),
			}
			if config.Collection != nil {
				model.WeaviateCloud.Collection = types.StringValue(*config.Collection)
			}
		}

	default:
		diagnostics.AddError(
			"Unsupported destination type",
			"Destination type '"+destination.Type+"' is not supported",
		)
	}

	return model
}
