package service

import (
	"context"
	"fmt"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
)

type makeTransference struct {
	repo entity.AccountRepository
}

func NewMakeTransference(repo entity.AccountRepository) entity.MakeTransference {
	return &makeTransference{
		repo: repo,
	}
}

func (m *makeTransference) MakeTransference(ctx context.Context, transference entity.Transfer) error {

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
		Type:        entity.Debit,
		SubType:     entity.Transference,
		Description: transference.Details.Description,
		Amount:      transference.Details.Amount + transference.Details.FromFee,
	}

	fromAccount.Movements = append(fromAccount.Movements, fromDetails)
	fromAccount.Balance += (fromDetails.Type.Modifier() * fromDetails.Amount)
	if fromAccount.Balance < 0 {
		return fmt.Errorf("Insufficient founds")
	}
	accounts = append(accounts, *fromAccount)

	if toAccount.Movements == nil {
		toAccount.Movements = make([]entity.TrxDetails, 0)
	}
	toDetails := entity.TrxDetails{
		Date:        transference.Details.Date,
		Type:        entity.Credit,
		SubType:     entity.Transference,
		Description: transference.Details.Description,
		Amount:      (transference.Details.Amount - transference.Details.ToFee) * transference.Details.Rate,
	}

	toAccount.Movements = append(toAccount.Movements, toDetails)
	toAccount.Balance += (toDetails.Type.Modifier() * toDetails.Amount)
	accounts = append(accounts, *toAccount)

	return m.repo.UpdateMany(ctx, accounts)
}
