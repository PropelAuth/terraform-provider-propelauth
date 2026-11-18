package provider

import (
	"context"
	"fmt"
	"slices"
	"sort"

	"terraform-provider-propelauth/internal/propelauth"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &customRoleMappingResource{}
var _ resource.ResourceWithConfigure = &customRoleMappingResource{}
var _ resource.ResourceWithImportState = &customRoleMappingResource{}

func NewCustomRoleMappingResource() resource.Resource {
	return &customRoleMappingResource{}
}

// customRoleMappingResource defines the resource implementation.
type customRoleMappingResource struct {
	client *propelauth.PropelAuthClient
}

// customRoleMappingResourceModel describes the resource data model.
type customRoleMappingResourceModel struct {
	Id    types.String                      `tfsdk:"id"`
	Name  types.String                      `tfsdk:"name"`
	Roles map[string]customMappingRoleModel `tfsdk:"roles"`
}

type customMappingRoleModel struct {
	CanViewOtherMembers     types.Bool   `tfsdk:"can_view_other_members"`
	CanInvite               types.Bool   `tfsdk:"can_invite"`
	CanChangeRoles          types.Bool   `tfsdk:"can_change_roles"`
	CanManageApiKeys        types.Bool   `tfsdk:"can_manage_api_keys"`
	CanRemoveUsers          types.Bool   `tfsdk:"can_remove_users"`
	CanSetupSaml            types.Bool   `tfsdk:"can_setup_saml"`
	CanDeleteOrg            types.Bool   `tfsdk:"can_delete_org"`
	CanEditOrgAccess        types.Bool   `tfsdk:"can_edit_org_access"`
	CanUpdateOrgMetadata    types.Bool   `tfsdk:"can_update_org_metadata"`
	CanViewOrgAuditLog      types.Bool   `tfsdk:"can_view_org_audit_log"`
	CanManageMcpConnections types.Bool   `tfsdk:"can_manage_mcp_connections"`
	Permissions             types.List   `tfsdk:"permissions"`
	RolesCanManage          types.List   `tfsdk:"roles_can_manage"`
	McpOrgScopes            types.List   `tfsdk:"mcp_org_scopes"`
	IsInternal              types.Bool   `tfsdk:"is_internal"`
	Disabled                types.Bool   `tfsdk:"disabled"`
	Deprecated              types.Bool   `tfsdk:"deprecated"`
	Description             types.String `tfsdk:"description"`
	ShowDescriptionTooltip  types.Bool   `tfsdk:"show_description_tooltip"`
}

func (r *customRoleMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_role_mapping"
}

func (r *customRoleMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Role attributes are Optional+Computed without defaults: the API fills in
	// values it owns (descriptions, role hierarchy, visibility) when config omits
	// them, and the state merge keeps them so plans don't churn.
	roleBool := func(description string) schema.BoolAttribute {
		return schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: description,
		}
	}
	stringList := func(description string) schema.ListAttribute {
		return schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Description: description,
		}
	}

	resp.Schema = schema.Schema{
		Description: "Custom Role Mapping resource. This is for configuring the roles in a non-default " +
			"role mapping. The default role mapping is managed by the propelauth_roles_and_permissions resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The mapping ID of the custom role mapping.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the custom role mapping.",
			},
			"roles": schema.MapNestedAttribute{
				Required:    true,
				Description: "A map of roles that can be assigned to members of an organization subscribed to this mapping.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"can_view_other_members": roleBool(
							"If true, users with this role in the org can view other members of the organization."),
						"can_invite": roleBool(
							"If true, users with this role in the org can invite other users to the organization."),
						"can_change_roles": roleBool(
							"If true, users with this role in the org can change the roles of other users in the organization."),
						"can_manage_api_keys": roleBool(
							"If true, users with this role in the org can manage API keys for the organization."),
						"can_remove_users": roleBool(
							"If true, users with this role in the org can remove other users from the organization."),
						"can_setup_saml": roleBool(
							"If true, users with this role in the org can setup enterprise SSO for the organization."),
						"can_delete_org": roleBool(
							"If true, users with this role in the org can delete the organization."),
						"can_edit_org_access": roleBool(
							"If true, users with this role in the org can edit the organization's access settings."),
						"can_update_org_metadata": roleBool(
							"If true, users with this role in the org can update the organization's metadata."),
						"can_view_org_audit_log": roleBool(
							"If true, users with this role in the org can view the organization's audit log."),
						"can_manage_mcp_connections": roleBool(
							"If true, users with this role in the org can manage the organization's MCP connections."),
						"permissions": stringList(
							"A list of permissions specific to your application that are assigned to this role."),
						"roles_can_manage": stringList(
							"A list of roles that this role can manage in a multi-role project."),
						"mcp_org_scopes": stringList(
							"A list of MCP scopes assigned to this role."),
						"is_internal": roleBool(
							"If true, this role is an internal role and cannot be assigned to or viewed by end users."),
						"disabled": roleBool(
							"If true, this role is disabled and cannot be assigned to users."),
						"deprecated": roleBool(
							"If true, this role is deprecated."),
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "A human-readable description of the role, inherited from the project role.",
						},
						"show_description_tooltip": schema.BoolAttribute{
							Computed:    true,
							Description: "If true, the role description is shown to end users as a tooltip. Inherited from the project role.",
						},
					},
				},
			},
		},
	}
}

func (r *customRoleMappingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *customRoleMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan customRoleMappingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapping, err := r.client.CreateCustomRoleMapping(propelauth.CustomRoleMapping{
		Name:  plan.Name.ValueString(),
		Roles: customMappingRolesFromState(ctx, plan.Roles),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating PropelAuth Custom Role Mapping",
			"Could not create custom role mapping, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(mapping.MappingId)
	plan.Roles = customMappingRolesToState(ctx, plan.Roles, mapping.Roles)

	tflog.Trace(ctx, "created a propelauth_custom_role_mapping resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *customRoleMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state customRoleMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapping, err := r.client.GetCustomRoleMapping(state.Id.ValueString())
	if err != nil {
		if propelauth.IsPropelAuthNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading PropelAuth Custom Role Mapping",
			"Could not read custom role mapping, unexpected error: "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(mapping.Name)
	state.Roles = customMappingRolesToState(ctx, state.Roles, mapping.Roles)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *customRoleMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan customRoleMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapping, err := r.client.UpdateCustomRoleMapping(plan.Id.ValueString(), propelauth.CustomRoleMapping{
		Name:  plan.Name.ValueString(),
		Roles: customMappingRolesFromState(ctx, plan.Roles),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating PropelAuth Custom Role Mapping",
			"Could not update custom role mapping, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Roles = customMappingRolesToState(ctx, plan.Roles, mapping.Roles)

	tflog.Trace(ctx, "updated a propelauth_custom_role_mapping resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *customRoleMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state customRoleMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteCustomRoleMapping(state.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting PropelAuth Custom Role Mapping",
			"Could not delete custom role mapping, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted a propelauth_custom_role_mapping resource")
}

func (r *customRoleMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func customMappingRolesFromState(ctx context.Context, roles map[string]customMappingRoleModel) []propelauth.CustomRoleMappingRole {
	result := make([]propelauth.CustomRoleMappingRole, 0, len(roles))
	for roleName, role := range roles {
		result = append(result, propelauth.CustomRoleMappingRole{
			Name:                    roleName,
			CanViewOtherMembers:     role.CanViewOtherMembers.ValueBool(),
			CanInvite:               role.CanInvite.ValueBool(),
			CanChangeRoles:          role.CanChangeRoles.ValueBool(),
			CanManageApiKeys:        role.CanManageApiKeys.ValueBool(),
			CanRemoveUsers:          role.CanRemoveUsers.ValueBool(),
			CanSetupSaml:            role.CanSetupSaml.ValueBool(),
			CanDeleteOrg:            role.CanDeleteOrg.ValueBool(),
			CanEditOrgAccess:        role.CanEditOrgAccess.ValueBool(),
			CanUpdateOrgMetadata:    role.CanUpdateOrgMetadata.ValueBool(),
			CanViewOrgAuditLog:      role.CanViewOrgAuditLog.ValueBool(),
			CanManageMcpConnections: role.CanManageMcpConnections.ValueBool(),
			RolesCanManage:          stringSliceFromList(ctx, role.RolesCanManage),
			Disabled:                role.Disabled.ValueBool(),
			Deprecated:              role.Deprecated.ValueBool(),
			Description:             role.Description.ValueStringPointer(),
			ShowDescriptionTooltip:  role.ShowDescriptionTooltip.ValueBoolPointer(),
			ExternalPermissions:     stringSliceFromList(ctx, role.Permissions),
			McpOrgScopes:            stringSliceFromList(ctx, role.McpOrgScopes),
			IsVisibleToEndUser:      !role.IsInternal.ValueBool(),
		})
	}
	return result
}

// customMappingRolesToState builds the roles map for state from the API
// response. Roles already in prior (plan or state) keep their list order when
// the contents match, since the API re-sorts roles_can_manage and similar
// fields on every write.
func customMappingRolesToState(ctx context.Context, prior map[string]customMappingRoleModel, roles []propelauth.CustomRoleMappingRole) map[string]customMappingRoleModel {
	result := make(map[string]customMappingRoleModel, len(roles))
	for i := range roles {
		result[roles[i].Name] = mergeCustomMappingRole(ctx, prior, &roles[i])
	}
	return result
}

// mergeCustomMappingRole merges the source role into the prior role entry,
// keeping list fields the API re-sorts (roles_can_manage, permissions,
// mcp_org_scopes) in prior order when the contents match so plans don't churn.
func mergeCustomMappingRole(ctx context.Context, prior map[string]customMappingRoleModel, role *propelauth.CustomRoleMappingRole) customMappingRoleModel {
	roleInState, ok := prior[role.Name]
	if !ok {
		return customMappingRoleToState(ctx, role)
	}

	roleInState.CanViewOtherMembers = types.BoolValue(role.CanViewOtherMembers)
	roleInState.CanInvite = types.BoolValue(role.CanInvite)
	roleInState.CanChangeRoles = types.BoolValue(role.CanChangeRoles)
	roleInState.CanManageApiKeys = types.BoolValue(role.CanManageApiKeys)
	roleInState.CanRemoveUsers = types.BoolValue(role.CanRemoveUsers)
	roleInState.CanSetupSaml = types.BoolValue(role.CanSetupSaml)
	roleInState.CanDeleteOrg = types.BoolValue(role.CanDeleteOrg)
	roleInState.CanEditOrgAccess = types.BoolValue(role.CanEditOrgAccess)
	roleInState.CanUpdateOrgMetadata = types.BoolValue(role.CanUpdateOrgMetadata)
	roleInState.CanViewOrgAuditLog = types.BoolValue(role.CanViewOrgAuditLog)
	roleInState.CanManageMcpConnections = types.BoolValue(role.CanManageMcpConnections)
	if !listMatchesIgnoreOrder(ctx, roleInState.Permissions, role.ExternalPermissions) {
		roleInState.Permissions = listValueFromStrings(ctx, role.ExternalPermissions)
	}
	if !listMatchesIgnoreOrder(ctx, roleInState.RolesCanManage, role.RolesCanManage) {
		roleInState.RolesCanManage = listValueFromStrings(ctx, role.RolesCanManage)
	}
	if !listMatchesIgnoreOrder(ctx, roleInState.McpOrgScopes, role.McpOrgScopes) {
		roleInState.McpOrgScopes = listValueFromStrings(ctx, role.McpOrgScopes)
	}
	roleInState.IsInternal = types.BoolValue(!role.IsVisibleToEndUser)
	roleInState.Disabled = types.BoolValue(role.Disabled)
	roleInState.Deprecated = types.BoolValue(role.Deprecated)
	roleInState.Description = types.StringPointerValue(role.Description)
	roleInState.ShowDescriptionTooltip = types.BoolPointerValue(role.ShowDescriptionTooltip)
	return roleInState
}

func customMappingRoleToState(ctx context.Context, role *propelauth.CustomRoleMappingRole) customMappingRoleModel {
	return customMappingRoleModel{
		CanViewOtherMembers:     types.BoolValue(role.CanViewOtherMembers),
		CanInvite:               types.BoolValue(role.CanInvite),
		CanChangeRoles:          types.BoolValue(role.CanChangeRoles),
		CanManageApiKeys:        types.BoolValue(role.CanManageApiKeys),
		CanRemoveUsers:          types.BoolValue(role.CanRemoveUsers),
		CanSetupSaml:            types.BoolValue(role.CanSetupSaml),
		CanDeleteOrg:            types.BoolValue(role.CanDeleteOrg),
		CanEditOrgAccess:        types.BoolValue(role.CanEditOrgAccess),
		CanUpdateOrgMetadata:    types.BoolValue(role.CanUpdateOrgMetadata),
		CanViewOrgAuditLog:      types.BoolValue(role.CanViewOrgAuditLog),
		CanManageMcpConnections: types.BoolValue(role.CanManageMcpConnections),
		Permissions:             listValueFromStrings(ctx, role.ExternalPermissions),
		RolesCanManage:          listValueFromStrings(ctx, role.RolesCanManage),
		McpOrgScopes:            listValueFromStrings(ctx, role.McpOrgScopes),
		IsInternal:              types.BoolValue(!role.IsVisibleToEndUser),
		Disabled:                types.BoolValue(role.Disabled),
		Deprecated:              types.BoolValue(role.Deprecated),
		Description:             types.StringPointerValue(role.Description),
		ShowDescriptionTooltip:  types.BoolPointerValue(role.ShowDescriptionTooltip),
	}
}

// stringSliceFromList returns an empty slice rather than nil for unset lists,
// since the API rejects null in place of an array.
func stringSliceFromList(ctx context.Context, list types.List) []string {
	result := []string{}
	if list.IsNull() || list.IsUnknown() {
		return result
	}
	list.ElementsAs(ctx, &result, false)
	return result
}

func listValueFromStrings(ctx context.Context, values []string) types.List {
	list, diags := types.ListValueFrom(ctx, types.StringType, values)
	if diags.HasError() {
		return types.ListNull(types.StringType)
	}
	return list
}

func listMatchesIgnoreOrder(ctx context.Context, listInPrior types.List, valuesFromSource []string) bool {
	if listInPrior.IsNull() || listInPrior.IsUnknown() {
		return false
	}
	prior := stringSliceFromList(ctx, listInPrior)
	source := append([]string{}, valuesFromSource...)
	sort.Strings(prior)
	sort.Strings(source)
	return slices.Equal(prior, source)
}
