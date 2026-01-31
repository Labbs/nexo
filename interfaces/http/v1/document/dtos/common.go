package dtos

import (
	"time"

	spaceDtos "github.com/labbs/nexo/interfaces/http/v1/space/dtos"
)

type Document struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`

	ParentId *string   `json:"parent_id,omitempty"`
	Parent   *Document `json:"parent,omitempty"`

	SpaceId string          `json:"space_id"`
	Space   spaceDtos.Space `json:"space"`

	Content []Block `json:"content"`

	Config   DocumentConfig `json:"config"`
	Metadata map[string]any `json:"metadata"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DocumentConfig struct {
	FullWidth        bool   `json:"full_width"`
	Icon             string `json:"icon"`
	Lock             bool   `json:"lock"`
	HeaderBackground string `json:"header_background"`
}

// Block représente un bloc BlockNote
type Block struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Props    map[string]any  `json:"props"`
	Content  []InlineContent `json:"content"`
	Children []Block         `json:"children"`
}

// InlineContent représente le contenu inline (texte, liens, etc.)
type InlineContent struct {
	Type   string          `json:"type"`
	Text   string          `json:"text,omitempty"`
	Href   string          `json:"href,omitempty"`
	Styles map[string]bool `json:"styles"`
}
