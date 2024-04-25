package service

import (
	"context"
	"errors"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
)

type MakeTransaction interface {
	MakeTransaction(accountRecord entity.Transaction) error
}

type makeTransaction struct {
	repo entity.AccountRepository
}

func NewMakeTransaction(repo entity.AccountRepository) MakeTransaction {
	return &makeTransaction{
		repo: repo,
	}
}

func (s *makeTransaction) MakeTransaction(accountRecord entity.Transaction) error {

	ctx := context.Background()

	account, err := s.repo.GetByID(ctx, accountRecord.AccountID)
	if err != nil {
		return err
	}

	if account.Movements == nil {
		account.Movements = make([]entity.TrxDetails, 0)
	}
	account.Movements = append(account.Movements, accountRecord.Details)

	account.Balance += (accountRecord.Details.Kind.Modifier() * accountRecord.Details.Amount)

	if account.Balance < 0 {
		return errors.New("Insufficient founds")
	}

	return s.repo.Update(context.Background(), *account)
}
