package document

import (
	"fmt"

	"github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
)

func (a *DocumentApp) ReorderDocuments(input dto.ReorderDocumentsInput) error {
	logger := a.Logger.With().Str("component", "application.document.reorder_documents").Logger()

	if len(input.Items) == 0 {
		return fmt.Errorf("no items to reorder")
	}

	// Convert application DTOs to domain items
	domainItems := make([]domain.ReorderItem, len(input.Items))
	for i, item := range input.Items {
		domainItems[i] = domain.ReorderItem{
			Id:       item.Id,
			Position: item.Position,
		}
	}

	if err := a.DocumentPers.Reorder(input.SpaceId, domainItems, input.UserId); err != nil {
		logger.Error().Err(err).Msg("failed to reorder documents")
		return err
	}

	return nil
}
