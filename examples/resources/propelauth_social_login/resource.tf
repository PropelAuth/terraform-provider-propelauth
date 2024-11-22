
variable "github_client_secret" {
  type      = string
  sensitive = true
}

resource "propelauth_social_login" "github_sso" {
  social_provider = "GitHub"
  client_id       = "my-client-id"
  client_secret   = var.github_client_secret
}
