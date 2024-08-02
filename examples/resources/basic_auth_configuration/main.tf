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

resource "propelauth_basic_auth_configuration" "my_auth_configuration" {
  has_password_login = true
  has_passwordless_login = true
}

output "auth_configuration_result" {
  value = propelauth_basic_auth_configuration.my_auth_configuration
}