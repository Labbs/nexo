package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDocumentVersion, downDocumentVersion)
}

func upDocumentVersion(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		CREATE TABLE IF NOT EXISTS document_version (
			id TEXT PRIMARY KEY,
			document_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			version INTEGER NOT NULL,
			name TEXT NOT NULL,
			content TEXT,
			config TEXT,
			description TEXT,
			created_at TIMESTAMP NOT NULL,
			FOREIGN KEY (document_id) REFERENCES document(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES user(id)
		);
		CREATE INDEX IF NOT EXISTS idx_document_version_document_id ON document_version(document_id);
		CREATE INDEX IF NOT EXISTS idx_document_version_created_at ON document_version(created_at);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_document_version_doc_version ON document_version(document_id, version);
		`
	case "postgres":
		query = `
		CREATE TABLE IF NOT EXISTS document_version (
			id UUID PRIMARY KEY,
			document_id UUID NOT NULL REFERENCES document(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES "user"(id),
			version INTEGER NOT NULL,
			name TEXT NOT NULL,
			content JSONB,
			config JSONB,
			description TEXT,
			created_at TIMESTAMPTZ NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_document_version_document_id ON document_version(document_id);
		CREATE INDEX IF NOT EXISTS idx_document_version_created_at ON document_version(created_at);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_document_version_doc_version ON document_version(document_id, version);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downDocumentVersion(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS document_version;`)
	return err
}
