package category

import (
	"context"

	"github.com/marcos-nsantos/aluraflix-api/internal/entity"
)

func NewService(repo Repor) Service {
	return &service{repo: repo}
}

type service struct {
	repo Repor
}

func (s *service) Post(ctx context.Context, category *entity.Category) error {
	return s.repo.Insert(ctx, category)
}

func (s *service) GetAll(ctx context.Context) ([]*entity.Category, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id uint64) (*entity.Category, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, category *entity.Category) error {
	return s.repo.Update(ctx, category)
}

func (s *service) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
