package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDocumentTemplate, downDocumentTemplate)
}

func upDocumentTemplate(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		ALTER TABLE document ADD COLUMN is_template BOOLEAN NOT NULL DEFAULT FALSE;
		ALTER TABLE document ADD COLUMN template_category TEXT NOT NULL DEFAULT '';
		CREATE INDEX IF NOT EXISTS idx_document_is_template ON document(is_template);
		`
	case "postgres":
		query = `
		ALTER TABLE document ADD COLUMN IF NOT EXISTS is_template BOOLEAN NOT NULL DEFAULT FALSE;
		ALTER TABLE document ADD COLUMN IF NOT EXISTS template_category TEXT NOT NULL DEFAULT '';
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_document_is_template ON document(is_template);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downDocumentTemplate(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE document DROP COLUMN IF EXISTS is_template;
		ALTER TABLE document DROP COLUMN IF EXISTS template_category;
	`)
	return err
}
