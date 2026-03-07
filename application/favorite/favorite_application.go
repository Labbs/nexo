package favorite

import (
	"github.com/labbs/nexo/application/ports"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type FavoriteApplication struct {
	Config              config.Config
	Logger              zerolog.Logger
	FavoritePers        domain.FavoritePers
	DocumentApplication ports.DocumentPort
}

func NewFavoriteApplication(config config.Config, logger zerolog.Logger, favoritePers domain.FavoritePers) *FavoriteApplication {
	return &FavoriteApplication{
		Config:       config,
		Logger:       logger,
		FavoritePers: favoritePers,
	}
}
