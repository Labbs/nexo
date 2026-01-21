package document

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberoapi "github.com/labbs/fiber-oapi"
	docDto "github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/helpers/mapper"
	"github.com/labbs/nexo/infrastructure/helpers/validator"
	"github.com/labbs/nexo/interfaces/http/v1/document/dtos"
)

func (ctrl *Controller) GetDocumentsFromSpace(ctx *fiber.Ctx, req dtos.GetDocumentsFromSpaceRequest) (*dtos.GetDocumentsFromSpaceResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.get_documents_from_space").Logger()

	// Get the authenticated user context
	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: "Authentication required",
			Type:    "AUTHENTICATION_REQUIRED",
		}
	}

	// Convert empty string to nil pointer for optional parentId
	var parentId *string
	if req.ParentId != "" {
		parentId = &req.ParentId
	}

	result, err := ctrl.DocumentApp.GetDocumentsFromSpaceWithUserPermissions(docDto.GetDocumentsFromSpaceInput{
		SpaceId:  req.SpaceId,
		UserId:   authCtx.UserID,
		ParentId: parentId,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get documents from space")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to retrieve documents",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// Map domain documents to DTO documents
	resp := &dtos.GetDocumentsFromSpaceResponse{
		Documents: make([]dtos.Document, len(result.Documents)),
	}

	for i, doc := range result.Documents {
		// Map the document
		err = mapper.MapStructByFieldNames(&doc, &resp.Documents[i])
		if err != nil {
			logger.Error().Err(err).Msg("failed to map document to response DTO")
			return nil, &fiberoapi.ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Details: "Failed to process document",
				Type:    "INTERNAL_SERVER_ERROR",
			}
		}

		// Map the parent if it exists
		if doc.Parent != nil {
			resp.Documents[i].Parent = &dtos.Document{}
			err = mapper.MapStructByFieldNames(doc.Parent, resp.Documents[i].Parent)
			if err != nil {
				logger.Error().Err(err).Msg("failed to map parent document to response DTO")
				return nil, &fiberoapi.ErrorResponse{
					Code:    fiber.StatusInternalServerError,
					Details: "Failed to process parent document",
					Type:    "INTERNAL_SERVER_ERROR",
				}
			}
		}

		// Map the space
		err = mapper.MapStructByFieldNames(&doc.Space, &resp.Documents[i].Space)
		if err != nil {
			logger.Error().Err(err).Msg("failed to map space to response DTO")
			return nil, &fiberoapi.ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Details: "Failed to process space",
				Type:    "INTERNAL_SERVER_ERROR",
			}
		}
	}

	return resp, nil
}

func (ctrl *Controller) GetDocument(ctx *fiber.Ctx, req dtos.GetDocumentRequest) (*dtos.GetDocumentResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.get_document").Logger()

	// Get the authenticated user context
	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: "Authentication required",
			Type:    "AUTHENTICATION_REQUIRED",
		}
	}

	var id *string
	var slug *string

	if validator.IsValidUUID(req.Identifier) {
		id = &req.Identifier
	} else {
		slug = &req.Identifier
	}

	result, err := ctrl.DocumentApp.GetDocumentWithSpace(docDto.GetDocumentWithSpaceInput{
		UserId:     authCtx.UserID,
		SpaceId:    req.SpaceId,
		DocumentId: id,
		Slug:       slug,
	})

	logger.Debug().Interface("result", result).Msg("GetDocumentWithSpace result")

	if err != nil || result.Document == nil {
		logger.Error().Err(err).Msg("failed to get document or document is nil")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusNotFound,
			Details: "Document not found",
			Type:    "DOCUMENT_NOT_FOUND",
		}
	}

	resp := &dtos.GetDocumentResponse{}
	err = mapper.MapStructByFieldNames(result.Document, resp)
	if err != nil {
		logger.Error().Err(err).Msg("failed to map document to response DTO")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to process document",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// Manual conversion of Content from datatypes.JSON to []Block
	if len(result.Document.Content) > 0 {
		appBlocks := docDto.JSONToBlocks(result.Document.Content)
		resp.Content = make([]dtos.Block, len(appBlocks))
		for i, b := range appBlocks {
			resp.Content[i] = dtos.Block{
				ID:       b.ID,
				Type:     b.Type,
				Props:    b.Props,
				Content:  convertToHttpInlineContent(b.Content),
				Children: convertToHttpBlocks(b.Children),
			}
		}
	}

	logger.Debug().Interface("resp.Content", resp.Content).Msg("Mapped document content")

	if result.Document.Parent != nil {
		resp.Parent = &dtos.Document{}
		err = mapper.MapStructByFieldNames(result.Document.Parent, resp.Parent)
		if err != nil {
			logger.Error().Err(err).Msg("failed to map parent document to response DTO")
			return nil, &fiberoapi.ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Details: "Failed to process parent document",
				Type:    "INTERNAL_SERVER_ERROR",
			}
		}
		// Pour éviter la récursion infinie, on ne mappe pas le parent du parent
		resp.Parent.Parent = nil
	}

	// Map the space if it exists (avoid nil pointer dereference)
	if result.Document.Space.Id != "" {
		err = mapper.MapStructByFieldNames(&result.Document.Space, &resp.Space)
		if err != nil {
			logger.Error().Err(err).Msg("failed to map space to response DTO")
			return nil, &fiberoapi.ErrorResponse{
				Code:    fiber.StatusInternalServerError,
				Details: "Failed to process space",
				Type:    "INTERNAL_SERVER_ERROR",
			}
		}
	}

	logger.Debug().Interface("resp", resp).Msg("Mapped document")

	return resp, nil
}

func (ctrl *Controller) CreateDocument(ctx *fiber.Ctx, req dtos.CreateDocumentRequest) (*dtos.CreateDocumentResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.create_document").Logger()

	// Get the authenticated user context
	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: "Authentication required",
			Type:    "AUTHENTICATION_REQUIRED",
		}
	}

	result, err := ctrl.DocumentApp.CreateDocument(docDto.CreateDocumentInput{
		Name:     "New Document",
		UserId:   authCtx.UserID,
		SpaceId:  req.SpaceId,
		Content:  []docDto.Block{},
		ParentId: req.ParentId,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create document")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to create document",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	if result.Document == nil {
		logger.Error().Msg("document creation returned nil document")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to create document",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	resp := &dtos.CreateDocumentResponse{}
	resp.Id = result.Document.Id
	resp.Name = result.Document.Name
	resp.Slug = result.Document.Slug
	resp.ParentId = result.Document.ParentId
	resp.SpaceId = result.Document.SpaceId
	resp.Content = []dtos.Block{}
	resp.Config = dtos.DocumentConfig{}
	resp.Metadata = map[string]any{}
	resp.CreatedAt = result.Document.CreatedAt
	resp.UpdatedAt = result.Document.UpdatedAt

	// Map the space if it exists
	if result.Document.Space.Id != "" {
		resp.Space.Id = result.Document.Space.Id
		resp.Space.Name = result.Document.Space.Name
		resp.Space.Slug = result.Document.Space.Slug
		resp.Space.Icon = result.Document.Space.Icon
		resp.Space.IconColor = result.Document.Space.IconColor
		resp.Space.Type = string(result.Document.Space.Type)
		resp.Space.CreatedAt = result.Document.Space.CreatedAt
		resp.Space.UpdatedAt = result.Document.Space.UpdatedAt
	}

	return resp, nil
}

func (ctrl *Controller) UpdateDocument(ctx *fiber.Ctx, req dtos.UpdateDocumentRequest) (*dtos.UpdateDocumentResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.update_document").Logger()

	// Get the authenticated user context
	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: "Authentication required",
			Type:    "AUTHENTICATION_REQUIRED",
		}
	}

	// Convert DTO config to domain config if provided
	var domainConfig *domain.DocumentConfig
	if req.Config != nil {
		dc := domain.DocumentConfig{
			FullWidth:        req.Config.FullWidth,
			Icon:             req.Config.Icon,
			Lock:             req.Config.Lock,
			HeaderBackground: req.Config.HeaderBackground,
		}
		domainConfig = &dc
	}

	// Convert metadata map to JSONB if provided
	var domainMetadata *domain.JSONB
	if req.Metadata != nil {
		jsonb := domain.JSONB(*req.Metadata)
		domainMetadata = &jsonb
	}

	// Convert HTTP blocks to application blocks if provided
	var appContent *[]docDto.Block
	if req.Content != nil {
		blocks := make([]docDto.Block, len(*req.Content))
		for i, b := range *req.Content {
			blocks[i] = docDto.Block{
				ID:       b.ID,
				Type:     b.Type,
				Props:    b.Props,
				Content:  convertInlineContent(b.Content),
				Children: convertBlocks(b.Children),
			}
		}
		appContent = &blocks
	}

	result, err := ctrl.DocumentApp.UpdateDocument(docDto.UpdateDocumentInput{
		UserId:     authCtx.UserID,
		SpaceId:    req.SpaceId,
		DocumentId: req.Id,
		Name:       req.Name,
		Content:    appContent,
		ParentId:   req.ParentId,
		Config:     domainConfig,
		Metadata:   domainMetadata,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to update document")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to update document",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	resp := &dtos.UpdateDocumentResponse{}
	err = mapper.MapStructByFieldNames(result.Document, resp)
	if err != nil {
		logger.Error().Err(err).Msg("failed to map updated document to response DTO")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Details: "Failed to process updated document",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return resp, nil
}

func (ctrl *Controller) DeleteDocument(ctx *fiber.Ctx, req dtos.DeleteDocumentRequest) (*dtos.DeleteDocumentResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.delete_document").Logger()

	// Get the authenticated user context
	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Details: "Authentication required",
			Type:    "AUTHENTICATION_REQUIRED",
		}
	}

	var id *string
	var slug *string
	if validator.IsValidUUID(req.Identifier) {
		id = &req.Identifier
	} else {
		slug = &req.Identifier
	}

	if err := ctrl.DocumentApp.DeleteDocument(docDto.DeleteDocumentInput{
		UserId:     authCtx.UserID,
		SpaceId:    req.SpaceId,
		DocumentId: id,
		Slug:       slug,
	}); err != nil {
		// Best-effort error typing
		details := err.Error()
		switch {
		case strings.Contains(details, "insufficient permissions") || strings.Contains(details, "permission"):
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		case strings.Contains(details, "child") || strings.Contains(details, "children"):
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusConflict, Details: "Document has child pages", Type: "CONFLICT"}
		case strings.Contains(details, "not found"):
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Document not found", Type: "DOCUMENT_NOT_FOUND"}
		default:
			logger.Error().Err(err).Msg("failed to delete document")
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete document", Type: "INTERNAL_SERVER_ERROR"}
		}
	}

	return &dtos.DeleteDocumentResponse{Message: "Document deleted"}, nil
}

func (ctrl *Controller) MoveDocument(ctx *fiber.Ctx, req dtos.MoveDocumentRequest) (*dtos.MoveDocumentResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.move_document").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.DocumentApp.MoveDocument(docDto.MoveDocumentInput{
		UserId:      authCtx.UserID,
		SpaceId:     req.SpaceId,
		DocumentId:  req.Id,
		NewParentId: req.ParentId,
	})
	if err != nil {
		details := err.Error()
		switch {
		case strings.Contains(details, "insufficient permissions"):
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		case strings.Contains(details, "invalid move"):
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: details, Type: "BAD_REQUEST"}
		case strings.Contains(details, "not found"):
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Document or parent not found", Type: "NOT_FOUND"}
		default:
			logger.Error().Err(err).Msg("failed to move document")
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to move document", Type: "INTERNAL_SERVER_ERROR"}
		}
	}

	return &dtos.MoveDocumentResponse{Id: result.Document.Id, ParentId: result.Document.ParentId, SpaceId: result.Document.SpaceId}, nil
}

// convertInlineContent converts HTTP DTO inline content to application DTO
func convertInlineContent(content []dtos.InlineContent) []docDto.InlineContent {
	result := make([]docDto.InlineContent, len(content))
	for i, c := range content {
		result[i] = docDto.InlineContent{
			Type:   c.Type,
			Text:   c.Text,
			Href:   c.Href,
			Styles: c.Styles,
		}
	}
	return result
}

// convertBlocks recursively converts HTTP DTO blocks to application DTO blocks
func convertBlocks(blocks []dtos.Block) []docDto.Block {
	result := make([]docDto.Block, len(blocks))
	for i, b := range blocks {
		result[i] = docDto.Block{
			ID:       b.ID,
			Type:     b.Type,
			Props:    b.Props,
			Content:  convertInlineContent(b.Content),
			Children: convertBlocks(b.Children),
		}
	}
	return result
}

// Permission handlers

func (ctrl *Controller) ListDocumentPermissions(ctx *fiber.Ctx, req dtos.ListDocumentPermissionsRequest) (*dtos.ListDocumentPermissionsResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.list_permissions").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.DocumentApp.ListDocumentPermissions(docDto.ListDocumentPermissionsInput{
		RequesterId: authCtx.UserID,
		SpaceId:     req.SpaceId,
		DocumentId:  req.DocumentId,
	})
	if err != nil {
		switch err.Error() {
		case "forbidden":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		case "not_found":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Document not found", Type: "DOCUMENT_NOT_FOUND"}
		default:
			logger.Error().Err(err).Msg("failed to list document permissions")
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to list permissions", Type: "INTERNAL_SERVER_ERROR"}
		}
	}

	resp := &dtos.ListDocumentPermissionsResponse{Permissions: make([]dtos.DocumentPermission, len(result.Permissions))}
	for i, p := range result.Permissions {
		perm := dtos.DocumentPermission{
			Id:     p.Id,
			UserId: p.UserId,
			Role:   string(p.Role),
		}
		if p.User != nil {
			perm.Username = &p.User.Username
		}
		resp.Permissions[i] = perm
	}
	return resp, nil
}

func (ctrl *Controller) UpsertDocumentUserPermission(ctx *fiber.Ctx, req dtos.UpsertDocumentUserPermissionRequest) (*dtos.UpsertDocumentUserPermissionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.upsert_permission").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	var role domain.PermissionRole
	switch req.Role {
	case string(domain.PermissionRoleOwner):
		role = domain.PermissionRoleOwner
	case string(domain.PermissionRoleEditor):
		role = domain.PermissionRoleEditor
	case string(domain.PermissionRoleDenied):
		role = domain.PermissionRoleDenied
	default:
		role = domain.PermissionRoleViewer
	}

	if err := ctrl.DocumentApp.UpsertDocumentUserPermission(docDto.UpsertDocumentUserPermissionInput{
		RequesterId:  authCtx.UserID,
		SpaceId:      req.SpaceId,
		DocumentId:   req.DocumentId,
		TargetUserId: req.UserId,
		Role:         role,
	}); err != nil {
		switch err.Error() {
		case "forbidden":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		case "not_found":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Document not found", Type: "DOCUMENT_NOT_FOUND"}
		default:
			logger.Error().Err(err).Msg("failed to upsert document permission")
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to upsert permission", Type: "INTERNAL_SERVER_ERROR"}
		}
	}

	return &dtos.UpsertDocumentUserPermissionResponse{Message: "permission updated"}, nil
}

func (ctrl *Controller) DeleteDocumentUserPermission(ctx *fiber.Ctx, req dtos.DeleteDocumentUserPermissionRequest) (*dtos.DeleteDocumentUserPermissionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.delete_permission").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	if err := ctrl.DocumentApp.DeleteDocumentUserPermission(docDto.DeleteDocumentUserPermissionInput{
		RequesterId:  authCtx.UserID,
		SpaceId:      req.SpaceId,
		DocumentId:   req.DocumentId,
		TargetUserId: req.UserId,
	}); err != nil {
		switch err.Error() {
		case "forbidden":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		case "not_found":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Document not found", Type: "DOCUMENT_NOT_FOUND"}
		default:
			logger.Error().Err(err).Msg("failed to delete document permission")
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete permission", Type: "INTERNAL_SERVER_ERROR"}
		}
	}

	return &dtos.DeleteDocumentUserPermissionResponse{Message: "permission deleted"}, nil
}

// Trash handlers

func (ctrl *Controller) GetTrash(ctx *fiber.Ctx, req dtos.GetTrashRequest) (*dtos.GetTrashResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.get_trash").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.DocumentApp.GetTrash(docDto.GetTrashInput{
		UserId:  authCtx.UserID,
		SpaceId: req.SpaceId,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get trash")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to get trash", Type: "INTERNAL_SERVER_ERROR"}
	}

	resp := &dtos.GetTrashResponse{Documents: make([]dtos.TrashDocument, len(result.Documents))}
	for i, doc := range result.Documents {
		resp.Documents[i] = dtos.TrashDocument{
			Id:        doc.Id,
			Name:      doc.Name,
			Slug:      doc.Slug,
			DeletedAt: doc.DeletedAt.Time,
		}
	}

	return resp, nil
}

func (ctrl *Controller) RestoreDocument(ctx *fiber.Ctx, req dtos.RestoreDocumentRequest) (*dtos.RestoreDocumentResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.restore_document").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.DocumentApp.RestoreDocument(docDto.RestoreDocumentInput{
		UserId:     authCtx.UserID,
		SpaceId:    req.SpaceId,
		DocumentId: req.DocumentId,
	})
	if err != nil {
		switch {
		case err.Error() == "access denied: insufficient permissions to restore document":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		case err.Error() == "document is not deleted":
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Document is not in trash", Type: "BAD_REQUEST"}
		default:
			logger.Error().Err(err).Msg("failed to restore document")
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to restore document", Type: "INTERNAL_SERVER_ERROR"}
		}
	}

	return &dtos.RestoreDocumentResponse{Message: "Document restored successfully"}, nil
}

// Public sharing handlers

func (ctrl *Controller) SetPublic(ctx *fiber.Ctx, req dtos.SetPublicRequest) (*dtos.SetPublicResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.set_public").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.DocumentApp.SetPublic(docDto.SetPublicInput{
		UserId:     authCtx.UserID,
		SpaceId:    req.SpaceId,
		DocumentId: req.DocumentId,
		Public:     req.Public,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		logger.Error().Err(err).Msg("failed to set document public status")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to update document", Type: "INTERNAL_SERVER_ERROR"}
	}

	message := "Document is now private"
	if req.Public {
		message = "Document is now public"
	}

	return &dtos.SetPublicResponse{Message: message, Public: req.Public}, nil
}

func (ctrl *Controller) GetPublicDocument(ctx *fiber.Ctx, req dtos.GetPublicDocumentRequest) (*dtos.GetDocumentResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.get_public_document").Logger()

	var id *string
	var slug *string

	if validator.IsValidUUID(req.Identifier) {
		id = &req.Identifier
	} else {
		slug = &req.Identifier
	}

	result, err := ctrl.DocumentApp.GetPublicDocument(docDto.GetPublicDocumentInput{
		SpaceId:    req.SpaceId,
		DocumentId: id,
		Slug:       slug,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get public document")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Document not found or not public", Type: "DOCUMENT_NOT_FOUND"}
	}

	resp := &dtos.GetDocumentResponse{}
	err = mapper.MapStructByFieldNames(result.Document, resp)
	if err != nil {
		logger.Error().Err(err).Msg("failed to map document to response DTO")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to process document", Type: "INTERNAL_SERVER_ERROR"}
	}

	return resp, nil
}

// Comment handlers

func (ctrl *Controller) GetComments(ctx *fiber.Ctx, req dtos.GetCommentsRequest) (*dtos.GetCommentsResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.get_comments").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.DocumentApp.GetComments(docDto.GetCommentsInput{
		UserId:     authCtx.UserID,
		DocumentId: req.DocumentId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		logger.Error().Err(err).Msg("failed to get comments")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to get comments", Type: "INTERNAL_SERVER_ERROR"}
	}

	resp := &dtos.GetCommentsResponse{Comments: make([]dtos.CommentResponse, len(result.Comments))}
	for i, c := range result.Comments {
		resp.Comments[i] = dtos.CommentResponse{
			Id:        c.Id,
			UserId:    c.UserId,
			UserName:  c.UserName,
			ParentId:  c.ParentId,
			Content:   c.Content,
			BlockId:   c.BlockId,
			Resolved:  c.Resolved,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}

	return resp, nil
}

func (ctrl *Controller) CreateComment(ctx *fiber.Ctx, req dtos.CreateCommentRequestWithParams) (*dtos.CreateCommentResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.create_comment").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.DocumentApp.CreateComment(docDto.CreateCommentInput{
		UserId:     authCtx.UserID,
		DocumentId: req.DocumentId,
		ParentId:   req.ParentId,
		Content:    req.Content,
		BlockId:    req.BlockId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		logger.Error().Err(err).Msg("failed to create comment")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to create comment", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.CreateCommentResponse{CommentId: result.CommentId}, nil
}

func (ctrl *Controller) UpdateComment(ctx *fiber.Ctx, req dtos.UpdateCommentRequestWithParams) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.update_comment").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.DocumentApp.UpdateComment(docDto.UpdateCommentInput{
		UserId:    authCtx.UserID,
		CommentId: req.CommentId,
		Content:   req.Content,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Only the comment author can update it", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Comment not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to update comment")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to update comment", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "Comment updated"}, nil
}

func (ctrl *Controller) DeleteComment(ctx *fiber.Ctx, req dtos.DeleteCommentRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.delete_comment").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.DocumentApp.DeleteComment(docDto.DeleteCommentInput{
		UserId:    authCtx.UserID,
		CommentId: req.CommentId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Only the comment author can delete it", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Comment not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to delete comment")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete comment", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "Comment deleted"}, nil
}

func (ctrl *Controller) ResolveComment(ctx *fiber.Ctx, req dtos.ResolveCommentRequestWithParams) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.resolve_comment").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.DocumentApp.ResolveComment(docDto.ResolveCommentInput{
		UserId:    authCtx.UserID,
		CommentId: req.CommentId,
		Resolved:  req.Resolved,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Comment not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to resolve comment")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to resolve comment", Type: "INTERNAL_SERVER_ERROR"}
	}

	message := "Comment unresolved"
	if req.Resolved {
		message = "Comment resolved"
	}

	return &dtos.MessageResponse{Message: message}, nil
}

// Search handler

func (ctrl *Controller) SearchDocuments(ctx *fiber.Ctx, req dtos.SearchRequest) (*dtos.SearchResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.search").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	if len(req.Query) < 2 {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Query must be at least 2 characters", Type: "BAD_REQUEST"}
	}

	result, err := ctrl.DocumentApp.Search(docDto.SearchInput{
		UserId:  authCtx.UserID,
		Query:   req.Query,
		SpaceId: req.SpaceId,
		Limit:   req.Limit,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to search documents")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to search documents", Type: "INTERNAL_SERVER_ERROR"}
	}

	resp := &dtos.SearchResponse{Results: make([]dtos.SearchResultItem, len(result.Results))}
	for i, r := range result.Results {
		resp.Results[i] = dtos.SearchResultItem{
			Id:        r.Id,
			Name:      r.Name,
			Slug:      r.Slug,
			SpaceId:   r.SpaceId,
			SpaceName: r.SpaceName,
			Icon:      r.Icon,
			UpdatedAt: r.UpdatedAt,
		}
	}

	return resp, nil
}

// Version handlers

func (ctrl *Controller) ListVersions(ctx *fiber.Ctx, req dtos.ListVersionsRequest) (*dtos.ListVersionsResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.list_versions").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	logger.Debug().Str("spaceId", req.SpaceId).Str("documentId", req.DocumentId).Str("userId", authCtx.UserID).Msg("ListVersions called")

	result, err := ctrl.DocumentApp.ListVersions(docDto.ListVersionsInput{
		UserId:     authCtx.UserID,
		SpaceId:    req.SpaceId,
		DocumentId: req.DocumentId,
		Limit:      req.Limit,
		Offset:     req.Offset,
	})
	if err != nil {
		logger.Error().Err(err).Str("spaceId", req.SpaceId).Str("documentId", req.DocumentId).Msg("failed to list versions")
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to list versions", Type: "INTERNAL_SERVER_ERROR"}
	}

	resp := &dtos.ListVersionsResponse{
		Versions:   make([]dtos.VersionItem, len(result.Versions)),
		TotalCount: result.TotalCount,
	}
	for i, v := range result.Versions {
		resp.Versions[i] = dtos.VersionItem{
			Id:          v.Id,
			Version:     v.Version,
			Name:        v.Name,
			Description: v.Description,
			UserId:      v.UserId,
			UserName:    v.UserName,
			CreatedAt:   v.CreatedAt,
		}
	}

	return resp, nil
}

func (ctrl *Controller) GetVersion(ctx *fiber.Ctx, req dtos.GetVersionRequest) (*dtos.GetVersionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.get_version").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.DocumentApp.GetVersion(docDto.GetVersionInput{
		UserId:    authCtx.UserID,
		VersionId: req.VersionId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Version not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to get version")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to get version", Type: "INTERNAL_SERVER_ERROR"}
	}

	// Convert content blocks
	content := make([]dtos.Block, len(result.Content))
	for i, b := range result.Content {
		content[i] = dtos.Block{
			ID:       b.ID,
			Type:     b.Type,
			Props:    b.Props,
			Content:  convertToHttpInlineContent(b.Content),
			Children: convertToHttpBlocks(b.Children),
		}
	}

	return &dtos.GetVersionResponse{
		Id:         result.Id,
		Version:    result.Version,
		DocumentId: result.DocumentId,
		Name:       result.Name,
		Content:    content,
		Config: dtos.VersionConfig{
			FullWidth:        result.Config.FullWidth,
			Icon:             result.Config.Icon,
			Lock:             result.Config.Lock,
			HeaderBackground: result.Config.HeaderBackground,
		},
		Description: result.Description,
		UserId:      result.UserId,
		UserName:    result.UserName,
		CreatedAt:   result.CreatedAt,
	}, nil
}

func (ctrl *Controller) RestoreVersion(ctx *fiber.Ctx, req dtos.RestoreVersionRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.restore_version").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.DocumentApp.RestoreVersion(docDto.RestoreVersionInput{
		UserId:    authCtx.UserID,
		VersionId: req.VersionId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Version not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to restore version")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to restore version", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "Version restored successfully"}, nil
}

func (ctrl *Controller) CreateVersion(ctx *fiber.Ctx, req dtos.CreateVersionRequest) (*dtos.CreateVersionResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.document.create_version").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.DocumentApp.CreateVersion(docDto.CreateVersionInput{
		UserId:      authCtx.UserID,
		DocumentId:  req.DocumentId,
		Description: req.Description,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		logger.Error().Err(err).Msg("failed to create version")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to create version", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.CreateVersionResponse{
		VersionId: result.VersionId,
		Version:   result.Version,
	}, nil
}

// Helper functions for converting blocks
func convertToHttpInlineContent(content []docDto.InlineContent) []dtos.InlineContent {
	result := make([]dtos.InlineContent, len(content))
	for i, c := range content {
		result[i] = dtos.InlineContent{
			Type:   c.Type,
			Text:   c.Text,
			Href:   c.Href,
			Styles: c.Styles,
		}
	}
	return result
}

func convertToHttpBlocks(blocks []docDto.Block) []dtos.Block {
	result := make([]dtos.Block, len(blocks))
	for i, b := range blocks {
		result[i] = dtos.Block{
			ID:       b.ID,
			Type:     b.Type,
			Props:    b.Props,
			Content:  convertToHttpInlineContent(b.Content),
			Children: convertToHttpBlocks(b.Children),
		}
	}
	return result
}
