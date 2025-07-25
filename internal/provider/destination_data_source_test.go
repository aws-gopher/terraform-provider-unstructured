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

func TestAccDestinationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccDestinationDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.unstructured_destination.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("1c684f00-0e7f-460b-ac31-21c062551997"),
					),
					statecheck.ExpectKnownValue(
						"data.unstructured_destination.test",
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

const testAccDestinationDataSourceConfig = `
data "unstructured_destination" "test" {
  id = "1c684f00-0e7f-460b-ac31-21c062551997"
}
`
