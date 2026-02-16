package dto

type CreateVersionInput struct {
	UserId      string
	DocumentId  string
	Description string
}

type CreateVersionOutput struct {
	VersionId string
	Version   int
}
