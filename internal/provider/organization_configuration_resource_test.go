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
				Config: testAccOrganizationConfigurationResourceConfig(true, false, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_organization_configuration.test", "users_can_delete_their_own_orgs", "true"),
					resource.TestCheckResourceAttr("propelauth_organization_configuration.test", "orgs_can_setup_saml", "false"),
				),
			},
			// Update and Read testing
			{
				Config: testAccOrganizationConfigurationResourceConfig(false, true, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_organization_configuration.test", "users_can_delete_their_own_orgs", "false"),
					resource.TestCheckResourceAttr("propelauth_organization_configuration.test", "orgs_can_setup_saml", "true"),
					resource.TestCheckResourceAttr("propelauth_organization_configuration.test", "skip_saml_role_mapping_step", "false"),
				),
			},
			{
				Config: testAccOrganizationConfigurationResourceConfig(false, true, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_organization_configuration.test", "users_can_delete_their_own_orgs", "false"),
					resource.TestCheckResourceAttr("propelauth_organization_configuration.test", "orgs_can_setup_saml", "true"),
					resource.TestCheckResourceAttr("propelauth_organization_configuration.test", "skip_saml_role_mapping_step", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccOrganizationConfigurationResourceConfig(
	usersCanDeleteTheirOwnOrgs bool,
	orgsCanViewTheirAuditLog bool,
	orgsCanSetupSaml bool,
	skipSamlRoleMappingStep bool,
) string {
	if orgsCanViewTheirAuditLog {
		return providerConfig + fmt.Sprintf(`
resource "propelauth_organization_configuration" "test" {
  users_can_delete_their_own_orgs = %[1]t
  orgs_can_setup_saml = %[2]t
  skip_saml_role_mapping_step = %[3]t
  customer_org_audit_log_settings = {
	enabled = true
	all_orgs_can_view_their_audit_log = false
	include_impersonation = true
	include_employee_actions = false
	include_api_key_actions = false
  }
}
`, usersCanDeleteTheirOwnOrgs, orgsCanSetupSaml, skipSamlRoleMappingStep)
	} else {
		return providerConfig + fmt.Sprintf(`
resource "propelauth_organization_configuration" "test" {
  users_can_delete_their_own_orgs = %[1]t
  orgs_can_setup_saml = %[2]t
  skip_saml_role_mapping_step = %[3]t
}
`, usersCanDeleteTheirOwnOrgs, orgsCanSetupSaml, skipSamlRoleMappingStep)
	}
}
