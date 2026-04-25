package dtos

type ListTemplatesRequest struct {
	SpaceId string `query:"space_id" validate:"omitempty,uuid4"`
}

type TemplateItem struct {
	Id               string        `json:"id"`
	Name             string        `json:"name"`
	Slug             string        `json:"slug"`
	TemplateCategory string        `json:"template_category,omitempty"`
	Config           DocumentConfig `json:"config"`
	SpaceId          string        `json:"space_id"`
	SpaceName        string        `json:"space_name"`
}

type ListTemplatesResponse struct {
	Templates []TemplateItem `json:"templates"`
}

type ToggleTemplateRequest struct {
	SpaceId    string `path:"space_id" validate:"required,uuid4"`
	DocumentId string `path:"document_id" validate:"required,uuid4"`
	IsTemplate bool   `json:"is_template"`
	Category   string `json:"category,omitempty"`
}

type ToggleTemplateResponse struct {
	Document
}
