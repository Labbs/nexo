package dtos

import (
	"time"

	spaceDtos "github.com/labbs/nexo/interfaces/http/v1/space/dtos"
)

type GetDocumentRequest struct {
	SpaceId    string `path:"space_id" validate:"required,uuid4"`
	Identifier string `path:"identifier" validate:"required"`
}

type GetDocumentResponse struct {
	Id   string `json:"document"`
	Name string `json:"name"`
	Slug string `json:"slug"`

	ParentId *string   `json:"parent_id,omitempty"`
	Parent   *Document `json:"parent,omitempty"`

	SpaceId string          `json:"space_id"`
	Space   spaceDtos.Space `json:"space"`

	Content []Block `json:"content"`

	Config   DocumentConfig `json:"config"`
	Metadata map[string]any `json:"metadata"`

	Public bool `json:"public"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
