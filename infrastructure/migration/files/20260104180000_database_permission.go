package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDatabasePermission, downDatabasePermission)
}

func upDatabasePermission(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		CREATE TABLE IF NOT EXISTS database_permission (
			id TEXT PRIMARY KEY,
			database_id TEXT NOT NULL,
			user_id TEXT,
			group_id TEXT,
			role TEXT CHECK(role IN ('editor', 'viewer', 'denied')) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP,
			CHECK((user_id IS NOT NULL AND group_id IS NULL) OR (user_id IS NULL AND group_id IS NOT NULL))
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_database_permission_db_user ON database_permission(database_id, user_id) WHERE user_id IS NOT NULL AND deleted_at IS NULL;
		CREATE UNIQUE INDEX IF NOT EXISTS idx_database_permission_db_group ON database_permission(database_id, group_id) WHERE group_id IS NOT NULL AND deleted_at IS NULL;
		CREATE INDEX IF NOT EXISTS idx_database_permission_database_id ON database_permission(database_id);
		CREATE INDEX IF NOT EXISTS idx_database_permission_user_id ON database_permission(user_id) WHERE user_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_database_permission_group_id ON database_permission(group_id) WHERE group_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_database_permission_deleted_at ON database_permission(deleted_at);
		`
	case "postgres":
		query = `
		CREATE TABLE IF NOT EXISTS database_permission (
			id UUID PRIMARY KEY,
			database_id UUID NOT NULL,
			user_id UUID,
			group_id UUID,
			role TEXT CHECK(role IN ('editor', 'viewer', 'denied')) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ,
			CHECK((user_id IS NOT NULL AND group_id IS NULL) OR (user_id IS NULL AND group_id IS NOT NULL))
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_database_permission_db_user ON database_permission(database_id, user_id) WHERE user_id IS NOT NULL AND deleted_at IS NULL;
		CREATE UNIQUE INDEX IF NOT EXISTS idx_database_permission_db_group ON database_permission(database_id, group_id) WHERE group_id IS NOT NULL AND deleted_at IS NULL;
		CREATE INDEX IF NOT EXISTS idx_database_permission_database_id ON database_permission(database_id);
		CREATE INDEX IF NOT EXISTS idx_database_permission_user_id ON database_permission(user_id) WHERE user_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_database_permission_group_id ON database_permission(group_id) WHERE group_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_database_permission_deleted_at ON database_permission(deleted_at);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downDatabasePermission(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS database_permission;`)
	return err
}
