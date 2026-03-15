package dto

// DrawingDetail contains the drawing data needed by other applications
type DrawingDetail struct {
	Id      string
	SpaceId string
}

type GetDrawingByIdInput struct {
	DrawingId string
}

type GetDrawingByIdOutput struct {
	Drawing *DrawingDetail
}
