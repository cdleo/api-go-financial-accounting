package entity

import (
	"context"
)

type AccountSource string

const (
	SourceAll        AccountSource = "All"
	SourceBudget     AccountSource = "Budget"
	SourceStandalone AccountSource = "Standalone"
)

var AccountSources = []AccountSource{
	SourceAll,
	SourceBudget,
	SourceStandalone,
}

type AccountHeader struct {
	Type     AccountType  `json:"type"`           /* Bank, Cash, etc.*/
	Name     string       `json:"name,omitempty"` /* Optional: Descriptive text*/
	Currency CurrencyCode `json:"currency"`       /* Account currency */
}

type Account struct {
	ID       string      `json:"-"`                /* ObjectID */
	ParentID string      `json:"-"`                /* Budget Id or nil value */
	Budget   *BudgetInfo `json:"budget,omitempty"` /* Info of the parent budget */
	AccountHeader
	Balance   float32      `json:"balance"`   /* Account balance (calculated from movements) */
	Movements []TrxDetails `json:"movements"` /* Actual transactions in the account */
}

type AccountInfo struct {
	ID     string      `json:"id,omitempty"`     /* ObjectID */
	Budget *BudgetInfo `json:"budget,omitempty"` /* Info of the parent budget */
	AccountHeader
	Balance float32 `json:"balance"` /* Account balance (calculated from movements) */
}

type AccountCreate interface {
	CreateAccount(ctx context.Context, value *Account) error
}

type AccountRetrieve interface {
	GetAccounts(ctx context.Context) ([]*AccountInfo, error)
	GetAccountByID(ctx context.Context, accountId string) (*Account, error)
}

type AccountUpdate interface {
	UpdateAccount(ctx context.Context, value Account) error
}

//go:generate mockgen -package mocks -destination mocks/accountingRecord.go . AccountingRecordRepository
type AccountRepository interface {
	GetByID(ctx context.Context, id string) (*Account, error)
	GetByHeader(ctx context.Context, header AccountHeader, parentId string) (*Account, error)
	GetInfo(ctx context.Context, budgetRepo BudgetRepository) ([]*AccountInfo, error)
	Save(ctx context.Context, value *Account) error
	Update(ctx context.Context, value Account) error
	UpdateMany(ctx context.Context, values []Account) error
}
