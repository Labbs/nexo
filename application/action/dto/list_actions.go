package dto

type ListActionsInput struct {
	UserId string
}

type ListActionsOutput struct {
	Actions []ActionItem
}
