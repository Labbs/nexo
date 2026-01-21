package document

import (
	"fmt"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/gosimple/slug"
	"github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/helpers/shortuuid"
)

func (a *DocumentApp) CreateDocument(input dto.CreateDocumentInput) (*dto.CreateDocumentOutput, error) {
	logger := a.Logger.With().Str("component", "application.document.create_document").Logger()

	var space *domain.Space

	if input.ParentId != nil {
		parent, err := a.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, input.ParentId, nil, input.UserId)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get parent document")
			return nil, fmt.Errorf("failed to get parent document: %w", err)
		}

		if !parent.HasPermission(input.UserId, domain.PermissionRoleEditor) {
			logger.Error().Msg("user does not have permission to create document under the specified parent")
			return nil, fmt.Errorf("user does not have permission to create document under the specified parent")
		}

		// Get the space for the response
		space, err = a.SpacePers.GetSpaceById(input.SpaceId)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get space for document creation")
			return nil, fmt.Errorf("failed to get space for document creation: %w", err)
		}
	} else {
		var err error
		space, err = a.SpacePers.GetSpaceById(input.SpaceId)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get space for document creation")
			return nil, fmt.Errorf("failed to get space for document creation: %w", err)
		}

		if !space.HasPermission(input.UserId, domain.PermissionRoleEditor) {
			logger.Error().Msg("user does not have permission to create document in the specified space")
			return nil, fmt.Errorf("user does not have permission to create document in the specified space")
		}
	}

	document := &domain.Document{
		Id:       utils.UUIDv4(),
		Name:     input.Name,
		Slug:     slug.Make(input.Name + "-" + shortuuid.GenerateShortUUID()),
		SpaceId:  input.SpaceId,
		Content:  dto.BlocksToJSON(input.Content),
		ParentId: input.ParentId,
		Public:   false,
	}

	err := a.DocumentPers.Create(document, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create document")
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	// Auto-create owner permission for the creator
	// This ensures they retain access and can manage permissions even if their space role is downgraded
	if err := a.PermissionPers.UpsertUser(domain.PermissionTypeDocument, document.Id, input.UserId, domain.PermissionRoleOwner); err != nil {
		// Log but don't fail - the document is already created
		logger.Warn().Err(err).Str("document_id", document.Id).Str("user_id", input.UserId).Msg("failed to create creator permission")
	}

	// Assign the space to the document for the response
	if space != nil {
		document.Space = *space
	}

	return &dto.CreateDocumentOutput{Document: document}, nil
}
