package dto

type MoveDrawingInput struct {
	UserId     string
	DrawingId  string
	DocumentId *string // nil = move to root (no parent document)
}

type MoveDrawingOutput struct {
	Id         string
	DocumentId *string
}
