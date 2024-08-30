package provider

import (
	"context"
	"fmt"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &rolePermissionsOverrideResource{}
var _ resource.ResourceWithConfigure   = &rolePermissionsOverrideResource{}

func NewRolePermissionsOverrideResource() resource.Resource {
	return &rolePermissionsOverrideResource{}
}

// rolePermissionsOverrideResource defines the resource implementation.
type rolePermissionsOverrideResource struct {
	client *propelauth.PropelAuthClient
}

// rolePermissionsOverrideResourceModel describes the resource data model.
type rolePermissionsOverrideResourceModel struct {
	MappingId types.String `tfsdk:"mapping_id"`
	Name types.String `tfsdk:"name"`
	Roles map[string]roleOverrideModel `tfsdk:"roles"`
}

type roleOverrideModel struct {
	// the optional overrides managed by the resource
	CanViewOtherMembers types.Bool `tfsdk:"can_view_other_members"`
	CanInvite types.Bool `tfsdk:"can_invite"`
	CanChangeRoles types.Bool `tfsdk:"can_change_roles"`
	CanManageApiKeys types.Bool `tfsdk:"can_manage_api_keys"`
	CanRemoveUsers types.Bool `tfsdk:"can_remove_users"`
	CanSetupSaml types.Bool `tfsdk:"can_setup_saml"`
	CanDeleteOrg types.Bool `tfsdk:"can_delete_org"`
	CanEditOrgAccess types.Bool `tfsdk:"can_edit_org_access"`
	CanUpdateOrgMetadata types.Bool `tfsdk:"can_update_org_metadata"`
	Permissions []types.String `tfsdk:"permissions"`
	Disabled types.Bool `tfsdk:"disabled"`
	// computed fields for tracking changes to the default settings
	DefaultCanViewOtherMembers types.Bool `tfsdk:"default_can_view_other_members"`
	DefaultCanInvite types.Bool `tfsdk:"default_can_invite"`
	DefaultCanChangeRoles types.Bool `tfsdk:"default_can_change_roles"`
	DefaultCanManageApiKeys types.Bool `tfsdk:"default_can_manage_api_keys"`
	DefaultCanRemoveUsers types.Bool `tfsdk:"default_can_remove_users"`
	DefaultCanSetupSaml types.Bool `tfsdk:"default_can_setup_saml"`
	DefaultCanDeleteOrg types.Bool `tfsdk:"default_can_delete_org"`
	DefaultCanEditOrgAccess types.Bool `tfsdk:"default_can_edit_org_access"`
	DefaultCanUpdateOrgMetadata types.Bool `tfsdk:"default_can_update_org_metadata"`
	DefaultPermissions []types.String `tfsdk:"default_permissions"`
	DefaultRolesCanManage []types.String `tfsdk:"default_roles_can_manage"`
	DefaultIsInternal types.Bool `tfsdk:"default_is_internal"`
	DefaultDisabled types.Bool `tfsdk:"default_disabled"`
}


func (r *rolePermissionsOverrideResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_permissions_override"
}

func (r *rolePermissionsOverrideResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Role-Permissions Override resource. This is used to expand and customize the permissions specific organizations have " +
			"access to. You can assign organizations to an override in the PropelAuth dashboard or programatically through a BE integration.",
		Attributes: map[string]schema.Attribute{
			"mapping_id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Description: "The name of the role-permission override (aka \"Mapping\" in the PropelAuth dashboard).",
			},
			"roles": schema.MapNestedAttribute{
				Optional: true,
				Computed: true,
				Description: "A map of roles with permission overrides. This resource depends on the `propelauth_roles_and_permissions_resource` " +
					"to be created first, as all the overrides defined here are applied on top of the base definitions of roles in the " +
					"parent resource.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"can_view_other_members": schema.BoolAttribute{
							Optional: true,
							Description: "If true, users with this role in the org can view other members of the organization.",
						},
						"can_invite": schema.BoolAttribute{
							Optional: true,
							Description: "If true, users with this role in the org can invite other users to the organization.",
						},
						"can_change_roles": schema.BoolAttribute{
							Optional: true,
							Description: "If true, users with this role in the org can change the roles of other users in the organization.",
						},	
						"can_manage_api_keys": schema.BoolAttribute{
							Optional: true,
							Description: "If true, users with this role in the org can manage API keys for the organization.",
						},
						"can_remove_users": schema.BoolAttribute{
							Optional: true,
							Description: "If true, users with this role in the org can remove other users from the organization.",
						},
						"can_setup_saml": schema.BoolAttribute{
							Optional: true,
							Description: "If true, users with this role in the org can setup enterprise SSO for the organization.",
						},
						"can_delete_org": schema.BoolAttribute{
							Optional: true,
							Description: "If true, users with this role in the org can delete the organization.",
						},
						"can_edit_org_access": schema.BoolAttribute{
							Optional: true,
							Description: "If true, users with this role in the org can edit the organization's access settings. " +
								"These settings incluede what email domains are included and whether 2FA is enforced for the org.",
						},
						"can_update_org_metadata": schema.BoolAttribute{
							Optional: true,
							Description: "If true, users with this role in the org can update the organization's metadata. " +
								"This includes changing the name of the organization.",
						},
						"permissions": schema.ListAttribute{
							ElementType: types.StringType,
							Optional: true,
							Computed: true,
							Default: listdefault.StaticValue(types.ListValueMust(
								types.StringType,
								[]attr.Value{types.StringValue("propelauth::null_permission")},
							)),
							Description: "A list of permissions specific to your application that are assigned to this role. These must first " +
								"exist in the `propelauth_roles_and_permissions_resource.permissions` attribute.",
						},
						"disabled": schema.BoolAttribute{
							Optional: true,
							Description: "If true, this role is disabled and cannot be assigned to users.",
						},
						"default_can_view_other_members": schema.BoolAttribute{
							Computed: true,
						},
						"default_can_invite": schema.BoolAttribute{
							Computed: true,
						},
						"default_can_change_roles": schema.BoolAttribute{
							Computed: true,
						},	
						"default_can_manage_api_keys": schema.BoolAttribute{
							Computed: true,
						},
						"default_can_remove_users": schema.BoolAttribute{
							Computed: true,
						},
						"default_can_setup_saml": schema.BoolAttribute{
							Computed: true,
						},
						"default_can_delete_org": schema.BoolAttribute{
							Computed: true,
						},
						"default_can_edit_org_access": schema.BoolAttribute{
							Computed: true,
						},
						"default_can_update_org_metadata": schema.BoolAttribute{
							Computed: true,
						},
						"default_permissions": schema.ListAttribute{
							ElementType: types.StringType,
							Computed: true,
						},
						"default_roles_can_manage": schema.ListAttribute{
							ElementType: types.StringType,
							Computed: true,
						},
						"default_disabled": schema.BoolAttribute{
							Computed: true,
						},
						"default_is_internal": schema.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (r *rolePermissionsOverrideResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*propelauth.PropelAuthClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *propelauth.PropelAuthClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *rolePermissionsOverrideResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan rolePermissionsOverrideResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

	// Create the role-permissions override
	// get the default roles and permissions to track changes and for constructing the full mapping
	defaultRolesAndPerms, err := r.client.GetRolesAndPermissions()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating PropelAuth Role-Permissions Override",
			"Could not read default PropelAuth Roles and Permissions: " + err.Error(),
		)
		return
	}

	updateDefaultPermissionsInState(&plan, defaultRolesAndPerms)
	overrideToCreate := constructOverrideFromState(&plan)

    mappingId, err := r.client.CreateRolePermissionsOverride(overrideToCreate)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Creating PropelAuth Role-Permissions Override",
            "Could not create the role-permissions override in PropelAuth, unexpected error: "+err.Error(),
        )
        return
    }

	// update the mapping id in the state
	plan.MappingId = types.StringPointerValue(mappingId)

	// log that the resource was created
	tflog.Trace(ctx, "created a propelauth_role_permissions_override resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *rolePermissionsOverrideResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state and read it into the model
	var state rolePermissionsOverrideResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the roles and permissions from PropelAuth
	_, err := r.client.GetRolePermissionsOverride(state.MappingId.ValueString())
	if err != nil {
        resp.Diagnostics.AddError(
            "Error Reading PropelAuth Role-Permissions Override",
            "Could not read PropelAuth Role-Permissions Override: " + err.Error(),
        )
        return
    }

	// update state
	// TKTK: update state with the retrieved roles-permissions override

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *rolePermissionsOverrideResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan rolePermissionsOverrideResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the roles and permissions
	// get the default roles and permissions to track changes and for constructing the full mapping
	defaultRolesAndPerms, err := r.client.GetRolesAndPermissions()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating PropelAuth Role-Permissions Override",
			"Could not read default PropelAuth Roles and Permissions: " + err.Error(),
		)
		return
	}

	updateDefaultPermissionsInState(&plan, defaultRolesAndPerms)
	overrideUpdate := constructOverrideFromState(&plan)

    _, err = r.client.UpdateRolePermissionsOverride(plan.MappingId.ValueString(), overrideUpdate)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Updating PropelAuth Role-Permissions Override",
            "Could not create the role-permissions override in PropelAuth, unexpected error: "+err.Error(),
        )
        return
    }

	tflog.Trace(ctx, "updated a propelauth_role_permissions_override resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *rolePermissionsOverrideResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state and read it into the model
	var state rolePermissionsOverrideResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the role-permissions override
	err := r.client.DeleteRolePermissionsOverride(state.MappingId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting PropelAuth Role-Permissions Override",
			"Could not delete the role-permissions override in PropelAuth, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted a propelauth_role_permissions_override resource")
}

func updateDefaultPermissionsInState(state *rolePermissionsOverrideResourceModel, defaultRolesAndPerms *propelauth.RolesAndPermissions) {
	for _, roleFromSourceDefaults := range defaultRolesAndPerms.Roles {
		roleInState, ok := state.Roles[roleFromSourceDefaults.Name]
		if !ok {
			state.Roles[roleFromSourceDefaults.Name] = roleOverrideModel{
				DefaultCanViewOtherMembers: types.BoolValue(roleFromSourceDefaults.CanViewOtherMembers),
				DefaultCanInvite: types.BoolValue(roleFromSourceDefaults.CanInvite),
				DefaultCanChangeRoles: types.BoolValue(roleFromSourceDefaults.CanChangeRoles),
				DefaultCanManageApiKeys: types.BoolValue(roleFromSourceDefaults.CanManageApiKeys),
				DefaultCanRemoveUsers: types.BoolValue(roleFromSourceDefaults.CanRemoveUsers),
				DefaultCanSetupSaml: types.BoolValue(roleFromSourceDefaults.CanSetupSaml),
				DefaultCanDeleteOrg: types.BoolValue(roleFromSourceDefaults.CanDeleteOrg),
				DefaultCanEditOrgAccess: types.BoolValue(roleFromSourceDefaults.CanEditOrgAccess),
				DefaultCanUpdateOrgMetadata: types.BoolValue(roleFromSourceDefaults.CanUpdateOrgMetadata),
				DefaultPermissions: convertArrayOfStringsForState(roleFromSourceDefaults.ExternalPermissions),
				DefaultRolesCanManage: convertArrayOfStringsForState(roleFromSourceDefaults.RolesCanManage),
				DefaultDisabled: types.BoolValue(roleFromSourceDefaults.Disabled),
				DefaultIsInternal: types.BoolValue(!roleFromSourceDefaults.IsVisibleToEndUser),
			}
		} else {
			roleInState.DefaultCanViewOtherMembers = types.BoolValue(roleFromSourceDefaults.CanViewOtherMembers)
			roleInState.DefaultCanInvite = types.BoolValue(roleFromSourceDefaults.CanInvite)
			roleInState.DefaultCanChangeRoles = types.BoolValue(roleFromSourceDefaults.CanChangeRoles)
			roleInState.DefaultCanManageApiKeys = types.BoolValue(roleFromSourceDefaults.CanManageApiKeys)
			roleInState.DefaultCanRemoveUsers = types.BoolValue(roleFromSourceDefaults.CanRemoveUsers)
			roleInState.DefaultCanSetupSaml = types.BoolValue(roleFromSourceDefaults.CanSetupSaml)
			roleInState.DefaultCanDeleteOrg = types.BoolValue(roleFromSourceDefaults.CanDeleteOrg)
			roleInState.DefaultCanEditOrgAccess = types.BoolValue(roleFromSourceDefaults.CanEditOrgAccess)
			roleInState.DefaultCanUpdateOrgMetadata = types.BoolValue(roleFromSourceDefaults.CanUpdateOrgMetadata)
			if !arraysMatchIgnoreOrder(roleInState.DefaultPermissions, roleFromSourceDefaults.ExternalPermissions) {
				defaultExternalPermissions := convertArrayOfStringsForState(roleFromSourceDefaults.ExternalPermissions)
				roleInState.DefaultPermissions = defaultExternalPermissions
			}
			if !arraysMatchIgnoreOrder(roleInState.DefaultRolesCanManage, roleFromSourceDefaults.RolesCanManage) {
				defaultRolesCanManage := convertArrayOfStringsForState(roleFromSourceDefaults.RolesCanManage)
				roleInState.DefaultRolesCanManage = defaultRolesCanManage
			}
			roleInState.DefaultDisabled = types.BoolValue(roleFromSourceDefaults.Disabled)
			roleInState.DefaultIsInternal = types.BoolValue(!roleFromSourceDefaults.IsVisibleToEndUser)

			state.Roles[roleFromSourceDefaults.Name] = roleInState
		}
	}
}

func constructOverrideFromState(state *rolePermissionsOverrideResourceModel) propelauth.RolePermissionsOverride {
	rolesOverrideForSource := propelauth.RolePermissionsOverride{
		Name: state.Name.ValueString(),
		Roles: []propelauth.RoleDefinition{},
	}

	for roleName, roleData := range state.Roles {
		roleForSource := propelauth.RoleDefinition{
			Name: roleName,
			CanViewOtherMembers: roleData.DefaultCanViewOtherMembers.ValueBool(),
			CanInvite: roleData.DefaultCanInvite.ValueBool(),
			CanChangeRoles: roleData.DefaultCanChangeRoles.ValueBool(),
			CanManageApiKeys: roleData.DefaultCanManageApiKeys.ValueBool(),
			CanRemoveUsers: roleData.DefaultCanRemoveUsers.ValueBool(),
			CanSetupSaml: roleData.DefaultCanSetupSaml.ValueBool(),
			CanDeleteOrg: roleData.DefaultCanDeleteOrg.ValueBool(),
			CanEditOrgAccess: roleData.DefaultCanEditOrgAccess.ValueBool(),
			CanUpdateOrgMetadata: roleData.DefaultCanUpdateOrgMetadata.ValueBool(),
			ExternalPermissions: convertArrayOfStringsForSource(roleData.DefaultPermissions),
			RolesCanManage: convertArrayOfStringsForSource(roleData.DefaultRolesCanManage),
			Disabled: roleData.DefaultDisabled.ValueBool(),
			IsVisibleToEndUser: !roleData.DefaultIsInternal.ValueBool(),
		}
		if !roleData.CanViewOtherMembers.IsNull() {
			roleForSource.CanViewOtherMembers = roleData.CanViewOtherMembers.ValueBool()
		}
		if !roleData.CanInvite.IsNull() {
			roleForSource.CanInvite = roleData.CanInvite.ValueBool()
		}
		if !roleData.CanChangeRoles.IsNull() {
			roleForSource.CanChangeRoles = roleData.CanChangeRoles.ValueBool()
		}
		if !roleData.CanManageApiKeys.IsNull() {
			roleForSource.CanManageApiKeys = roleData.CanManageApiKeys.ValueBool()
		}
		if !roleData.CanRemoveUsers.IsNull() {
			roleForSource.CanRemoveUsers = roleData.CanRemoveUsers.ValueBool()
		}
		if !roleData.CanSetupSaml.IsNull() {
			roleForSource.CanSetupSaml = roleData.CanSetupSaml.ValueBool()
		}
		if !roleData.CanDeleteOrg.IsNull() {
			roleForSource.CanDeleteOrg = roleData.CanDeleteOrg.ValueBool()
		}
		if !roleData.CanEditOrgAccess.IsNull() {
			roleForSource.CanEditOrgAccess = roleData.CanEditOrgAccess.ValueBool()
		}
		if !roleData.CanUpdateOrgMetadata.IsNull() {
			roleForSource.CanUpdateOrgMetadata = roleData.CanUpdateOrgMetadata.ValueBool()
		}
		if !roleData.Disabled.IsNull() {
			roleForSource.Disabled = roleData.Disabled.ValueBool()
		}
		convertedExternalPermissions := convertArrayOfStringsForSource(roleData.Permissions)
		if !propelauth.Contains(convertedExternalPermissions, "propelauth::null_permission") {
			roleForSource.ExternalPermissions = convertedExternalPermissions
		}
		rolesOverrideForSource.Roles = append(rolesOverrideForSource.Roles, roleForSource)
	}

	return rolesOverrideForSource
}
