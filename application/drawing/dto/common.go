package dto

import "time"

type DrawingItem struct {
	Id         string
	DocumentId *string
	Name       string
	Icon       string
	Thumbnail  string
	CreatedBy  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
