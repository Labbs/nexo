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
	Properties    map[string]interface{}
	Content       map[string]interface{}
	ShowInSidebar bool
	CreatedBy     string
	CreatedByUser *UserInfo
	UpdatedBy     string
	UpdatedByUser *UserInfo
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
