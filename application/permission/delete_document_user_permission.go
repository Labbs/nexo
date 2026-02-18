package permission

import (
	"fmt"

	dto "github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
)

// DeleteDocumentUserPermission removes a user permission from a document.
// The requester must be able to manage permissions (owner of document OR admin of space).
func (app *PermissionApplication) DeleteDocumentUserPermission(input dto.DeleteDocumentUserPermissionInput) error {
	// Get the document with space to check permissions
	doc, err := app.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, &input.DocumentId, nil, input.RequesterId)
	if err != nil || doc == nil {
		return fmt.Errorf("not_found")
	}

	// User must be able to manage permissions (owner of document OR admin of space)
	if !doc.CanManagePermissions(input.RequesterId) {
		return fmt.Errorf("forbidden")
	}

	// Prevent removing the space owner from documents in personal/private spaces
	if (doc.Space.Type == domain.SpaceTypePersonal || doc.Space.Type == domain.SpaceTypePrivate) &&
		doc.Space.OwnerId != nil && *doc.Space.OwnerId == input.TargetUserId {
		return fmt.Errorf("cannot_remove_owner")
	}

	return app.PermissionPers.DeleteUser(domain.PermissionTypeDocument, input.DocumentId, input.TargetUserId)
}
