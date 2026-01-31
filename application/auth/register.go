package auth

import (
	"fmt"

	"github.com/labbs/nexo/application/auth/dto"
	d "github.com/labbs/nexo/application/document/dto"
	s "github.com/labbs/nexo/application/space/dto"
	u "github.com/labbs/nexo/application/user/dto"
	"github.com/labbs/nexo/domain"
	"golang.org/x/crypto/bcrypt"
)

func (c *AuthApp) Register(input dto.RegisterInput) error {
	logger := c.Logger.With().Str("component", "application.auth.register").Logger()

	// check if the email is already in use
	_, err := c.UserApp.GetByEmail(u.GetByEmailInput{Email: input.Email})
	if err == nil {
		return fmt.Errorf("email is already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(err).Str("email", input.Email).Msg("failed to hash password")
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := domain.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Active:   true,
	}

	createdUser, err := c.UserApp.Create(u.CreateUserInput{User: user})
	if err != nil {
		logger.Error().Err(err).Str("email", input.Email).Msg("failed to create user")
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Create a private space for the user
	fmt.Println("Creating private space for user", createdUser.User.Id)
	space, err := c.SpaceApp.CreatePrivateSpaceForUser(s.CreatePrivateSpaceForUserInput{UserId: createdUser.User.Id})
	if err != nil {
		logger.Error().Err(err).Str("user_id", createdUser.User.Id).Msg("failed to create private space for user")
		return fmt.Errorf("failed to create private space for user: %w", err)
	}

	// Create a welcome document in the user's private space
	welcomeContent := []d.Block{
		{
			ID:   "welcome-1",
			Type: d.BlockTypeParagraph,
			Props: map[string]any{
				"textColor":       "default",
				"backgroundColor": "default",
				"textAlignment":   "left",
			},
			Content: []d.InlineContent{
				{
					Type:   "text",
					Text:   "This is your private space. Start adding your notes and documents here!",
					Styles: map[string]bool{},
				},
			},
			Children: []d.Block{},
		},
	}

	_, err = c.DocumentApp.CreateDocument(d.CreateDocumentInput{
		Name:    "Welcome to Your Private Space",
		UserId:  createdUser.User.Id,
		SpaceId: space.Space.Id,
		Content: welcomeContent,
	})
	if err != nil {
		logger.Error().Err(err).Str("space_id", space.Space.Id).Str("user_id", createdUser.User.Id).Msg("failed to create welcome document")
		return fmt.Errorf("failed to create welcome document: %w", err)
	}

	return nil
}
