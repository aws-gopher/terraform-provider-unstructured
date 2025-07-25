resource "unstructured_source" "example" {
  name = "example_s3_source"

  s3 {
    remote_url = "s3://my-source-bucket/"
    recursive  = true
    anonymous  = true
  }
}
