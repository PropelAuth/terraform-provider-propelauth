# As there is only one theme per project there's no need to specify the id,
# but terraform import requires an id to be specified, so we can use an arbitrary string here.
terraform import propelauth_theme.example arbitrary_string_here
# Note: The propelauth_theme resource has many default values for attributes set by the provider
# if you do not provide them. Carefully review the plan for any unexpected changes before applying.