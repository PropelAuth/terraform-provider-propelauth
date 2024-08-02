package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBasicAuthConfigurationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBasicAuthConfigurationResourceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_basic_auth_configuration.test", "has_passwordless_login", "true"),
				),
			},
			// Update and Read testing
			{
				Config: testAccBasicAuthConfigurationResourceConfig(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_basic_auth_configuration.test", "has_passwordless_login", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccBasicAuthConfigurationResourceConfig(has_passwordless_login bool) string {
	return providerConfig + fmt.Sprintf(`
resource "propelauth_basic_auth_configuration" "test" {
  has_passwordless_login = %v
}
`, has_passwordless_login)
}
