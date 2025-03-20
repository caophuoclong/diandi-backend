package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"

	"diandi-backend/domains"
)

// OAuthService defines the interface for OAuth operations
type OAuthService interface {
	// Configuration
	GetAuthURL(provider domains.OAuthProvider, state string) (string, error)

	// OAuth Flow
	ExchangeCode(ctx context.Context, provider domains.OAuthProvider, code string) (*domains.OAuthToken, error)
	GetUserProfile(ctx context.Context, provider domains.OAuthProvider, token *domains.OAuthToken) (*domains.OAuthProfile, error)

	// Token Management
	RefreshToken(ctx context.Context, token *domains.OAuthToken) (*domains.OAuthToken, error)
	RevokeToken(ctx context.Context, token *domains.OAuthToken) error

	// User Management
	LinkAccount(ctx context.Context, userID string, profile *domains.OAuthProfile, token *domains.OAuthToken) error
	UnlinkAccount(ctx context.Context, userID string, provider domains.OAuthProvider) error
}

// oauthService implements OAuthService
type oauthService struct {
	configs map[domains.OAuthProvider]*oauth2.Config
	repo    OAuthRepository
}

// OAuthRepository defines the interface for OAuth data persistence
type OAuthRepository interface {
	SaveToken(ctx context.Context, token *domains.OAuthToken) error
	GetToken(ctx context.Context, userID string, provider domains.OAuthProvider) (*domains.OAuthToken, error)
	DeleteToken(ctx context.Context, userID string, provider domains.OAuthProvider) error
	SaveProfile(ctx context.Context, profile *domains.OAuthProfile) error
	GetProfile(ctx context.Context, providerID string, provider domains.OAuthProvider) (*domains.OAuthProfile, error)
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(repo OAuthRepository, configs map[domains.OAuthProvider]*domains.OAuthConfig) OAuthService {
	oauthConfigs := make(map[domains.OAuthProvider]*oauth2.Config)

	for provider, config := range configs {
		var endpoint oauth2.Endpoint
		switch provider {
		case domains.GoogleOAuthProvider:
			endpoint = google.Endpoint
		case domains.FacebookOAuthProvider:
			endpoint = facebook.Endpoint
		}

		oauthConfigs[provider] = &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,
			Scopes:       config.Scopes,
			Endpoint:     endpoint,
		}
	}

	return &oauthService{
		configs: oauthConfigs,
		repo:    repo,
	}
}

func (s *oauthService) GetAuthURL(provider domains.OAuthProvider, state string) (string, error) {
	config, ok := s.configs[provider]
	if !ok {
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}

	return config.AuthCodeURL(state), nil
}

func (s *oauthService) ExchangeCode(ctx context.Context, provider domains.OAuthProvider, code string) (*domains.OAuthToken, error) {
	config, ok := s.configs[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	oauthToken := &domains.OAuthToken{
		Provider:     provider,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int(token.Expiry.Sub(time.Now()).Seconds()),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return oauthToken, nil
}

func (s *oauthService) GetUserProfile(ctx context.Context, provider domains.OAuthProvider, token *domains.OAuthToken) (*domains.OAuthProfile, error) {
	var profileURL string
	switch provider {
	case domains.GoogleOAuthProvider:
		profileURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	case domains.FacebookOAuthProvider:
		profileURL = "https://graph.facebook.com/v12.0/me?fields=id,email,first_name,last_name,picture,locale"
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", profileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user profile: %s", string(body))
	}

	var profile domains.OAuthProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse profile data: %w", err)
	}

	profile.Provider = provider
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	return &profile, nil
}

func (s *oauthService) RefreshToken(ctx context.Context, token *domains.OAuthToken) (*domains.OAuthToken, error) {
	config, ok := s.configs[token.Provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", token.Provider)
	}

	t := &oauth2.Token{
		RefreshToken: token.RefreshToken,
	}

	newToken, err := config.TokenSource(ctx, t).Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	token.AccessToken = newToken.AccessToken
	token.TokenType = newToken.TokenType
	token.ExpiresIn = int(newToken.Expiry.Sub(time.Now()).Seconds())
	token.UpdatedAt = time.Now()

	if err := s.repo.SaveToken(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to save refreshed token: %w", err)
	}

	return token, nil
}

func (s *oauthService) RevokeToken(ctx context.Context, token *domains.OAuthToken) error {
	var revokeURL string
	switch token.Provider {
	case domains.GoogleOAuthProvider:
		revokeURL = fmt.Sprintf("https://accounts.google.com/o/oauth2/revoke?token=%s", token.AccessToken)
	case domains.FacebookOAuthProvider:
		revokeURL = fmt.Sprintf("https://graph.facebook.com/v12.0/me/permissions?access_token=%s", token.AccessToken)
	default:
		return fmt.Errorf("unsupported provider: %s", token.Provider)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", revokeURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create revoke request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to revoke token: %s", string(body))
	}

	return s.repo.DeleteToken(ctx, token.UserID, token.Provider)
}

func (s *oauthService) LinkAccount(ctx context.Context, userID string, profile *domains.OAuthProfile, token *domains.OAuthToken) error {
	token.UserID = userID

	if err := s.repo.SaveProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to save profile: %w", err)
	}

	if err := s.repo.SaveToken(ctx, token); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

func (s *oauthService) UnlinkAccount(ctx context.Context, userID string, provider domains.OAuthProvider) error {
	token, err := s.repo.GetToken(ctx, userID, provider)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	if err := s.RevokeToken(ctx, token); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}
