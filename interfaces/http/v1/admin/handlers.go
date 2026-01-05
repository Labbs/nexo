package admin

import (
	"github.com/gofiber/fiber/v2"
	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/interfaces/http/v1/admin/dtos"
)

// checkAdmin verifies the user has admin role
func (ctrl *Controller) checkAdmin(ctx *fiber.Ctx) (*fiberoapi.AuthContext, *fiberoapi.ErrorResponse) {
	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: "Authentication required",
			Type:    "AUTHENTICATION_REQUIRED",
		}
	}

	// Get user to check role
	user, err := ctrl.UserApp.GetByUserId(struct{ UserId string }{UserId: authCtx.UserID})
	if err != nil {
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to retrieve user",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	if user.User.Role != domain.RoleAdmin {
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusForbidden,
			Details: "Admin access required",
			Type:    "FORBIDDEN",
		}
	}

	return authCtx, nil
}

// Users

func (ctrl *Controller) ListUsers(ctx *fiber.Ctx, req dtos.ListUsersRequest) (*dtos.ListUsersResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.list_users").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	// Default pagination
	limit := req.Limit
	offset := req.Offset
	if limit == 0 {
		limit = 50
	}

	users, total, err := ctrl.UserApp.GetAllUsers(limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get users")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to retrieve users",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	userItems := make([]dtos.UserItem, len(users))
	for i, u := range users {
		userItems[i] = dtos.UserItem{
			Id:        u.Id,
			Username:  u.Username,
			Email:     u.Email,
			AvatarUrl: u.AvatarUrl,
			Role:      string(u.Role),
			Active:    u.Active,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
	}

	return &dtos.ListUsersResponse{
		Users:      userItems,
		TotalCount: total,
	}, nil
}

func (ctrl *Controller) UpdateUserRole(ctx *fiber.Ctx, req dtos.UpdateUserRoleRequest) (*dtos.UpdateUserRoleResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.update_user_role").Logger()

	authCtx, errResp := ctrl.checkAdmin(ctx)
	if errResp != nil {
		return nil, errResp
	}

	// Prevent admin from changing their own role
	if req.UserId == authCtx.UserID {
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Details: "Cannot change your own role",
			Type:    "CANNOT_CHANGE_OWN_ROLE",
		}
	}

	role := domain.Role(req.Role)
	err := ctrl.UserApp.UpdateRole(req.UserId, role)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update user role")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to update user role",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.UpdateUserRoleResponse{
		Message: "User role updated successfully",
	}, nil
}

func (ctrl *Controller) UpdateUserActive(ctx *fiber.Ctx, req dtos.UpdateUserActiveRequest) (*dtos.UpdateUserActiveResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.update_user_active").Logger()

	authCtx, errResp := ctrl.checkAdmin(ctx)
	if errResp != nil {
		return nil, errResp
	}

	// Prevent admin from deactivating themselves
	if req.UserId == authCtx.UserID {
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Details: "Cannot deactivate your own account",
			Type:    "CANNOT_DEACTIVATE_SELF",
		}
	}

	err := ctrl.UserApp.UpdateActive(req.UserId, req.Active)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update user active status")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to update user status",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.UpdateUserActiveResponse{
		Message: "User status updated successfully",
	}, nil
}

func (ctrl *Controller) DeleteUser(ctx *fiber.Ctx, req dtos.DeleteUserRequest) (*dtos.DeleteUserResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.delete_user").Logger()

	authCtx, errResp := ctrl.checkAdmin(ctx)
	if errResp != nil {
		return nil, errResp
	}

	// Prevent admin from deleting themselves
	if req.UserId == authCtx.UserID {
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Details: "Cannot delete your own account",
			Type:    "CANNOT_DELETE_SELF",
		}
	}

	err := ctrl.UserApp.DeleteUser(req.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete user")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to delete user",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}

func (ctrl *Controller) InviteUser(ctx *fiber.Ctx, req dtos.InviteUserRequest) (*dtos.InviteUserResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.invite_user").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	// TODO: Implement user invitation (create user with temporary password and send email)
	// For now, just return a placeholder response
	logger.Info().Str("email", req.Email).Str("role", req.Role).Msg("user invitation requested")

	return &dtos.InviteUserResponse{
		Message: "User invitation feature coming soon",
	}, nil
}

// Spaces

func (ctrl *Controller) ListAllSpaces(ctx *fiber.Ctx, req dtos.ListAllSpacesRequest) (*dtos.ListAllSpacesResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.list_all_spaces").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	// Default pagination
	limit := req.Limit
	offset := req.Offset
	if limit == 0 {
		limit = 50
	}

	spaces, total, err := ctrl.SpaceApp.GetAllSpaces(limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get spaces")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to retrieve spaces",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	spaceItems := make([]dtos.SpaceItem, len(spaces))
	for i, s := range spaces {
		spaceItems[i] = dtos.SpaceItem{
			Id:          s.Id,
			Name:        s.Name,
			Description: "", // Add if available
			Icon:        s.Icon,
			Type:        string(s.Type),
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		}
		if s.OwnerId != nil {
			spaceItems[i].OwnerId = *s.OwnerId
		}
		if s.Owner != nil {
			spaceItems[i].OwnerName = s.Owner.Username
		}
	}

	return &dtos.ListAllSpacesResponse{
		Spaces:     spaceItems,
		TotalCount: total,
	}, nil
}

// API Keys

func (ctrl *Controller) ListAllApiKeys(ctx *fiber.Ctx, req dtos.ListAllApiKeysRequest) (*dtos.ListAllApiKeysResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.list_all_apikeys").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	// Default pagination
	limit := req.Limit
	offset := req.Offset
	if limit == 0 {
		limit = 50
	}

	apiKeys, total, err := ctrl.ApiKeyApp.GetAllApiKeys(limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get api keys")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to retrieve API keys",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	apiKeyItems := make([]dtos.ApiKeyItem, len(apiKeys))
	for i, k := range apiKeys {
		// Extract permissions as string array
		var permissions []string
		if k.Permissions != nil {
			if scopes, ok := k.Permissions["scopes"].([]interface{}); ok {
				for _, s := range scopes {
					if str, ok := s.(string); ok {
						permissions = append(permissions, str)
					}
				}
			}
		}

		apiKeyItems[i] = dtos.ApiKeyItem{
			Id:          k.Id,
			Name:        k.Name,
			KeyPrefix:   k.KeyPrefix,
			UserId:      k.UserId,
			Permissions: permissions,
			CreatedAt:   k.CreatedAt,
		}
		if k.User.Id != "" {
			apiKeyItems[i].Username = k.User.Username
		}
		if k.ExpiresAt != nil {
			apiKeyItems[i].ExpiresAt = *k.ExpiresAt
		}
		if k.LastUsedAt != nil {
			apiKeyItems[i].LastUsedAt = *k.LastUsedAt
		}
	}

	return &dtos.ListAllApiKeysResponse{
		ApiKeys:    apiKeyItems,
		TotalCount: total,
	}, nil
}

func (ctrl *Controller) RevokeApiKey(ctx *fiber.Ctx, req dtos.RevokeApiKeyRequest) (*dtos.RevokeApiKeyResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.revoke_apikey").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	err := ctrl.ApiKeyApp.AdminDeleteApiKey(req.ApiKeyId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to revoke api key")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to revoke API key",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.RevokeApiKeyResponse{
		Message: "API key revoked successfully",
	}, nil
}

// Groups

func (ctrl *Controller) ListGroups(ctx *fiber.Ctx, req dtos.ListGroupsRequest) (*dtos.ListGroupsResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.list_groups").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	// Default pagination
	limit := req.Limit
	offset := req.Offset
	if limit == 0 {
		limit = 50
	}

	groups, total, err := ctrl.GroupApp.GetAllGroups(limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get groups")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to retrieve groups",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	groupItems := make([]dtos.GroupItem, len(groups))
	for i, g := range groups {
		members := make([]dtos.MemberItem, len(g.Members))
		for j, m := range g.Members {
			members[j] = dtos.MemberItem{
				Id:        m.Id,
				Username:  m.Username,
				Email:     m.Email,
				AvatarUrl: m.AvatarUrl,
			}
		}

		groupItems[i] = dtos.GroupItem{
			Id:          g.Id,
			Name:        g.Name,
			Description: g.Description,
			Role:        string(g.Role),
			OwnerId:     g.OwnerId,
			MemberCount: len(g.Members),
			Members:     members,
			CreatedAt:   g.CreatedAt,
			UpdatedAt:   g.UpdatedAt,
		}
		if g.OwnerUser.Id != "" {
			groupItems[i].OwnerName = g.OwnerUser.Username
		}
	}

	return &dtos.ListGroupsResponse{
		Groups:     groupItems,
		TotalCount: total,
	}, nil
}

func (ctrl *Controller) CreateGroup(ctx *fiber.Ctx, req dtos.CreateGroupRequest) (*dtos.CreateGroupResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.create_group").Logger()

	authCtx, errResp := ctrl.checkAdmin(ctx)
	if errResp != nil {
		return nil, errResp
	}

	group, err := ctrl.GroupApp.CreateGroup(req.Name, req.Description, authCtx.UserID, domain.Role(req.Role))
	if err != nil {
		logger.Error().Err(err).Msg("failed to create group")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to create group",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.CreateGroupResponse{
		Id:      group.Id,
		Message: "Group created successfully",
	}, nil
}

func (ctrl *Controller) UpdateGroup(ctx *fiber.Ctx, req dtos.UpdateGroupRequest) (*dtos.UpdateGroupResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.update_group").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	err := ctrl.GroupApp.UpdateGroup(req.GroupId, req.Name, req.Description, domain.Role(req.Role))
	if err != nil {
		logger.Error().Err(err).Msg("failed to update group")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to update group",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.UpdateGroupResponse{
		Message: "Group updated successfully",
	}, nil
}

func (ctrl *Controller) DeleteGroup(ctx *fiber.Ctx, req dtos.DeleteGroupRequest) (*dtos.DeleteGroupResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.delete_group").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	err := ctrl.GroupApp.DeleteGroup(req.GroupId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete group")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to delete group",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.DeleteGroupResponse{
		Message: "Group deleted successfully",
	}, nil
}

func (ctrl *Controller) GetGroupMembers(ctx *fiber.Ctx, req dtos.GetGroupMembersRequest) (*dtos.GetGroupMembersResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.get_group_members").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	members, err := ctrl.GroupApp.GetMembers(req.GroupId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get group members")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to retrieve group members",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	memberItems := make([]dtos.MemberItem, len(members))
	for i, m := range members {
		memberItems[i] = dtos.MemberItem{
			Id:        m.Id,
			Username:  m.Username,
			Email:     m.Email,
			AvatarUrl: m.AvatarUrl,
		}
	}

	return &dtos.GetGroupMembersResponse{
		Members: memberItems,
	}, nil
}

func (ctrl *Controller) AddGroupMember(ctx *fiber.Ctx, req dtos.AddGroupMemberRequest) (*dtos.AddGroupMemberResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.add_group_member").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	err := ctrl.GroupApp.AddMember(req.GroupId, req.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to add member to group")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to add member to group",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.AddGroupMemberResponse{
		Message: "Member added successfully",
	}, nil
}

func (ctrl *Controller) RemoveGroupMember(ctx *fiber.Ctx, req dtos.RemoveGroupMemberRequest) (*dtos.RemoveGroupMemberResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.remove_group_member").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	err := ctrl.GroupApp.RemoveMember(req.GroupId, req.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to remove member from group")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to remove member from group",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.RemoveGroupMemberResponse{
		Message: "Member removed successfully",
	}, nil
}

// Spaces (Admin)

func (ctrl *Controller) AdminCreateSpace(ctx *fiber.Ctx, req dtos.AdminCreateSpaceRequest) (*dtos.AdminCreateSpaceResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.create_space").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	spaceType := domain.SpaceType(req.Type)
	space, err := ctrl.SpaceApp.AdminCreateSpace(req.Name, req.Icon, req.IconColor, spaceType, req.OwnerId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create space")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to create space",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.AdminCreateSpaceResponse{
		Id:      space.Id,
		Message: "Space created successfully",
	}, nil
}

func (ctrl *Controller) AdminUpdateSpace(ctx *fiber.Ctx, req dtos.AdminUpdateSpaceRequest) (*dtos.AdminUpdateSpaceResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.update_space").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	spaceType := domain.SpaceType(req.Type)
	err := ctrl.SpaceApp.AdminUpdateSpace(req.SpaceId, req.Name, req.Icon, req.IconColor, spaceType, req.OwnerId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update space")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to update space",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.AdminUpdateSpaceResponse{
		Message: "Space updated successfully",
	}, nil
}

func (ctrl *Controller) AdminDeleteSpace(ctx *fiber.Ctx, req dtos.AdminDeleteSpaceRequest) (*dtos.AdminDeleteSpaceResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.delete_space").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	err := ctrl.SpaceApp.AdminDeleteSpace(req.SpaceId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete space")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to delete space",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.AdminDeleteSpaceResponse{
		Message: "Space deleted successfully",
	}, nil
}

func (ctrl *Controller) AdminListSpacePermissions(ctx *fiber.Ctx, req dtos.AdminListSpacePermissionsRequest) (*dtos.AdminListSpacePermissionsResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.list_space_permissions").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	permissions, err := ctrl.SpaceApp.AdminListSpacePermissions(req.SpaceId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list space permissions")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to list space permissions",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	items := make([]dtos.SpacePermissionItem, len(permissions))
	for i, p := range permissions {
		items[i] = dtos.SpacePermissionItem{
			Id:        p.Id,
			UserId:    p.UserId,
			GroupId:   p.GroupId,
			Role:      string(p.Role),
			CreatedAt: p.CreatedAt,
		}
		if p.User != nil {
			items[i].Username = p.User.Username
		}
		if p.Group != nil {
			items[i].GroupName = p.Group.Name
		}
	}

	return &dtos.AdminListSpacePermissionsResponse{
		Permissions: items,
	}, nil
}

func (ctrl *Controller) AdminAddSpaceUserPermission(ctx *fiber.Ctx, req dtos.AdminAddSpaceUserPermissionRequest) (*dtos.AdminAddSpaceUserPermissionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.add_space_user_permission").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	role := domain.SpaceRole(req.Role)
	err := ctrl.SpaceApp.AdminAddSpaceUserPermission(req.SpaceId, req.UserId, role)
	if err != nil {
		logger.Error().Err(err).Msg("failed to add user permission")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to add user permission",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.AdminAddSpaceUserPermissionResponse{
		Message: "User permission added successfully",
	}, nil
}

func (ctrl *Controller) AdminRemoveSpaceUserPermission(ctx *fiber.Ctx, req dtos.AdminRemoveSpaceUserPermissionRequest) (*dtos.AdminRemoveSpaceUserPermissionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.remove_space_user_permission").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	err := ctrl.SpaceApp.AdminRemoveSpaceUserPermission(req.SpaceId, req.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to remove user permission")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to remove user permission",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.AdminRemoveSpaceUserPermissionResponse{
		Message: "User permission removed successfully",
	}, nil
}

func (ctrl *Controller) AdminAddSpaceGroupPermission(ctx *fiber.Ctx, req dtos.AdminAddSpaceGroupPermissionRequest) (*dtos.AdminAddSpaceGroupPermissionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.add_space_group_permission").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	role := domain.SpaceRole(req.Role)
	err := ctrl.SpaceApp.AdminAddSpaceGroupPermission(req.SpaceId, req.GroupId, role)
	if err != nil {
		logger.Error().Err(err).Msg("failed to add group permission")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to add group permission",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.AdminAddSpaceGroupPermissionResponse{
		Message: "Group permission added successfully",
	}, nil
}

func (ctrl *Controller) AdminRemoveSpaceGroupPermission(ctx *fiber.Ctx, req dtos.AdminRemoveSpaceGroupPermissionRequest) (*dtos.AdminRemoveSpaceGroupPermissionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.admin.remove_space_group_permission").Logger()

	if _, errResp := ctrl.checkAdmin(ctx); errResp != nil {
		return nil, errResp
	}

	err := ctrl.SpaceApp.AdminRemoveSpaceGroupPermission(req.SpaceId, req.GroupId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to remove group permission")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to remove group permission",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.AdminRemoveSpaceGroupPermissionResponse{
		Message: "Group permission removed successfully",
	}, nil
}
