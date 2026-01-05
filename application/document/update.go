package document

import (
	"fmt"

	"github.com/gosimple/slug"
	"github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/helpers/shortuuid"
)

func (a *DocumentApp) UpdateDocument(input dto.UpdateDocumentInput) (*dto.UpdateDocumentOutput, error) {
	logger := a.Logger.With().Str("component", "application.document.update_document").Logger()

	document, err := a.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, &input.DocumentId, nil, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get document for update")
		return nil, fmt.Errorf("failed to get document for update: %w", err)
	}

	if !document.HasPermission(input.UserId, domain.DocumentRoleEditor) {
		logger.Error().Msg("user does not have permission to update document")
		return nil, fmt.Errorf("user does not have permission to update document")
	}

	// Update name only if provided
	if input.Name != nil && *input.Name != "" && document.Name != *input.Name {
		document.Name = *input.Name
		document.Slug = slug.Make(*input.Name + "-" + shortuuid.GenerateShortUUID())
	}

	// Update content only if provided
	if input.Content != nil {
		document.Content = dto.BlocksToJSON(*input.Content)
	}

	// Update parentId if provided
	if input.ParentId != nil {
		// If parentId is empty string, set to nil (move to root)
		if *input.ParentId == "" {
			document.ParentId = nil
		} else {
			document.ParentId = input.ParentId
		}
	}

	// Update config if provided
	if input.Config != nil {
		document.Config = *input.Config
	}

	// Update metadata if provided
	if input.Metadata != nil {
		document.Metadata = domain.JSONB(*input.Metadata)
	}

	// Create a version snapshot before updating (for content changes)
	if input.Content != nil {
		a.CreateVersionOnUpdate(document, input.UserId)
	}

	err = a.DocumentPers.Update(document, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update document")
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	return &dto.UpdateDocumentOutput{Document: document}, nil
}
