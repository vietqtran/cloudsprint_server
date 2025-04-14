package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"cloud-sprint/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GitHubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type GitHubRepository struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	Private     bool      `json:"private"`
	HTMLURL     string    `json:"html_url"`
	CloneURL    string    `json:"clone_url"`
	Language    string    `json:"language"`
	Fork        bool      `json:"fork"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GitHubService struct {
	config config.Config
}

func NewGitHubService(config config.Config) *GitHubService {
	return &GitHubService{
		config: config,
	}
}

func (s *GitHubService) GetOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.config.OAuth.GitHubClientID,
		ClientSecret: s.config.OAuth.GitHubClientSecret,
		RedirectURL:  s.config.OAuth.GitHubRedirectURL,
		Scopes:       []string{"user:email", "read:user", "repo"},
		Endpoint:     github.Endpoint,
	}
}

func (s *GitHubService) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	oauthConfig := s.GetOAuthConfig()
	return oauthConfig.Exchange(ctx, code)
}

func (s *GitHubService) GetUserInfo(token *oauth2.Token) (*GitHubUserInfo, error) {
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GitHubUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	if userInfo.Email == "" {
		userInfo.Email, err = s.getPrimaryEmail(client)
		if err != nil {
			return nil, fmt.Errorf("failed to get primary email: %w", err)
		}
	}

	return &userInfo, nil
}

func (s *GitHubService) GetUserRepositories(token *oauth2.Token) ([]GitHubRepository, error) {
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))

	repos := []GitHubRepository{}
	page := 1
	const perPage = 100

	for {
		url := fmt.Sprintf("https://api.github.com/user/repos?page=%d&per_page=%d&sort=updated", page, perPage)
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to get repositories: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			if err := resp.Body.Close(); err != nil {
				return nil, fmt.Errorf("failed to close response body: %w", err)
			}
			return nil, fmt.Errorf("GitHub API returned non-200 status: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err := resp.Body.Close(); err != nil {
			return nil, fmt.Errorf("failed to close response body: %w", err)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var pageRepos []GitHubRepository
		if err := json.Unmarshal(body, &pageRepos); err != nil {
			return nil, fmt.Errorf("failed to unmarshal repositories: %w", err)
		}

		if len(pageRepos) == 0 {
			break
		}

		repos = append(repos, pageRepos...)

		linkHeader := resp.Header.Get("Link")
		if linkHeader == "" || !containsNextPage(linkHeader) {
			break
		}

		page++
	}

	return repos, nil
}

func (s *GitHubService) getPrimaryEmail(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type GitHubEmail struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	var emails []GitHubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no verified email found")
}

func containsNextPage(linkHeader string) bool {
	return len(linkHeader) > 0 && linkHeader != "" && linkHeader != "undefined"
}
