package document

import (
	"fmt"

	"github.com/labbs/nexo/application/document/dto"
)

func (a *DocumentApplication) ToggleTemplate(input dto.ToggleTemplateInput) (*dto.ToggleTemplateOutput, error) {
	logger := a.Logger.With().Str("component", "application.document.toggle_template").Logger()

	if err := a.DocumentPers.SetTemplate(input.DocumentId, input.IsTemplate, input.Category, input.UserId); err != nil {
		logger.Error().Err(err).Str("document_id", input.DocumentId).Msg("failed to set template flag")
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	doc, err := a.DocumentPers.GetDocumentWithPermissions(input.DocumentId, input.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated document: %w", err)
	}

	return &dto.ToggleTemplateOutput{Document: doc}, nil
}
