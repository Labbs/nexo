package dto

import "time"

type CreateRowInput struct {
	UserId        string
	DatabaseId    string
	Properties    map[string]interface{}
	Content       map[string]interface{}
	ShowInSidebar bool
}

type CreateRowOutput struct {
	Id         string
	Properties map[string]interface{}
	CreatedAt  time.Time
}
