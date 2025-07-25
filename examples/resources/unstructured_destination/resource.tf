resource "unstructured_destination" "example" {
  name = "example_s3_destination"

  s3 {
    remote_url = "s3://my-destination-bucket/"
    key        = "aws-access-key-id"
    secret     = "aws-secret-access-key"
  }
}
