package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSocialLoginResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSocialLoginResourceConfig("client-id", "SECRET"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_social_login.test", "client_id", "client-id"),
					resource.TestCheckResourceAttr("propelauth_social_login.test", "client_secret", "SECRET"),
				),
			},
			// Update and Read testing
			{
				Config: testAccSocialLoginResourceConfig("client-id-v2", "SECRET2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_social_login.test", "client_id", "client-id-v2"),
					resource.TestCheckResourceAttr("propelauth_social_login.test", "client_secret", "SECRET2"),
				),
			},
			{
				Config: testAccSocialLoginResourceConfig("client-id-v2", "SECRET3"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_social_login.test", "client_id", "client-id-v2"),
					resource.TestCheckResourceAttr("propelauth_social_login.test", "client_secret", "SECRET3"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccSocialLoginResourceConfig(clientId string, clientSecret string) string {
	return providerConfig + fmt.Sprintf(`
resource "propelauth_social_login" "test" {
	social_provider = "GitHub"
	client_id = %[1]q
	client_secret = %[2]q
}	  
`, clientId, clientSecret)
}
