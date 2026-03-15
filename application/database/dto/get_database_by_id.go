package dto

// DatabaseDetail contains the database data needed by other applications
type DatabaseDetail struct {
	Id        string
	SpaceId   string
	CreatedBy string
}

type GetDatabaseByIdInput struct {
	DatabaseId string
}

type GetDatabaseByIdOutput struct {
	Database *DatabaseDetail
}
