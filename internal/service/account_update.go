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

func (s *accountUpdate) UpdateAccount(value entity.Account) error {

	var balance float32
	for _, trx := range value.Movements {

		if trx.Kind == entity.Debit {
			balance -= trx.Amount
		} else if trx.Kind == entity.Credit {
			balance += trx.Amount
		} else {
			return fmt.Errorf("record kind [%s] is not valid", trx.Kind)
		}

	}
	value.Balance = balance

	return s.repo.Update(context.Background(), value)
}
