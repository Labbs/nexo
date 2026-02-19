package document

import (
	"github.com/labbs/nexo/application/ports"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type DocumentApplication struct {
	Config              config.Config
	Logger              zerolog.Logger
	DocumentPers        domain.DocumentPers
	CommentPers         domain.CommentPers
	DocumentVersionPers domain.DocumentVersionPers
	SpaceApp            ports.SpacePort
	PermissionApp       ports.PermissionPort
}

func NewDocumentApplication(config config.Config, logger zerolog.Logger, documentPers domain.DocumentPers, commentPers domain.CommentPers, documentVersionPers domain.DocumentVersionPers) *DocumentApplication {
	return &DocumentApplication{
		Config:              config,
		Logger:              logger,
		DocumentPers:        documentPers,
		CommentPers:         commentPers,
		DocumentVersionPers: documentVersionPers,
	}
}
