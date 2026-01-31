package user

import (
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type UserApp struct {
	Config       config.Config
	Logger       zerolog.Logger
	UserPres     domain.UserPers
	GroupPres    domain.GroupPers
	FavoritePers domain.FavoritePers
}

func NewUserApp(config config.Config, logger zerolog.Logger, userPers domain.UserPers, groupPers domain.GroupPers, favoritePers domain.FavoritePers) *UserApp {
	return &UserApp{
		Config:       config,
		Logger:       logger,
		UserPres:     userPers,
		GroupPres:    groupPers,
		FavoritePers: favoritePers,
	}
}
