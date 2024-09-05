resource "propelauth_custom_domain" "my_custom_domain" {
  environment = "Prod"
  domain      = "example.com"
}

# AWS Route53 Example
resource "aws_route53_zone" "primary" {
  name = "example.com"
}

resource "aws_route53_record" "txt_record_for_propelauth" {
  zone_id    = aws_route53_zone.primary.zone_id
  name       = propelauth_custom_domain.my_custom_domain.txt_record_key
  type       = "TXT"
  ttl        = 300
  records    = [propelauth_custom_domain.my_custom_domain.txt_record_value]
  depends_on = [propelauth_custom_domain.my_custom_domain]
}

resource "aws_route53_record" "cname_record_for_propelauth" {
  zone_id    = aws_route53_zone.primary.zone_id
  name       = propelauth_custom_domain.my_custom_domain.cname_record_key
  type       = "CNAME"
  ttl        = 300
  records    = [propelauth_custom_domain.my_custom_domain.cname_record_value]
  depends_on = [propelauth_custom_domain.my_custom_domain]
}

# This resource will verify the domain once your DNS records have been set up. See the above
# for how to do that with AWS Route 53 and the output of the "propelauth_custom_domain" resource 
# to set up the DNS records.
resource "propelauth_custom_domain_verification" "my_custom_domain_verification" {
  depends_on = [
    propelauth_custom_domain.my_custom_domain,
    aws_route53_record.txt_record_for_propelauth,
    aws_route53_record.cname_record_for_propelauth
  ]
  environment = propelauth_custom_domain.my_custom_domain.environment
  domain      = propelauth_custom_domain.my_custom_domain.domain
}
