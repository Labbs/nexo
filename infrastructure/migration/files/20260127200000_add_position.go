package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddPosition, downAddPosition)
}

func upAddPosition(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		ALTER TABLE document ADD COLUMN position INTEGER DEFAULT 0;
		CREATE INDEX IF NOT EXISTS idx_document_position ON document(position);

		ALTER TABLE database ADD COLUMN position INTEGER DEFAULT 0;
		CREATE INDEX IF NOT EXISTS idx_database_position ON database(position);

		ALTER TABLE drawing ADD COLUMN position INTEGER DEFAULT 0;
		CREATE INDEX IF NOT EXISTS idx_drawing_position ON drawing(position);
		`
	case "postgres":
		query = `
		ALTER TABLE document ADD COLUMN position INTEGER DEFAULT 0;
		CREATE INDEX IF NOT EXISTS idx_document_position ON document(position);

		ALTER TABLE "database" ADD COLUMN position INTEGER DEFAULT 0;
		CREATE INDEX IF NOT EXISTS idx_database_position ON "database"(position);

		ALTER TABLE drawing ADD COLUMN position INTEGER DEFAULT 0;
		CREATE INDEX IF NOT EXISTS idx_drawing_position ON drawing(position);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	// Initialize positions for existing documents based on created_at order
	// Group by parent_id + space_id to set position within each sibling group
	initQueries := []string{
		`UPDATE document SET position = sub.rn FROM (
			SELECT id, ROW_NUMBER() OVER (PARTITION BY COALESCE(parent_id, ''), space_id ORDER BY created_at ASC) - 1 AS rn
			FROM document WHERE deleted_at IS NULL
		) sub WHERE document.id = sub.id`,
		`UPDATE "database" SET position = sub.rn FROM (
			SELECT id, ROW_NUMBER() OVER (PARTITION BY COALESCE(document_id, ''), space_id ORDER BY created_at ASC) - 1 AS rn
			FROM "database" WHERE deleted_at IS NULL
		) sub WHERE "database".id = sub.id`,
		`UPDATE drawing SET position = sub.rn FROM (
			SELECT id, ROW_NUMBER() OVER (PARTITION BY COALESCE(document_id, ''), space_id ORDER BY created_at ASC) - 1 AS rn
			FROM drawing WHERE deleted_at IS NULL
		) sub WHERE drawing.id = sub.id`,
	}

	if dialect == "sqlite" {
		// SQLite doesn't support UPDATE ... FROM with window functions
		// Use a simpler approach: positions will be 0 for all existing records (default)
		// This is acceptable as existing records had no explicit ordering
		return nil
	}

	for _, q := range initQueries {
		if _, err := tx.ExecContext(ctx, q); err != nil {
			return err
		}
	}

	return nil
}

func downAddPosition(ctx context.Context, tx *sql.Tx) error {
	dialect, _ := ctx.Value("dbDialect").(string)
	var queries []string
	switch dialect {
	case "sqlite":
		// SQLite doesn't support DROP COLUMN before 3.35.0
		// We'll just leave the columns as-is for down migration safety
		return nil
	case "postgres":
		queries = []string{
			`ALTER TABLE document DROP COLUMN IF EXISTS position`,
			`ALTER TABLE "database" DROP COLUMN IF EXISTS position`,
			`ALTER TABLE drawing DROP COLUMN IF EXISTS position`,
		}
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	for _, q := range queries {
		if _, err := tx.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}
