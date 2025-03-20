package config

import (
	"os"

	"diandi-backend/domains"
)

// OAuthConfigs holds all OAuth provider configurations
type OAuthConfigs struct {
	Google   *domains.OAuthConfig
	Facebook *domains.OAuthConfig
}

// LoadOAuthConfigs loads OAuth configurations from environment variables
func LoadOAuthConfigs() *OAuthConfigs {
	return &OAuthConfigs{
		Google: &domains.OAuthConfig{
			Provider:     domains.GoogleOAuthProvider,
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		},
		Facebook: &domains.OAuthConfig{
			Provider:     domains.FacebookOAuthProvider,
			ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
			ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("FACEBOOK_REDIRECT_URL"),
			Scopes: []string{
				"email",
				"public_profile",
			},
		},
	}
}

// GetConfig returns the configuration for a specific provider
func (c *OAuthConfigs) GetConfig(provider domains.OAuthProvider) *domains.OAuthConfig {
	switch provider {
	case domains.GoogleOAuthProvider:
		return c.Google
	case domains.FacebookOAuthProvider:
		return c.Facebook
	default:
		return nil
	}
}

// GetAllConfigs returns a map of all OAuth configurations
func (c *OAuthConfigs) GetAllConfigs() map[domains.OAuthProvider]*domains.OAuthConfig {
	return map[domains.OAuthProvider]*domains.OAuthConfig{
		domains.GoogleOAuthProvider:   c.Google,
		domains.FacebookOAuthProvider: c.Facebook,
	}
}
