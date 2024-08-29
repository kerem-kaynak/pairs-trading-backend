package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"pairs-trading-backend/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func GetGoogleOauthConfig(cfg *config.Config) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GetGoogleUserInfo(code string, cfg *config.Config) (*GoogleUser, error) {
	googleConfig := GetGoogleOauthConfig(cfg)
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	var user GoogleUser
	err = json.Unmarshal(contents, &user)
	if err != nil {
		return nil, fmt.Errorf("failed parsing user info: %s", err.Error())
	}

	return &user, nil
}
