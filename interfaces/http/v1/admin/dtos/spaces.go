package dtos

import "time"

// Space list for admin

type ListAllSpacesRequest struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

type SpaceItem struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Type        string    `json:"type"`
	OwnerId     string    `json:"owner_id"`
	OwnerName   string    `json:"owner_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListAllSpacesResponse struct {
	Spaces     []SpaceItem `json:"spaces"`
	TotalCount int64       `json:"total_count"`
}

// Create space (admin)

type AdminCreateSpaceRequest struct {
	Name      string  `json:"name" validate:"required,min=1,max=100"`
	Icon      string  `json:"icon"`
	IconColor string  `json:"icon_color"`
	Type      string  `json:"type" validate:"required,oneof=public private restricted"`
	OwnerId   *string `json:"owner_id"` // Optional - can be nil for unowned spaces
}

type AdminCreateSpaceResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// Update space (admin)

type AdminUpdateSpaceRequest struct {
	SpaceId   string  `path:"space_id" validate:"required"`
	Name      string  `json:"name" validate:"required,min=1,max=100"`
	Icon      string  `json:"icon"`
	IconColor string  `json:"icon_color"`
	Type      string  `json:"type" validate:"required,oneof=public private restricted"`
	OwnerId   *string `json:"owner_id"`
}

type AdminUpdateSpaceResponse struct {
	Message string `json:"message"`
}

// Delete space (admin)

type AdminDeleteSpaceRequest struct {
	SpaceId string `path:"space_id" validate:"required"`
}

type AdminDeleteSpaceResponse struct {
	Message string `json:"message"`
}

// Space permissions (admin)

type AdminListSpacePermissionsRequest struct {
	SpaceId string `path:"space_id" validate:"required"`
}

type SpacePermissionItem struct {
	Id        string    `json:"id"`
	UserId    *string   `json:"user_id,omitempty"`
	Username  string    `json:"username,omitempty"`
	GroupId   *string   `json:"group_id,omitempty"`
	GroupName string    `json:"group_name,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminListSpacePermissionsResponse struct {
	Permissions []SpacePermissionItem `json:"permissions"`
}

type AdminAddSpaceUserPermissionRequest struct {
	SpaceId string `path:"space_id" validate:"required"`
	UserId  string `json:"user_id" validate:"required"`
	Role    string `json:"role" validate:"required,oneof=viewer editor admin"`
}

type AdminAddSpaceUserPermissionResponse struct {
	Message string `json:"message"`
}

type AdminRemoveSpaceUserPermissionRequest struct {
	SpaceId string `path:"space_id" validate:"required"`
	UserId  string `path:"user_id" validate:"required"`
}

type AdminRemoveSpaceUserPermissionResponse struct {
	Message string `json:"message"`
}

type AdminAddSpaceGroupPermissionRequest struct {
	SpaceId string `path:"space_id" validate:"required"`
	GroupId string `json:"group_id" validate:"required"`
	Role    string `json:"role" validate:"required,oneof=viewer editor admin"`
}

type AdminAddSpaceGroupPermissionResponse struct {
	Message string `json:"message"`
}

type AdminRemoveSpaceGroupPermissionRequest struct {
	SpaceId string `path:"space_id" validate:"required"`
	GroupId string `path:"group_id" validate:"required"`
}

type AdminRemoveSpaceGroupPermissionResponse struct {
	Message string `json:"message"`
}
