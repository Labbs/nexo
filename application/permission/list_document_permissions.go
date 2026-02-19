package permission

import (
	"fmt"

	documentDto "github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
)

// ListDocumentPermissions returns all permissions for a document.
// The requester must have at least viewer access to the document.
func (app *PermissionApplication) ListDocumentPermissions(input documentDto.ListDocumentPermissionsInput) (*documentDto.ListDocumentPermissionsOutput, error) {
	// Get the document with space to check permissions
	docResult, err := app.DocumentApp.GetDocumentByIdOrSlugWithUserPermissions(documentDto.GetDocumentByIdOrSlugWithUserPermissionsInput{
		SpaceId:    input.SpaceId,
		DocumentId: &input.DocumentId,
		UserId:     input.RequesterId,
	})
	if err != nil || docResult.Document == nil {
		return nil, fmt.Errorf("not_found")
	}

	// User must have at least viewer access to the document to see permissions
	if !docResult.Document.HasPermission(input.RequesterId, "viewer") {
		return nil, fmt.Errorf("forbidden")
	}

	permissions, err := app.PermissionPers.ListByResource(domain.PermissionTypeDocument, input.DocumentId)
	if err != nil {
		return nil, err
	}

	return &documentDto.ListDocumentPermissionsOutput{Permissions: permissions}, nil
}
