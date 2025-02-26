# Create a PropelAuth OAuth client for integrating your application.
resource "propelauth_oauth_client" "example" {
  # for clients in Prod and Staging environments, the environment must first be live
  depends_on  = [propelauth_custom_domain_verification.my_prod_domain_verification]
  environment = "Prod"
  redirect_uris = [
    "https://*.example.com/oauth/callback"
  ]
}

output "oauth_client_id_result" {
  value     = propelauth_oauth_client.example.client_id
  sensitive = false
}

output "oauth_client_secret_result" {
  value     = propelauth_oauth_client.example.client_secret
  sensitive = true
}