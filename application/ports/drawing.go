package ports

import (
	"github.com/labbs/nexo/application/drawing/dto"
)

type DrawingPort interface {
	CreateDrawing(input dto.CreateDrawingInput) (*dto.CreateDrawingOutput, error)
	ListDrawings(input dto.ListDrawingsInput) (*dto.ListDrawingsOutput, error)
	GetDrawing(input dto.GetDrawingInput) (*dto.GetDrawingOutput, error)
	UpdateDrawing(input dto.UpdateDrawingInput) error
	MoveDrawing(input dto.MoveDrawingInput) (*dto.MoveDrawingOutput, error)
	DeleteDrawing(input dto.DeleteDrawingInput) error
	ListDrawingPermissions(input dto.ListDrawingPermissionsInput) (*dto.ListDrawingPermissionsOutput, error)
	UpsertDrawingUserPermission(input dto.UpsertDrawingUserPermissionInput) error
	DeleteDrawingUserPermission(input dto.DeleteDrawingUserPermissionInput) error
}
