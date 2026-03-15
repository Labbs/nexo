package apikey

import (
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type ApiKeyApplication struct {
	Config     config.Config
	Logger     zerolog.Logger
	ApiKeyPers domain.ApiKeyPers
}

func NewApiKeyApplication(config config.Config, logger zerolog.Logger, apiKeyPers domain.ApiKeyPers) *ApiKeyApplication {
	return &ApiKeyApplication{
		Config:     config,
		Logger:     logger,
		ApiKeyPers: apiKeyPers,
	}
}
