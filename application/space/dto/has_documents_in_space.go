package dto

type HasDocumentsInSpaceInput struct {
	SpaceId string
	UserId  string
}

type HasDocumentsInSpaceOutput struct {
	HasDocuments bool
}
