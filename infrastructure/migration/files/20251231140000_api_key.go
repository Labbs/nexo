package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upApiKey, downApiKey)
}

func upApiKey(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		CREATE TABLE IF NOT EXISTS api_key (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			name TEXT NOT NULL,
			key_hash TEXT NOT NULL UNIQUE,
			key_prefix TEXT NOT NULL,
			permissions TEXT,
			last_used_at TIMESTAMP,
			expires_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES user(id)
		);
		CREATE INDEX IF NOT EXISTS idx_api_key_user_id ON api_key(user_id);
		CREATE INDEX IF NOT EXISTS idx_api_key_key_hash ON api_key(key_hash);
		CREATE INDEX IF NOT EXISTS idx_api_key_deleted_at ON api_key(deleted_at);
		`
	case "postgres":
		query = `
		CREATE TABLE IF NOT EXISTS api_key (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES "user"(id),
			name TEXT NOT NULL,
			key_hash TEXT NOT NULL UNIQUE,
			key_prefix TEXT NOT NULL,
			permissions JSONB,
			last_used_at TIMESTAMPTZ,
			expires_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ
		);
		CREATE INDEX IF NOT EXISTS idx_api_key_user_id ON api_key(user_id);
		CREATE INDEX IF NOT EXISTS idx_api_key_key_hash ON api_key(key_hash);
		CREATE INDEX IF NOT EXISTS idx_api_key_deleted_at ON api_key(deleted_at);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downApiKey(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS api_key;`)
	return err
}
