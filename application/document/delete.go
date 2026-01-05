package document

import (
	"fmt"

	"github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
)

func (a *DocumentApp) DeleteDocument(input dto.DeleteDocumentInput) error {
	logger := a.Logger.With().Str("component", "application.document.delete_document").Logger()

	if input.DocumentId == nil && input.Slug == nil {
		return fmt.Errorf("either documentId or slug must be provided")
	}

	document, err := a.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, input.DocumentId, input.Slug, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get document for delete")
		return fmt.Errorf("failed to get document for delete: %w", err)
	}

	if !document.HasPermission(input.UserId, domain.DocumentRoleEditor) {
		logger.Error().Msg("user does not have permission to delete document")
		return fmt.Errorf("user does not have permission to delete document")
	}

	if err := a.DocumentPers.Delete(document.Id, input.UserId); err != nil {
		logger.Error().Err(err).Msg("failed to delete document")
		return err
	}

	return nil
}
