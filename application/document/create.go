package document

import (
	"fmt"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/gosimple/slug"
	"github.com/labbs/nexo/application/document/dto"
	permissionDto "github.com/labbs/nexo/application/permission/dto"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/helpers/shortuuid"
)

func (a *DocumentApplication) CreateDocument(input dto.CreateDocumentInput) (*dto.CreateDocumentOutput, error) {
	logger := a.Logger.With().Str("component", "application.document.create_document").Logger()

	// Get the space via port (returns SpaceDetail DTO)
	spaceResult, err := a.SpaceApplication.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: input.SpaceId})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get space for document creation")
		return nil, fmt.Errorf("failed to get space for document creation: %w", err)
	}
	spaceDetail := spaceResult.Space

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
	} else {
		if !spaceDetail.HasPermission(input.UserId, "editor") {
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

	err = a.DocumentPers.Create(document, input.UserId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create document")
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	// Auto-create owner permission for the creator
	if err := a.PermissionApplication.AssignOwnerPermission(permissionDto.AssignOwnerPermissionInput{
		ResourceType: "document",
		ResourceId:   document.Id,
		UserId:       input.UserId,
		Role:         "owner",
	}); err != nil {
		logger.Warn().Err(err).Str("document_id", document.Id).Str("user_id", input.UserId).Msg("failed to create creator permission")
	}

	// Map SpaceDetail DTO to domain.Space for the response
	document.Space = domain.Space{
		Id:        spaceDetail.Id,
		Name:      spaceDetail.Name,
		Slug:      spaceDetail.Slug,
		Icon:      spaceDetail.Icon,
		IconColor: spaceDetail.IconColor,
		Type:      domain.SpaceType(spaceDetail.Type),
		OwnerId:   spaceDetail.OwnerId,
		CreatedAt: spaceDetail.CreatedAt,
		UpdatedAt: spaceDetail.UpdatedAt,
	}

	return &dto.CreateDocumentOutput{Document: document}, nil
}
