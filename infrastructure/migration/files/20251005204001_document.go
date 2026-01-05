package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDocument, downDocument)
}

func upDocument(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		CREATE TABLE IF NOT EXISTS document (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			config JSONB,
			metadata JSONB,
			parent_id TEXT,
			space_id TEXT NOT NULL,
			public BOOLEAN DEFAULT FALSE,
			content JSONB DEFAULT '[]',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP,
			FOREIGN KEY (space_id) REFERENCES space(id) ON DELETE CASCADE,
			FOREIGN KEY (parent_id) REFERENCES document(id) ON DELETE SET NULL
		);
		CREATE INDEX IF NOT EXISTS idx_document_space_id ON document(space_id);
		CREATE INDEX IF NOT EXISTS idx_document_parent_id ON document(parent_id);
		CREATE INDEX IF NOT EXISTS idx_document_public ON document(public);
		CREATE INDEX IF NOT EXISTS idx_document_deleted_at ON document(deleted_at);
		`
	case "postgres":
		query = `
		CREATE TABLE IF NOT EXISTS document (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			config JSONB,
			metadata JSONB,
			parent_id UUID,
			space_id UUID NOT NULL,
			public BOOLEAN DEFAULT FALSE,
			content JSONB DEFAULT '[]',
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ,
			FOREIGN KEY (space_id) REFERENCES space(id) ON DELETE CASCADE,
			FOREIGN KEY (parent_id) REFERENCES document(id) ON DELETE SET NULL
		);
		CREATE INDEX IF NOT EXISTS idx_document_space_id ON document(space_id);
		CREATE INDEX IF NOT EXISTS idx_document_parent_id ON document(parent_id);
		CREATE INDEX IF NOT EXISTS idx_document_public ON document(public);
		CREATE INDEX IF NOT EXISTS idx_document_deleted_at ON document(deleted_at);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downDocument(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS document;`)
	return err
}
