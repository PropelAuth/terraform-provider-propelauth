# As there is only one default roles_and_permissions per project there's no need to specify the id,
# but terraform import requires an id to be specified, so we can use an arbitrary string here.
terraform import propelauth_roles_and_permissions.example arbitrary_string_here