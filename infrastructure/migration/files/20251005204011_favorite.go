package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upFavorite, downFavorite)
}

func upFavorite(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		CREATE TABLE IF NOT EXISTS favorite (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			document_id TEXT NOT NULL,
			space_id TEXT NOT NULL,
			position INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_favorite_user_id ON favorite(user_id);
		CREATE INDEX IF NOT EXISTS idx_favorite_document_id ON favorite(document_id);
		CREATE INDEX IF NOT EXISTS idx_favorite_space_id ON favorite(space_id);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_favorite_user_document ON favorite(user_id, document_id, space_id);
		`
	case "postgres":
		query = `
		CREATE TABLE IF NOT EXISTS favorite (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			document_id UUID NOT NULL,
			space_id UUID NOT NULL,
			position INTEGER NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		);
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_favorite_user_id ON favorite(user_id);
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_favorite_document_id ON favorite(document_id);
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_favorite_space_id ON favorite(space_id);
		CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_favorite_user_document ON favorite(user_id, document_id, space_id);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downFavorite(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS favorite;`)
	return err
}
