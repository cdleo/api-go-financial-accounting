package service

import (
	"context"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
)

type budgetRetrieve struct {
	repo           entity.BudgetRepository
	acountRetrieve entity.AccountRetrieve
}

func NewBudgetRetrieve(repo entity.BudgetRepository, acountRetrieve entity.AccountRetrieve) entity.BudgetRetrieve {
	return &budgetRetrieve{
		repo:           repo,
		acountRetrieve: acountRetrieve,
	}
}

func (s *budgetRetrieve) GetBudgets() ([]*entity.Budget, error) {
	return s.repo.GetAll(context.Background())
}

func (s *budgetRetrieve) GetBudgetByDate(month int, year int) (*entity.Budget, error) {
	return s.repo.GetByDate(context.Background(), month, year)
}

func (s *budgetRetrieve) GetBudgetById(id string) (*entity.Budget, error) {
	return s.repo.GetById(context.Background(), id)
}
