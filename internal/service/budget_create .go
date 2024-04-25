package service

import (
	"context"
	"errors"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type Budgetcreate struct {
	repo entity.BudgetRepository
}

func NewBudgetCreate(repo entity.BudgetRepository) entity.BudgetCreate {
	return &Budgetcreate{
		repo: repo,
	}
}

func (s *Budgetcreate) CreateBudget(value *entity.Budget) error {

	ctx := context.Background()

	budget, err := s.repo.GetByDate(ctx, value.Month, value.Year)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if budget != nil {
		return errors.New("already created")
	}

	return s.repo.Save(ctx, value)
}
