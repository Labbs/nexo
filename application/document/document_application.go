package document

import (
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type DocumentApplication struct {
	Config              config.Config
	Logger              zerolog.Logger
	DocumentPers        domain.DocumentPers
	SpacePers           domain.SpacePers
	PermissionPers      domain.PermissionPers
	CommentPers         domain.CommentPers
	DocumentVersionPers domain.DocumentVersionPers
}

func NewDocumentApplication(config config.Config, logger zerolog.Logger, documentPers domain.DocumentPers, spacePers domain.SpacePers, permissionPers domain.PermissionPers, commentPers domain.CommentPers, documentVersionPers domain.DocumentVersionPers) *DocumentApplication {
	return &DocumentApplication{
		Config:              config,
		Logger:              logger,
		DocumentPers:        documentPers,
		SpacePers:           spacePers,
		PermissionPers:      permissionPers,
		CommentPers:         commentPers,
		DocumentVersionPers: documentVersionPers,
	}
}
