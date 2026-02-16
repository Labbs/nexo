package dto

import fiberoapi "github.com/labbs/fiber-oapi"

type HasScopeInput struct {
	Context *fiberoapi.AuthContext
	Scope   string
}
