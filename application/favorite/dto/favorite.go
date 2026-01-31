package dto

import "time"

// Favorite represents a favorite document for the application layer
type Favorite struct {
	Id         string
	UserId     string
	DocumentId string
	SpaceId    string
	Position   int
	CreatedAt  time.Time

	// Document info for display
	Document *FavoriteDocument
}

// FavoriteDocument contains minimal document info for favorites list
type FavoriteDocument struct {
	Id     string
	Name   string
	Slug   string
	Icon   string
}
