package domains

import "time"

// OAuthProvider represents supported OAuth providers
type OAuthProvider string

const (
	GoogleOAuthProvider   OAuthProvider = "google"
	FacebookOAuthProvider OAuthProvider = "facebook"
)

// OAuthProfile represents user profile data from OAuth providers
type OAuthProfile struct {
	ID            string        `json:"id" bson:"_id,omitempty"`
	Provider      OAuthProvider `json:"provider" bson:"provider"`
	ProviderID    string        `json:"providerId" bson:"providerId"`
	Email         string        `json:"email" bson:"email"`
	EmailVerified bool          `json:"emailVerified" bson:"emailVerified"`
	Name          string        `json:"name" bson:"name"`
	FirstName     string        `json:"firstName" bson:"firstName"`
	LastName      string        `json:"lastName" bson:"lastName"`
	Picture       string        `json:"picture" bson:"picture"`
	Locale        string        `json:"locale" bson:"locale"`
	CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// OAuthToken represents OAuth access tokens and refresh tokens
type OAuthToken struct {
	ID           string        `json:"id" bson:"_id,omitempty"`
	UserID       string        `json:"userId" bson:"userId"`
	Provider     OAuthProvider `json:"provider" bson:"provider"`
	AccessToken  string        `json:"accessToken" bson:"accessToken"`
	TokenType    string        `json:"tokenType" bson:"tokenType"`
	RefreshToken string        `json:"refreshToken" bson:"refreshToken"`
	ExpiresIn    int           `json:"expiresIn" bson:"expiresIn"`
	Scope        string        `json:"scope" bson:"scope"`
	CreatedAt    time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// OAuthConfig represents OAuth provider configuration
type OAuthConfig struct {
	Provider     OAuthProvider `json:"provider" bson:"provider"`
	ClientID     string        `json:"clientId" bson:"clientId"`
	ClientSecret string        `json:"clientSecret" bson:"clientSecret"`
	RedirectURL  string        `json:"redirectUrl" bson:"redirectUrl"`
	Scopes       []string      `json:"scopes" bson:"scopes"`
}
