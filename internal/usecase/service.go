package usecase

import "mcp-gitlab-mng/internal/domain"

type Service struct {
	repo Repository
}

type Repository interface {
	GetEntity(id string) (*domain.Entity, error)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetEntity(id string) (*domain.Entity, error) {
	return s.repo.GetEntity(id)
}