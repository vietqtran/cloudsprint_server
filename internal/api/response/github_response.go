package response

import "cloud-sprint/internal/service"

type GitHubRepositoryResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	URL         string `json:"url"`
	CloneURL    string `json:"clone_url"`
	Language    string `json:"language"`
	Fork        bool   `json:"fork"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func NewGitHubRepositoryResponse(repo service.GitHubRepository) GitHubRepositoryResponse {
	return GitHubRepositoryResponse{
		ID:          repo.ID,
		Name:        repo.Name,
		FullName:    repo.FullName,
		Description: repo.Description,
		Private:     repo.Private,
		URL:         repo.HTMLURL,
		CloneURL:    repo.CloneURL,
		Language:    repo.Language,
		Fork:        repo.Fork,
		CreatedAt:   repo.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   repo.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func NewGitHubRepositoriesResponse(repos []service.GitHubRepository) []GitHubRepositoryResponse {
	response := make([]GitHubRepositoryResponse, len(repos))
	for i, repo := range repos {
		response[i] = NewGitHubRepositoryResponse(repo)
	}
	return response
}

type GitHubUserResponse struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func NewGitHubUserResponse(userInfo service.GitHubUserInfo) GitHubUserResponse {
	return GitHubUserResponse{
		ID:        userInfo.ID,
		Login:     userInfo.Login,
		Name:      userInfo.Name,
		Email:     userInfo.Email,
		AvatarURL: userInfo.AvatarURL,
	}
}
