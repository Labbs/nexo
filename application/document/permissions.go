package document

import (
	"fmt"

	"github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
)

func (c *DocumentApp) ListDocumentPermissions(input dto.ListDocumentPermissionsInput) (*dto.ListDocumentPermissionsOutput, error) {
	// Get the document with space to check permissions
	doc, err := c.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, &input.DocumentId, nil, input.RequesterId)
	if err != nil || doc == nil {
		return nil, fmt.Errorf("not_found")
	}

	// User must have at least viewer access to the document to see permissions
	if !doc.HasPermission(input.RequesterId, domain.DocumentRoleViewer) {
		return nil, fmt.Errorf("forbidden")
	}

	permissions, err := c.DocumentPermissionPers.ListByDocumentId(input.DocumentId)
	if err != nil {
		return nil, err
	}

	return &dto.ListDocumentPermissionsOutput{Permissions: permissions}, nil
}

func (c *DocumentApp) UpsertDocumentUserPermission(input dto.UpsertDocumentUserPermissionInput) error {
	// Get the document with space to check permissions
	doc, err := c.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, &input.DocumentId, nil, input.RequesterId)
	if err != nil || doc == nil {
		return fmt.Errorf("not_found")
	}

	// User must be able to manage permissions (owner of document OR admin of space)
	if !doc.CanManagePermissions(input.RequesterId) {
		return fmt.Errorf("forbidden")
	}

	// Prevent changing the owner's role on documents in personal/private spaces
	if (doc.Space.Type == domain.SpaceTypePersonal || doc.Space.Type == domain.SpaceTypePrivate) &&
		doc.Space.OwnerId != nil && *doc.Space.OwnerId == input.TargetUserId && input.Role != domain.DocumentRoleOwner {
		return fmt.Errorf("cannot_change_owner_role")
	}

	return c.DocumentPermissionPers.Upsert(input.DocumentId, input.TargetUserId, input.Role)
}

func (c *DocumentApp) DeleteDocumentUserPermission(input dto.DeleteDocumentUserPermissionInput) error {
	// Get the document with space to check permissions
	doc, err := c.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, &input.DocumentId, nil, input.RequesterId)
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

	return c.DocumentPermissionPers.Delete(input.DocumentId, input.TargetUserId)
}
