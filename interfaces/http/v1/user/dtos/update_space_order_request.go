package dtos

type UpdateSpaceOrderRequest struct {
	SpaceIds []string `json:"space_ids" validate:"required"`
}

type UpdateSpaceOrderResponse struct {
	SpaceIds []string `json:"space_ids"`
}
