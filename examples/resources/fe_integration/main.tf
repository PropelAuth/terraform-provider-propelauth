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

resource "propelauth_fe_integration" "my_test_fe_integration" {
  environment          = "Test"
  application_url      = "http://localhost:3001"
  login_redirect_path  = "/home"
  logout_redirect_path = "/goodbye"
  additional_fe_locations = [
    {
      domain              = "http://localhost:3002"
      allow_any_subdomain = false
    }
  ]
}

output "test_fe_integration_result" {
  value = propelauth_fe_integration.my_test_fe_integration
}

resource "propelauth_fe_integration" "my_prod_fe_integration" {
  environment          = "Prod"
  application_url      = "https://app.sharpenchess.com"
  login_redirect_path  = "/home"
  logout_redirect_path = "/goodbye"
  additional_fe_locations = [
    {
      domain              = "https://worker.sharpenchess.com"
      allow_any_subdomain = false
    }
  ]
}

output "prod_fe_integration_result" {
  value = propelauth_fe_integration.my_prod_fe_integration
}
