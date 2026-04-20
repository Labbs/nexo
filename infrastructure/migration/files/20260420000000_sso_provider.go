package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upOAuthProvider, downOAuthProvider)
}

func upOAuthProvider(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		CREATE TABLE IF NOT EXISTS oauth_provider (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			provider TEXT NOT NULL,
			provider_user_id TEXT NOT NULL,
			email TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
			UNIQUE(provider, provider_user_id)
		);
		CREATE INDEX IF NOT EXISTS idx_oauth_provider_user_id ON oauth_provider(user_id);
		`
	case "postgres":
		query = `
		CREATE TABLE IF NOT EXISTS oauth_provider (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			provider TEXT NOT NULL,
			provider_user_id TEXT NOT NULL,
			email TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(provider, provider_user_id)
		);
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_oauth_provider_user_id ON oauth_provider(user_id);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downOAuthProvider(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS oauth_provider;`)
	return err
}
