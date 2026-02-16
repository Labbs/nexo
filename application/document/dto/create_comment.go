package dto

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
