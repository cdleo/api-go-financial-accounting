package service

import (
	"context"
	"fmt"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
)

type accountUpdate struct {
	repo entity.AccountRepository
}

func NewAccountUpdate(repo entity.AccountRepository) entity.AccountUpdate {
	return &accountUpdate{
		repo: repo,
	}
}

func (s *accountUpdate) UpdateAccount(ctx context.Context, value entity.Account) error {

	accountDB, err := s.repo.GetByID(ctx, value.ID)
	if err != nil {
		return err
	}
	value.ParentID = accountDB.ParentID

	var balance float32
	for _, trx := range value.Movements {

		if trx.Type == entity.Debit {
			balance -= trx.Amount
		} else if trx.Type == entity.Credit {
			balance += trx.Amount
		} else {
			return fmt.Errorf("record type [%s] is not valid", trx.Type)
		}

	}
	value.Balance = balance

	return s.repo.Update(context.Background(), value)
}
