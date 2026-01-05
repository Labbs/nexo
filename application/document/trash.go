package document

import (
	"fmt"

	"github.com/labbs/nexo/application/document/dto"
)

func (c *DocumentApp) GetTrash(input dto.GetTrashInput) (*dto.GetTrashOutput, error) {
	logger := c.Logger.With().Str("component", "application.document.get_trash").Logger()

	docs, err := c.DocumentPers.GetDeletedDocuments(input.SpaceId, input.UserId)
	if err != nil {
		logger.Error().Err(err).Str("space_id", input.SpaceId).Msg("failed to get deleted documents")
		return nil, fmt.Errorf("failed to get trash")
	}

	return &dto.GetTrashOutput{Documents: docs}, nil
}

func (c *DocumentApp) RestoreDocument(input dto.RestoreDocumentInput) error {
	logger := c.Logger.With().Str("component", "application.document.restore_document").Logger()

	err := c.DocumentPers.Restore(input.DocumentId, input.UserId)
	if err != nil {
		logger.Error().Err(err).Str("document_id", input.DocumentId).Msg("failed to restore document")
		return err
	}

	return nil
}
