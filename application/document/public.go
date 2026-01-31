package document

import (
	"fmt"

	"github.com/labbs/nexo/application/document/dto"
)

func (c *DocumentApp) SetPublic(input dto.SetPublicInput) error {
	logger := c.Logger.With().Str("component", "application.document.set_public").Logger()

	err := c.DocumentPers.SetPublic(input.DocumentId, input.Public, input.UserId)
	if err != nil {
		logger.Error().Err(err).Str("document_id", input.DocumentId).Msg("failed to set document public status")
		return err
	}

	return nil
}

func (c *DocumentApp) GetPublicDocument(input dto.GetPublicDocumentInput) (*dto.GetPublicDocumentOutput, error) {
	logger := c.Logger.With().Str("component", "application.document.get_public_document").Logger()

	doc, err := c.DocumentPers.GetPublicDocument(input.SpaceId, input.DocumentId, input.Slug)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get public document")
		return nil, fmt.Errorf("document not found or not public")
	}

	return &dto.GetPublicDocumentOutput{Document: doc}, nil
}
