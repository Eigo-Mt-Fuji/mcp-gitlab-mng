package domain

import "time"

type Repository struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	NameWithNamespace string `json:"name_with_namespace"`
	Path          string    `json:"path"`
	PathWithNamespace string `json:"path_with_namespace"`
	HTTPURLToRepo string    `json:"http_url_to_repo"`
	SSHURLToRepo  string    `json:"ssh_url_to_repo"`
	WebURL        string    `json:"web_url"`
	Description   string    `json:"description"`
	DefaultBranch string    `json:"default_branch"`
	Visibility    string    `json:"visibility"`
	CreatedAt     time.Time `json:"created_at"`
	LastActivityAt time.Time `json:"last_activity_at"`
}

type RepositoryService interface {
	ListRepositories(groupPath string) ([]*Repository, error)
	GetRepository(projectID int) (*Repository, error)
}