package space

import (
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type SpaceApp struct {
	Config       config.Config
	Logger       zerolog.Logger
	SpacePres    domain.SpacePers
	DocumentPers domain.DocumentPers
}

func NewSpaceApp(config config.Config, logger zerolog.Logger, spacePers domain.SpacePers, documentPers domain.DocumentPers) *SpaceApp {
	return &SpaceApp{
		Config:       config,
		Logger:       logger,
		SpacePres:    spacePers,
		DocumentPers: documentPers,
	}
}
