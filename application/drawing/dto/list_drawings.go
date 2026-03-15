package dto

type ListDrawingsInput struct {
	UserId  string
	SpaceId string
}

type ListDrawingsOutput struct {
	Drawings []DrawingItem
}
