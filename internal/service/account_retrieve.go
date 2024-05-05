package service

import (
	"context"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
)

type accountRetrieve struct {
	accountRepo entity.AccountRepository
	budgetRepo  entity.BudgetRepository
}

func NewAccountRetrieve(accountRepo entity.AccountRepository, budgetRepo entity.BudgetRepository) entity.AccountRetrieve {
	return &accountRetrieve{
		accountRepo: accountRepo,
		budgetRepo:  budgetRepo,
	}
}

func (s *accountRetrieve) GetAccounts(ctx context.Context) ([]*entity.AccountInfo, error) {
	return s.accountRepo.GetInfo(ctx, s.budgetRepo)
}

func (s *accountRetrieve) GetAccountByID(ctx context.Context, accountId string) (*entity.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, accountId)
	if err != nil {
		return nil, err
	}
	if len(account.ParentID) > 0 {
		account.Budget, err = s.budgetRepo.GetInfoById(context.Background(), account.ParentID)
	}
	return account, err
}
