package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upGroupMembers, downGroupMembers)
}

func upGroupMembers(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		-- Add deleted_at column to group table
		ALTER TABLE "group" ADD COLUMN deleted_at TIMESTAMP;

		-- Create group_members join table
		CREATE TABLE IF NOT EXISTS group_members (
			group_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (group_id, user_id),
			FOREIGN KEY (group_id) REFERENCES "group"(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members(group_id);
		CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON group_members(user_id);
		`
	case "postgres":
		query = `
		-- Add deleted_at column to group table
		ALTER TABLE "group" ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

		-- Create group_members join table
		CREATE TABLE IF NOT EXISTS group_members (
			group_id UUID NOT NULL,
			user_id UUID NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (group_id, user_id),
			FOREIGN KEY (group_id) REFERENCES "group"(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members(group_id);
		CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON group_members(user_id);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downGroupMembers(ctx context.Context, tx *sql.Tx) error {
	query := `
	DROP TABLE IF EXISTS group_members;
	ALTER TABLE "group" DROP COLUMN IF EXISTS deleted_at;
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}
