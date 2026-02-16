package dto

type GetDocumentsFromSpaceInput struct {
	SpaceId  string
	UserId   string
	ParentId *string
}

type GetDocumentsFromSpaceOutput struct {
	Documents []Document
}
