package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	authDto "github.com/labbs/nexo/application/auth/dto"
	"github.com/labbs/nexo/interfaces/http/v1/auth/dtos"
)

func (ctrl Controller) Login(ctx *fiber.Ctx, req dtos.LoginRequest) (*dtos.LoginResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.auth.login").Logger()

	resp, err := ctrl.AuthApplication.Authenticate(authDto.AuthenticateInput{
		Email:    req.Email,
		Password: req.Password,
		Context:  ctx,
	})
	if err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("failed to authenticate user")
		if errors.Is(err, apperrors.ErrInvalidCredentials) || errors.Is(err, apperrors.ErrUserNotActive) {
			return nil, &fiberoapi.ErrorResponse{
				Code:    fiber.StatusUnauthorized,
				Details: err.Error(),
				Type:    "AUTHENTICATION_FAILED",
			}
		}
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Authentication failed",
			Type:    "INTERNAL_SERVER_ERROR",
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

	sessionId, _ := authCtx.Claims["session_id"].(string)
	err = ctrl.AuthApplication.Logout(authDto.LogoutInput{SessionId: sessionId})
	if err != nil {
		logger.Error().Err(err).Str("session_id", sessionId).Msg("failed to logout user")
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

	err := ctrl.AuthApplication.Register(authDto.RegisterInput{
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

func (ctrl Controller) SSORedirect(ctx *fiber.Ctx, input struct{}) (*dtos.SSORedirectResponse, *fiberoapi.ErrorResponse) {
	out, err := ctrl.AuthApplication.SSORedirect()
	if err != nil {
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Details: err.Error(),
			Type:    "SSO_DISABLED",
		}
	}
	return &dtos.SSORedirectResponse{URL: out.URL, State: out.State}, nil
}

func (ctrl Controller) SSOCallback(ctx *fiber.Ctx, req dtos.SSOCallbackRequest) (*dtos.SSOCallbackResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.auth.sso_callback").Logger()

	out, err := ctrl.AuthApplication.SSOCallback(authDto.SSOCallbackInput{
		Code:    req.Code,
		State:   req.State,
		Context: ctx,
	})
	if err != nil {
		logger.Error().Err(err).Msg("SSO callback failed")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: err.Error(),
			Type:    "SSO_CALLBACK_FAILED",
		}
	}
	return &dtos.SSOCallbackResponse{Token: out.Token}, nil
}

//TODO: implement password reset, email verification, ...
