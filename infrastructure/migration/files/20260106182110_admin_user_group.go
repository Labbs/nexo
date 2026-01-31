package migrations

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	goose.AddMigrationContext(upAdminUser, downAdminUser)
}

func generateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func upAdminUser(ctx context.Context, tx *sql.Tx) error {
	// Generate UUIDs for admin user and group
	adminUserId := uuid.New().String()
	adminGroupId := uuid.New().String()

	// Generate random password
	randomPassword, err := generateRandomPassword(16)
	if err != nil {
		return fmt.Errorf("failed to generate random password: %w", err)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Log the generated password
	logger, _ := ctx.Value("logger").(zerolog.Logger)
	logger.Info().Str("username", "admin").Str("email", "admin@nexo.local").Str("password", randomPassword).Msg("Admin user created")

	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		// Insert admin user
		_, err = tx.ExecContext(ctx, `
			INSERT INTO user (id, username, email, password, role, active, created_at, updated_at)
			VALUES (?, 'admin', 'admin@nexo.local', ?, 'admin', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
		`, adminUserId, string(hashedPassword))
		if err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}

		// Insert admin group
		_, err = tx.ExecContext(ctx, `
			INSERT INTO "group" (id, name, description, role, owner_id, created_at, updated_at)
			VALUES (?, 'Administrators', 'Default administrator group', 'admin', ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
		`, adminGroupId, adminUserId)
		if err != nil {
			return fmt.Errorf("failed to create admin group: %w", err)
		}

		// Add admin user to admin group
		_, err = tx.ExecContext(ctx, `
			INSERT INTO group_members (group_id, user_id, created_at)
			VALUES (?, ?, CURRENT_TIMESTAMP);
		`, adminGroupId, adminUserId)
		if err != nil {
			return fmt.Errorf("failed to add admin user to admin group: %w", err)
		}

	case "postgres":
		// Insert admin user
		_, err = tx.ExecContext(ctx, `
			INSERT INTO "user" (id, username, email, password, role, active, created_at, updated_at)
			VALUES ($1, 'admin', 'admin@nexo.local', $2, 'admin', true, NOW(), NOW());
		`, adminUserId, string(hashedPassword))
		if err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}

		// Insert admin group
		_, err = tx.ExecContext(ctx, `
			INSERT INTO "group" (id, name, description, role, owner_id, created_at, updated_at)
			VALUES ($1, 'Administrators', 'Default administrator group', 'admin', $2, NOW(), NOW());
		`, adminGroupId, adminUserId)
		if err != nil {
			return fmt.Errorf("failed to create admin group: %w", err)
		}

		// Add admin user to admin group
		_, err = tx.ExecContext(ctx, `
			INSERT INTO group_members (group_id, user_id, created_at)
			VALUES ($1, $2, NOW());
		`, adminGroupId, adminUserId)
		if err != nil {
			return fmt.Errorf("failed to add admin user to admin group: %w", err)
		}

	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}

	return nil
}

func downAdminUser(ctx context.Context, tx *sql.Tx) error {
	dialect, _ := ctx.Value("dbDialect").(string)
	switch dialect {
	case "sqlite":
		_, err := tx.ExecContext(ctx, `
			DELETE FROM group_members WHERE group_id IN (SELECT id FROM "group" WHERE name = 'Administrators');
			DELETE FROM "group" WHERE name = 'Administrators';
			DELETE FROM user WHERE username = 'admin';
		`)
		return err
	case "postgres":
		_, err := tx.ExecContext(ctx, `
			DELETE FROM group_members WHERE group_id IN (SELECT id FROM "group" WHERE name = 'Administrators');
			DELETE FROM "group" WHERE name = 'Administrators';
			DELETE FROM "user" WHERE username = 'admin';
		`)
		return err
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}
}
