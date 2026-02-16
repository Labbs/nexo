package dto

type GetRunsInput struct {
	UserId   string
	ActionId string
	Limit    int
}

type GetRunsOutput struct {
	Runs []RunItem
}
