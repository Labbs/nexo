package user

import (
	"github.com/gofiber/fiber/v2"
	fiberoapi "github.com/labbs/fiber-oapi"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	userDto "github.com/labbs/nexo/application/user/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/helpers/mapper"
	spaceDtos "github.com/labbs/nexo/interfaces/http/v1/space/dtos"
	"github.com/labbs/nexo/interfaces/http/v1/user/dtos"
)

func (ctrl *Controller) GetProfile(ctx *fiber.Ctx, input struct{}) (*dtos.ProfileResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.user.get_profile").Logger()

	// Get the authenticated user context
	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: "Authentication required",
			Type:    "AUTHENTICATION_REQUIRED",
		}
	}

	result, err := ctrl.UserApp.GetByUserId(userDto.GetByUserIdInput{UserId: authCtx.UserID})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get user by id")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to retrieve user",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	profile := dtos.ProfileResponse{}
	err = mapper.MapStructByFieldNames(result.User, &profile)
	if err != nil {
		logger.Error().Err(err).Msg("failed to map user to profile")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to retrieve profile",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// Add role (not auto-mapped because it's a custom type)
	profile.Role = string(result.User.Role)

	// Add preferences (not auto-mapped because domain.JSONB != map[string]any for mapper)
	if result.User.Preferences != nil {
		profile.Preferences = map[string]any(result.User.Preferences)
	}

	return &profile, nil
}

func (ctrl *Controller) GetMySpaces(ctx *fiber.Ctx, input struct{}) (*dtos.GetMySpacesResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.user.get_my_spaces").Logger()

	// Get the authenticated user context
	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: "Authentication required",
			Type:    "AUTHENTICATION_REQUIRED",
		}
	}

	result, err := ctrl.SpaceApp.GetSpacesForUser(spaceDto.GetSpacesForUserInput{UserId: authCtx.UserID})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get spaces for user")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to retrieve spaces",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// Map domain spaces to DTO spaces
	spaceDtoList := make([]spaceDtos.Space, len(result.Spaces))
	for i, space := range result.Spaces {
		var spaceItem spaceDtos.Space
		err := mapper.MapStructByFieldNames(&space, &spaceItem)
		if err != nil {
			logger.Error().Err(err).Msg("failed to map space to DTO")
			return nil, &fiberoapi.ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Details: "Failed to process spaces",
				Type:    "INTERNAL_SERVER_ERROR",
			}
		}

		// Handle the Type field conversion (SpaceType to string)
		spaceItem.Type = string(space.Type)

		// Get the user's role in this space
		userRole := space.GetUserRole(authCtx.UserID)
		if userRole != nil {
			spaceItem.MyRole = string(*userRole)
		} else if space.Type == domain.SpaceTypePublic {
			// For public spaces without explicit permission, user is a viewer
			spaceItem.MyRole = string(domain.PermissionRoleViewer)
		}

		spaceDtoList[i] = spaceItem
	}

	response := &dtos.GetMySpacesResponse{
		Spaces: spaceDtoList,
	}

	return response, nil
}

func (ctrl *Controller) GetMyFavorites(ctx *fiber.Ctx, input struct{}) (*dtos.GetMyFavoritesResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.user.get_my_favorites").Logger()

	// Get the authenticated user context
	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: "Authentication required",
			Type:    "AUTHENTICATION_REQUIRED",
		}
	}

	result, err := ctrl.UserApp.GetMyFavorites(userDto.GetMyFavoritesInput{UserId: authCtx.UserID})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get my favorites")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to retrieve favorites",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// Map domain favorites to DTO favorites
	favoriteDtoList := make([]dtos.Favorite, len(result.Favorites))
	for i, favorite := range result.Favorites {
		favoriteDto := dtos.Favorite{
			Id:       favorite.Id,
			SpaceId:  favorite.SpaceId,
			Position: favorite.Position,
		}

		// Manually map the document
		if favorite.Document.Id != "" {
			err := mapper.MapStructByFieldNames(&favorite.Document, &favoriteDto.Document)
			if err != nil {
				logger.Error().Err(err).Msg("failed to map favorite document to DTO")
				return nil, &fiberoapi.ErrorResponse{
					Code:    fiber.StatusInternalServerError,
					Details: "Failed to process favorites",
					Type:    "INTERNAL_SERVER_ERROR",
				}
			}
		}

		favoriteDtoList[i] = favoriteDto
	}

	response := &dtos.GetMyFavoritesResponse{
		Favorites: favoriteDtoList,
	}

	return response, nil
}

func (ctrl *Controller) AddFavorite(ctx *fiber.Ctx, req dtos.AddFavoriteRequest) (*dtos.AddFavoriteResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.user.add_favorite").Logger()

	// Get the authenticated user context
	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: "Authentication required",
			Type:    "AUTHENTICATION_REQUIRED",
		}
	}

	err = ctrl.UserApp.CreateFavorite(userDto.CreateFavoriteInput{
		DocumentId: req.DocumentId,
		SpaceId:    req.SpaceId,
		UserId:     authCtx.UserID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to add favorite")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to add favorite",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.AddFavoriteResponse{
		Message: "favorite added",
	}, nil
}

func (ctrl *Controller) RemoveFavorite(ctx *fiber.Ctx, req dtos.RemoveFavoriteRequest) (*dtos.RemoveFavoriteResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.user.remove_favorite").Logger()

	// Get the authenticated user context
	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: "Authentication required",
			Type:    "AUTHENTICATION_REQUIRED",
		}
	}

	result, err := ctrl.UserApp.GetFavoriteByIdAndUserId(userDto.GetFavoriteByIdAndUserIdInput{
		FavoriteId: req.FavoriteId,
		UserId:     authCtx.UserID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get favorite by id and user id")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to retrieve favorite",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	err = ctrl.UserApp.DeleteFavorite(userDto.DeleteFavoriteInput{
		DocumentId: result.Favorite.DocumentId,
		UserId:     authCtx.UserID,
		SpaceId:    result.Favorite.SpaceId,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to remove favorite")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to remove favorite",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return &dtos.RemoveFavoriteResponse{
		Message: "favorite removed",
	}, nil
}

func (ctrl *Controller) UpdateFavoritePosition(ctx *fiber.Ctx, req dtos.UpdateFavoritePositionRequest) (*dtos.UpdateFavoritePositionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.user.update_favorite_position").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.UserApp.UpdateFavoritePosition(userDto.UpdateFavoritePositionInput{
		UserId:      authCtx.UserID,
		FavoriteId:  req.FavoriteId,
		NewPosition: req.Position,
	})
	if err != nil {
		if err.Error() == "not_found" {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Favorite not found", Type: "FAVORITE_NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to update favorite position")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to update favorite position", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.UpdateFavoritePositionResponse{FavoriteId: result.Favorite.Id, Position: result.Favorite.Position}, nil
}

func (ctrl *Controller) UpdateProfile(ctx *fiber.Ctx, req dtos.UpdateProfileRequest) (*dtos.UpdateProfileResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.user.update_profile").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	// Convert preferences if provided
	var prefs *domain.JSONB
	if req.Preferences != nil {
		jsonb := domain.JSONB(*req.Preferences)
		prefs = &jsonb
	}

	result, err := ctrl.UserApp.UpdateProfile(userDto.UpdateProfileInput{
		UserId:      authCtx.UserID,
		Username:    req.Username,
		AvatarUrl:   req.AvatarUrl,
		Preferences: prefs,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to update profile")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: err.Error(), Type: "UPDATE_PROFILE_FAILED"}
	}

	return &dtos.UpdateProfileResponse{
		Id:          result.User.Id,
		Username:    result.User.Username,
		Email:       result.User.Email,
		AvatarUrl:   result.User.AvatarUrl,
		Preferences: result.User.Preferences,
	}, nil
}

func (ctrl *Controller) ChangePassword(ctx *fiber.Ctx, req dtos.ChangePasswordRequest) (*dtos.ChangePasswordResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.user.change_password").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.UserApp.ChangePassword(userDto.ChangePasswordInput{
		UserId:          authCtx.UserID,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		if err.Error() == "invalid current password" {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Invalid current password", Type: "INVALID_PASSWORD"}
		}
		logger.Error().Err(err).Msg("failed to change password")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: err.Error(), Type: "CHANGE_PASSWORD_FAILED"}
	}

	return &dtos.ChangePasswordResponse{Message: "Password changed successfully"}, nil
}

func (ctrl *Controller) UpdateSpaceOrder(ctx *fiber.Ctx, req dtos.UpdateSpaceOrderRequest) (*dtos.UpdateSpaceOrderResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.user.update_space_order").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.UserApp.UpdateSpaceOrder(userDto.UpdateSpaceOrderInput{
		UserId:   authCtx.UserID,
		SpaceIds: req.SpaceIds,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to update space order")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to update space order", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.UpdateSpaceOrderResponse{SpaceIds: result.SpaceIds}, nil
}

func (ctrl *Controller) ListUsers(ctx *fiber.Ctx, req dtos.ListUsersRequest) (*dtos.ListUsersResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.user.list_users").Logger()

	// Get the authenticated user context (just to ensure user is authenticated)
	_, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: "Authentication required",
			Type:    "AUTHENTICATION_REQUIRED",
		}
	}

	// Set default limit if not provided
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}

	// Get users from persistence layer
	users, totalCount, err := ctrl.UserApp.UserPres.GetAll(limit, req.Offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get users")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to retrieve users",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// Map to simplified DTO (only id, username, avatar)
	userItems := make([]dtos.UserListItem, len(users))
	for i, user := range users {
		userItems[i] = dtos.UserListItem{
			Id:        user.Id,
			Username:  user.Username,
			AvatarUrl: user.AvatarUrl,
		}
	}

	return &dtos.ListUsersResponse{
		Users:      userItems,
		TotalCount: totalCount,
	}, nil
}
