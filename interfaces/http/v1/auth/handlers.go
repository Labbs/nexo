package auth

import (
	"github.com/gofiber/fiber/v2"
	fiberoapi "github.com/labbs/fiber-oapi"
	authDto "github.com/labbs/nexo/application/auth/dto"
	"github.com/labbs/nexo/interfaces/http/v1/auth/dtos"
)

func (ctrl Controller) Login(ctx *fiber.Ctx, req dtos.LoginRequest) (*dtos.LoginResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.auth.login").Logger()

	resp, err := ctrl.AuthApp.Authenticate(authDto.AuthenticateInput{
		Email:    req.Email,
		Password: req.Password,
		Context:  ctx,
	})
	if err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("failed to authenticate user")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: err.Error(),
			Type:    "AUTHENTICATION_FAILED",
		}
	}
	return &dtos.LoginResponse{Token: resp.Token}, nil
}

func (ctrl Controller) Logout(ctx *fiber.Ctx, input struct{}) (*dtos.LogoutResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.auth.logout").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: "Authentication required",
			Type:    "AUTHENTICATION_REQUIRED",
		}
	}

	err = ctrl.AuthApp.Logout(authDto.LogoutInput{SessionId: authCtx.Claims["session_id"].(string)})
	if err != nil {
		logger.Error().Err(err).Str("session_id", authCtx.Claims["session_id"].(string)).Msg("failed to logout user")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: err.Error(),
			Type:    "LOGOUT_FAILED",
		}
	}

	return &dtos.LogoutResponse{
		Message: "Logged out successfully",
	}, nil
}

func (ctrl Controller) Register(ctx *fiber.Ctx, req dtos.RegisterRequest) (*dtos.RegisterResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.auth.register").Logger()

	err := ctrl.AuthApp.Register(authDto.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("failed to register user")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Details: err.Error(),
			Type:    "REGISTRATION_FAILED",
		}
	}

	return &dtos.RegisterResponse{
		Message: "User registered successfully",
	}, nil
}

//TODO: implement password reset, email verification, ...
