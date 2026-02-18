package ports

import (
	"github.com/labbs/nexo/application/database/dto"
)

type DatabasePort interface {
	// Database CRUD
	CreateDatabase(input dto.CreateDatabaseInput) (*dto.CreateDatabaseOutput, error)
	ListDatabases(input dto.ListDatabasesInput) (*dto.ListDatabasesOutput, error)
	GetDatabase(input dto.GetDatabaseInput) (*dto.GetDatabaseOutput, error)
	UpdateDatabase(input dto.UpdateDatabaseInput) error
	DeleteDatabase(input dto.DeleteDatabaseInput) error
	MoveDatabase(input dto.MoveDatabaseInput) (*dto.MoveDatabaseOutput, error)
	Search(input dto.SearchDatabasesInput) (*dto.SearchDatabasesOutput, error)

	// Views
	CreateView(input dto.CreateViewInput) (*dto.CreateViewOutput, error)
	UpdateView(input dto.UpdateViewInput) error
	DeleteView(input dto.DeleteViewInput) error

	// Rows
	CreateRow(input dto.CreateRowInput) (*dto.CreateRowOutput, error)
	ListRows(input dto.ListRowsInput) (*dto.ListRowsOutput, error)
	GetRow(input dto.GetRowInput) (*dto.GetRowOutput, error)
	UpdateRow(input dto.UpdateRowInput) error
	DeleteRow(input dto.DeleteRowInput) error
	BulkDeleteRows(input dto.BulkDeleteRowsInput) error
}
