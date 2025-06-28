package repository

import (
	"context"
	"mcp-gitlab-mng/internal/domain"

	"gitlab.com/gitlab-org/api/client-go"
)

type GitLabRepository struct {
	client *gitlab.Client
}

func NewGitLabRepository(token, baseURL string) (*GitLabRepository, error) {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	if err != nil {
		return nil, err
	}

	return &GitLabRepository{
		client: client,
	}, nil
}

func (r *GitLabRepository) ListRepositories(groupPath string) ([]*domain.Repository, error) {
	var projects []*gitlab.Project
	var err error

	if groupPath != "" {
		// List projects for specific group
		opts := &gitlab.ListGroupProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
			},
		}
		projects, _, err = r.client.Groups.ListGroupProjects(groupPath, opts)
	} else {
		// List all projects accessible to user
		opts := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
			},
		}
		projects, _, err = r.client.Projects.ListProjects(opts)
	}

	if err != nil {
		return nil, err
	}

	repositories := make([]*domain.Repository, len(projects))
	for i, project := range projects {
		repositories[i] = convertProjectToRepository(project)
	}

	return repositories, nil
}

func (r *GitLabRepository) ListRepository(ctx context.Context, onlyPrivate bool) ([]*domain.Repository, error) {
	opts := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
	}

	if onlyPrivate {
		visibility := gitlab.PrivateVisibility
		opts.Visibility = &visibility
	}

	projects, _, err := r.client.Projects.ListProjects(opts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	repositories := make([]*domain.Repository, len(projects))
	for i, project := range projects {
		repositories[i] = convertProjectToRepository(project)
	}

	return repositories, nil
}

func (r *GitLabRepository) GetRepository(projectID int) (*domain.Repository, error) {
	project, _, err := r.client.Projects.GetProject(projectID, nil)
	if err != nil {
		return nil, err
	}

	return convertProjectToRepository(project), nil
}

func convertProjectToRepository(project *gitlab.Project) *domain.Repository {
	repo := &domain.Repository{
		ID:                project.ID,
		Name:              project.Name,
		NameWithNamespace: project.NameWithNamespace,
		Path:              project.Path,
		PathWithNamespace: project.PathWithNamespace,
		HTTPURLToRepo:     project.HTTPURLToRepo,
		SSHURLToRepo:      project.SSHURLToRepo,
		WebURL:            project.WebURL,
		Description:       project.Description,
		DefaultBranch:     project.DefaultBranch,
		Visibility:        string(project.Visibility),
	}

	if project.CreatedAt != nil {
		repo.CreatedAt = *project.CreatedAt
	}
	if project.LastActivityAt != nil {
		repo.LastActivityAt = *project.LastActivityAt
	}

	return repo
}
