package dtos

type ListDocumentPermissionsRequest struct {
	SpaceId    string `path:"space_id" validate:"required,uuid4"`
	DocumentId string `path:"document_id" validate:"required,uuid4"`
}

type DocumentPermission struct {
	Id       string  `json:"id"`
	UserId   *string `json:"user_id,omitempty"`
	Username *string `json:"username,omitempty"`
	Role     string  `json:"role"`
}

type ListDocumentPermissionsResponse struct {
	Permissions []DocumentPermission `json:"permissions"`
}

type UpsertDocumentUserPermissionRequest struct {
	SpaceId    string `path:"space_id" validate:"required,uuid4"`
	DocumentId string `path:"document_id" validate:"required,uuid4"`
	UserId     string `json:"user_id" validate:"required,uuid4"`
	Role       string `json:"role" validate:"required,oneof=owner editor viewer denied"`
}

type UpsertDocumentUserPermissionResponse struct {
	Message string `json:"message"`
}

type DeleteDocumentUserPermissionRequest struct {
	SpaceId    string `path:"space_id" validate:"required,uuid4"`
	DocumentId string `path:"document_id" validate:"required,uuid4"`
	UserId     string `path:"user_id" validate:"required,uuid4"`
}

type DeleteDocumentUserPermissionResponse struct {
	Message string `json:"message"`
}
