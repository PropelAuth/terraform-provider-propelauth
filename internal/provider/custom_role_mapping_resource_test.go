package provider

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomRoleMappingResource(t *testing.T) {
	// Creating a custom role mapping requires a role for every role defined in
	// the project, so the test fixture is sourced from the project itself.
	defaultRole, rolesConfig := testAccCustomRoleMappingRoles(t, "")
	_, updateRolesConfig := testAccCustomRoleMappingRoles(t, "can_invite = true")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCustomRoleMappingResourceConfig("TF Test Mapping", rolesConfig),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_custom_role_mapping.test", "name", "TF Test Mapping"),
					resource.TestCheckResourceAttrSet("propelauth_custom_role_mapping.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccCustomRoleMappingResourceConfig("TF Test Mapping", updateRolesConfig),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_custom_role_mapping.test",
						"roles."+defaultRole+".can_invite", "true"),
				),
			},
			{
				Config: testAccCustomRoleMappingResourceConfig("TF Test Mapping Renamed", rolesConfig),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("propelauth_custom_role_mapping.test", "name", "TF Test Mapping Renamed"),
				),
			},
			// Import testing
			{
				ResourceName:      "propelauth_custom_role_mapping.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// testAccCustomRoleMappingRoles returns the project's default role and a roles
// map literal covering every project role. Role visibility must match the
// project definition, so is_internal is sourced from the live roles.
func testAccCustomRoleMappingRoles(t *testing.T, defaultRoleExtra string) (string, string) {
	t.Helper()

	tenantID := os.Getenv("PROPELAUTH_TENANT_ID")
	projectID := os.Getenv("PROPELAUTH_PROJECT_ID")
	apiKey := os.Getenv("PROPELAUTH_API_KEY")
	client, err := propelauth.NewClient(&tenantID, &projectID, &apiKey)
	if err != nil {
		t.Fatalf("failed to create PropelAuth client: %s", err)
	}

	rolesAndPermissions, err := client.GetRolesAndPermissions()
	if err != nil {
		t.Fatalf("failed to fetch roles and permissions: %s", err)
	}

	var roles strings.Builder
	for _, role := range rolesAndPermissions.Roles {
		fmt.Fprintf(&roles, "\t\t%q = {\n\t\t\tis_internal = %t\n", role.Name, !role.IsVisibleToEndUser)
		if role.Name == rolesAndPermissions.DefaultRole && defaultRoleExtra != "" {
			fmt.Fprintf(&roles, "\t\t\t%s\n", defaultRoleExtra)
		}
		roles.WriteString("\t\t}\n")
	}

	return rolesAndPermissions.DefaultRole, roles.String()
}

func testAccCustomRoleMappingResourceConfig(name string, rolesConfig string) string {
	return providerConfig + fmt.Sprintf(`
resource "propelauth_custom_role_mapping" "test" {
  name = %[1]q
  roles = {
%[2]s  }
}
`, name, rolesConfig)
}
