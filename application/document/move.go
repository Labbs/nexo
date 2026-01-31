package document

import (
	"fmt"

	"github.com/labbs/nexo/application/document/dto"
)

func (a *DocumentApp) MoveDocument(input dto.MoveDocumentInput) (*dto.MoveDocumentOutput, error) {
	logger := a.Logger.With().Str("component", "application.document.move_document").Logger()

	// Ensure the document belongs to the provided space and user has access
	doc, err := a.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, &input.DocumentId, nil, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get document for move")
		return nil, fmt.Errorf("failed to get document for move: %w", err)
	}
	// Delegate to persistence Move (permission checks and save)
	moved, err := a.DocumentPers.Move(doc.Id, input.NewParentId, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to move document")
		return nil, err
	}
	return &dto.MoveDocumentOutput{Document: moved}, nil
}
