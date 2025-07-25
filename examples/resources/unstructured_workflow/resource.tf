resource "unstructured_source" "example" {
  name = "example_s3_source"

  s3 {
    remote_url = "s3://my-source-bucket/"
    recursive  = true
    anonymous  = true
  }
}

resource "unstructured_destination" "example" {
  name = "example_s3_destination"

  s3 {
    remote_url = "s3://my-destination-bucket/"
    key        = "aws-access-key-id"
    secret     = "aws-secret-access-key"
  }
}

resource "unstructured_workflow" "example" {
  name          = "example_workflow"
  workflow_type = "basic"

  source_id      = unstructured_source.example.id
  destination_id = unstructured_destination.example.id

  workflow_nodes {
    name    = "Partitioner"
    type    = "partition"
    subtype = "vlm"
    settings = {
      "provider"           = "anthropic"
      "provider_api_key"   = null
      "model"              = "claude-3-5-sonnet-20241022"
      "output_format"      = "text/html"
      "prompt"             = null
      "format_html"        = true
      "unique_element_ids" = true
      "is_dynamic"         = true
      "allow_fast"         = true
    }
  }

  workflow_nodes {
    name    = "Image summarizer"
    type    = "prompter"
    subtype = "openai_image_description"
  }

  workflow_nodes {
    name    = "Table summarizer"
    type    = "prompter"
    subtype = "anthropic_table_description"
  }

  workflow_nodes {
    name    = "Chunker"
    type    = "chunk"
    subtype = "chunk_by_title"
    settings = {
      "unstructured_api_url"         = null
      "unstructured_api_key"         = null
      "multipage_sections"           = false
      "combine_text_under_n_chars"   = 0
      "include_orig_elements"        = false
      "new_after_n_chars"            = 1500
      "max_characters"               = 2048
      "overlap"                      = 160
      "overlap_all"                  = false
      "contextual_chunking_strategy" = "v1"
    }
  }

  workflow_nodes {
    name    = "Embedder"
    type    = "embed"
    subtype = "azure_openai"
    settings = {
      "model_name" = "text-embedding-3-large"
    }
  }
}

