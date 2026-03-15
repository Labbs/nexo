package dto

import "time"

type CreateRowInput struct {
	UserId        string
	DatabaseId    string
	Properties    map[string]any
	Content       map[string]any
	ShowInSidebar bool
}

type CreateRowOutput struct {
	Id         string
	Properties map[string]any
	CreatedAt  time.Time
}
