package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"diandi-backend/domains"
	"diandi-backend/services"

	"github.com/gin-gonic/gin"
)

// OAuthHandler handles OAuth-related HTTP requests
type OAuthHandler struct {
	oauthService services.OAuthService
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(oauthService services.OAuthService) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
	}
}

// RegisterRoutes registers the OAuth routes
func (h *OAuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	oauth := router.Group("/oauth")
	{
		oauth.GET("/login/:provider", h.HandleOAuthLogin)
		oauth.GET("/callback/:provider", h.HandleOAuthCallback)
		oauth.POST("/unlink/:provider", h.HandleUnlinkAccount)
	}
}

// HandleOAuthLogin initiates the OAuth flow
func (h *OAuthHandler) HandleOAuthLogin(c *gin.Context) {
	provider := domains.OAuthProvider(c.Param("provider"))

	// Generate random state
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}
	state := base64.URLEncoding.EncodeToString(b)

	// Store state in session or cookie
	c.SetCookie("oauth_state", state, 3600, "/", "", false, true)

	// Get authorization URL
	url, err := h.oauthService.GetAuthURL(provider, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Redirect to provider's consent page
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleOAuthCallback processes the OAuth callback
func (h *OAuthHandler) HandleOAuthCallback(c *gin.Context) {
	provider := domains.OAuthProvider(c.Param("provider"))
	code := c.Query("code")
	state := c.Query("state")

	// Verify state
	savedState, err := c.Cookie("oauth_state")
	if err != nil || state != savedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	// Clear state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// Exchange code for token
	token, err := h.oauthService.ExchangeCode(c.Request.Context(), provider, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get user profile
	profile, err := h.oauthService.GetUserProfile(c.Request.Context(), provider, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from session or create new user
	userID := c.GetString("user_id") // Assuming user ID is set in middleware
	if userID == "" {
		// Create new user logic here
		// This should be handled by your user service
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation not implemented"})
		return
	}

	// Link account
	if err := h.oauthService.LinkAccount(c.Request.Context(), userID, profile, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Redirect to frontend with success
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully linked account",
		"profile": profile,
	})
}

// HandleUnlinkAccount unlinks a social account from the user
func (h *OAuthHandler) HandleUnlinkAccount(c *gin.Context) {
	provider := domains.OAuthProvider(c.Param("provider"))
	userID := c.GetString("user_id") // Assuming user ID is set in middleware

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.oauthService.UnlinkAccount(c.Request.Context(), userID, provider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unlinked account"})
}
