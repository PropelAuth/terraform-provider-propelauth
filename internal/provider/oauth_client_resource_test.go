package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOauthClientResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOauthClientResourceConfig("https://*.redirect.com/path", "Test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_oauth_client.test", "redirect_uris.0", "https://*.redirect.com/path"),
				),
			},
			// Update and Read testing
			{
				Config: testAccOauthClientResourceConfig("https://*.redirect.com/path/v2", "Test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_oauth_client.test", "redirect_uris.0", "https://*.redirect.com/path/v2"),
				),
			},
			{
				Config: testAccOauthClientResourceConfig("https://exact.redirect.com/path/v2", "Test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_oauth_client.test", "redirect_uris.0", "https://exact.redirect.com/path/v2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccOauthClientResourceConfig(redirect_uri string, environment string) string {
	return providerConfig + fmt.Sprintf(`
resource "propelauth_oauth_client" "test" {
  environment   = %[1]q
  redirect_uris = [ %[2]q ]
}
`, environment, redirect_uri)
}
