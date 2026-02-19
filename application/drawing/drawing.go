package drawing

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/drawing/dto"
	permissionDto "github.com/labbs/nexo/application/permission/dto"
	"github.com/labbs/nexo/application/ports"
	spaceDto "github.com/labbs/nexo/application/space/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type DrawingApplication struct {
	Config        config.Config
	Logger        zerolog.Logger
	DrawingPers   domain.DrawingPers
	SpaceApp      ports.SpacePort
	PermissionApp ports.PermissionPort
}

func NewDrawingApplication(config config.Config, logger zerolog.Logger, drawingPers domain.DrawingPers) *DrawingApplication {
	return &DrawingApplication{
		Config:      config,
		Logger:      logger,
		DrawingPers: drawingPers,
	}
}

func (app *DrawingApplication) CreateDrawing(input dto.CreateDrawingInput) (*dto.CreateDrawingOutput, error) {
	// Verify user has access to the space
	spaceResult, err := app.SpaceApp.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: input.SpaceId})
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if spaceResult.Space.GetUserRole(input.UserId) == nil {
		return nil, fmt.Errorf("access denied")
	}

	// Convert elements to JSONBArray
	var elements domain.JSONBArray
	if input.Elements != nil {
		elementsJSON, err := json.Marshal(input.Elements)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal elements: %w", err)
		}
		json.Unmarshal(elementsJSON, &elements)
	} else {
		elements = domain.JSONBArray{}
	}

	// Convert appState to JSONB
	var appState domain.JSONB
	if input.AppState != nil {
		appStateJSON, err := json.Marshal(input.AppState)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal appState: %w", err)
		}
		json.Unmarshal(appStateJSON, &appState)
	} else {
		appState = domain.JSONB{}
	}

	// Convert files to JSONB
	var files domain.JSONB
	if input.Files != nil {
		filesJSON, err := json.Marshal(input.Files)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal files: %w", err)
		}
		json.Unmarshal(filesJSON, &files)
	} else {
		files = domain.JSONB{}
	}

	drawing := &domain.Drawing{
		Id:         uuid.New().String(),
		SpaceId:    input.SpaceId,
		DocumentId: input.DocumentId,
		Name:       input.Name,
		Icon:       input.Icon,
		Elements:   elements,
		AppState:   appState,
		Files:      files,
		Thumbnail:  input.Thumbnail,
		CreatedBy:  input.UserId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := app.DrawingPers.Create(drawing); err != nil {
		return nil, fmt.Errorf("failed to create drawing: %w", err)
	}

	// Auto-create owner permission for the creator
	// This ensures they retain access and can manage permissions even if their space role is downgraded
	if err := app.PermissionApp.AssignOwnerPermission(permissionDto.AssignOwnerPermissionInput{
		ResourceType: "drawing",
		ResourceId:   drawing.Id,
		UserId:       input.UserId,
		Role:         "owner",
	}); err != nil {
		// Log but don't fail - the drawing is already created
		app.Logger.Warn().Err(err).Str("drawing_id", drawing.Id).Str("user_id", input.UserId).Msg("failed to create creator permission")
	}

	return &dto.CreateDrawingOutput{
		Id:        drawing.Id,
		Name:      drawing.Name,
		CreatedAt: drawing.CreatedAt,
	}, nil
}

func (app *DrawingApplication) ListDrawings(input dto.ListDrawingsInput) (*dto.ListDrawingsOutput, error) {
	// Verify user has access to the space
	spaceResult, err := app.SpaceApp.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: input.SpaceId})
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if spaceResult.Space.GetUserRole(input.UserId) == nil {
		return nil, fmt.Errorf("access denied")
	}

	drawings, err := app.DrawingPers.GetBySpaceId(input.SpaceId)
	if err != nil {
		return nil, fmt.Errorf("failed to list drawings: %w", err)
	}

	output := &dto.ListDrawingsOutput{
		Drawings: make([]dto.DrawingItem, len(drawings)),
	}

	for i, d := range drawings {
		output.Drawings[i] = dto.DrawingItem{
			Id:         d.Id,
			DocumentId: d.DocumentId,
			Name:       d.Name,
			Icon:       d.Icon,
			Thumbnail:  d.Thumbnail,
			CreatedBy:  d.User.Username,
			CreatedAt:  d.CreatedAt,
			UpdatedAt:  d.UpdatedAt,
		}
	}

	return output, nil
}

func (app *DrawingApplication) GetDrawing(input dto.GetDrawingInput) (*dto.GetDrawingOutput, error) {
	drawing, err := app.DrawingPers.GetById(input.DrawingId)
	if err != nil {
		return nil, fmt.Errorf("drawing not found: %w", err)
	}

	// Verify user has access to the space
	spaceResult, err := app.SpaceApp.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: drawing.SpaceId})
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if spaceResult.Space.GetUserRole(input.UserId) == nil {
		return nil, fmt.Errorf("access denied")
	}

	// Convert JSONB to slices/maps
	var elements []interface{}
	if drawing.Elements != nil {
		elementsJSON, _ := json.Marshal(drawing.Elements)
		json.Unmarshal(elementsJSON, &elements)
	}

	var appState map[string]interface{}
	if drawing.AppState != nil {
		appStateJSON, _ := json.Marshal(drawing.AppState)
		json.Unmarshal(appStateJSON, &appState)
	}

	var files map[string]interface{}
	if drawing.Files != nil {
		filesJSON, _ := json.Marshal(drawing.Files)
		json.Unmarshal(filesJSON, &files)
	}

	return &dto.GetDrawingOutput{
		Id:         drawing.Id,
		SpaceId:    drawing.SpaceId,
		DocumentId: drawing.DocumentId,
		Name:       drawing.Name,
		Icon:       drawing.Icon,
		Elements:   elements,
		AppState:   appState,
		Files:      files,
		Thumbnail:  drawing.Thumbnail,
		CreatedBy:  drawing.User.Username,
		CreatedAt:  drawing.CreatedAt,
		UpdatedAt:  drawing.UpdatedAt,
	}, nil
}

func (app *DrawingApplication) UpdateDrawing(input dto.UpdateDrawingInput) error {
	drawing, err := app.DrawingPers.GetById(input.DrawingId)
	if err != nil {
		return fmt.Errorf("drawing not found: %w", err)
	}

	// Verify user has access to the space
	spaceResult, err := app.SpaceApp.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: drawing.SpaceId})
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	if spaceResult.Space.GetUserRole(input.UserId) == nil {
		return fmt.Errorf("access denied")
	}

	if input.Name != nil {
		drawing.Name = *input.Name
	}

	if input.Icon != nil {
		drawing.Icon = *input.Icon
	}

	if input.Elements != nil {
		elementsJSON, err := json.Marshal(input.Elements)
		if err != nil {
			return fmt.Errorf("failed to marshal elements: %w", err)
		}
		var elements domain.JSONBArray
		json.Unmarshal(elementsJSON, &elements)
		drawing.Elements = elements
	}

	if input.AppState != nil {
		appStateJSON, err := json.Marshal(input.AppState)
		if err != nil {
			return fmt.Errorf("failed to marshal appState: %w", err)
		}
		var appState domain.JSONB
		json.Unmarshal(appStateJSON, &appState)
		drawing.AppState = appState
	}

	if input.Files != nil {
		filesJSON, err := json.Marshal(input.Files)
		if err != nil {
			return fmt.Errorf("failed to marshal files: %w", err)
		}
		var files domain.JSONB
		json.Unmarshal(filesJSON, &files)
		drawing.Files = files
	}

	if input.Thumbnail != nil {
		drawing.Thumbnail = *input.Thumbnail
	}

	drawing.UpdatedAt = time.Now()

	if err := app.DrawingPers.Update(drawing); err != nil {
		return fmt.Errorf("failed to update drawing: %w", err)
	}

	return nil
}

func (app *DrawingApplication) MoveDrawing(input dto.MoveDrawingInput) (*dto.MoveDrawingOutput, error) {
	drawing, err := app.DrawingPers.GetById(input.DrawingId)
	if err != nil {
		return nil, fmt.Errorf("drawing not found: %w", err)
	}

	// Verify user has access to the space
	spaceResult, err := app.SpaceApp.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: drawing.SpaceId})
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if spaceResult.Space.GetUserRole(input.UserId) == nil {
		return nil, fmt.Errorf("access denied")
	}

	drawing.DocumentId = input.DocumentId
	drawing.UpdatedAt = time.Now()

	if err := app.DrawingPers.Update(drawing); err != nil {
		return nil, fmt.Errorf("failed to move drawing: %w", err)
	}

	return &dto.MoveDrawingOutput{
		Id:         drawing.Id,
		DocumentId: drawing.DocumentId,
	}, nil
}

func (app *DrawingApplication) DeleteDrawing(input dto.DeleteDrawingInput) error {
	drawing, err := app.DrawingPers.GetById(input.DrawingId)
	if err != nil {
		return fmt.Errorf("drawing not found: %w", err)
	}

	// Verify user has access to the space
	spaceResult, err := app.SpaceApp.GetSpaceById(spaceDto.GetSpaceByIdInput{SpaceId: drawing.SpaceId})
	if err != nil {
		return fmt.Errorf("space not found: %w", err)
	}

	if spaceResult.Space.GetUserRole(input.UserId) == nil {
		return fmt.Errorf("access denied")
	}

	if err := app.DrawingPers.Delete(input.DrawingId); err != nil {
		return fmt.Errorf("failed to delete drawing: %w", err)
	}

	return nil
}
