package permission

import (
	"fmt"

	documentDto "github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
)

// DeleteDocumentUserPermission removes a user permission from a document.
// The requester must be able to manage permissions (owner of document OR admin of space).
func (app *PermissionApplication) DeleteDocumentUserPermission(input documentDto.DeleteDocumentUserPermissionInput) error {
	// Get the document with space to check permissions
	docResult, err := app.DocumentApp.GetDocumentByIdOrSlugWithUserPermissions(documentDto.GetDocumentByIdOrSlugWithUserPermissionsInput{
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

	// Prevent removing the space owner from documents in personal/private spaces
	if (docResult.Document.Space.Type == "personal" || docResult.Document.Space.Type == "private") &&
		docResult.Document.Space.OwnerId != nil && *docResult.Document.Space.OwnerId == input.TargetUserId {
		return fmt.Errorf("cannot_remove_owner")
	}

	return app.PermissionPers.DeleteUser(domain.PermissionTypeDocument, input.DocumentId, input.TargetUserId)
}
