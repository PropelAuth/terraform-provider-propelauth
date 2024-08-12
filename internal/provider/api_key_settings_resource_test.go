package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApiKeySettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApiKeySettingsResourceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_api_key_settings.test", "org_api_keys_enabled", "true"),
				),
			},
			// Update and Read testing
			{
				Config: testAccApiKeySettingsResourceConfig(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_api_key_settings.test", "org_api_keys_enabled", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccApiKeySettingsResourceConfig(org_api_keys_enabled bool) string {
	return providerConfig + fmt.Sprintf(`
resource "propelauth_api_key_settings" "test" {
  personal_api_keys_enabled = true
  org_api_keys_enabled = %v
  api_key_config = {
    expiration_options = {
      options = [
        "TwoWeeks",
        "OneMonth",
        "ThreeMonths",
      ]
      default = "ThreeMonths"
    }
  }
}
`, org_api_keys_enabled)
}
