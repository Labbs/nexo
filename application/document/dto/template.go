package dto

import "github.com/labbs/nexo/domain"

type ListTemplatesInput struct {
	UserId   string
	SpaceIds []string // if empty, all accessible spaces
}

type ListTemplatesOutput struct {
	Templates []domain.Document
}

type ToggleTemplateInput struct {
	DocumentId string
	SpaceId    string
	UserId     string
	IsTemplate bool
	Category   string
}

type ToggleTemplateOutput struct {
	Document *domain.Document
}
