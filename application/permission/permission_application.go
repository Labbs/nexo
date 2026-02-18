package permission

import (
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type PermissionApplication struct {
	Config         config.Config
	Logger         zerolog.Logger
	PermissionPers domain.PermissionPers
	SpacePers      domain.SpacePers
	DrawingPers    domain.DrawingPers
	DocumentPers   domain.DocumentPers
	DatabasePers   domain.DatabasePers
}

func NewPermissionApplication(
	config config.Config,
	logger zerolog.Logger,
	permissionPers domain.PermissionPers,
	spacePers domain.SpacePers,
	drawingPers domain.DrawingPers,
	documentPers domain.DocumentPers,
	databasePers domain.DatabasePers,
) *PermissionApplication {
	return &PermissionApplication{
		Config:         config,
		Logger:         logger,
		PermissionPers: permissionPers,
		SpacePers:      spacePers,
		DrawingPers:    drawingPers,
		DocumentPers:   documentPers,
		DatabasePers:   databasePers,
	}
}
