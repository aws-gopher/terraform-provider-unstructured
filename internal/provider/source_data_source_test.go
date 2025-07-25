// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccSourceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccSourceDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.unstructured_source.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("f49d1dcc-99f4-4be1-9afe-ce2495654ffb"),
					),
					statecheck.ExpectKnownValue(
						"data.unstructured_source.test",
						tfjsonpath.New("name"),
						knownvalue.StringFunc(func(v string) error {
							if v == "" {
								return fmt.Errorf("expected name to not be empty")
							}
							return nil
						}),
					),
				},
			},
		},
	})
}

const testAccSourceDataSourceConfig = `
data "unstructured_source" "test" {
  id = "f49d1dcc-99f4-4be1-9afe-ce2495654ffb"
}
`
