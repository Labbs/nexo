package ports

import (
	"github.com/labbs/nexo/application/document/dto"
)

type DocumentPort interface {
	CreateDocument(input dto.CreateDocumentInput) (*dto.CreateDocumentOutput, error)
	GetDocumentWithSpace(input dto.GetDocumentWithSpaceInput) (*dto.GetDocumentWithSpaceOutput, error)
	GetDocumentsFromSpaceWithUserPermissions(input dto.GetDocumentsFromSpaceInput) (*dto.GetDocumentsFromSpaceOutput, error)
	UpdateDocument(input dto.UpdateDocumentInput) (*dto.UpdateDocumentOutput, error)
	MoveDocument(input dto.MoveDocumentInput) (*dto.MoveDocumentOutput, error)
	DeleteDocument(input dto.DeleteDocumentInput) error
}
