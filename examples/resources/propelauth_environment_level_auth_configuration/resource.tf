resource "propelauth_environment_level_auth_configuration" "test_example" {
  environment                = "Test"
  require_email_confirmation = false
  allow_public_signups       = true
}

resource "propelauth_custom_domain_verification" "my_custom_domain_verification" {
  # Fields are incomplete here for simplicity.
  # See the documentation for the "propelauth_custom_domain_verification" resource for more information
  environment = "Prod"
}

# Prod and Staging environments don't exist until a domain has been verified for them,
# so we need to depend on the verification of the domain before creating the environment-level auth configuration
resource "propelauth_environment_level_auth_configuration" "prod_example" {
  depends_on = [
    propelauth_custom_domain_verification.my_custom_domain_verification
  ]
  environment                = "Prod"
  require_email_confirmation = true
  allow_public_signups       = false
}
