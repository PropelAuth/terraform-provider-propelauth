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

resource "propelauth_be_api_key" "my_be_api_key" {
  environment = "Test"
  name        = "My API Key Updated"
  read_only   = true
}

output "be_api_key_result" {
  value     = propelauth_be_api_key.my_be_api_key
  sensitive = true
}
