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
