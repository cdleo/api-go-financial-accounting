package service

import (
	"context"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
)

type accountRetrieve struct {
	repo entity.AccountRepository
}

func NewAccountRetrieve(repo entity.AccountRepository) entity.AccountRetrieve {
	return &accountRetrieve{
		repo: repo,
	}
}

func (s *accountRetrieve) GetAccounts() ([]*entity.AccountSummary, error) {
	return s.repo.GetAll(context.Background())
}

func (s *accountRetrieve) GetAccountByID(accountId string) (*entity.Account, error) {
	return s.repo.GetByID(context.Background(), accountId)
}
