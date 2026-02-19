package ports

import (
	"github.com/labbs/nexo/application/drawing/dto"
)

type DrawingPort interface {
	CreateDrawing(input dto.CreateDrawingInput) (*dto.CreateDrawingOutput, error)
	ListDrawings(input dto.ListDrawingsInput) (*dto.ListDrawingsOutput, error)
	GetDrawing(input dto.GetDrawingInput) (*dto.GetDrawingOutput, error)
	GetDrawingById(input dto.GetDrawingByIdInput) (*dto.GetDrawingByIdOutput, error)
	UpdateDrawing(input dto.UpdateDrawingInput) error
	MoveDrawing(input dto.MoveDrawingInput) (*dto.MoveDrawingOutput, error)
	DeleteDrawing(input dto.DeleteDrawingInput) error
}
