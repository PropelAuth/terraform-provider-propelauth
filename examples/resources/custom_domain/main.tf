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
  tenant_id  = "<PROPELAUTH_TENANT_ID>"  # or PROPELAUTH_TENANT_ID environment variable
  project_id = "<PROPELAUTH_PROJECT_ID>" # or PROPELAUTH_PROJECT_ID environment variable
  api_key    = "<PROPELAUTH_API_KEY>"    # or PROPELAUTH_API_KEY environment variable
}

# This resource just sets up the process of verifying the domain. 
# It will return the TXT and CNAME records that you need to add to your DNS settings.
# You will need to add these records to your DNS settings manually or using Terraform.
# Then, the "propelauth_custom_domain_verification" resource will verify the domain.
resource "propelauth_custom_domain" "my_custom_domain" {
  # The environment can either be "Prod" or "Staging"
  environment = "Prod" 

  # The domain to use. Your auth domain will be `auth.<domain>`. 
  # You can also specify a subdomain like prod.example.com which will result in auth.prod.example.com
  domain      = "example.com"

  # The subdomain where your application is hosted. 
  # This is optional, but recommended, as it will allow PropelAuth to automatically redirect users to your application after they login.
  # subdomain   = "app"
}

# This resource will verify the domain. You must set up the DNS records first.
# You can use the output of the "propelauth_custom_domain" resource to set up the DNS records.
resource "propelauth_custom_domain_verification" "my_custom_domain_verification" {
  depends_on = [
    propelauth_custom_domain.my_custom_domain,

    # See below for an example of using Route53 for the custom domain
    aws_route53_record.txt_record_for_propelauth,
    aws_route53_record.cname_record_for_propelauth
  ]
  environment = propelauth_custom_domain.my_custom_domain.environment
  domain      = propelauth_custom_domain.my_custom_domain.domain
}

output "project_custom_domain_result" {
  value = propelauth_custom_domain.my_custom_domain
}



# AWS Route53 Example
resource "aws_route53_zone" "primary" {
  name = "example.com"
}

resource "aws_route53_record" "txt_record_for_propelauth" {
  zone_id = aws_route53_zone.primary.zone_id
  name    = propelauth_custom_domain.my_custom_domain.txt_record_key
  type    = "TXT"
  ttl     = 300
  records = [propelauth_custom_domain.my_custom_domain.txt_record_value]
  depends_on = [propelauth_custom_domain.my_custom_domain]
}

resource "aws_route53_record" "cname_record_for_propelauth" {
  zone_id = aws_route53_zone.primary.zone_id
  name    = propelauth_custom_domain.my_custom_domain.cname_record_key
  type    = "CNAME"
  ttl     = 300
  records = [propelauth_custom_domain.my_custom_domain.cname_record_value]
  depends_on = [propelauth_custom_domain.my_custom_domain]
}
