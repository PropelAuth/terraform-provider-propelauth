# As there is only one user_property_settings per project there's no need to specify the id,
# but terraform import requires an id to be specified, so we can use an arbitrary string here.
terraform import propelauth_user_property_settings.example arbitrary_string_here