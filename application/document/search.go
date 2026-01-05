package document

import (
	"fmt"

	"github.com/labbs/nexo/application/document/dto"
)

func (app *DocumentApp) Search(input dto.SearchInput) (*dto.SearchOutput, error) {
	if len(input.Query) < 2 {
		return nil, fmt.Errorf("query must be at least 2 characters")
	}

	docs, err := app.DocumentPers.Search(input.Query, input.UserId, input.SpaceId, input.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}

	output := &dto.SearchOutput{
		Results: make([]dto.SearchResultItem, len(docs)),
	}

	for i, doc := range docs {
		output.Results[i] = dto.SearchResultItem{
			Id:        doc.Id,
			Name:      doc.Name,
			Slug:      doc.Slug,
			SpaceId:   doc.SpaceId,
			SpaceName: doc.Space.Name,
			Icon:      doc.Config.Icon,
			UpdatedAt: doc.UpdatedAt,
		}
	}

	return output, nil
}
