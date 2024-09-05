package provider

import (
	"context"
	"fmt"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &rolesAndPermissionsResource{}
var _ resource.ResourceWithConfigure = &rolesAndPermissionsResource{}
var _ resource.ResourceWithValidateConfig = &rolesAndPermissionsResource{}

func NewRolesAndPermissionsResource() resource.Resource {
	return &rolesAndPermissionsResource{}
}

// rolesAndPermissionsResource defines the resource implementation.
type rolesAndPermissionsResource struct {
	client *propelauth.PropelAuthClient
}

// rolesAndPermissionsResourceModel describes the resource data model.
type rolesAndPermissionsResourceModel struct {
	MultipleRolesPerUser types.Bool           `tfsdk:"multiple_roles_per_user"`
	Permissions          []permissionModel    `tfsdk:"permissions"`
	Roles                map[string]roleModel `tfsdk:"roles"`
	RoleHierarchy        []types.String       `tfsdk:"role_hierarchy"`
	DefaultRole          types.String         `tfsdk:"default_role"`
	DefaultOwnerRole     types.String         `tfsdk:"default_owner_role"`
}

type permissionModel struct {
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
}

type roleModel struct {
	CanViewOtherMembers  types.Bool     `tfsdk:"can_view_other_members"`
	CanInvite            types.Bool     `tfsdk:"can_invite"`
	CanChangeRoles       types.Bool     `tfsdk:"can_change_roles"`
	CanManageApiKeys     types.Bool     `tfsdk:"can_manage_api_keys"`
	CanRemoveUsers       types.Bool     `tfsdk:"can_remove_users"`
	CanSetupSaml         types.Bool     `tfsdk:"can_setup_saml"`
	CanDeleteOrg         types.Bool     `tfsdk:"can_delete_org"`
	CanEditOrgAccess     types.Bool     `tfsdk:"can_edit_org_access"`
	CanUpdateOrgMetadata types.Bool     `tfsdk:"can_update_org_metadata"`
	Permissions          []types.String `tfsdk:"permissions"`
	RolesCanManage       []types.String `tfsdk:"roles_can_manage"`
	IsInternal           types.Bool     `tfsdk:"is_internal"`
	Disabled             types.Bool     `tfsdk:"disabled"`
	Description          types.String   `tfsdk:"description"`
	ReplacingRole        types.String   `tfsdk:"replacing_role"`
}

func (r *rolesAndPermissionsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_roles_and_permissions"
}

func (r *rolesAndPermissionsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Roles and Permissions resource. This is for configuring the basic roles and permissions information in PropelAuth.",
		Attributes: map[string]schema.Attribute{
			"multiple_roles_per_user": schema.BoolAttribute{
				Computed: true,
				Description: "If true, than each member of an organization can have multiple roles and their is no hierarchy between roles. " +
					"Instead, the relationship between roles is defined by the `roles_can_manage` field on each individual role definition. " +
					"A single-role project can be migrated to multi-role, but not the other way around. Because of this, " +
					"this can only be set in the PropelAuth dashboard.",
			},
			"permissions": schema.ListNestedAttribute{
				Optional:    true,
				Description: "A list of permissions that are specific to your application and can be assigned to individual roles.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The name of the permission. This should be a unique identifier for the permission.",
						},
						"display_name": schema.StringAttribute{
							Optional: true,
							Description: "The display name of the permission. This is the human readable name of the permission. " +
								"If not provided, the `name` will be used.",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "A description of the permission. This is a human readable description of what the permission allows.",
						},
					},
				},
			},
			"default_role": schema.StringAttribute{
				Required: true,
				Description: "The `default_role` is the role assigned to a user if they join an organization and no other role is assigned to them. " +
					"It is also the fallback role in the instance their role is deleted from the configuration without a replacement.",
			},
			"default_owner_role": schema.StringAttribute{
				Required:    true,
				Description: "The `default_owner_role` is the role automatically assigned to the user who creates the organization.",
			},
			"roles": schema.MapNestedAttribute{
				Required: true,
				Description: "A map of roles that can be assigned to members of an organization. This provides the the permissions " +
					"in the default mapping. For overrides (ie custom mappings) that can be applied on top of this, see the " +
					"`propelauth_role_permissions_override` resource. in the default mapping. For overrides (ie custom mappings) " +
					"that can be applied on top of this, see the `propelauth_role_permissions_override` resource.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"can_view_other_members": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(true),
							Description: "If true, users with this role in the org can view other members of the organization. " +
								"The default is true.",
						},
						"can_invite": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
							Description: "If true, users with this role in the org can invite other users to the organization. " +
								"The default is false.",
						},
						"can_change_roles": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
							Description: "If true, users with this role in the org can change the roles of other users in the organization. " +
								"The default is false.",
						},
						"can_manage_api_keys": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
							Description: "If true, users with this role in the org can manage API keys for the organization. " +
								"The default is false.",
						},
						"can_remove_users": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
							Description: "If true, users with this role in the org can remove other users from the organization. " +
								"The default is false.",
						},
						"can_setup_saml": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
							Description: "If true, users with this role in the org can setup enterprise SSO for the organization. " +
								"The default is false.",
						},
						"can_delete_org": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
							Description: "If true, users with this role in the org can delete the organization. " +
								"The default is false.",
						},
						"can_edit_org_access": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
							Description: "If true, users with this role in the org can edit the organization's access settings. " +
								"These settings incluede what email domains are included and whether 2FA is enforced for the org. " +
								"The default is false.",
						},
						"can_update_org_metadata": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
							Description: "If true, users with this role in the org can update the organization's metadata. " +
								"This includes changing the name of the organization. The default is false.",
						},
						"permissions": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
							Default: listdefault.StaticValue(types.ListValueMust(
								types.StringType,
								[]attr.Value{},
							)),
							Description: "A list of permissions specific to your application that are assigned to this role.",
						},
						"roles_can_manage": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
							Default: listdefault.StaticValue(types.ListValueMust(
								types.StringType,
								[]attr.Value{},
							)),
							Description: "A list of roles that this role can manage. This is only relevant if `multiple_roles_per_user` " +
								"is true. If `multiple_roles_per_user` is false, the other roles that a role can manage is defined by " +
								"the order in `role_hierarchy` where the first role is able to manage every other role including itself.",
						},
						"is_internal": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
							Description: "If true, this role is an internal role and cannot be assigned to or viewed by end users. " +
								"The default is false.",
						},
						"disabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
							Description: "If true, this role is disabled and cannot be assigned to users. It is only useful if you " +
								"intend to use the role in non-default role mappings exclusively. The default is false.",
						},
						"replacing_role": schema.StringAttribute{
							Optional: true,
							Description: "The name of a role that no longer exists but this role is replacing. This should only be used " +
								"if you are attempting to change the name of an existing role and want to ensure that users who had the old role " +
								"now have this role. The `replacing_role` should not exist in the `roles` map.",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "A human-readable description of the role.",
						},
					},
				},
			},
			"role_hierarchy": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default: listdefault.StaticValue(types.ListValueMust(
					types.StringType,
					[]attr.Value{},
				)),
				Description: "A list of roles in order of hierarchy. The first role in the list is the highest role and " +
					"the last role is the lowest role. This is only relevant if `multiple_roles_per_user` is false. " +
					"If `multiple_roles_per_user` is true, the roles that a role can manage is defined by the `roles_can_manage` " +
					"field on each individual role definition.",
			},
		},
	}
}

func (r *rolesAndPermissionsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var plan rolesAndPermissionsResourceModel

	// Read Terraform plan data into the model
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the plan data

	// Verify that all roles in the hierarchy are defined
	for _, roleName := range plan.RoleHierarchy {
		if _, ok := plan.Roles[roleName.ValueString()]; !ok {
			resp.Diagnostics.AddAttributeError(
				path.Root("role_hierarchy"),
				"PropelAuth Role in hierarchy not defined",
				fmt.Sprintf("Role %s is in the role hierarchy but is not defined in the roles map.", roleName.ValueString()),
			)
			return
		}
		if len(plan.RoleHierarchy) < len(plan.Roles) {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("role_hierarchy"),
				"Not all defined roles are included in the hierarchy",
				"If multiple_roles_per_user = true, you can ignore this. If not, any roles excluded from the hierarchy "+
					"will be omitted in the source system.",
			)
			return
		}
	}
}

func (r *rolesAndPermissionsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *rolesAndPermissionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan rolesAndPermissionsResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the roles and permissions
	updateBuilder := propelauth.NewRolesAndPermissionsUpdateBuilder()

	updateBuilder = updateBuilder.
		SetDefaultRole(plan.DefaultRole.ValueString()).
		SetDefaultOwnerRole(plan.DefaultOwnerRole.ValueString())

	for _, permission := range plan.Permissions {
		updateBuilder = updateBuilder.InsertPermission(propelauth.Permission{
			Name:        permission.Name.ValueString(),
			DisplayName: permission.DisplayName.ValueStringPointer(),
			Description: permission.Description.ValueStringPointer(),
		})
	}

	for roleName, role := range plan.Roles {
		updateBuilder = updateBuilder.InsertRole(roleName, convertRoleFromState(roleName, &role))
		if role.ReplacingRole.ValueString() != "" {
			updateBuilder = updateBuilder.InsertOldToNewRoleMapping(role.ReplacingRole.ValueString(), roleName)
		}
	}

	updateBuilder.SetRoleHierarchy(convertArrayOfStringsForSource(plan.RoleHierarchy))

	// get the old roles and permissions to track changes/deletions in role names
	oldRolesAndPermissions, err := r.client.GetRolesAndPermissions()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating PropelAuth Roles and Permissions",
			"Could not read old PropelAuth Roles and Permissions: "+err.Error(),
		)
		return
	}

	updateBuilder.SetMultipleRolesPerUser(oldRolesAndPermissions.IsMultiRole())
	plan.MultipleRolesPerUser = types.BoolValue(oldRolesAndPermissions.IsMultiRole())

	for _, oldRole := range oldRolesAndPermissions.Roles {
		updateBuilder.InsertOldRoleName(oldRole.Name)
	}

	_, err = r.client.UpdateRolesAndPermissions(updateBuilder.Build())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting roles and permissions",
			"Could not set roles and permissions, unexpected error: "+err.Error(),
		)
		return
	}

	// log that the resource was created
	tflog.Trace(ctx, "created a propelauth_roles_and_permissions resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *rolesAndPermissionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state and read it into the model
	var state rolesAndPermissionsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the roles and permissions from PropelAuth
	rolesAndPermissions, err := r.client.GetRolesAndPermissions()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth Roles and Permissions",
			"Could not read PropelAuth Roles and Permissions: "+err.Error(),
		)
		return
	}

	// update state
	// easy ones first
	state.MultipleRolesPerUser = types.BoolValue(rolesAndPermissions.IsMultiRole())
	state.DefaultRole = types.StringValue(rolesAndPermissions.DefaultRole)
	state.DefaultOwnerRole = types.StringValue(rolesAndPermissions.DefaultOwnerRole)
	// role definitions
	for _, role := range rolesAndPermissions.Roles {
		updateStateForRole(&state, &role)
	}
	// permissions
	reconcilePermissions(&state, rolesAndPermissions)
	// role hierarchy
	sourceRoleHierarchy := rolesAndPermissions.GetHierarchy()
	if !state.MultipleRolesPerUser.ValueBool() && !arraysMatch(state.RoleHierarchy, sourceRoleHierarchy) {
		state.RoleHierarchy = convertArrayOfStringsForState(sourceRoleHierarchy)
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *rolesAndPermissionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan data into the model
	var plan rolesAndPermissionsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the roles and permissions
	// Update the roles and permissions
	updateBuilder := propelauth.NewRolesAndPermissionsUpdateBuilder()

	updateBuilder = updateBuilder.
		SetDefaultRole(plan.DefaultRole.ValueString()).
		SetDefaultOwnerRole(plan.DefaultOwnerRole.ValueString())

	for _, permission := range plan.Permissions {
		updateBuilder = updateBuilder.InsertPermission(propelauth.Permission{
			Name:        permission.Name.ValueString(),
			DisplayName: permission.DisplayName.ValueStringPointer(),
			Description: permission.Description.ValueStringPointer(),
		})
	}

	for roleName, role := range plan.Roles {
		updateBuilder = updateBuilder.InsertRole(roleName, convertRoleFromState(roleName, &role))
		if role.ReplacingRole.ValueString() != "" {
			updateBuilder = updateBuilder.InsertOldToNewRoleMapping(role.ReplacingRole.ValueString(), roleName)
		}
	}

	updateBuilder.SetRoleHierarchy(convertArrayOfStringsForSource(plan.RoleHierarchy))

	// get the old roles and permissions to track changes/deletions in role names
	oldRolesAndPermissions, err := r.client.GetRolesAndPermissions()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating PropelAuth Roles and Permissions",
			"Could not read current PropelAuth Roles and Permissions: "+err.Error(),
		)
		return
	}

	updateBuilder.SetMultipleRolesPerUser(oldRolesAndPermissions.IsMultiRole())
	plan.MultipleRolesPerUser = types.BoolValue(oldRolesAndPermissions.IsMultiRole())

	for _, oldRole := range oldRolesAndPermissions.Roles {
		updateBuilder.InsertOldRoleName(oldRole.Name)
	}

	_, err = r.client.UpdateRolesAndPermissions(updateBuilder.Build())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting roles and permissions",
			"Could not set roles and permissions, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "updated a propelauth_roles_and_permissions resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *rolesAndPermissionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "deleted a propelauth_roles_and_permissions resource")
}

func convertRoleFromState(roleName string, role *roleModel) propelauth.RoleDefinition {
	return propelauth.RoleDefinition{
		Name:                 roleName,
		CanViewOtherMembers:  role.CanViewOtherMembers.ValueBool(),
		CanInvite:            role.CanInvite.ValueBool(),
		CanChangeRoles:       role.CanChangeRoles.ValueBool(),
		CanManageApiKeys:     role.CanManageApiKeys.ValueBool(),
		CanRemoveUsers:       role.CanRemoveUsers.ValueBool(),
		CanSetupSaml:         role.CanSetupSaml.ValueBool(),
		CanDeleteOrg:         role.CanDeleteOrg.ValueBool(),
		CanEditOrgAccess:     role.CanEditOrgAccess.ValueBool(),
		CanUpdateOrgMetadata: role.CanUpdateOrgMetadata.ValueBool(),
		ExternalPermissions:  convertArrayOfStringsForSource(role.Permissions),
		RolesCanManage:       convertArrayOfStringsForSource(role.RolesCanManage),
		IsVisibleToEndUser:   !role.IsInternal.ValueBool(),
		Disabled:             role.Disabled.ValueBool(),
		Description:          role.Description.ValueStringPointer(),
	}
}

func updateStateForRole(state *rolesAndPermissionsResourceModel, role *propelauth.RoleDefinition) {
	roleInState, ok := state.Roles[role.Name]
	if ok {
		roleInState.CanViewOtherMembers = types.BoolValue(role.CanViewOtherMembers)
		roleInState.CanInvite = types.BoolValue(role.CanInvite)
		roleInState.CanChangeRoles = types.BoolValue(role.CanChangeRoles)
		roleInState.CanManageApiKeys = types.BoolValue(role.CanManageApiKeys)
		roleInState.CanRemoveUsers = types.BoolValue(role.CanRemoveUsers)
		roleInState.CanSetupSaml = types.BoolValue(role.CanSetupSaml)
		roleInState.CanDeleteOrg = types.BoolValue(role.CanDeleteOrg)
		roleInState.CanEditOrgAccess = types.BoolValue(role.CanEditOrgAccess)
		roleInState.CanUpdateOrgMetadata = types.BoolValue(role.CanUpdateOrgMetadata)
		if !arraysMatchIgnoreOrder(roleInState.Permissions, role.ExternalPermissions) {
			roleInState.Permissions = convertArrayOfStringsForState(role.ExternalPermissions)
		}
		if !arraysMatchIgnoreOrder(roleInState.RolesCanManage, role.RolesCanManage) {
			roleInState.RolesCanManage = convertArrayOfStringsForState(role.RolesCanManage)
		}
		roleInState.IsInternal = types.BoolValue(!role.IsVisibleToEndUser)
		roleInState.Disabled = types.BoolValue(role.Disabled)
		roleInState.Description = types.StringPointerValue(role.Description)
		state.Roles[role.Name] = roleInState
	} else {
		state.Roles[role.Name] = convertRoleToState(role)
	}
}

func arraysMatchIgnoreOrder(arrayInState []types.String, arrayFromSource []string) bool {
	if len(arrayInState) != len(arrayFromSource) {
		return false
	}

	for _, value := range arrayInState {
		if !propelauth.Contains(arrayFromSource, value.ValueString()) {
			return false
		}
	}

	return true
}

func arraysMatch(arrayInState []types.String, arrayFromSource []string) bool {
	if len(arrayInState) != len(arrayFromSource) {
		return false
	}

	for i := range arrayInState {
		if arrayInState[i].ValueString() != arrayFromSource[i] {
			return false
		}
	}

	return true
}

func convertArrayOfStringsForState(array []string) []types.String {
	result := make([]types.String, len(array))
	for i, value := range array {
		result[i] = types.StringValue(value)
	}
	return result
}

func convertArrayOfStringsForSource(array []types.String) []string {
	result := make([]string, len(array))
	for i, value := range array {
		result[i] = value.ValueString()
	}
	return result
}

func convertRoleToState(role *propelauth.RoleDefinition) roleModel {
	roleInState := roleModel{
		CanViewOtherMembers:  types.BoolValue(role.CanViewOtherMembers),
		CanInvite:            types.BoolValue(role.CanInvite),
		CanChangeRoles:       types.BoolValue(role.CanChangeRoles),
		CanManageApiKeys:     types.BoolValue(role.CanManageApiKeys),
		CanRemoveUsers:       types.BoolValue(role.CanRemoveUsers),
		CanSetupSaml:         types.BoolValue(role.CanSetupSaml),
		CanDeleteOrg:         types.BoolValue(role.CanDeleteOrg),
		CanEditOrgAccess:     types.BoolValue(role.CanEditOrgAccess),
		CanUpdateOrgMetadata: types.BoolValue(role.CanUpdateOrgMetadata),
		Permissions:          make([]types.String, len(role.ExternalPermissions)),
		RolesCanManage:       make([]types.String, len(role.RolesCanManage)),
		IsInternal:           types.BoolValue(!role.IsVisibleToEndUser),
		Disabled:             types.BoolValue(role.Disabled),
		Description:          types.StringPointerValue(role.Description),
	}

	for i, permission := range role.ExternalPermissions {
		roleInState.Permissions[i] = types.StringValue(permission)
	}

	for i, roleName := range role.RolesCanManage {
		roleInState.RolesCanManage[i] = types.StringValue(roleName)
	}

	return roleInState
}

func reconcilePermissions(state *rolesAndPermissionsResourceModel, rolesAndPermissions *propelauth.RolesAndPermissions) {
	for _, permissionInState := range state.Permissions {
		permission, ok := rolesAndPermissions.GetPermission(permissionInState.Name.ValueString())
		if ok {
			permissionInState.DisplayName = types.StringPointerValue(permission.DisplayName)
			permissionInState.Description = types.StringPointerValue(permission.Description)
		} else {
			permissionInState.Name = types.StringPointerValue(nil)
			permissionInState.DisplayName = types.StringPointerValue(nil)
			permissionInState.Description = types.StringPointerValue(nil)
		}
	}

	for _, permissionInRolesAndPermissions := range rolesAndPermissions.Permissions {
		exists := state.PermissionExists(permissionInRolesAndPermissions.Name)
		if !exists {
			state.Permissions = append(state.Permissions, permissionModel{
				Name:        types.StringValue(permissionInRolesAndPermissions.Name),
				DisplayName: types.StringPointerValue(permissionInRolesAndPermissions.DisplayName),
				Description: types.StringPointerValue(permissionInRolesAndPermissions.Description),
			})
		}
	}
}

func (r *rolesAndPermissionsResourceModel) PermissionExists(permissionName string) bool {
	for _, permission := range r.Permissions {
		if permission.Name.ValueString() == permissionName {
			return true
		}
	}
	return false
}
