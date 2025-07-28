package main

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name unstructured

// Format Terraform code for use in documentation.
//go:generate terraform fmt -recursive ./examples/
