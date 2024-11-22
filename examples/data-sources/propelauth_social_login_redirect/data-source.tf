# Retrieve the redirect url needed for configuring a google social login 
# in your test environment.
data "propelauth_social_login_redirect" "google_test_redirect" {
  environment = "Test"
  provider    = "Google"
}

output "google_test_redirect_result" {
  value = propelauth_social_login_redirect.google_test_redirect.redirect_url
}

# To do the same for a production or staging environment,
# you must first verify a domain for the target environment,
# or no social login redirect URLs will be available.
data "propelauth_social_login_redirect" "google_prod_redirect" {
  environment = "Prod"
  provider    = "Google"
  depends_on  = [data.propelauth_custom_domain_verification.my_prod_custom_domain_verification]
}

output "google_prod_redirect_result" {
  value = propelauth_social_login_redirect.google_prod_redirect.redirect_url
}