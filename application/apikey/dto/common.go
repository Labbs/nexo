package dto

import "time"

type ApiKeyItem struct {
	Id         string
	Name       string
	KeyPrefix  string
	Scopes     []string
	LastUsedAt *time.Time
	ExpiresAt  *time.Time
	CreatedAt  time.Time
}
