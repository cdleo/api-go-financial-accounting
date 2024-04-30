package service

import (
	"context"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
)

type budgetUpdate struct {
	repo          entity.BudgetRepository
	accountUpdate entity.AccountUpdate
}

func NewBudgetUpdate(repo entity.BudgetRepository, accountUpdate entity.AccountUpdate) entity.BudgetUpdate {
	return &budgetUpdate{
		repo:          repo,
		accountUpdate: accountUpdate,
	}
}

func (s *budgetUpdate) UpdateBudget(value entity.Budget) error {

	ctx := context.Background()

	_, err := s.repo.GetById(ctx, value.ID)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, &value)
}
