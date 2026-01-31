package dtos

type ListDrawingPermissionsRequest struct {
	DrawingId string `path:"drawing_id" validate:"required,uuid4"`
}

type DrawingPermission struct {
	Id       string  `json:"id"`
	UserId   *string `json:"user_id,omitempty"`
	Username *string `json:"username,omitempty"`
	Role     string  `json:"role"`
}

type ListDrawingPermissionsResponse struct {
	Permissions []DrawingPermission `json:"permissions"`
}

type UpsertDrawingUserPermissionRequest struct {
	DrawingId string `path:"drawing_id" validate:"required,uuid4"`
	UserId    string `json:"user_id" validate:"required,uuid4"`
	Role      string `json:"role" validate:"required,oneof=owner editor viewer denied"`
}

type UpsertDrawingUserPermissionResponse struct {
	Message string `json:"message"`
}

type DeleteDrawingUserPermissionRequest struct {
	DrawingId string `path:"drawing_id" validate:"required,uuid4"`
	UserId    string `path:"user_id" validate:"required,uuid4"`
}

type DeleteDrawingUserPermissionResponse struct {
	Message string `json:"message"`
}
