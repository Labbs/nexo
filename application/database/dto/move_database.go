package dto

type MoveDatabaseInput struct {
	UserId     string
	DatabaseId string
	DocumentId *string // nil = move to root (no parent document)
}

type MoveDatabaseOutput struct {
	Id         string
	DocumentId *string
}
