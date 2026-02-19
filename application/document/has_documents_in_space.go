package document

import (
	"github.com/labbs/nexo/application/document/dto"
)

func (a *DocumentApplication) HasDocumentsInSpace(input dto.HasDocumentsInSpaceInput) (*dto.HasDocumentsInSpaceOutput, error) {
	docs, err := a.DocumentPers.GetRootDocumentsFromSpaceWithUserPermissions(input.SpaceId, input.UserId)
	if err != nil {
		return &dto.HasDocumentsInSpaceOutput{HasDocuments: false}, nil
	}

	return &dto.HasDocumentsInSpaceOutput{HasDocuments: len(docs) > 0}, nil
}
