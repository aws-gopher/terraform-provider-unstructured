// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccWorkflowResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkflowResourceConfig("Terraform Test One"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"unstructured_workflow.test",
						tfjsonpath.New("id"),
						knownvalue.StringFunc(func(v string) error {
							_, err := uuid.ParseUUID(v)
							return err
						}),
					),
					statecheck.ExpectKnownValue(
						"unstructured_workflow.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Terraform Test One"),
					),
					statecheck.ExpectKnownValue(
						"unstructured_workflow.test",
						tfjsonpath.New("workflow_type"),
						knownvalue.StringExact("basic"),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "unstructured_workflow.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccWorkflowResourceConfig("Terraform Test Two"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"unstructured_workflow.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Terraform Test Two"),
					),
					statecheck.ExpectKnownValue(
						"unstructured_workflow.test",
						tfjsonpath.New("id"),
						knownvalue.StringFunc(func(v string) error {
							_, err := uuid.ParseUUID(v)
							return err
						}),
					),
					statecheck.ExpectKnownValue(
						"unstructured_workflow.test",
						tfjsonpath.New("workflow_type"),
						knownvalue.StringExact("basic"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWorkflowResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "unstructured_workflow" "test" {
  name = %[1]q
  workflow_type = "basic"
}
`, configurableAttribute)
}
