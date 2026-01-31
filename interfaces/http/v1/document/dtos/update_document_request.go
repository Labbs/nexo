package dtos

type UpdateDocumentRequest struct {
	SpaceId  string          `path:"space_id" validate:"required,uuid4"`
	Id       string          `path:"id" validate:"required"`
	Name     *string         `json:"name,omitempty" validate:"omitempty,min=1"`
	Content  *[]Block        `json:"content,omitempty"`
	ParentId *string         `json:"parent_id,omitempty" validate:"omitempty"`
	Config   *DocumentConfig `json:"config,omitempty"`
	Metadata *map[string]any `json:"metadata,omitempty"`
	Public   *bool           `json:"public,omitempty"`
}

type UpdateDocumentResponse struct {
	Id       string         `json:"id"`
	Name     string         `json:"name"`
	Slug     string         `json:"slug"`
	ParentId *string        `json:"parent_id,omitempty"`
	SpaceId  string         `json:"space_id"`
	Content  []Block        `json:"content"`
	Config   DocumentConfig `json:"config"`
	Metadata map[string]any `json:"metadata"`
}
