package document

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/document/dto"
	"github.com/labbs/nexo/domain"
)

const MaxVersionsPerDocument = 50

func (app *DocumentApp) ListVersions(input dto.ListVersionsInput) (*dto.ListVersionsOutput, error) {
	// Verify user has access to the document (accept ID or slug)
	doc, err := app.DocumentPers.GetDocumentByIdOrSlugWithUserPermissions(input.SpaceId, &input.DocumentId, &input.DocumentId, input.UserId)
	if err != nil {
		return nil, fmt.Errorf("document not found or access denied: %w", err)
	}

	versions, err := app.DocumentVersionPers.GetByDocumentId(doc.Id, input.Limit, input.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	totalCount, err := app.DocumentVersionPers.GetVersionCount(doc.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get version count: %w", err)
	}

	output := &dto.ListVersionsOutput{
		Versions:   make([]dto.VersionItem, len(versions)),
		TotalCount: totalCount,
	}

	for i, v := range versions {
		output.Versions[i] = dto.VersionItem{
			Id:          v.Id,
			Version:     v.Version,
			Name:        v.Name,
			Description: v.Description,
			UserId:      v.UserId,
			UserName:    v.User.Username,
			CreatedAt:   v.CreatedAt,
		}
	}

	return output, nil
}

func (app *DocumentApp) GetVersion(input dto.GetVersionInput) (*dto.GetVersionOutput, error) {
	version, err := app.DocumentVersionPers.GetById(input.VersionId)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	// Verify user has access to the document
	_, err = app.DocumentPers.GetDocumentWithPermissions(version.DocumentId, input.UserId)
	if err != nil {
		return nil, fmt.Errorf("access denied: %w", err)
	}

	// Parse content
	var blocks []dto.Block
	if version.Content != nil {
		if err := json.Unmarshal(version.Content, &blocks); err != nil {
			blocks = []dto.Block{}
		}
	}

	return &dto.GetVersionOutput{
		Id:         version.Id,
		Version:    version.Version,
		DocumentId: version.DocumentId,
		Name:       version.Name,
		Content:    blocks,
		Config: dto.DocumentConfig{
			FullWidth:        version.Config.FullWidth,
			Icon:             version.Config.Icon,
			Lock:             version.Config.Lock,
			HeaderBackground: version.Config.HeaderBackground,
		},
		Description: version.Description,
		UserId:      version.UserId,
		UserName:    version.User.Username,
		CreatedAt:   version.CreatedAt,
	}, nil
}

func (app *DocumentApp) RestoreVersion(input dto.RestoreVersionInput) error {
	version, err := app.DocumentVersionPers.GetById(input.VersionId)
	if err != nil {
		return fmt.Errorf("version not found: %w", err)
	}

	// Verify user has editor access to the document
	doc, err := app.DocumentPers.GetDocumentWithPermissions(version.DocumentId, input.UserId)
	if err != nil {
		return fmt.Errorf("access denied: %w", err)
	}

	if !doc.HasPermission(input.UserId, domain.PermissionRoleEditor) {
		return fmt.Errorf("access denied: insufficient permissions")
	}

	// Create a new version before restoring (to preserve current state)
	if _, err := app.createVersionFromDocument(doc, input.UserId, fmt.Sprintf("Before restore to version %d", version.Version)); err != nil {
		return fmt.Errorf("failed to create backup version: %w", err)
	}

	// Restore the document to the selected version
	doc.Name = version.Name
	doc.Content = version.Content
	doc.Config = version.Config

	if err := app.DocumentPers.Update(doc, input.UserId); err != nil {
		return fmt.Errorf("failed to restore document: %w", err)
	}

	return nil
}

func (app *DocumentApp) CreateVersion(input dto.CreateVersionInput) (*dto.CreateVersionOutput, error) {
	// Verify user has editor access to the document
	doc, err := app.DocumentPers.GetDocumentWithPermissions(input.DocumentId, input.UserId)
	if err != nil {
		return nil, fmt.Errorf("document not found or access denied: %w", err)
	}

	if !doc.HasPermission(input.UserId, domain.PermissionRoleEditor) {
		return nil, fmt.Errorf("access denied: insufficient permissions")
	}

	version, err := app.createVersionFromDocument(doc, input.UserId, input.Description)
	if err != nil {
		return nil, err
	}

	return &dto.CreateVersionOutput{
		VersionId: version.Id,
		Version:   version.Version,
	}, nil
}

// createVersionFromDocument creates a new version snapshot of a document
func (app *DocumentApp) createVersionFromDocument(doc *domain.Document, userId string, description string) (*domain.DocumentVersion, error) {
	// Get the next version number
	latestVersion, err := app.DocumentVersionPers.GetLatestVersion(doc.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version: %w", err)
	}

	nextVersion := 1
	if latestVersion != nil {
		nextVersion = latestVersion.Version + 1
	}

	version := &domain.DocumentVersion{
		Id:          uuid.New().String(),
		DocumentId:  doc.Id,
		UserId:      userId,
		Version:     nextVersion,
		Name:        doc.Name,
		Content:     doc.Content,
		Config:      doc.Config,
		Description: description,
		CreatedAt:   time.Now(),
	}

	if err := app.DocumentVersionPers.Create(version); err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	// Clean up old versions if we exceed the limit
	count, err := app.DocumentVersionPers.GetVersionCount(doc.Id)
	if err == nil && count > MaxVersionsPerDocument {
		_ = app.DocumentVersionPers.DeleteOldVersions(doc.Id, MaxVersionsPerDocument)
	}

	return version, nil
}

// CreateVersionOnUpdate should be called when a document is updated to create an automatic version
func (app *DocumentApp) CreateVersionOnUpdate(doc *domain.Document, userId string) {
	logger := app.Logger.With().Str("component", "application.document.version_on_update").Logger()
	// Create version silently - don't fail the update if versioning fails
	_, err := app.createVersionFromDocument(doc, userId, "Auto-save")
	if err != nil {
		logger.Error().Err(err).Str("document_id", doc.Id).Msg("failed to create version on update")
	}
}
