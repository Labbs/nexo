package dto

import "time"

type GetCommentsInput struct {
	UserId     string
	DocumentId string
}

type CommentOutput struct {
	Id        string
	UserId    string
	UserName  string
	ParentId  *string
	Content   string
	BlockId   *string
	Resolved  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetCommentsOutput struct {
	Comments []CommentOutput
}
