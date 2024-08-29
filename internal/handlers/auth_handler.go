package handlers

import (
	"net/http"
	"net/url"

	"pairs-trading-backend/internal/auth"
	"pairs-trading-backend/internal/config"
	"pairs-trading-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	frontendRedirectURL := c.Query("redirect_url")
	if frontendRedirectURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing redirect_url parameter"})
		return
	}

	state := url.QueryEscape(frontendRedirectURL)

	url := auth.GetGoogleOauthConfig(h.cfg).AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	frontendRedirectURL, err := url.QueryUnescape(state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	googleUser, err := auth.GetGoogleUserInfo(code, h.cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info from Google"})
		return
	}

	user, err := models.UpsertUser(h.db, googleUser.Email, googleUser.Name, googleUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert user"})
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.Email, h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	redirectURL, err := url.Parse(frontendRedirectURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid redirect URL"})
		return
	}

	query := redirectURL.Query()
	query.Set("token", token)
	redirectURL.RawQuery = query.Encode()

	c.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
}
