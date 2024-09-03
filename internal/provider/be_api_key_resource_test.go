package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBeApiKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBeApiKeyResourceConfig("First Name", false, "Test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_be_api_key.test", "name", "First Name"),
					resource.TestCheckResourceAttr("propelauth_be_api_key.test", "read_only", "false"),
				),
			},
			// Update and Read testing
			{
				Config: testAccBeApiKeyResourceConfig("Second Name", false, "Test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_be_api_key.test", "name", "Second Name"),
					resource.TestCheckResourceAttr("propelauth_be_api_key.test", "read_only", "false"),
				),
			},
			{
				Config: testAccBeApiKeyResourceConfig("Second Name", true, "Test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_be_api_key.test", "name", "Second Name"),
					resource.TestCheckResourceAttr("propelauth_be_api_key.test", "read_only", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccBeApiKeyResourceConfig(name string, read_only bool, environment string) string {
	return providerConfig + fmt.Sprintf(`
resource "propelauth_be_api_key" "test" {
  environment = %[1]q
  name        = %[2]q
  read_only   = %[3]t
}
`, environment, name, read_only)
}
