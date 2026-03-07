package space

import (
	"github.com/labbs/nexo/application/ports"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type SpaceApplication struct {
	Config                config.Config
	Logger                zerolog.Logger
	SpacePres             domain.SpacePers
	DocumentApplication   ports.DocumentPort
	PermissionApplication ports.PermissionPort
}

func NewSpaceApplication(config config.Config, logger zerolog.Logger, spacePers domain.SpacePers) *SpaceApplication {
	return &SpaceApplication{
		Config:    config,
		Logger:    logger,
		SpacePres: spacePers,
	}
}

// SpaceApp is a type alias for backward compatibility
type SpaceApp = SpaceApplication
