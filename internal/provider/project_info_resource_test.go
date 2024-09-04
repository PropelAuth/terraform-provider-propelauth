package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectInfoResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectInfoResourceConfig("Terraform Acceptance Testing"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_project_info.test", "name", "Terraform Acceptance Testing"),
				),
			},
			// Update and Read testing
			{
				Config: testAccProjectInfoResourceConfig("Terraform Testing"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_project_info.test", "name", "Terraform Testing"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProjectInfoResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "propelauth_project_info" "test" {
  name = %[1]q
}
`, name)
}
