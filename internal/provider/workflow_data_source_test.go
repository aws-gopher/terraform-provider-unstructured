// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccWorkflowDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccWorkflowDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.unstructured_workflow.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("f57569fe-36e4-4f3f-9ab3-447dc9419d2f"),
					),
					statecheck.ExpectKnownValue(
						"data.unstructured_workflow.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Project Bronco"),
					),
				},
			},
		},
	})
}

const testAccWorkflowDataSourceConfig = `
data "unstructured_workflow" "test" {
  id = "f57569fe-36e4-4f3f-9ab3-447dc9419d2f"
}
`
