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

func (s *GitLabService) ListRepository(ctx context.Context, onlyPrivate bool) ([]*domain.Repository, error) {
	return s.repo.ListRepository(ctx, onlyPrivate)
}

func (s *GitLabService) ListRepositories(groupPath string) ([]*domain.Repository, error) {
	return s.repo.ListRepositories(groupPath)
}

func (s *GitLabService) GetRepository(projectID int) (*domain.Repository, error) {
	return s.repo.GetRepository(projectID)
}
