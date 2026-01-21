package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upPermission, downPermission)
}

func upPermission(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		CREATE TABLE IF NOT EXISTS permission (
			id TEXT PRIMARY KEY,
			type TEXT CHECK(type IN ('space', 'document', 'database', 'drawing')) NOT NULL,
			space_id TEXT,
			document_id TEXT,
			database_id TEXT,
			drawing_id TEXT,
			user_id TEXT,
			group_id TEXT,
			role TEXT CHECK(role IN ('owner', 'admin', 'editor', 'viewer', 'denied')) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP,
			FOREIGN KEY (space_id) REFERENCES space(id) ON DELETE CASCADE,
			FOREIGN KEY (document_id) REFERENCES document(id) ON DELETE CASCADE,
			FOREIGN KEY (database_id) REFERENCES database(id) ON DELETE CASCADE,
			FOREIGN KEY (drawing_id) REFERENCES drawing(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
			FOREIGN KEY (group_id) REFERENCES "group"(id) ON DELETE CASCADE,
			CHECK((user_id IS NOT NULL AND group_id IS NULL) OR (user_id IS NULL AND group_id IS NOT NULL))
		);
		CREATE INDEX IF NOT EXISTS idx_permission_type ON permission(type);
		CREATE INDEX IF NOT EXISTS idx_permission_space_id ON permission(space_id) WHERE space_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_permission_document_id ON permission(document_id) WHERE document_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_permission_database_id ON permission(database_id) WHERE database_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_permission_drawing_id ON permission(drawing_id) WHERE drawing_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_permission_user_id ON permission(user_id) WHERE user_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_permission_group_id ON permission(group_id) WHERE group_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_permission_deleted_at ON permission(deleted_at);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_permission_resource_user ON permission(type, space_id, document_id, database_id, drawing_id, user_id) WHERE user_id IS NOT NULL AND deleted_at IS NULL;
		CREATE UNIQUE INDEX IF NOT EXISTS idx_permission_resource_group ON permission(type, space_id, document_id, database_id, drawing_id, group_id) WHERE group_id IS NOT NULL AND deleted_at IS NULL;
		`
	case "postgres":
		query = `
		CREATE TABLE IF NOT EXISTS permission (
			id UUID PRIMARY KEY,
			type TEXT CHECK(type IN ('space', 'document', 'database', 'drawing')) NOT NULL,
			space_id UUID REFERENCES space(id) ON DELETE CASCADE,
			document_id UUID REFERENCES document(id) ON DELETE CASCADE,
			database_id UUID REFERENCES database(id) ON DELETE CASCADE,
			drawing_id UUID REFERENCES drawing(id) ON DELETE CASCADE,
			user_id UUID REFERENCES "user"(id) ON DELETE CASCADE,
			group_id UUID REFERENCES "group"(id) ON DELETE CASCADE,
			role TEXT CHECK(role IN ('owner', 'admin', 'editor', 'viewer', 'denied')) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ,
			CHECK((user_id IS NOT NULL AND group_id IS NULL) OR (user_id IS NULL AND group_id IS NOT NULL))
		);
		CREATE INDEX IF NOT EXISTS idx_permission_type ON permission(type);
		CREATE INDEX IF NOT EXISTS idx_permission_space_id ON permission(space_id) WHERE space_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_permission_document_id ON permission(document_id) WHERE document_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_permission_database_id ON permission(database_id) WHERE database_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_permission_drawing_id ON permission(drawing_id) WHERE drawing_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_permission_user_id ON permission(user_id) WHERE user_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_permission_group_id ON permission(group_id) WHERE group_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_permission_deleted_at ON permission(deleted_at);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_permission_resource_user ON permission(type, COALESCE(space_id, '00000000-0000-0000-0000-000000000000'), COALESCE(document_id, '00000000-0000-0000-0000-000000000000'), COALESCE(database_id, '00000000-0000-0000-0000-000000000000'), COALESCE(drawing_id, '00000000-0000-0000-0000-000000000000'), user_id) WHERE user_id IS NOT NULL AND deleted_at IS NULL;
		CREATE UNIQUE INDEX IF NOT EXISTS idx_permission_resource_group ON permission(type, COALESCE(space_id, '00000000-0000-0000-0000-000000000000'), COALESCE(document_id, '00000000-0000-0000-0000-000000000000'), COALESCE(database_id, '00000000-0000-0000-0000-000000000000'), COALESCE(drawing_id, '00000000-0000-0000-0000-000000000000'), group_id) WHERE group_id IS NOT NULL AND deleted_at IS NULL;
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downPermission(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS permission;`)
	return err
}
