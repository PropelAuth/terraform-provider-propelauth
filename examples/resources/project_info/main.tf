terraform {
  required_providers {
    propelauth = {
      source = "registry.terraform.io/propelauth/propelauth"
    }
  }
}

provider "propelauth" {
  # tenant_id  = "<PROPELAUTH_TENANT_ID>"  # or PROPELAUTH_TENANT_ID environment variable
  # project_id = "<PROPELAUTH_PROJECT_ID>" # or PROPELAUTH_PROJECT_ID environment variable
  # api_key    = "<PROPELAUTH_API_KEY>"    # or PROPELAUTH_API_KEY environment variable
}

resource "propelauth_project_info" "my_project_info" {
  name = "name-set-by-terraform"
}

output "project_info_result" {
  value = propelauth_project_info.my_project_info
}