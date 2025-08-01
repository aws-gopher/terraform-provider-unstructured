---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "unstructured_destination Resource - unstructured"
subcategory: ""
description: |-
  
---

# unstructured_destination (Resource)



## Example Usage

```terraform
resource "unstructured_destination" "example" {
  name = "example_s3_destination"

  s3 = {
    remote_url = "s3://my-destination-bucket/"
    key        = "aws-access-key-id"
    secret     = "aws-secret-access-key"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)

### Optional

- `astradb` (Attributes) (see [below for nested schema](#nestedatt--astradb))
- `azure_ai_search` (Attributes) (see [below for nested schema](#nestedatt--azure_ai_search))
- `couchbase` (Attributes) (see [below for nested schema](#nestedatt--couchbase))
- `databricks_volume_delta_tables` (Attributes) (see [below for nested schema](#nestedatt--databricks_volume_delta_tables))
- `databricks_volumes` (Attributes) (see [below for nested schema](#nestedatt--databricks_volumes))
- `delta_table` (Attributes) (see [below for nested schema](#nestedatt--delta_table))
- `elasticsearch` (Attributes) (see [below for nested schema](#nestedatt--elasticsearch))
- `gcs` (Attributes) (see [below for nested schema](#nestedatt--gcs))
- `ibm_watsonx_s3` (Attributes) (see [below for nested schema](#nestedatt--ibm_watsonx_s3))
- `kafka_cloud` (Attributes) (see [below for nested schema](#nestedatt--kafka_cloud))
- `milvus` (Attributes) (see [below for nested schema](#nestedatt--milvus))
- `mongodb` (Attributes) (see [below for nested schema](#nestedatt--mongodb))
- `motherduck` (Attributes) (see [below for nested schema](#nestedatt--motherduck))
- `neo4j` (Attributes) (see [below for nested schema](#nestedatt--neo4j))
- `onedrive` (Attributes) (see [below for nested schema](#nestedatt--onedrive))
- `pinecone` (Attributes) (see [below for nested schema](#nestedatt--pinecone))
- `postgres` (Attributes) (see [below for nested schema](#nestedatt--postgres))
- `qdrant_cloud` (Attributes) (see [below for nested schema](#nestedatt--qdrant_cloud))
- `redis` (Attributes) (see [below for nested schema](#nestedatt--redis))
- `s3` (Attributes) (see [below for nested schema](#nestedatt--s3))
- `snowflake` (Attributes) (see [below for nested schema](#nestedatt--snowflake))
- `weaviate_cloud` (Attributes) (see [below for nested schema](#nestedatt--weaviate_cloud))

### Read-Only

- `created_at` (String)
- `id` (String) The ID of this resource.
- `updated_at` (String)

<a id="nestedatt--astradb"></a>
### Nested Schema for `astradb`

Required:

- `api_endpoint` (String)
- `batch_size` (Number)
- `collection_name` (String)
- `token` (String)

Optional:

- `keyspace` (String)


<a id="nestedatt--azure_ai_search"></a>
### Nested Schema for `azure_ai_search`

Required:

- `endpoint` (String)
- `index` (String)
- `key` (String)


<a id="nestedatt--couchbase"></a>
### Nested Schema for `couchbase`

Required:

- `batch_size` (Number)
- `bucket` (String)
- `collection_id` (String)
- `connection_string` (String)
- `password` (String)
- `username` (String)

Optional:

- `collection` (String)
- `scope` (String)


<a id="nestedatt--databricks_volume_delta_tables"></a>
### Nested Schema for `databricks_volume_delta_tables`

Required:

- `catalog` (String)
- `http_path` (String)
- `schema` (String)
- `server_hostname` (String)
- `volume` (String)

Optional:

- `client_id` (String)
- `client_secret` (String)
- `database` (String)
- `table_name` (String)
- `token` (String)
- `volume_path` (String)


<a id="nestedatt--databricks_volumes"></a>
### Nested Schema for `databricks_volumes`

Required:

- `catalog` (String)
- `client_id` (String)
- `client_secret` (String)
- `host` (String)
- `volume` (String)
- `volume_path` (String)

Optional:

- `schema` (String)


<a id="nestedatt--delta_table"></a>
### Nested Schema for `delta_table`

Required:

- `aws_access_key_id` (String)
- `aws_region` (String)
- `aws_secret_access_key` (String)
- `table_uri` (String)


<a id="nestedatt--elasticsearch"></a>
### Nested Schema for `elasticsearch`

Required:

- `es_api_key` (String)
- `hosts` (List of String)
- `index_name` (String)


<a id="nestedatt--gcs"></a>
### Nested Schema for `gcs`

Required:

- `remote_url` (String)
- `service_account_key` (String)


<a id="nestedatt--ibm_watsonx_s3"></a>
### Nested Schema for `ibm_watsonx_s3`

Required:

- `access_key_id` (String)
- `catalog` (String)
- `iam_api_key` (String)
- `iceberg_endpoint` (String)
- `namespace` (String)
- `object_storage_endpoint` (String)
- `object_storage_region` (String)
- `secret_access_key` (String)
- `table` (String)

Optional:

- `max_retries` (Number)
- `max_retries_connection` (Number)
- `record_id_key` (String)


<a id="nestedatt--kafka_cloud"></a>
### Nested Schema for `kafka_cloud`

Required:

- `bootstrap_servers` (String)
- `kafka_api_key` (String)
- `secret` (String)
- `topic` (String)

Optional:

- `batch_size` (Number)
- `group_id` (String)
- `port` (Number)


<a id="nestedatt--milvus"></a>
### Nested Schema for `milvus`

Required:

- `collection_name` (String)
- `record_id_key` (String)
- `uri` (String)

Optional:

- `db_name` (String)
- `password` (String)
- `token` (String)
- `user` (String)


<a id="nestedatt--mongodb"></a>
### Nested Schema for `mongodb`

Required:

- `collection` (String)
- `database` (String)
- `uri` (String)


<a id="nestedatt--motherduck"></a>
### Nested Schema for `motherduck`

Required:

- `account` (String)
- `database` (String)
- `host` (String)
- `password` (String)
- `role` (String)
- `user` (String)

Optional:

- `batch_size` (Number)
- `port` (Number)
- `record_id_key` (String)
- `schema` (String)
- `table_name` (String)


<a id="nestedatt--neo4j"></a>
### Nested Schema for `neo4j`

Required:

- `database` (String)
- `password` (String)
- `uri` (String)
- `username` (String)

Optional:

- `batch_size` (Number)


<a id="nestedatt--onedrive"></a>
### Nested Schema for `onedrive`

Required:

- `authority_url` (String)
- `client_cred` (String)
- `client_id` (String)
- `path` (String)
- `tenant` (String)
- `user_pname` (String)

Optional:

- `recursive` (Boolean)


<a id="nestedatt--pinecone"></a>
### Nested Schema for `pinecone`

Required:

- `api_key` (String)
- `index_name` (String)
- `namespace` (String)

Optional:

- `batch_size` (Number)


<a id="nestedatt--postgres"></a>
### Nested Schema for `postgres`

Required:

- `batch_size` (Number)
- `database` (String)
- `host` (String)
- `password` (String)
- `port` (Number)
- `table_name` (String)
- `username` (String)


<a id="nestedatt--qdrant_cloud"></a>
### Nested Schema for `qdrant_cloud`

Required:

- `api_key` (String)
- `collection_name` (String)
- `url` (String)

Optional:

- `batch_size` (Number)


<a id="nestedatt--redis"></a>
### Nested Schema for `redis`

Required:

- `host` (String)

Optional:

- `batch_size` (Number)
- `database` (Number)
- `password` (String)
- `port` (Number)
- `ssl` (Boolean)
- `uri` (String)
- `username` (String)


<a id="nestedatt--s3"></a>
### Nested Schema for `s3`

Required:

- `remote_url` (String)

Optional:

- `anonymous` (Boolean)
- `endpoint_url` (String)
- `key` (String)
- `recursive` (Boolean)
- `secret` (String)
- `token` (String)


<a id="nestedatt--snowflake"></a>
### Nested Schema for `snowflake`

Required:

- `account` (String)
- `database` (String)
- `host` (String)
- `password` (String)
- `role` (String)
- `table_name` (String)
- `user` (String)

Optional:

- `batch_size` (Number)
- `port` (Number)
- `record_id_key` (String)
- `schema` (String)


<a id="nestedatt--weaviate_cloud"></a>
### Nested Schema for `weaviate_cloud`

Required:

- `api_key` (String)
- `cluster_url` (String)

Optional:

- `collection` (String)

## Import

Import is supported using the following syntax:

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used, for example:

```shell
terraform import unstructured_destination.example "bc15ad07-a46e-4cc4-893a-8173047b7165"
```
