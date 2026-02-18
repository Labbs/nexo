package permission

import (
	"fmt"

	dto "github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
)

// ListDocumentPermissions returns all permissions for a document.
// The requester must have at least viewer access to the document.
func (app *PermissionApplication) ListDocumentPermissions(input dto.ListDocumentPermissionsInput) (*dto.ListDocumentPermissionsOutput, error) {
	// Get the document with space to check permissions
	doc, err := app.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, &input.DocumentId, nil, input.RequesterId)
	if err != nil || doc == nil {
		return nil, fmt.Errorf("not_found")
	}

	// User must have at least viewer access to the document to see permissions
	if !doc.HasPermission(input.RequesterId, domain.PermissionRoleViewer) {
		return nil, fmt.Errorf("forbidden")
	}

	permissions, err := app.PermissionPers.ListByResource(domain.PermissionTypeDocument, input.DocumentId)
	if err != nil {
		return nil, err
	}

	return &dto.ListDocumentPermissionsOutput{Permissions: permissions}, nil
}
