package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDrawing, downDrawing)
}

func upDrawing(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		CREATE TABLE IF NOT EXISTS drawing (
			id TEXT PRIMARY KEY,
			space_id TEXT NOT NULL,
			document_id TEXT,
			name TEXT NOT NULL,
			icon TEXT DEFAULT '',
			elements TEXT DEFAULT '[]',
			app_state TEXT DEFAULT '{}',
			files TEXT DEFAULT '{}',
			thumbnail TEXT,
			created_by TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP,
			FOREIGN KEY (space_id) REFERENCES space(id) ON DELETE CASCADE,
			FOREIGN KEY (document_id) REFERENCES document(id) ON DELETE SET NULL,
			FOREIGN KEY (created_by) REFERENCES user(id)
		);
		CREATE INDEX IF NOT EXISTS idx_drawing_space_id ON drawing(space_id);
		CREATE INDEX IF NOT EXISTS idx_drawing_document_id ON drawing(document_id);
		CREATE INDEX IF NOT EXISTS idx_drawing_deleted_at ON drawing(deleted_at);
		`
	case "postgres":
		query = `
		CREATE TABLE IF NOT EXISTS drawing (
			id UUID PRIMARY KEY,
			space_id UUID NOT NULL REFERENCES space(id) ON DELETE CASCADE,
			document_id UUID REFERENCES document(id) ON DELETE SET NULL,
			name TEXT NOT NULL,
			icon TEXT DEFAULT '',
			elements JSONB DEFAULT '[]',
			app_state JSONB DEFAULT '{}',
			files JSONB DEFAULT '{}',
			thumbnail TEXT,
			created_by UUID NOT NULL REFERENCES "user"(id),
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ
		);
		CREATE INDEX IF NOT EXISTS idx_drawing_space_id ON drawing(space_id);
		CREATE INDEX IF NOT EXISTS idx_drawing_document_id ON drawing(document_id);
		CREATE INDEX IF NOT EXISTS idx_drawing_deleted_at ON drawing(deleted_at);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downDrawing(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS drawing;`)
	return err
}
