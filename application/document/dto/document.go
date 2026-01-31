package dto

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// Document représente un document dans la couche application
type Document struct {
	Id       string
	Name     string
	Slug     string
	ParentId *string
	SpaceId  string
	Public   bool
	Content  []Block
	Config   DocumentConfig
	Metadata map[string]any

	CreatedAt time.Time
	UpdatedAt time.Time
}

// DocumentConfig représente la configuration d'un document
type DocumentConfig struct {
	FullWidth        bool   `json:"fullWidth"`
	Icon             string `json:"icon"`
	Lock             bool   `json:"lock"`
	HeaderBackground string `json:"headerBackground"`
}

// Block représente un bloc BlockNote
type Block struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"` // paragraph, heading, bulletListItem, table, etc.
	Props    map[string]any  `json:"props"`
	Content  []InlineContent `json:"content"`
	Children []Block         `json:"children"`
}

// InlineContent représente le contenu inline (texte, liens, etc.)
type InlineContent struct {
	Type   string          `json:"type"` // "text", "link"
	Text   string          `json:"text,omitempty"`
	Href   string          `json:"href,omitempty"`
	Styles map[string]bool `json:"styles"`
}

// Types de blocs BlockNote supportés
const (
	BlockTypeParagraph        = "paragraph"
	BlockTypeHeading          = "heading"
	BlockTypeBulletListItem   = "bulletListItem"
	BlockTypeNumberedListItem = "numberedListItem"
	BlockTypeCheckListItem    = "checkListItem"
	BlockTypeTable            = "table"
	BlockTypeImage            = "image"
	BlockTypeVideo            = "video"
	BlockTypeAudio            = "audio"
	BlockTypeFile             = "file"
	BlockTypeCodeBlock        = "codeBlock"
	BlockTypeColumn           = "column"
	BlockTypeColumnList       = "columnList"
)

// BlocksToJSON convertit []Block en datatypes.JSON pour le stockage
func BlocksToJSON(blocks []Block) datatypes.JSON {
	if blocks == nil {
		return datatypes.JSON("[]")
	}
	data, err := json.Marshal(blocks)
	if err != nil {
		return datatypes.JSON("[]")
	}
	return datatypes.JSON(data)
}

// JSONToBlocks convertit datatypes.JSON en []Block
func JSONToBlocks(data datatypes.JSON) []Block {
	if data == nil || len(data) == 0 {
		return []Block{}
	}
	var blocks []Block
	if err := json.Unmarshal(data, &blocks); err != nil {
		return []Block{}
	}
	return blocks
}
