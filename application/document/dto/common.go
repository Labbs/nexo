package dto

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// DocumentSpace contains space info embedded in a document response
type DocumentSpace struct {
	Id        string
	Name      string
	Slug      string
	Icon      string
	IconColor string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Document represents a document in the application layer
type Document struct {
	Id       string
	Name     string
	Slug     string
	ParentId *string
	Parent   *Document
	SpaceId  string
	Space    DocumentSpace
	Public   bool
	Content  []Block
	Config   DocumentConfig
	Metadata map[string]any

	CreatedAt time.Time
	UpdatedAt time.Time
}

// DocumentConfig represents the configuration of a document
type DocumentConfig struct {
	FullWidth        bool   `json:"fullWidth"`
	Icon             string `json:"icon"`
	Lock             bool   `json:"lock"`
	HeaderBackground string `json:"headerBackground"`
}

// Block represents a BlockNote block
type Block struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"` // paragraph, heading, bulletListItem, table, etc.
	Props    map[string]any  `json:"props"`
	Content  []InlineContent `json:"content"`
	Children []Block         `json:"children"`
}

// InlineContent represents inline content (text, links, etc.)
type InlineContent struct {
	Type   string          `json:"type"` // "text", "link"
	Text   string          `json:"text,omitempty"`
	Href   string          `json:"href,omitempty"`
	Styles map[string]bool `json:"styles"`
}

// BlockNote block types
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

// BlocksToJSON converts []Block to datatypes.JSON for storage
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

// JSONToBlocks converts datatypes.JSON to []Block
func JSONToBlocks(data datatypes.JSON) []Block {
	if len(data) == 0 {
		return []Block{}
	}
	var blocks []Block
	if err := json.Unmarshal(data, &blocks); err != nil {
		return []Block{}
	}
	return blocks
}

type VersionItem struct {
	Id          string
	Version     int
	Name        string
	Description string
	UserId      string
	UserName    string
	CreatedAt   time.Time
}

type ReorderItem struct {
	Id       string
	Position int
}

type SearchResultItem struct {
	Id        string
	Name      string
	Slug      string
	SpaceId   string
	SpaceName string
	Icon      string
	UpdatedAt time.Time
}
