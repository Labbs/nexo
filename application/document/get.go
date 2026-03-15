package document

import (
	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
	"github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
)

func (a *DocumentApplication) GetDocumentWithSpace(input dto.GetDocumentWithSpaceInput) (*dto.GetDocumentWithSpaceOutput, error) {
	logger := a.Logger.With().Str("component", "application.document.get_document").Logger()

	if input.DocumentId == nil && input.Slug == nil {
		return nil, apperrors.ErrInvalidInput
	}

	document, err := a.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, input.DocumentId, input.Slug, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get document")
		return nil, err
	}

	doc := &dto.Document{
		Id:       document.Id,
		Name:     document.Name,
		Slug:     document.Slug,
		SpaceId:  document.SpaceId,
		ParentId: document.ParentId,
		Public:   document.Public,
		Content:  dto.JSONToBlocks(document.Content),
		Config: dto.DocumentConfig{
			Icon:             document.Config.Icon,
			FullWidth:        document.Config.FullWidth,
			Lock:             document.Config.Lock,
			HeaderBackground: document.Config.HeaderBackground,
		},
		CreatedAt: document.CreatedAt,
		UpdatedAt: document.UpdatedAt,
	}

	// Map parent if loaded
	if document.Parent != nil {
		doc.Parent = &dto.Document{
			Id:       document.Parent.Id,
			Name:     document.Parent.Name,
			Slug:     document.Parent.Slug,
			SpaceId:  document.Parent.SpaceId,
			ParentId: document.Parent.ParentId,
			Config: dto.DocumentConfig{
				Icon:             document.Parent.Config.Icon,
				FullWidth:        document.Parent.Config.FullWidth,
				Lock:             document.Parent.Config.Lock,
				HeaderBackground: document.Parent.Config.HeaderBackground,
			},
		}
	}

	// Map space if loaded
	if document.Space.Id != "" {
		doc.Space = dto.DocumentSpace{
			Id:        document.Space.Id,
			Name:      document.Space.Name,
			Slug:      document.Space.Slug,
			Icon:      document.Space.Icon,
			IconColor: document.Space.IconColor,
			Type:      string(document.Space.Type),
			CreatedAt: document.Space.CreatedAt,
			UpdatedAt: document.Space.UpdatedAt,
		}
	}

	return &dto.GetDocumentWithSpaceOutput{Document: doc}, nil
}

func (a DocumentApplication) GetDocumentsFromSpaceWithUserPermissions(input dto.GetDocumentsFromSpaceInput) (*dto.GetDocumentsFromSpaceOutput, error) {
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
