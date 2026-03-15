package ports

import (
	databaseDto "github.com/labbs/nexo/application/database/dto"
	documentDto "github.com/labbs/nexo/application/document/dto"
	drawingDto "github.com/labbs/nexo/application/drawing/dto"
	permissionDto "github.com/labbs/nexo/application/permission/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
)

type PermissionPort interface {
	// Internal permission assignment (no authorization checks)
	AssignOwnerPermission(input permissionDto.AssignOwnerPermissionInput) error

	// Space permissions
	ListSpacePermissions(input spaceDto.ListSpacePermissionsInput) (*spaceDto.ListSpacePermissionsOutput, error)
	UpsertSpaceUserPermission(input spaceDto.UpsertSpaceUserPermissionInput) error
	DeleteSpaceUserPermission(input spaceDto.DeleteSpaceUserPermissionInput) error
	UpsertSpaceGroupPermission(input spaceDto.UpsertSpaceGroupPermissionInput) error
	DeleteSpaceGroupPermission(input spaceDto.DeleteSpaceGroupPermissionInput) error

	// Drawing permissions
	ListDrawingPermissions(input drawingDto.ListDrawingPermissionsInput) (*drawingDto.ListDrawingPermissionsOutput, error)
	UpsertDrawingUserPermission(input drawingDto.UpsertDrawingUserPermissionInput) error
	DeleteDrawingUserPermission(input drawingDto.DeleteDrawingUserPermissionInput) error

	// Document permissions
	ListDocumentPermissions(input documentDto.ListDocumentPermissionsInput) (*documentDto.ListDocumentPermissionsOutput, error)
	UpsertDocumentUserPermission(input documentDto.UpsertDocumentUserPermissionInput) error
	DeleteDocumentUserPermission(input documentDto.DeleteDocumentUserPermissionInput) error

	// Database permissions
	ListDatabasePermissions(input databaseDto.ListDatabasePermissionsInput) (*databaseDto.ListDatabasePermissionsOutput, error)
	UpsertDatabasePermission(input databaseDto.UpsertDatabasePermissionInput) error
	DeleteDatabasePermission(input databaseDto.DeleteDatabasePermissionInput) error
}
