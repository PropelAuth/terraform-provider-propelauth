# As there is only one basic_auth_configuration per project there's no need to specify the id,
# but terraform import requires an id to be specified, so we can use an arbitrary string here.
terraform import propelauth_basic_auth_configuration.example arbitrary_string_here