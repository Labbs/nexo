package dto

import fiberoapi "github.com/labbs/fiber-oapi"

type HasRoleInput struct {
	Context *fiberoapi.AuthContext
	Role    string
}
