# PropelAuth Terraform Provider

Use terraform to manage your [PropelAuth](https://www.propelauth.com/) integration for authentication, B2B authorization, and user management.

## Using the provider

To use the provider you must first have a [PropelAuth project](https://app.propelauth.com/) where you can setup you can generate the necessary credentials. Find the latest documentation on how to use the provider [here](https://registry.terraform.io/providers/PropelAuth/propelauth/latest/docs).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
