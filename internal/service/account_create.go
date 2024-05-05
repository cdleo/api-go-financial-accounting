package service

import (
	"context"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
)

type accountcreate struct {
	repo entity.AccountRepository
}

func NewAccountCreate(repo entity.AccountRepository) entity.AccountCreate {
	return &accountcreate{
		repo: repo,
	}
}

func (s *accountcreate) CreateAccount(ctx context.Context, value *entity.Account) error {
	return s.repo.Save(ctx, value)
}
