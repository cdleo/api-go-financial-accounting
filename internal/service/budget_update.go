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

	_, err := s.repo.GetByDate(ctx, value.Month, value.Year)
	if err != nil {
		return err
	}

	for i := 0; i < len(value.Accounts); i++ {
		if err := s.accountUpdate.UpdateAccount(value.Accounts[i].Account); err != nil {
			return err
		}
	}

	return s.repo.Update(ctx, &value)
}
