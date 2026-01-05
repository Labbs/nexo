package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upWebhook, downWebhook)
}

func upWebhook(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		CREATE TABLE IF NOT EXISTS webhook (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			space_id TEXT,
			name TEXT NOT NULL,
			url TEXT NOT NULL,
			secret TEXT NOT NULL,
			events TEXT,
			active INTEGER DEFAULT 1,
			last_error TEXT,
			last_error_at TIMESTAMP,
			success_count INTEGER DEFAULT 0,
			failure_count INTEGER DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES user(id),
			FOREIGN KEY (space_id) REFERENCES space(id)
		);
		CREATE INDEX IF NOT EXISTS idx_webhook_user_id ON webhook(user_id);
		CREATE INDEX IF NOT EXISTS idx_webhook_space_id ON webhook(space_id);
		CREATE INDEX IF NOT EXISTS idx_webhook_active ON webhook(active);
		CREATE INDEX IF NOT EXISTS idx_webhook_deleted_at ON webhook(deleted_at);

		CREATE TABLE IF NOT EXISTS webhook_delivery (
			id TEXT PRIMARY KEY,
			webhook_id TEXT NOT NULL,
			event TEXT NOT NULL,
			payload TEXT,
			status_code INTEGER,
			response TEXT,
			duration INTEGER,
			success INTEGER,
			created_at TIMESTAMP NOT NULL,
			FOREIGN KEY (webhook_id) REFERENCES webhook(id) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_webhook_delivery_webhook_id ON webhook_delivery(webhook_id);
		CREATE INDEX IF NOT EXISTS idx_webhook_delivery_created_at ON webhook_delivery(created_at);
		`
	case "postgres":
		query = `
		CREATE TABLE IF NOT EXISTS webhook (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES "user"(id),
			space_id UUID REFERENCES space(id),
			name TEXT NOT NULL,
			url TEXT NOT NULL,
			secret TEXT NOT NULL,
			events JSONB,
			active BOOLEAN DEFAULT TRUE,
			last_error TEXT,
			last_error_at TIMESTAMPTZ,
			success_count INTEGER DEFAULT 0,
			failure_count INTEGER DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ
		);
		CREATE INDEX IF NOT EXISTS idx_webhook_user_id ON webhook(user_id);
		CREATE INDEX IF NOT EXISTS idx_webhook_space_id ON webhook(space_id);
		CREATE INDEX IF NOT EXISTS idx_webhook_active ON webhook(active);
		CREATE INDEX IF NOT EXISTS idx_webhook_deleted_at ON webhook(deleted_at);

		CREATE TABLE IF NOT EXISTS webhook_delivery (
			id UUID PRIMARY KEY,
			webhook_id UUID NOT NULL REFERENCES webhook(id) ON DELETE CASCADE,
			event TEXT NOT NULL,
			payload JSONB,
			status_code INTEGER,
			response TEXT,
			duration INTEGER,
			success BOOLEAN,
			created_at TIMESTAMPTZ NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_webhook_delivery_webhook_id ON webhook_delivery(webhook_id);
		CREATE INDEX IF NOT EXISTS idx_webhook_delivery_created_at ON webhook_delivery(created_at);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downWebhook(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS webhook_delivery;
		DROP TABLE IF EXISTS webhook;
	`)
	return err
}
