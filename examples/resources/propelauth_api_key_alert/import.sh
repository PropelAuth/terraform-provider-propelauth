# As there is only one api_key_alert per project there's no need to specify the id,
# but terraform import requires an id to be specified, so we can use an arbitrary string here.
terraform import propelauth_api_key_alert.example arbitrary_string_here