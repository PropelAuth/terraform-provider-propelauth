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

resource "propelauth_custom_domain" "my_custom_domain" {
  environment = "Staging"
  domain      = "example.com"
  # subdomain   = "app" # Optional
}

output "project_custom_domain_result" {
  value = propelauth_custom_domain.my_custom_domain
}
