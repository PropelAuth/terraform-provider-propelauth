package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationConfigurationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOrganizationConfigurationResourceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_organization_configuration.test", "users_can_delete_their_own_orgs", "true"),
				),
			},
			// Update and Read testing
			{
				Config: testAccOrganizationConfigurationResourceConfig(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_organization_configuration.test", "users_can_delete_their_own_orgs", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccOrganizationConfigurationResourceConfig(users_can_delete_their_own_orgs bool) string {
	return providerConfig + fmt.Sprintf(`
resource "propelauth_organization_configuration" "test" {
  users_can_delete_their_own_orgs = %v
}
`, users_can_delete_their_own_orgs)
}
