package domain

import "time"

// OAuthProvider stores the link between a user and an external SSO identity.
type OAuthProvider struct {
	Id             string
	UserId         string
	Provider       string // e.g. "oidc"
	ProviderUserId string // subject (sub) from the provider
	Email          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (o *OAuthProvider) TableName() string {
	return "oauth_provider"
}

type OAuthProviderPers interface {
	FindByProviderAndSubject(provider, subject string) (OAuthProvider, error)
	FindByUserId(userId string) ([]OAuthProvider, error)
	Create(op OAuthProvider) (OAuthProvider, error)
}
