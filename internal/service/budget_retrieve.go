package service

import (
	"context"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
)

type budgetRetrieve struct {
	repo entity.BudgetRepository
}

func NewBudgetRetrieve(repo entity.BudgetRepository) entity.BudgetRetrieve {
	return &budgetRetrieve{
		repo: repo,
	}
}

func (s *budgetRetrieve) GetBudgetInfo(ctx context.Context) ([]*entity.BudgetInfo, error) {
	return s.repo.GetInfo(ctx)
}

func (s *budgetRetrieve) GetBudgetByDate(ctx context.Context, month int, year int) (*entity.Budget, error) {
	return s.repo.GetByDate(ctx, month, year)
}

func (s *budgetRetrieve) GetBudgetById(ctx context.Context, id string) (*entity.Budget, error) {
	return s.repo.GetById(ctx, id)
}
