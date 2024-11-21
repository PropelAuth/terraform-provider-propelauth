# As there is only one project_info per project there's no need to specify the id,
# but terraform import requires an id to be specified, so we can use an arbitrary string here.
terraform import propelauth_project_info.example arbitrary_string_here