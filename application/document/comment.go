package document

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
)

func (app *DocumentApp) CreateComment(input dto.CreateCommentInput) (*dto.CreateCommentOutput, error) {
	// Verify user has access to the document
	_, err := app.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions("", &input.DocumentId, nil, input.UserId)
	if err != nil {
		return nil, fmt.Errorf("document not found or access denied: %w", err)
	}

	comment := &domain.Comment{
		Id:         uuid.New().String(),
		DocumentId: input.DocumentId,
		UserId:     input.UserId,
		ParentId:   input.ParentId,
		Content:    input.Content,
		BlockId:    input.BlockId,
		Resolved:   false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := app.CommentPers.Create(comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &dto.CreateCommentOutput{
		CommentId: comment.Id,
	}, nil
}

func (app *DocumentApp) GetComments(input dto.GetCommentsInput) (*dto.GetCommentsOutput, error) {
	// Verify user has access to the document
	_, err := app.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions("", &input.DocumentId, nil, input.UserId)
	if err != nil {
		return nil, fmt.Errorf("document not found or access denied: %w", err)
	}

	comments, err := app.CommentPers.GetByDocumentId(input.DocumentId)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	output := &dto.GetCommentsOutput{
		Comments: make([]dto.CommentOutput, len(comments)),
	}

	for i, c := range comments {
		output.Comments[i] = dto.CommentOutput{
			Id:        c.Id,
			UserId:    c.UserId,
			UserName:  c.User.Username,
			ParentId:  c.ParentId,
			Content:   c.Content,
			BlockId:   c.BlockId,
			Resolved:  c.Resolved,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}

	return output, nil
}

func (app *DocumentApp) UpdateComment(input dto.UpdateCommentInput) error {
	comment, err := app.CommentPers.GetById(input.CommentId)
	if err != nil {
		return fmt.Errorf("comment not found: %w", err)
	}

	// Only the comment author can update it
	if comment.UserId != input.UserId {
		return fmt.Errorf("access denied: only the comment author can update it")
	}

	comment.Content = input.Content
	comment.UpdatedAt = time.Now()

	if err := app.CommentPers.Update(comment); err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

func (app *DocumentApp) DeleteComment(input dto.DeleteCommentInput) error {
	comment, err := app.CommentPers.GetById(input.CommentId)
	if err != nil {
		return fmt.Errorf("comment not found: %w", err)
	}

	// Only the comment author can delete it
	if comment.UserId != input.UserId {
		return fmt.Errorf("access denied: only the comment author can delete it")
	}

	if err := app.CommentPers.Delete(input.CommentId); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

func (app *DocumentApp) ResolveComment(input dto.ResolveCommentInput) error {
	comment, err := app.CommentPers.GetById(input.CommentId)
	if err != nil {
		return fmt.Errorf("comment not found: %w", err)
	}

	// Verify user has access to the document (any user with document access can resolve)
	_, err = app.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions("", &comment.DocumentId, nil, input.UserId)
	if err != nil {
		return fmt.Errorf("access denied: %w", err)
	}

	if err := app.CommentPers.Resolve(input.CommentId, input.Resolved); err != nil {
		return fmt.Errorf("failed to resolve comment: %w", err)
	}

	return nil
}
