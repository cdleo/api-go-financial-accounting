package entity

import (
	"context"
)

type AccountHeader struct {
	Type     AccountType  `json:"type"`           /* Bank, Cash, etc.*/
	Name     string       `json:"name,omitempty"` /* Optional: Descriptive text*/
	Currency CurrencyCode `json:"currency"`       /* Account currency */
}

type Account struct {
	ID       string `json:"-"` /* ObjectID */
	ParentID string `json:"-"` /* Budget Id or nil value */
	AccountHeader
	Balance   float32      `json:"balance"`   /* Account balance (calculated from movements) */
	Movements []TrxDetails `json:"movements"` /* Actual transactions in the account */
}

type AccountSummary struct {
	ID string `json:"id,omitempty"` /* ObjectID */
	AccountHeader
	Balance float32 `json:"balance"` /* Account balance (calculated from movements) */
}

type AccountCreate interface {
	CreateAccount(value *Account) error
}

type AccountRetrieve interface {
	GetAccounts() ([]*AccountSummary, error)
	GetAccountByID(accountId string) (*Account, error)
}

type AccountUpdate interface {
	UpdateAccount(value Account) error
}

//go:generate mockgen -package mocks -destination mocks/accountingRecord.go . AccountingRecordRepository
type AccountRepository interface {
	GetByID(ctx context.Context, id string) (*Account, error)
	GetByHeader(ctx context.Context, header AccountHeader, parentId string) (*Account, error)
	GetAll(ctx context.Context) ([]*AccountSummary, error)
	Save(ctx context.Context, value *Account) error
	Update(ctx context.Context, value Account) error
	UpdateMany(ctx context.Context, values []Account) error
}
