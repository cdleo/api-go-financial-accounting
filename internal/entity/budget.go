package entity

import (
	"context"
)

type BudgetType string

const (
	Monthly BudgetType = "Monthly"
	Travel  BudgetType = "Travel"
)

var BudgetTypes = []BudgetType{
	Monthly,
	Travel,
}

type BudgetInfo struct {
	ID          string `json:"id,omitempty"` /* ObjectID */
	Description string `json:"description"`  /* Optional: Descriptive text*/
}

type Budget struct {
	ID          string           `json:"-"`
	Description string           `json:"description"` /* Optional: Descriptive text*/
	Type        BudgetType       `json:"type"`        /* Monthly or Travel */
	Month       int              `json:"month"`       /* Month of the Budget */
	Year        int              `json:"year"`        /* Year of the Budget */
	Accounts    []*BudgetAccount `json:"accounts"`
}

type BudgetAccount struct {
	Account `json:",inline"`
	Planned []ExpensePlanned `json:"planned"`
}

type BudgetCreate interface {
	CreateBudget(ctx context.Context, value *Budget) error
}

type BudgetRetrieve interface {
	GetBudgetInfo(ctx context.Context) ([]*BudgetInfo, error)
	GetBudgetById(ctx context.Context, id string) (*Budget, error)
	GetBudgetByDate(ctx context.Context, month int, year int) (*Budget, error)
}

type BudgetUpdate interface {
	UpdateBudget(ctx context.Context, value Budget) error
}

type BudgetRepository interface {
	GetByDate(ctx context.Context, month int, year int) (*Budget, error)
	GetById(ctx context.Context, id string) (*Budget, error)
	GetInfo(ctx context.Context) ([]*BudgetInfo, error)
	GetInfoById(ctx context.Context, id string) (*BudgetInfo, error)
	Save(ctx context.Context, value *Budget) error
	Update(ctx context.Context, value *Budget) error
}
