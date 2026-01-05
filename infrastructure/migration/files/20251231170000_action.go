package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAction, downAction)
}

func upAction(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		CREATE TABLE IF NOT EXISTS action (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			space_id TEXT,
			database_id TEXT,
			name TEXT NOT NULL,
			description TEXT,
			trigger_type TEXT NOT NULL,
			trigger_config TEXT,
			steps TEXT,
			active INTEGER DEFAULT 1,
			last_run_at TIMESTAMP,
			last_error TEXT,
			run_count INTEGER DEFAULT 0,
			success_count INTEGER DEFAULT 0,
			failure_count INTEGER DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES user(id),
			FOREIGN KEY (space_id) REFERENCES space(id)
		);
		CREATE INDEX IF NOT EXISTS idx_action_user_id ON action(user_id);
		CREATE INDEX IF NOT EXISTS idx_action_space_id ON action(space_id);
		CREATE INDEX IF NOT EXISTS idx_action_database_id ON action(database_id);
		CREATE INDEX IF NOT EXISTS idx_action_trigger_type ON action(trigger_type);
		CREATE INDEX IF NOT EXISTS idx_action_active ON action(active);
		CREATE INDEX IF NOT EXISTS idx_action_deleted_at ON action(deleted_at);

		CREATE TABLE IF NOT EXISTS action_run (
			id TEXT PRIMARY KEY,
			action_id TEXT NOT NULL,
			trigger_data TEXT,
			steps_result TEXT,
			success INTEGER,
			error TEXT,
			duration INTEGER,
			created_at TIMESTAMP NOT NULL,
			FOREIGN KEY (action_id) REFERENCES action(id) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_action_run_action_id ON action_run(action_id);
		CREATE INDEX IF NOT EXISTS idx_action_run_created_at ON action_run(created_at);
		`
	case "postgres":
		query = `
		CREATE TABLE IF NOT EXISTS action (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES "user"(id),
			space_id UUID REFERENCES space(id),
			database_id UUID,
			name TEXT NOT NULL,
			description TEXT,
			trigger_type TEXT NOT NULL,
			trigger_config JSONB,
			steps JSONB,
			active BOOLEAN DEFAULT TRUE,
			last_run_at TIMESTAMPTZ,
			last_error TEXT,
			run_count INTEGER DEFAULT 0,
			success_count INTEGER DEFAULT 0,
			failure_count INTEGER DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ
		);
		CREATE INDEX IF NOT EXISTS idx_action_user_id ON action(user_id);
		CREATE INDEX IF NOT EXISTS idx_action_space_id ON action(space_id);
		CREATE INDEX IF NOT EXISTS idx_action_database_id ON action(database_id);
		CREATE INDEX IF NOT EXISTS idx_action_trigger_type ON action(trigger_type);
		CREATE INDEX IF NOT EXISTS idx_action_active ON action(active);
		CREATE INDEX IF NOT EXISTS idx_action_deleted_at ON action(deleted_at);

		CREATE TABLE IF NOT EXISTS action_run (
			id UUID PRIMARY KEY,
			action_id UUID NOT NULL REFERENCES action(id) ON DELETE CASCADE,
			trigger_data JSONB,
			steps_result JSONB,
			success BOOLEAN,
			error TEXT,
			duration INTEGER,
			created_at TIMESTAMPTZ NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_action_run_action_id ON action_run(action_id);
		CREATE INDEX IF NOT EXISTS idx_action_run_created_at ON action_run(created_at);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downAction(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS action_run;
		DROP TABLE IF EXISTS action;
	`)
	return err
}
