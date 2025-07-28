package main

// Generate the provider code from the specification files
//go:generate go run github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework generate all --input ./provider-code-spec.json --output ./internal

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name unstructured

// Format Terraform code for use in documentation.
//go:generate terraform fmt -recursive ./examples/
