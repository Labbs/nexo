package document

import (
	"fmt"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/gosimple/slug"
	"github.com/labbs/nexo/infrastructure/helpers/apperrors"
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
			return nil, apperrors.ErrAccessDenied
		}
	} else {
		if !spaceDetail.HasPermission(input.UserId, "editor") {
			logger.Error().Msg("user does not have permission to create document in the specified space")
			return nil, apperrors.ErrAccessDenied
		}
	}

	name := input.Name
	content := dto.BlocksToJSON(input.Content)
	var config domain.DocumentConfig

	// Clone content/name/config from template when requested
	if input.TemplateId != nil {
		tmpl, err := a.DocumentPers.GetDocumentWithPermissions(*input.TemplateId, input.UserId)
		if err != nil {
			logger.Error().Err(err).Str("template_id", *input.TemplateId).Msg("failed to get template document")
			return nil, fmt.Errorf("failed to get template: %w", err)
		}
		if name == "" || name == "New Document" {
			name = tmpl.Name
		}
		content = tmpl.Content
		config = domain.DocumentConfig{
			FullWidth:        tmpl.Config.FullWidth,
			Icon:             tmpl.Config.Icon,
			HeaderBackground: tmpl.Config.HeaderBackground,
			// Lock is intentionally not cloned — the new doc starts unlocked
		}
	}

	document := &domain.Document{
		Id:       utils.UUIDv4(),
		Name:     name,
		Slug:     slug.Make(name + "-" + shortuuid.GenerateShortUUID()),
		SpaceId:  input.SpaceId,
		Content:  content,
		Config:   config,
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
