package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRolesAndPermissionsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRolesAndPermissionsResourceConfig("ai::deploy", "Member", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"propelauth_roles_and_permissions.test",
						"roles.Support.permissions.3",
						"ai::deploy",
					),
					resource.TestCheckResourceAttr(
						"propelauth_roles_and_permissions.test",
						"default_role",
						"Member",
					),
					resource.TestCheckResourceAttr(
						"propelauth_roles_and_permissions.test",
						"roles.Admin.can_manage_api_keys",
						"true",
					),
				),
			},
			// Update and Read testing
			{
				Config: testAccRolesAndPermissionsResourceConfig("ai::deploy::v2", "Admin", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"propelauth_roles_and_permissions.test",
						"roles.Support.permissions.3",
						"ai::deploy::v2",
					),
					resource.TestCheckResourceAttr(
						"propelauth_roles_and_permissions.test",
						"default_role",
						"Admin",
					),
					resource.TestCheckResourceAttr(
						"propelauth_roles_and_permissions.test",
						"roles.Admin.can_manage_api_keys",
						"false",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRolesAndPermissionsResourceConfig(permission string, defaultRole string, adminCanManageApiKeys bool) string {
	return providerConfig + fmt.Sprintf(`
resource "propelauth_roles_and_permissions" "test" {
  permissions = [
	  {
	    name = "doc::read"
	    display_name = "Can read documents." # optional
	    description = "A description here." # optional
	  },
	  {
	    name = "doc::write"
	  },
	  {
	    name = "ticket::read"
	  },
	  {
	    name = "ticket::write"
	  },
	  {
	    name = %[1]q
	  }
	]
	# multiple_roles_per_user = false # default is false
	roles = {
    "Owner" = {
      can_view_other_members = true
      can_invite = true
      can_change_roles = true
      can_manage_api_keys = true
      can_remove_users = true
      can_setup_saml = true # will always be true for the default_owner_role
      can_delete_org = true # will always be true for the default_owner_role
      can_edit_org_access = true
      can_update_org_metadata = true
      permissions = [
        "doc::read",
        "doc::write",
        "ticket::read",
        "ticket::write"
      ]
      # roles_can_manage = ["Admin", "Member"] # only for multiple_roles_per_user = true
    }
    "Admin" = {
      can_view_other_members = true
      can_invite = true
      can_change_roles = true
      can_manage_api_keys = %[3]t
      can_remove_users = true
      can_setup_saml = false
      can_delete_org = false
      can_edit_org_access = false
      can_update_org_metadata = false
      permissions = [
        "doc::read",
        "doc::write",
        "ticket::read"
      ]
    # roles_can_manage = ["Member"] # only for multiple_roles_per_user = true
    }
    "Member" = {
      # the defaults for all PropelAuth permissions
      permissions = [
        "doc::read",
        "ticket::read"
      ]
    }
    "Support" = {
      is_internal = true
      can_view_other_members = true
      can_invite = false
      can_change_roles = false
      can_manage_api_keys = false
      can_remove_users = false
      can_setup_saml = false
      can_delete_org = false
      can_edit_org_access = false
      can_update_org_metadata = false
      permissions = [
        "doc::read",
        "doc::write",
        "ticket::read",
        %[1]q
      ]
	  }
	}
	role_hierarchy = [
	  "Owner",
	  "Admin",
	  "Support",
	  "Member"
	]
	default_role = %[2]q
	default_owner_role = "Owner"
}
`, permission, defaultRole, adminCanManageApiKeys)
}
