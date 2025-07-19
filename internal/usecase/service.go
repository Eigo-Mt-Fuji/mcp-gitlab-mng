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

func (s *GitLabService) ListCodeByKeyword(ctx context.Context, repositories []*domain.Repository, keyword string, regexPatterns map[string]string) ([]*domain.RepositoryCodeSearch, error) {
	return s.repo.ListCodeByKeyword(ctx, repositories, keyword, regexPatterns)
}
