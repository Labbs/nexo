package dto

type UpdateSpaceOrderInput struct {
	UserId   string
	SpaceIds []string
}

type UpdateSpaceOrderOutput struct {
	SpaceIds []string
}
