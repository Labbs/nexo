package dtos

import (
	documentDtos "github.com/labbs/nexo/interfaces/http/v1/document/dtos"
)

type Favorite struct {
	Id       string                `json:"id"`
	SpaceId  string                `json:"space_id"`
	Document documentDtos.Document `json:"document"`
	Position int                   `json:"position"`
}
