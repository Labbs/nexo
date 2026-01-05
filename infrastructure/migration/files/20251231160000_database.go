package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDatabase, downDatabase)
}

func upDatabase(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		CREATE TABLE IF NOT EXISTS database (
			id TEXT PRIMARY KEY,
			space_id TEXT NOT NULL,
			document_id TEXT,
			name TEXT NOT NULL,
			description TEXT,
			icon TEXT,
			schema TEXT,
			views TEXT,
			default_view TEXT DEFAULT 'table',
			type TEXT DEFAULT 'spreadsheet',
			created_by TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP,
			FOREIGN KEY (space_id) REFERENCES space(id),
			FOREIGN KEY (document_id) REFERENCES document(id),
			FOREIGN KEY (created_by) REFERENCES user(id)
		);
		CREATE INDEX IF NOT EXISTS idx_database_space_id ON database(space_id);
		CREATE INDEX IF NOT EXISTS idx_database_document_id ON database(document_id);
		CREATE INDEX IF NOT EXISTS idx_database_deleted_at ON database(deleted_at);

		CREATE TABLE IF NOT EXISTS database_row (
			id TEXT PRIMARY KEY,
			database_id TEXT NOT NULL,
			properties TEXT,
			content TEXT,
			show_in_sidebar INTEGER DEFAULT 0,
			created_by TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP,
			FOREIGN KEY (database_id) REFERENCES database(id) ON DELETE CASCADE,
			FOREIGN KEY (created_by) REFERENCES user(id)
		);
		CREATE INDEX IF NOT EXISTS idx_database_row_database_id ON database_row(database_id);
		CREATE INDEX IF NOT EXISTS idx_database_row_deleted_at ON database_row(deleted_at);
		`
	case "postgres":
		query = `
		CREATE TABLE IF NOT EXISTS database (
			id UUID PRIMARY KEY,
			space_id UUID NOT NULL REFERENCES space(id),
			document_id UUID REFERENCES document(id),
			name TEXT NOT NULL,
			description TEXT,
			icon TEXT,
			schema JSONB,
			views JSONB,
			default_view TEXT DEFAULT 'table',
			type TEXT DEFAULT 'spreadsheet',
			created_by UUID NOT NULL REFERENCES "user"(id),
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ
		);
		CREATE INDEX IF NOT EXISTS idx_database_space_id ON database(space_id);
		CREATE INDEX IF NOT EXISTS idx_database_document_id ON database(document_id);
		CREATE INDEX IF NOT EXISTS idx_database_deleted_at ON database(deleted_at);

		CREATE TABLE IF NOT EXISTS database_row (
			id UUID PRIMARY KEY,
			database_id UUID NOT NULL REFERENCES database(id) ON DELETE CASCADE,
			properties JSONB,
			content JSONB,
			show_in_sidebar BOOLEAN DEFAULT false,
			created_by UUID NOT NULL REFERENCES "user"(id),
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ
		);
		CREATE INDEX IF NOT EXISTS idx_database_row_database_id ON database_row(database_id);
		CREATE INDEX IF NOT EXISTS idx_database_row_deleted_at ON database_row(deleted_at);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downDatabase(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS database_row;
		DROP TABLE IF EXISTS database;
	`)
	return err
}
