package dto

import "time"

// Create comment
type CreateCommentInput struct {
	UserId     string
	DocumentId string
	ParentId   *string
	Content    string
	BlockId    *string
}

type CreateCommentOutput struct {
	CommentId string
}

// Get comments
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

// Update comment
type UpdateCommentInput struct {
	UserId    string
	CommentId string
	Content   string
}

// Delete comment
type DeleteCommentInput struct {
	UserId    string
	CommentId string
}

// Resolve comment
type ResolveCommentInput struct {
	UserId    string
	CommentId string
	Resolved  bool
}
