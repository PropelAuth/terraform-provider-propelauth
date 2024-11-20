package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnvironmentLevelAuthConfigurationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEnvironmentLevelAuthConfigurationResourceConfig(true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_environment_level_auth_configuration.test", "require_email_confirmation", "true"),
					resource.TestCheckResourceAttr("propelauth_environment_level_auth_configuration.test", "allow_public_signups", "false"),
				),
			},
			// Update and Read testing
			{
				Config: testAccEnvironmentLevelAuthConfigurationResourceConfig(false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_environment_level_auth_configuration.test", "require_email_confirmation", "false"),
					resource.TestCheckResourceAttr("propelauth_environment_level_auth_configuration.test", "allow_public_signups", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccEnvironmentLevelAuthConfigurationResourceConfig(require_email_confirmation bool, allow_public_signups bool) string {
	return providerConfig + fmt.Sprintf(`
resource "propelauth_environment_level_auth_configuration" "test" {
  environment = "Test"
  require_email_confirmation = %v
  allow_public_signups = %v
}
`, require_email_confirmation, allow_public_signups)
}
