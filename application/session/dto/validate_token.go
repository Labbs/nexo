package dto

import fiberoapi "github.com/labbs/fiber-oapi"

type ValidateTokenInput struct {
	Token string
}

type ValidateTokenOutput struct {
	AuthContext *fiberoapi.AuthContext
}
