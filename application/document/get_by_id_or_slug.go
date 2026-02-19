package document

import (
	"fmt"

	"github.com/labbs/nexo/application/document/dto"
)

func (a *DocumentApplication) GetDocumentByIdOrSlugWithUserPermissions(input dto.GetDocumentByIdOrSlugWithUserPermissionsInput) (*dto.GetDocumentByIdOrSlugWithUserPermissionsOutput, error) {
	logger := a.Logger.With().Str("component", "application.document.get_document_by_id_or_slug").Logger()

	doc, err := a.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, input.DocumentId, input.Slug, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get document by id or slug")
		return nil, fmt.Errorf("document not found: %w", err)
	}

	// Map domain.Document → DocumentDetail DTO
	permissions := make([]dto.DocumentPermission, len(doc.Permissions))
	for i, p := range doc.Permissions {
		permissions[i] = dto.DocumentPermission{
			UserId: p.UserId,
			Role:   string(p.Role),
		}
	}

	detail := &dto.DocumentDetail{
		Id:      doc.Id,
		Name:    doc.Name,
		Slug:    doc.Slug,
		SpaceId: doc.SpaceId,
		Config: dto.DocumentConfig{
			FullWidth:        doc.Config.FullWidth,
			Icon:             doc.Config.Icon,
			Lock:             doc.Config.Lock,
			HeaderBackground: doc.Config.HeaderBackground,
		},
		Space: dto.DocumentSpaceInfo{
			Type:    string(doc.Space.Type),
			OwnerId: doc.Space.OwnerId,
		},
		Permissions: permissions,
	}

	return &dto.GetDocumentByIdOrSlugWithUserPermissionsOutput{Document: detail}, nil
}
