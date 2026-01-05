package database

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberoapi "github.com/labbs/fiber-oapi"
	databaseDto "github.com/labbs/nexo/application/database/dto"
	"github.com/labbs/nexo/interfaces/http/v1/database/dtos"
)

func (ctrl *Controller) CreateDatabase(ctx *fiber.Ctx, req dtos.CreateDatabaseRequest) (*dtos.CreateDatabaseResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.create").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	if req.Name == "" {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Name is required", Type: "BAD_REQUEST"}
	}

	if req.SpaceId == "" {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Space ID is required", Type: "BAD_REQUEST"}
	}

	// Convert schema
	schema := make([]databaseDto.PropertySchema, len(req.Schema))
	for i, s := range req.Schema {
		schema[i] = databaseDto.PropertySchema{
			Id:      s.Id,
			Name:    s.Name,
			Type:    s.Type,
			Options: s.Options,
		}
	}

	result, err := ctrl.DatabaseApp.CreateDatabase(databaseDto.CreateDatabaseInput{
		UserId:      authCtx.UserID,
		SpaceId:     req.SpaceId,
		DocumentId:  req.DocumentId,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Schema:      schema,
		Type:        req.Type,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		logger.Error().Err(err).Msg("failed to create database")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to create database", Type: "INTERNAL_SERVER_ERROR"}
	}

	respSchema := make([]dtos.PropertySchema, len(result.Schema))
	for i, s := range result.Schema {
		respSchema[i] = dtos.PropertySchema{
			Id:      s.Id,
			Name:    s.Name,
			Type:    s.Type,
			Options: s.Options,
		}
	}

	return &dtos.CreateDatabaseResponse{
		Id:          result.Id,
		Name:        result.Name,
		Description: result.Description,
		Icon:        result.Icon,
		Schema:      respSchema,
		DefaultView: result.DefaultView,
		Type:        result.Type,
		CreatedAt:   result.CreatedAt,
	}, nil
}

func (ctrl *Controller) ListDatabases(ctx *fiber.Ctx, req dtos.ListDatabasesRequest) (*dtos.ListDatabasesResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.list").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	if req.SpaceId == "" {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Space ID is required", Type: "BAD_REQUEST"}
	}

	result, err := ctrl.DatabaseApp.ListDatabases(databaseDto.ListDatabasesInput{
		UserId:  authCtx.UserID,
		SpaceId: req.SpaceId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		logger.Error().Err(err).Msg("failed to list databases")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to list databases", Type: "INTERNAL_SERVER_ERROR"}
	}

	resp := &dtos.ListDatabasesResponse{Databases: make([]dtos.DatabaseItem, len(result.Databases))}
	for i, db := range result.Databases {
		resp.Databases[i] = dtos.DatabaseItem{
			Id:          db.Id,
			DocumentId:  db.DocumentId,
			Name:        db.Name,
			Description: db.Description,
			Icon:        db.Icon,
			Type:        db.Type,
			RowCount:    db.RowCount,
			CreatedBy:   db.CreatedBy,
			CreatedAt:   db.CreatedAt,
			UpdatedAt:   db.UpdatedAt,
		}
	}

	return resp, nil
}

func (ctrl *Controller) GetDatabase(ctx *fiber.Ctx, req dtos.GetDatabaseRequest) (*dtos.GetDatabaseResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.get").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.DatabaseApp.GetDatabase(databaseDto.GetDatabaseInput{
		UserId:     authCtx.UserID,
		DatabaseId: req.DatabaseId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Database not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to get database")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to get database", Type: "INTERNAL_SERVER_ERROR"}
	}

	schema := make([]dtos.PropertySchema, len(result.Schema))
	for i, s := range result.Schema {
		schema[i] = dtos.PropertySchema{
			Id:      s.Id,
			Name:    s.Name,
			Type:    s.Type,
			Options: s.Options,
		}
	}

	views := make([]dtos.ViewConfig, len(result.Views))
	for i, v := range result.Views {
		sort := make([]dtos.SortConfig, len(v.Sort))
		for j, s := range v.Sort {
			sort[j] = dtos.SortConfig{
				PropertyId: s.PropertyId,
				Direction:  s.Direction,
			}
		}
		views[i] = dtos.ViewConfig{
			Id:            v.Id,
			Name:          v.Name,
			Type:          v.Type,
			Filter:        v.Filter,
			Sort:          sort,
			Columns:       v.Columns,
			HiddenColumns: v.HiddenColumns,
		}
	}

	return &dtos.GetDatabaseResponse{
		Id:          result.Id,
		SpaceId:     result.SpaceId,
		DocumentId:  result.DocumentId,
		Name:        result.Name,
		Description: result.Description,
		Icon:        result.Icon,
		Schema:      schema,
		Views:       views,
		DefaultView: result.DefaultView,
		Type:        result.Type,
		CreatedBy:   result.CreatedBy,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (ctrl *Controller) UpdateDatabase(ctx *fiber.Ctx, req dtos.UpdateDatabaseRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.update").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	var schema *[]databaseDto.PropertySchema
	if req.Schema != nil {
		s := make([]databaseDto.PropertySchema, len(*req.Schema))
		for i, ps := range *req.Schema {
			s[i] = databaseDto.PropertySchema{
				Id:      ps.Id,
				Name:    ps.Name,
				Type:    ps.Type,
				Options: ps.Options,
			}
		}
		schema = &s
	}

	err = ctrl.DatabaseApp.UpdateDatabase(databaseDto.UpdateDatabaseInput{
		UserId:      authCtx.UserID,
		DatabaseId:  req.DatabaseId,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Schema:      schema,
		DefaultView: req.DefaultView,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Database not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to update database")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to update database", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "Database updated"}, nil
}

func (ctrl *Controller) DeleteDatabase(ctx *fiber.Ctx, req dtos.DeleteDatabaseRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.delete").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.DatabaseApp.DeleteDatabase(databaseDto.DeleteDatabaseInput{
		UserId:     authCtx.UserID,
		DatabaseId: req.DatabaseId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Database not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to delete database")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete database", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "Database deleted"}, nil
}

func (ctrl *Controller) GetAvailableTypes(ctx *fiber.Ctx, _ dtos.EmptyRequest) (*dtos.AvailableTypesResponse, *fiberoapi.ErrorResponse) {
	types := []dtos.TypeInfo{
		{Type: "title", Description: "Title property (required, one per database)"},
		{Type: "text", Description: "Plain text"},
		{Type: "number", Description: "Number with optional format"},
		{Type: "currency", Description: "Currency with symbol and locale formatting"},
		{Type: "select", Description: "Single select from options"},
		{Type: "multi_select", Description: "Multiple select from options"},
		{Type: "date", Description: "Date with optional time"},
		{Type: "checkbox", Description: "Boolean checkbox"},
		{Type: "url", Description: "URL link"},
		{Type: "email", Description: "Email address"},
		{Type: "phone", Description: "Phone number"},
		{Type: "image", Description: "Image URL with preview"},
		{Type: "relation", Description: "Relation to another database"},
		{Type: "rollup", Description: "Rollup from relation"},
		{Type: "formula", Description: "Calculated formula"},
		{Type: "created_time", Description: "Auto-populated creation time"},
		{Type: "updated_time", Description: "Auto-populated update time"},
		{Type: "created_by", Description: "Auto-populated creator"},
		{Type: "updated_by", Description: "Auto-populated last editor"},
		{Type: "files", Description: "File attachments"},
		{Type: "person", Description: "Person/user reference"},
	}

	return &dtos.AvailableTypesResponse{Types: types}, nil
}

// View handlers

func (ctrl *Controller) CreateView(ctx *fiber.Ctx, req dtos.CreateViewRequest) (*dtos.CreateViewResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.view.create").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	if req.Name == "" {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Name is required", Type: "BAD_REQUEST"}
	}
	if req.Type == "" {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Type is required", Type: "BAD_REQUEST"}
	}

	// Convert sort config
	sort := make([]databaseDto.SortConfig, len(req.Sort))
	for i, s := range req.Sort {
		sort[i] = databaseDto.SortConfig{
			PropertyId: s.PropertyId,
			Direction:  s.Direction,
		}
	}

	result, err := ctrl.DatabaseApp.CreateView(databaseDto.CreateViewInput{
		UserId:     authCtx.UserID,
		DatabaseId: req.DatabaseId,
		Name:       req.Name,
		Type:       req.Type,
		Filter:     req.Filter,
		Sort:       sort,
		Columns:    req.Columns,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Database not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to create view")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to create view", Type: "INTERNAL_SERVER_ERROR"}
	}

	respSort := make([]dtos.SortConfig, len(result.Sort))
	for i, s := range result.Sort {
		respSort[i] = dtos.SortConfig{
			PropertyId: s.PropertyId,
			Direction:  s.Direction,
		}
	}

	return &dtos.CreateViewResponse{
		Id:      result.Id,
		Name:    result.Name,
		Type:    result.Type,
		Filter:  result.Filter,
		Sort:    respSort,
		Columns: result.Columns,
	}, nil
}

func (ctrl *Controller) UpdateView(ctx *fiber.Ctx, req dtos.UpdateViewRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.view.update").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	// Convert sort config
	var sort []databaseDto.SortConfig
	if len(req.Sort) > 0 {
		sort = make([]databaseDto.SortConfig, len(req.Sort))
		for i, s := range req.Sort {
			sort[i] = databaseDto.SortConfig{
				PropertyId: s.PropertyId,
				Direction:  s.Direction,
			}
		}
	}

	err = ctrl.DatabaseApp.UpdateView(databaseDto.UpdateViewInput{
		UserId:        authCtx.UserID,
		DatabaseId:    req.DatabaseId,
		ViewId:        req.ViewId,
		Name:          req.Name,
		Filter:        req.Filter,
		Sort:          sort,
		Columns:       req.Columns,
		HiddenColumns: req.HiddenColumns,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "View not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to update view")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to update view", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "View updated"}, nil
}

func (ctrl *Controller) DeleteView(ctx *fiber.Ctx, req dtos.DeleteViewRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.view.delete").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.DatabaseApp.DeleteView(databaseDto.DeleteViewInput{
		UserId:     authCtx.UserID,
		DatabaseId: req.DatabaseId,
		ViewId:     req.ViewId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "View not found", Type: "NOT_FOUND"}
		}
		if strings.Contains(err.Error(), "cannot delete last view") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Cannot delete the last view", Type: "BAD_REQUEST"}
		}
		logger.Error().Err(err).Msg("failed to delete view")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete view", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "View deleted"}, nil
}

// Row handlers

func (ctrl *Controller) CreateRow(ctx *fiber.Ctx, req dtos.CreateRowRequest) (*dtos.CreateRowResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.row.create").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.DatabaseApp.CreateRow(databaseDto.CreateRowInput{
		UserId:        authCtx.UserID,
		DatabaseId:    req.DatabaseId,
		Properties:    req.Properties,
		Content:       req.Content,
		ShowInSidebar: req.ShowInSidebar,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Database not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to create row")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to create row", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.CreateRowResponse{
		Id:         result.Id,
		Properties: result.Properties,
		CreatedAt:  result.CreatedAt,
	}, nil
}

func (ctrl *Controller) ListRows(ctx *fiber.Ctx, req dtos.ListRowsRequest) (*dtos.ListRowsResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.row.list").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.DatabaseApp.ListRows(databaseDto.ListRowsInput{
		UserId:     authCtx.UserID,
		DatabaseId: req.DatabaseId,
		ViewId:     req.ViewId,
		Limit:      req.Limit,
		Offset:     req.Offset,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Database not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to list rows")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to list rows", Type: "INTERNAL_SERVER_ERROR"}
	}

	resp := &dtos.ListRowsResponse{
		Rows:       make([]dtos.RowItem, len(result.Rows)),
		TotalCount: result.TotalCount,
	}
	for i, row := range result.Rows {
		resp.Rows[i] = dtos.RowItem{
			Id:            row.Id,
			Properties:    row.Properties,
			ShowInSidebar: row.ShowInSidebar,
			CreatedBy:     row.CreatedBy,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		}
	}

	return resp, nil
}

func (ctrl *Controller) GetRow(ctx *fiber.Ctx, req dtos.GetRowRequest) (*dtos.GetRowResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.row.get").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	result, err := ctrl.DatabaseApp.GetRow(databaseDto.GetRowInput{
		UserId:     authCtx.UserID,
		DatabaseId: req.DatabaseId,
		RowId:      req.RowId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Row not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to get row")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to get row", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.GetRowResponse{
		Id:            result.Id,
		DatabaseId:    result.DatabaseId,
		Properties:    result.Properties,
		Content:       result.Content,
		ShowInSidebar: result.ShowInSidebar,
		CreatedBy:     result.CreatedBy,
		CreatedAt:     result.CreatedAt,
		UpdatedAt:     result.UpdatedAt,
	}, nil
}

func (ctrl *Controller) UpdateRow(ctx *fiber.Ctx, req dtos.UpdateRowRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.row.update").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.DatabaseApp.UpdateRow(databaseDto.UpdateRowInput{
		UserId:        authCtx.UserID,
		DatabaseId:    req.DatabaseId,
		RowId:         req.RowId,
		Properties:    req.Properties,
		Content:       req.Content,
		ShowInSidebar: req.ShowInSidebar,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Row not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to update row")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to update row", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "Row updated"}, nil
}

func (ctrl *Controller) DeleteRow(ctx *fiber.Ctx, req dtos.DeleteRowRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.row.delete").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	err = ctrl.DatabaseApp.DeleteRow(databaseDto.DeleteRowInput{
		UserId:     authCtx.UserID,
		DatabaseId: req.DatabaseId,
		RowId:      req.RowId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Row not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to delete row")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete row", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "Row deleted"}, nil
}

func (ctrl *Controller) BulkDeleteRows(ctx *fiber.Ctx, req dtos.BulkDeleteRowsRequest) (*dtos.MessageResponse, *fiberoapi.ErrorResponse) {
	requestId := ctx.Locals("requestid").(string)
	logger := ctrl.Logger.With().Str("request_id", requestId).Str("component", "http.api.v1.database.row.bulk_delete").Logger()

	authCtx, err := fiberoapi.GetAuthContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get auth context")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusUnauthorized, Details: "Authentication required", Type: "AUTHENTICATION_REQUIRED"}
	}

	if len(req.RowIds) == 0 {
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusBadRequest, Details: "Row IDs are required", Type: "BAD_REQUEST"}
	}

	err = ctrl.DatabaseApp.BulkDeleteRows(databaseDto.BulkDeleteRowsInput{
		UserId:     authCtx.UserID,
		DatabaseId: req.DatabaseId,
		RowIds:     req.RowIds,
	})
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusForbidden, Details: "Forbidden", Type: "FORBIDDEN"}
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusNotFound, Details: "Database not found", Type: "NOT_FOUND"}
		}
		logger.Error().Err(err).Msg("failed to delete rows")
		return nil, &fiberoapi.ErrorResponse{Code: fiber.StatusInternalServerError, Details: "Failed to delete rows", Type: "INTERNAL_SERVER_ERROR"}
	}

	return &dtos.MessageResponse{Message: "Rows deleted"}, nil
}
