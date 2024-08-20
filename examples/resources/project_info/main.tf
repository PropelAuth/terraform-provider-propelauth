terraform {
  required_providers {
    propelauth = {
      source = "registry.terraform.io/propelauth/propelauth"
    }
  }
}

provider "propelauth" {
  tenant_id  = "e1dc8461-5d8a-4bad-a929-19745de693f4"                                                             # or PROPELAUTH_TENANT_ID environment variable
  project_id = "5a5f7a4f-1a51-4312-bbbe-4126cceab59b"                                                             # or PROPELAUTH_PROJECT_ID environment variable
  api_key    = "c557308180b7da18d7e0e9cbd2ae3b36833c0165b5158c439efe59662df01701c2e23b00211b9c25b5223e51417f323b" # or PROPELAUTH_API_KEY environment variable
}

resource "propelauth_project_info" "my_project_info" {
  name = "name-set-by-terraform"
}

output "project_info_result" {
  value = propelauth_project_info.my_project_info
}
