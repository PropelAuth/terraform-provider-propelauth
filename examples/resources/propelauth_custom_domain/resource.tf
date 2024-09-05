# Initialize a custom domain with your PropelAuth environment. This does not verify the domain:
# You will need a DNS provider and the `propelauth_custom_domain_verification` resource to verify the domain.
resource "propelauth_custom_domain" "my_custom_domain" {
  environment = "Prod"
  domain      = "example.com"
  subdomain   = "app"
}