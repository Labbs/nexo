package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upSpace, downSpace)
}

func upSpace(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		CREATE TABLE IF NOT EXISTS space (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			icon TEXT,
			icon_color TEXT,
			type TEXT CHECK(type IN ('personal', 'public', 'restricted', 'private')) NOT NULL,
			owner_id TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_space_name ON space(name);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_space_slug ON space(slug);
		CREATE INDEX IF NOT EXISTS idx_space_type ON space(type);
		CREATE INDEX IF NOT EXISTS idx_space_deleted_at ON space(deleted_at);
		`
	case "postgres":
		query = `
		CREATE TABLE IF NOT EXISTS space (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			icon TEXT,
			icon_color TEXT,
			type TEXT CHECK(type IN ('personal', 'public', 'restricted', 'private')) NOT NULL,
			owner_id UUID,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ
		);
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_space_name ON space(name);
		CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_space_slug ON space(slug);
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_space_type ON space(type);
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_space_deleted_at ON space(deleted_at);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downSpace(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS space;`)
	return err
}
