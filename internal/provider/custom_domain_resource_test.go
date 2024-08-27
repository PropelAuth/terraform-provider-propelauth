package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomDomainResource(t *testing.T) {
	subdomain := "app"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCustomDomainResourceConfig("example.com", nil, "Prod"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_custom_domain.test", "environment", "Prod"),
					resource.TestCheckResourceAttr("propelauth_custom_domain.test", "domain", "example.com"),
					resource.TestCheckResourceAttrSet("propelauth_custom_domain.test", "txt_record_key"),
					resource.TestCheckResourceAttrSet("propelauth_custom_domain.test", "txt_record_value"),
					resource.TestCheckResourceAttrSet("propelauth_custom_domain.test", "cname_record_key"),
					resource.TestCheckResourceAttrSet("propelauth_custom_domain.test", "cname_record_value"),
					resource.TestCheckNoResourceAttr("propelauth_custom_domain.test", "subdomain"),
				),
			},
			// Update and Read testing
			{
				Config: testAccCustomDomainResourceConfig("example2.com", nil, "Prod"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_custom_domain.test", "environment", "Prod"),
					resource.TestCheckResourceAttr("propelauth_custom_domain.test", "domain", "example2.com"),
					resource.TestCheckResourceAttrSet("propelauth_custom_domain.test", "txt_record_key"),
					resource.TestCheckResourceAttrSet("propelauth_custom_domain.test", "txt_record_value"),
					resource.TestCheckResourceAttrSet("propelauth_custom_domain.test", "cname_record_key"),
					resource.TestCheckResourceAttrSet("propelauth_custom_domain.test", "cname_record_value"),
					resource.TestCheckNoResourceAttr("propelauth_custom_domain.test", "subdomain"),
				),
			},
			{
				Config: testAccCustomDomainResourceConfig("example.com", &subdomain, "Prod"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_custom_domain.test", "environment", "Prod"),
					resource.TestCheckResourceAttr("propelauth_custom_domain.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("propelauth_custom_domain.test", "subdomain", "app"),
					resource.TestCheckResourceAttrSet("propelauth_custom_domain.test", "txt_record_key"),
					resource.TestCheckResourceAttrSet("propelauth_custom_domain.test", "txt_record_value"),
					resource.TestCheckResourceAttrSet("propelauth_custom_domain.test", "cname_record_key"),
					resource.TestCheckResourceAttrSet("propelauth_custom_domain.test", "cname_record_value"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCustomDomainResourceConfig(domain string, subdomain *string, environment string) string {
	if subdomain != nil {
		return providerConfig + fmt.Sprintf(`
resource "propelauth_custom_domain" "test" {
  environment = %[1]q
  domain = %[2]q
  subdomain = %[3]q
}
`, environment, domain, *subdomain)
	}
	return providerConfig + fmt.Sprintf(`
resource "propelauth_custom_domain" "test" {
  environment = %[1]q
  domain = %[2]q
}
`, environment, domain)
}
