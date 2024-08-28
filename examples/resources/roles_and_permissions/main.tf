terraform {
  required_providers {
    propelauth = {
      source = "registry.terraform.io/propelauth/propelauth"
    }
  }
}

provider "propelauth" {
#   tenant_id  = "<PROPELAUTH_TENANT_ID>"  # or PROPELAUTH_TENANT_ID environment variable
#   project_id = "<PROPELAUTH_PROJECT_ID>" # or PROPELAUTH_PROJECT_ID environment variable
#   api_key    = "<PROPELAUTH_API_KEY>"    # or PROPELAUTH_API_KEY environment variable
}

resource "propelauth_roles_and_permissions" "my_roles_and_permissions" {
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
	    name = "ai::deploy"
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
      can_manage_api_keys = true
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
      replacing_role = "Member With Wrong Name"
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
        "ai::deploy"
      ]
	  }
	  # "New Role Name" = {
		#   replacing_role = "Old Role Name" # this field is for when you want to rename a role
	  # }
	}
	role_hierarchy = [ # only for multiple_roles_per_user = false
	  "Owner",
	  "Admin",
	  "Member",
	  "Support"
	]
	default_role = "Member"
	default_owner_role = "Owner"
	# default_mapping_name = "Free Tier"
}

output "test_fe_integration_result" {
  value = propelauth_roles_and_permissions.my_roles_and_permissions
}
