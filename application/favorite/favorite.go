package favorite

import (
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type FavoriteApp struct {
	Config       config.Config
	Logger       zerolog.Logger
	FavoritePers domain.FavoritePers
	DocumentPers domain.DocumentPers
}

func NewFavoriteApp(config config.Config, logger zerolog.Logger, favoritePers domain.FavoritePers, documentPers domain.DocumentPers) *FavoriteApp {
	return &FavoriteApp{
		Config:       config,
		Logger:       logger,
		FavoritePers: favoritePers,
		DocumentPers: documentPers,
	}
}
