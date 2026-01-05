package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upSpacePermission, downSpacePermission)
}

func upSpacePermission(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		CREATE TABLE IF NOT EXISTS space_permission (
			id TEXT PRIMARY KEY,
			space_id TEXT NOT NULL,
			user_id TEXT,
			group_id TEXT,
			role TEXT CHECK(role IN ('owner', 'admin', 'editor', 'viewer')) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP,
			CHECK((user_id IS NOT NULL AND group_id IS NULL) OR (user_id IS NULL AND group_id IS NOT NULL))
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_space_permission_space_user ON space_permission(space_id, user_id) WHERE user_id IS NOT NULL AND deleted_at IS NULL;
		CREATE UNIQUE INDEX IF NOT EXISTS idx_space_permission_space_group ON space_permission(space_id, group_id) WHERE group_id IS NOT NULL AND deleted_at IS NULL;
		CREATE INDEX IF NOT EXISTS idx_space_permission_space_id ON space_permission(space_id);
		CREATE INDEX IF NOT EXISTS idx_space_permission_user_id ON space_permission(user_id) WHERE user_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_space_permission_group_id ON space_permission(group_id) WHERE group_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_space_permission_deleted_at ON space_permission(deleted_at);
		`
	case "postgres":
		query = `
    CREATE TABLE IF NOT EXISTS space_permission (
			id UUID PRIMARY KEY,
			space_id UUID NOT NULL,
			user_id UUID,
			group_id UUID,
			role TEXT CHECK(role IN ('owner', 'admin', 'editor', 'viewer')) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ,
			CHECK((user_id IS NOT NULL AND group_id IS NULL) OR (user_id IS NULL AND group_id IS NOT NULL))
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_space_permission_space_user ON space_permission(space_id, user_id) WHERE user_id IS NOT NULL AND deleted_at IS NULL;
		CREATE UNIQUE INDEX IF NOT EXISTS idx_space_permission_space_group ON space_permission(space_id, group_id) WHERE group_id IS NOT NULL AND deleted_at IS NULL;
		CREATE INDEX IF NOT EXISTS idx_space_permission_space_id ON space_permission(space_id);
		CREATE INDEX IF NOT EXISTS idx_space_permission_user_id ON space_permission(user_id) WHERE user_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_space_permission_group_id ON space_permission(group_id) WHERE group_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_space_permission_deleted_at ON space_permission(deleted_at);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downSpacePermission(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS space_permission;`)
	return err
}
