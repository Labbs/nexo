package dtos

import "time"

// Request DTOs

type GetCommentsRequest struct {
	SpaceId    string `path:"space_id"`
	DocumentId string `path:"document_id"`
}

type CreateCommentRequest struct {
	ParentId *string `json:"parent_id,omitempty"`
	Content  string  `json:"content"`
	BlockId  *string `json:"block_id,omitempty"`
}

type CreateCommentRequestWithParams struct {
	SpaceId    string  `path:"space_id"`
	DocumentId string  `path:"document_id"`
	ParentId   *string `json:"parent_id,omitempty"`
	Content    string  `json:"content"`
	BlockId    *string `json:"block_id,omitempty"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

type UpdateCommentRequestWithParams struct {
	SpaceId    string `path:"space_id"`
	DocumentId string `path:"document_id"`
	CommentId  string `path:"comment_id"`
	Content    string `json:"content"`
}

type DeleteCommentRequest struct {
	SpaceId    string `path:"space_id"`
	DocumentId string `path:"document_id"`
	CommentId  string `path:"comment_id"`
}

type ResolveCommentRequest struct {
	Resolved bool `json:"resolved"`
}

type ResolveCommentRequestWithParams struct {
	SpaceId    string `path:"space_id"`
	DocumentId string `path:"document_id"`
	CommentId  string `path:"comment_id"`
	Resolved   bool   `json:"resolved"`
}

// Response DTOs

type CommentResponse struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	ParentId  *string   `json:"parent_id,omitempty"`
	Content   string    `json:"content"`
	BlockId   *string   `json:"block_id,omitempty"`
	Resolved  bool      `json:"resolved"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetCommentsResponse struct {
	Comments []CommentResponse `json:"comments"`
}

type CreateCommentResponse struct {
	CommentId string `json:"comment_id"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
