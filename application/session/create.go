package session

import (
	"github.com/gofiber/fiber/v2/utils"
	"github.com/labbs/nexo/application/session/dto"
	"github.com/labbs/nexo/domain"
)

func (c *SessionApp) Create(input dto.CreateSessionInput) (*dto.CreateSessionOutput, error) {
	logger := c.Logger.With().Str("component", "application.session.create").Logger()

	session := &domain.Session{
		Id:        utils.UUIDv4(),
		UserId:    input.UserId,
		UserAgent: input.UserAgent,
		IpAddress: input.IpAddress,
		ExpiresAt: input.ExpiresAt,
	}

	err := c.SessionPers.Create(session)
	if err != nil {
		logger.Error().Err(err).Str("session_id", session.Id).Str("user_id", session.UserId).Msg("failed to create session")
		return nil, err
	}

	return &dto.CreateSessionOutput{SessionId: session.Id}, nil
}

func (c *SessionApp) DeleteExpired() error {
	// logger := c.Logger.With().Str("component", "application.session.delete_expired").Logger()

	return nil
}
