package permission

import (
	"fmt"

	dto "github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
)

// UpsertDocumentUserPermission adds or updates a user permission for a document.
// The requester must be able to manage permissions (owner of document OR admin of space).
func (app *PermissionApplication) UpsertDocumentUserPermission(input dto.UpsertDocumentUserPermissionInput) error {
	// Get the document with space to check permissions
	doc, err := app.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, &input.DocumentId, nil, input.RequesterId)
	if err != nil || doc == nil {
		return fmt.Errorf("not_found")
	}

	// User must be able to manage permissions (owner of document OR admin of space)
	if !doc.CanManagePermissions(input.RequesterId) {
		return fmt.Errorf("forbidden")
	}

	// Prevent changing the owner's role on documents in personal/private spaces
	if (doc.Space.Type == domain.SpaceTypePersonal || doc.Space.Type == domain.SpaceTypePrivate) &&
		doc.Space.OwnerId != nil && *doc.Space.OwnerId == input.TargetUserId && input.Role != domain.PermissionRoleOwner {
		return fmt.Errorf("cannot_change_owner_role")
	}

	return app.PermissionPers.UpsertUser(domain.PermissionTypeDocument, input.DocumentId, input.TargetUserId, input.Role)
}
