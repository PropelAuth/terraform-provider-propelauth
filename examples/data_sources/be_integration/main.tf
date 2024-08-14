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

data "propelauth_be_integration" "my_test_be_integration" {
  environment = "Test"
}

output "test_be_integration_result" {
  value = data.propelauth_be_integration.my_test_be_integration
}

data "propelauth_be_integration" "my_prod_be_integration" {
  environment = "Prod"
}

output "prod_be_integration_result" {
  value = data.propelauth_be_integration.my_prod_be_integration
}
