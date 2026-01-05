package document

import (
	"fmt"

	"github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
)

func (a *DocumentApp) GetDocumentWithSpace(input dto.GetDocumentWithSpaceInput) (*dto.GetDocumentWithSpaceOutput, error) {
	logger := a.Logger.With().Str("component", "application.document.get_document").Logger()

	if input.DocumentId == nil && input.Slug == nil {
		return nil, fmt.Errorf("either documentId or slug must be provided")
	}

	document, err := a.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, input.DocumentId, input.Slug, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get document")
		return nil, err
	}

	return &dto.GetDocumentWithSpaceOutput{Document: document}, nil
}

func (a DocumentApp) GetDocumentsFromSpaceWithUserPermissions(input dto.GetDocumentsFromSpaceInput) (*dto.GetDocumentsFromSpaceOutput, error) {
	logger := a.Logger.With().Str("component", "application.document.get_documents_from_space").Logger()

	var documents []domain.Document
	var err error

	if input.ParentId != nil {
		documents, err = a.DocumentPers.GetChildDocumentsWithUserPermissions(*input.ParentId, input.UserId)
	} else {
		documents, err = a.DocumentPers.GetRootDocumentsFromSpaceWithUserPermissions(input.SpaceId, input.UserId)
	}

	if err != nil {
		logger.Error().Err(err).Msg("failed to get documents from space")
		return nil, err
	}

	return &dto.GetDocumentsFromSpaceOutput{Documents: documents}, nil
}
