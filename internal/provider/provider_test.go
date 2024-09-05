package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// PROPELAUTH_TENANT_ID=e1dc8461-5d8a-4bad-a929-19745de693f4
// PROPELAUTH_PROJECT_ID=5a5f7a4f-1a51-4312-bbbe-4126cceab59b
// PROPELAUTH_API_KEY=c557308180b7da18d7e0e9cbd2ae3b36833c0165b5158c439efe59662df01701c2e23b00211b9c25b5223e51417f323b

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the PropelAuth client is properly configured.
	// It is also possible to use the PROPELAUTH_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "propelauth" {
#   tenant_id = "<PROPELAUTH_TENANT_ID>"
#   project_id = "<PROPELAUTH_PROJECT_ID>"
#   api_key = "<PROPELAUTH_API_KEY>"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"propelauth": providerserver.NewProtocol6WithError(New("test")()),
	}
)
