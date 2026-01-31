package dtos

type ListSpacePermissionsRequest struct {
    SpaceId string `path:"space_id" validate:"required,uuid4"`
}

type SpacePermission struct {
    Id        string  `json:"id"`
    UserId    *string `json:"user_id,omitempty"`
    Username  string  `json:"username,omitempty"`
    GroupId   *string `json:"group_id,omitempty"`
    GroupName string  `json:"group_name,omitempty"`
    Role      string  `json:"role"`
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




