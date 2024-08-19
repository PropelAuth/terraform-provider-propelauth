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

resource "propelauth_organization_configuration" "my_org_configuration" {
  has_orgs                        = true
  orgs_metaname                   = "Company"
  max_num_orgs_users_can_be_in    = 1
  users_can_delete_their_own_orgs = true
  orgs_can_require_2fa            = true
}

output "org_configuration_result" {
  value = propelauth_organization_configuration.my_org_configuration
}