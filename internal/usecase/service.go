package usecase

import (
	"context"
	"mcp-gitlab-mng/internal/domain"
)

type GitLabService struct {
	repo domain.RepositoryService
}

func NewGitLabService(repo domain.RepositoryService) *GitLabService {
	return &GitLabService{
		repo: repo,
	}
}

func (s *GitLabService) ListRepositories(ctx context.Context, groupPath string, onlyPrivate bool) ([]*domain.Repository, error) {
	return s.repo.ListRepositories(ctx, groupPath, onlyPrivate)
}

func (s *GitLabService) GetRepository(projectID int) (*domain.Repository, error) {
	return s.repo.GetRepository(projectID)
}

func (s *GitLabService) ListTerraformVersions(ctx context.Context, repositories []*domain.Repository) ([]*domain.RepositoryTerraformVersions, error) {
	return s.repo.ListTerraformVersions(ctx, repositories)
}
