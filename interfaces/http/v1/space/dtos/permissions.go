package dtos

type ListSpacePermissionsRequest struct {
    SpaceId string `path:"space_id" validate:"required,uuid4"`
}

type SpacePermission struct {
    UserId *string `json:"user_id,omitempty"`
    Role   string  `json:"role"`
}

type ListSpacePermissionsResponse struct {
    Permissions []SpacePermission `json:"permissions"`
}

type UpsertSpaceUserPermissionRequest struct {
    SpaceId string `path:"space_id" validate:"required,uuid4"`
    UserId  string `json:"user_id" validate:"required,uuid4"`
    Role    string `json:"role" validate:"required,oneof=owner admin editor viewer"`
}

type UpsertSpaceUserPermissionResponse struct {
    Message string `json:"message"`
}

type DeleteSpaceUserPermissionRequest struct {
    SpaceId string `path:"space_id" validate:"required,uuid4"`
    UserId  string `path:"user_id" validate:"required,uuid4"`
}

type DeleteSpaceUserPermissionResponse struct {
    Message string `json:"message"`
}




