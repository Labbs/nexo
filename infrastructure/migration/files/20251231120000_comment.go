package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upComment, downComment)
}

func upComment(ctx context.Context, tx *sql.Tx) error {
	var query string
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		query = `
		CREATE TABLE IF NOT EXISTS comment (
			id TEXT PRIMARY KEY,
			document_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			parent_id TEXT,
			content TEXT NOT NULL,
			block_id TEXT,
			resolved INTEGER DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP,
			FOREIGN KEY (document_id) REFERENCES document(id),
			FOREIGN KEY (user_id) REFERENCES user(id),
			FOREIGN KEY (parent_id) REFERENCES comment(id)
		);
		CREATE INDEX IF NOT EXISTS idx_comment_document_id ON comment(document_id);
		CREATE INDEX IF NOT EXISTS idx_comment_user_id ON comment(user_id);
		CREATE INDEX IF NOT EXISTS idx_comment_parent_id ON comment(parent_id);
		CREATE INDEX IF NOT EXISTS idx_comment_block_id ON comment(block_id);
		CREATE INDEX IF NOT EXISTS idx_comment_deleted_at ON comment(deleted_at);
		`
	case "postgres":
		query = `
		CREATE TABLE IF NOT EXISTS comment (
			id UUID PRIMARY KEY,
			document_id UUID NOT NULL REFERENCES document(id),
			user_id UUID NOT NULL REFERENCES "user"(id),
			parent_id UUID REFERENCES comment(id),
			content TEXT NOT NULL,
			block_id TEXT,
			resolved BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ
		);
		CREATE INDEX IF NOT EXISTS idx_comment_document_id ON comment(document_id);
		CREATE INDEX IF NOT EXISTS idx_comment_user_id ON comment(user_id);
		CREATE INDEX IF NOT EXISTS idx_comment_parent_id ON comment(parent_id);
		CREATE INDEX IF NOT EXISTS idx_comment_block_id ON comment(block_id);
		CREATE INDEX IF NOT EXISTS idx_comment_deleted_at ON comment(deleted_at);
		`
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	_, err := tx.ExecContext(ctx, query)
	return err
}

func downComment(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS comment;`)
	return err
}
