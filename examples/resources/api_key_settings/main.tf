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

resource "propelauth_api_key_settings" "my_api_key_settings" {
  personal_api_keys_enabled                = true
  org_api_keys_enabled                     = true
  invalidate_org_api_key_upon_user_removal = true
  api_key_config = {
    expiration_options = {
      options = [
        "TwoWeeks",
        "OneMonth",
        "ThreeMonths",
      ]
      default = "ThreeMonths"
    }
  }
}

output "api_key_settings_result" {
  value = propelauth_api_key_settings.my_api_key_settings
}