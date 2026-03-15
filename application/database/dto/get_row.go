package dto

import "time"

type GetRowInput struct {
	UserId     string
	DatabaseId string
	RowId      string
}

type GetRowOutput struct {
	Id            string
	DatabaseId    string
	Properties    map[string]any
	Content       map[string]any
	ShowInSidebar bool
	CreatedBy     string
	CreatedByUser *UserInfo
	UpdatedBy     string
	UpdatedByUser *UserInfo
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
