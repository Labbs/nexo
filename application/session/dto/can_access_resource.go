package dto

import fiberoapi "github.com/labbs/fiber-oapi"

type CanAccessResourceInput struct {
	Context      *fiberoapi.AuthContext
	ResourceType string
	ResourceID   string
	Action       string
}
