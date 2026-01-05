package dtos

import (
	spaceDtos "github.com/labbs/nexo/interfaces/http/v1/space/dtos"
)

type GetMySpacesResponse struct {
	Spaces []spaceDtos.Space `json:"spaces"`
}
