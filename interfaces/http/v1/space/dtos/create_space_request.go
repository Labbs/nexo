package dtos

type CreateSpaceRequest struct {
	Name      string  `json:"name" validate:"required,min=3,max=100"`
	Icon      *string `json:"icon,omitempty"`
	IconColor *string `json:"icon_color,omitempty"`
	Type      *string `json:"type,omitempty" validate:"omitempty,oneof=public private"`
}

type CreateSpaceResponse struct {
	SpaceId string `json:"space_id"`
}

// My spaces list is defined under user DTOs (see interfaces/http/v1/user/dtos)

type UpdateSpaceRequest struct {
    SpaceId   string  `path:"space_id" validate:"required,uuid4"`
    Name      *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
    Icon      *string `json:"icon,omitempty"`
    IconColor *string `json:"icon_color,omitempty"`
}

type UpdateSpaceResponse struct {
    Space Space `json:"space"`
}

type DeleteSpaceRequest struct {
    SpaceId string `path:"space_id" validate:"required,uuid4"`
}

type DeleteSpaceResponse struct {
    SpaceId string `json:"space_id"`
}
