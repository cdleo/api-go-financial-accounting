package service

import (
	"context"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
)

type budgetUpdate struct {
	repo entity.BudgetRepository
}

func NewBudgetUpdate(repo entity.BudgetRepository) entity.BudgetUpdate {
	return &budgetUpdate{
		repo: repo,
	}
}

func (s *budgetUpdate) UpdateBudget(ctx context.Context, value entity.Budget) error {

	_, err := s.repo.GetById(ctx, value.ID)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, &value)
}
