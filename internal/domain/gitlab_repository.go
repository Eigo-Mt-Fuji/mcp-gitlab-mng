package domain

import (
	"context"
	"time"
)

type Repository struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	NameWithNamespace string    `json:"name_with_namespace"`
	Path              string    `json:"path"`
	PathWithNamespace string    `json:"path_with_namespace"`
	HTTPURLToRepo     string    `json:"http_url_to_repo"`
	SSHURLToRepo      string    `json:"ssh_url_to_repo"`
	WebURL            string    `json:"web_url"`
	Description       string    `json:"description"`
	DefaultBranch     string    `json:"default_branch"`
	Visibility        string    `json:"visibility"`
	CreatedAt         time.Time `json:"created_at"`
	LastActivityAt    time.Time `json:"last_activity_at"`
}

type TerraformComponent struct {
	Path            string `json:"path"`
	RequiredVersion string `json:"required_version"`
	FilePath        string `json:"file_path"`
}

type RepositoryTerraformVersions struct {
	Repository *Repository           `json:"repository"`
	Components []*TerraformComponent `json:"components"`
}

type CodeMatch struct {
	Path     string            `json:"path"`
	FilePath string            `json:"file_path"`
	Matches  map[string]string `json:"matches"`
}

type RepositoryCodeSearch struct {
	Repository *Repository  `json:"repository"`
	Results    []*CodeMatch `json:"results"`
}

type RepositoryService interface {
	ListRepositories(ctx context.Context, groupPath string, onlyPrivate bool) ([]*Repository, error)
	GetRepository(projectID int) (*Repository, error)
	ListTerraformVersions(ctx context.Context, repositories []*Repository) ([]*RepositoryTerraformVersions, error)
	ListCodeByKeyword(ctx context.Context, repositories []*Repository, keyword string, regexPatterns map[string]string) ([]*RepositoryCodeSearch, error)
}
