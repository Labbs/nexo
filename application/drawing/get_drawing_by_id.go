package drawing

import (
	"fmt"

	"github.com/labbs/nexo/application/drawing/dto"
)

// GetDrawingById returns a drawing by its ID without authorization checks.
// This is used internally by other applications that need to look up a drawing.
func (app *DrawingApplication) GetDrawingById(input dto.GetDrawingByIdInput) (*dto.GetDrawingByIdOutput, error) {
	drawing, err := app.DrawingPers.GetById(input.DrawingId)
	if err != nil {
		return nil, fmt.Errorf("drawing not found: %w", err)
	}

	detail := &dto.DrawingDetail{
		Id:      drawing.Id,
		SpaceId: drawing.SpaceId,
	}

	return &dto.GetDrawingByIdOutput{Drawing: detail}, nil
}
