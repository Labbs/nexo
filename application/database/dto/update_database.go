package dto

type UpdateDatabaseInput struct {
	UserId      string
	DatabaseId  string
	Name        *string
	Description *string
	Icon        *string
	Schema      *[]PropertySchema
	DefaultView *string
}
