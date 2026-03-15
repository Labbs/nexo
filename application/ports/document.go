package ports

import (
	"github.com/labbs/nexo/application/document/dto"
)

type DocumentPort interface {
	// CRUD
	CreateDocument(input dto.CreateDocumentInput) (*dto.CreateDocumentOutput, error)
	GetDocumentWithSpace(input dto.GetDocumentWithSpaceInput) (*dto.GetDocumentWithSpaceOutput, error)
	GetDocumentsFromSpaceWithUserPermissions(input dto.GetDocumentsFromSpaceInput) (*dto.GetDocumentsFromSpaceOutput, error)
	UpdateDocument(input dto.UpdateDocumentInput) (*dto.UpdateDocumentOutput, error)
	MoveDocument(input dto.MoveDocumentInput) (*dto.MoveDocumentOutput, error)
	DeleteDocument(input dto.DeleteDocumentInput) error
	GetDocumentByIdOrSlugWithUserPermissions(input dto.GetDocumentByIdOrSlugWithUserPermissionsInput) (*dto.GetDocumentByIdOrSlugWithUserPermissionsOutput, error)
	HasDocumentsInSpace(input dto.HasDocumentsInSpaceInput) (*dto.HasDocumentsInSpaceOutput, error)

	// Search
	Search(input dto.SearchInput) (*dto.SearchOutput, error)

	// Reorder
	ReorderDocuments(input dto.ReorderDocumentsInput) error

	// Trash
	GetTrash(input dto.GetTrashInput) (*dto.GetTrashOutput, error)
	RestoreDocument(input dto.RestoreDocumentInput) error

	// Public
	SetPublic(input dto.SetPublicInput) error
	GetPublicDocument(input dto.GetPublicDocumentInput) (*dto.GetPublicDocumentOutput, error)

	// Versions
	ListVersions(input dto.ListVersionsInput) (*dto.ListVersionsOutput, error)
	GetVersion(input dto.GetVersionInput) (*dto.GetVersionOutput, error)
	RestoreVersion(input dto.RestoreVersionInput) error
	CreateVersion(input dto.CreateVersionInput) (*dto.CreateVersionOutput, error)

	// Comments
	CreateComment(input dto.CreateCommentInput) (*dto.CreateCommentOutput, error)
	GetComments(input dto.GetCommentsInput) (*dto.GetCommentsOutput, error)
	UpdateComment(input dto.UpdateCommentInput) error
	DeleteComment(input dto.DeleteCommentInput) error
	ResolveComment(input dto.ResolveCommentInput) error
}
