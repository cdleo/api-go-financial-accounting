package service

import (
	"context"
	"fmt"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
)

type MakeTransference interface {
	MakeTransference(transference entity.Transference) error
}

type makeTransference struct {
	repo entity.AccountRepository
}

func NewMakeTransference(repo entity.AccountRepository) MakeTransference {
	return &makeTransference{
		repo: repo,
	}
}

func (m *makeTransference) MakeTransference(transference entity.Transference) error {

	ctx := context.Background()

	fromAccount, err := m.repo.GetByID(ctx, transference.FromID)
	if err != nil {
		return err
	}

	toAccount, err := m.repo.GetByID(ctx, transference.ToID)
	if err != nil {
		return err
	}

	if (fromAccount == nil) || (toAccount == nil) {
		return fmt.Errorf("Account does not exists")
	}

	accounts := []entity.Account{}

	if fromAccount.Movements == nil {
		fromAccount.Movements = make([]entity.TrxDetails, 0)
	}
	fromDetails := entity.TrxDetails{
		Date:        transference.Details.Date,
		Kind:        entity.Debit,
		Description: transference.Details.Description,
		Amount:      transference.Details.Amount + transference.Details.FromFee,
	}

	fromAccount.Movements = append(fromAccount.Movements, fromDetails)
	fromAccount.Balance += (fromDetails.Kind.Modifier() * fromDetails.Amount)
	if fromAccount.Balance < 0 {
		return fmt.Errorf("Insufficient founds")
	}
	accounts = append(accounts, *fromAccount)

	if toAccount.Movements == nil {
		toAccount.Movements = make([]entity.TrxDetails, 0)
	}
	toDetails := entity.TrxDetails{
		Date:        transference.Details.Date,
		Kind:        entity.Credit,
		Description: transference.Details.Description,
		Amount:      (transference.Details.Amount - transference.Details.ToFee) * transference.Details.Rate,
	}

	toAccount.Movements = append(toAccount.Movements, toDetails)
	toAccount.Balance += (toDetails.Kind.Modifier() * toDetails.Amount)
	accounts = append(accounts, *toAccount)

	return m.repo.UpdateMany(ctx, accounts)
}
