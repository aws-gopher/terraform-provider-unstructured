// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestSourceResourceSchema(t *testing.T) {
	// Test that the schema is properly defined
	sourceResource := NewSourceResource()

	// Test the schema method
	var resp frameworkresource.SchemaResponse
	sourceResource.Schema(t.Context(), frameworkresource.SchemaRequest{}, &resp)

	if resp.Schema.Attributes["s3"] == nil {
		t.Error("S3 attribute not found in schema")
	}

	if resp.Schema.Attributes["postgres"] == nil {
		t.Error("Postgres attribute not found in schema")
	}

	if resp.Schema.Attributes["name"] == nil {
		t.Error("Name attribute not found in schema")
	}
}

func TestAccSourceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSourceResourceConfig("Terraform Test Source One"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"unstructured_source.test",
						tfjsonpath.New("id"),
						knownvalue.StringFunc(func(v string) error {
							_, err := uuid.ParseUUID(v)
							return err
						}),
					),
					statecheck.ExpectKnownValue(
						"unstructured_source.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Terraform Test Source One"),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "unstructured_source.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSourceResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "unstructured_source" "test" {
  name = %[1]q
  
  s3 {
    remote_url = "s3://example-bucket/"
    anonymous  = true
  }
}
`, name)
}
