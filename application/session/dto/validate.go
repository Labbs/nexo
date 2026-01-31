package dto

import fiberoapi "github.com/labbs/fiber-oapi"

type ValidateTokenInput struct {
	Token string
}

type ValidateTokenOutput struct {
	AuthContext *fiberoapi.AuthContext
}

type HasRoleInput struct {
	Context *fiberoapi.AuthContext
	Role    string
}

type HasScopeInput struct {
	Context *fiberoapi.AuthContext
	Scope   string
}

type CanAccessResourceInput struct {
	Context      *fiberoapi.AuthContext
	ResourceType string
	ResourceID   string
	Action       string
}

type GetUserPermissionsInput struct {
	Context      *fiberoapi.AuthContext
	ResourceType string
	ResourceID   string
}

type GetUserPermissionsOutput struct {
	Permission *fiberoapi.ResourcePermission
}
