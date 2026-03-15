package dtos

import (
	ddtos "github.com/labbs/nexo/interfaces/http/v1/document/dtos"
)

type Favorite struct {
	Id       string         `json:"id"`
	SpaceId  string         `json:"space_id"`
	Document ddtos.Document `json:"document"`
	Position int            `json:"position"`
}
