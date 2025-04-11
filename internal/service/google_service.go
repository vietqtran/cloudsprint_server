package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"cloud-sprint/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type GoogleService struct {
	config       *oauth2.Config
	clientConfig config.Config
}

func NewGoogleService(clientConfig config.Config) *GoogleService {
	config := &oauth2.Config{
		ClientID:     clientConfig.OAuth.GoogleClientID,
		ClientSecret: clientConfig.OAuth.GoogleClientSecret,
		RedirectURL:  clientConfig.OAuth.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleService{
		config:       config,
		clientConfig: clientConfig,
	}
}

func (s *GoogleService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

func (s *GoogleService) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.config.Exchange(ctx, code)
}

func (s *GoogleService) GetUserInfo(token *oauth2.Token) (*GoogleUserInfo, error) {
	client := s.config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info from Google")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
