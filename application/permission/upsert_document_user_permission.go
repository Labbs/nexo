package permission

import (
	"fmt"

	documentDto "github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
)

// UpsertDocumentUserPermission adds or updates a user permission for a document.
// The requester must be able to manage permissions (owner of document OR admin of space).
func (app *PermissionApplication) UpsertDocumentUserPermission(input documentDto.UpsertDocumentUserPermissionInput) error {
	// Get the document with space to check permissions
	docResult, err := app.DocumentApplication.GetDocumentByIdOrSlugWithUserPermissions(documentDto.GetDocumentByIdOrSlugWithUserPermissionsInput{
		SpaceId:    input.SpaceId,
		DocumentId: &input.DocumentId,
		UserId:     input.RequesterId,
	})
	if err != nil || docResult.Document == nil {
		return fmt.Errorf("not_found")
	}

	// User must be able to manage permissions (owner of document OR admin of space)
	if !docResult.Document.CanManagePermissions(input.RequesterId) {
		return fmt.Errorf("forbidden")
	}

	// Prevent changing the owner's role on documents in personal/private spaces
	if (docResult.Document.Space.Type == "personal" || docResult.Document.Space.Type == "private") &&
		docResult.Document.Space.OwnerId != nil && *docResult.Document.Space.OwnerId == input.TargetUserId && input.Role != "owner" {
		return fmt.Errorf("cannot_change_owner_role")
	}

	return app.PermissionPers.UpsertUser(domain.PermissionTypeDocument, input.DocumentId, input.TargetUserId, domain.PermissionRole(input.Role))
}
