terraform {
  required_providers {
    propelauth = {
      source = "registry.terraform.io/propelauth/propelauth"
    }
    # aws = {
    #   source  = "hashicorp/aws"
    #   version = "~> 5.0"
    # }
  }
}

provider "propelauth" {
  # tenant_id  = "<PROPELAUTH_TENANT_ID>"  # or PROPELAUTH_TENANT_ID environment variable
  # project_id = "<PROPELAUTH_PROJECT_ID>" # or PROPELAUTH_PROJECT_ID environment variable
  # api_key    = "<PROPELAUTH_API_KEY>"    # or PROPELAUTH_API_KEY environment variable
}

resource "propelauth_custom_domain" "my_custom_domain" {
  environment = "Staging"
  domain      = "itailevi.com"
  # subdomain   = "app" # Optional
}

# resource "aws_route53_record" "txt_record_for_propelauth" {
#   zone_id = aws_route53_zone.primary.zone_id
#   name    = propelauth_custom_domain.my_custom_domain.txt_record_key
#   type    = "TXT"
#   ttl     = 300
#   records = [propelauth_custom_domin.my_custom_domain.txt_record_value]
# }

# resource "aws_route53_record" "cname_record_for_propelauth" {
#   zone_id = aws_route53_zone.primary.zone_id
#   name    = propelauth_custom_domain.my_custom_domain.cname_record_key
#   type    = "CNAME"
#   records = [propelauth_custom_domin.my_custom_domain.cname_record_value]
# }

resource "propelauth_custom_domain_verification" "my_custom_domain_verification" {
  depends_on = [
    propelauth_custom_domain.my_custom_domain,
  ]
  environment = propelauth_custom_domain.my_custom_domain.environment
  # timeouts { create = "15m" }
}

output "project_custom_domain_result" {
  value = propelauth_custom_domain.my_custom_domain
}
