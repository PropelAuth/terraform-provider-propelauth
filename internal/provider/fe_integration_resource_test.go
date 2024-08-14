package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFeIntegrationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFeIntegrationResourceConfig("3000", "3001"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"propelauth_fe_integration.test",
						"application_url",
						"http://localhost:3000",
					),
					resource.TestCheckResourceAttr(
						"propelauth_fe_integration.test",
						"login_redirect_path",
						"/home",
					),
					resource.TestCheckResourceAttr(
						"propelauth_fe_integration.test",
						"logout_redirect_path",
						"/goodbye",
					),
					resource.TestCheckResourceAttr(
						"propelauth_fe_integration.test",
						"additional_fe_locations.0.domain",
						"http://localhost:3001",
					),
				),
			},
			// Update and Read testing
			{
				Config: testAccFeIntegrationResourceConfig("3001", "3002"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"propelauth_fe_integration.test",
						"application_url",
						"http://localhost:3001",
					),
					resource.TestCheckResourceAttr(
						"propelauth_fe_integration.test",
						"login_redirect_path",
						"/home",
					),
					resource.TestCheckResourceAttr(
						"propelauth_fe_integration.test",
						"logout_redirect_path",
						"/goodbye",
					),
					resource.TestCheckResourceAttr(
						"propelauth_fe_integration.test",
						"additional_fe_locations.0.domain",
						"http://localhost:3002",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccFeIntegrationResourceConfig(port string, additionalPort string) string {
	return providerConfig + fmt.Sprintf(`
resource "propelauth_fe_integration" "test" {
  environment = "Test"
  application_url = "http://localhost:%v"
  login_redirect_path = "/home"
  logout_redirect_path = "/goodbye"
  additional_fe_locations = [
    {
      domain = "http://localhost:%v"
      allow_any_subdomain = false
    }
  ]
}
`, port, additionalPort)
}
