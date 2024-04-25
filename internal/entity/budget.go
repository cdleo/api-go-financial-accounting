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
	Account  `json:",inline"`
	Expected []ExpenseExpected `json:"expected"`
}

type BudgetCreate interface {
	CreateBudget(value *Budget) error
}

type BudgetRetrieve interface {
	GetBudgets() ([]*Budget, error)
	GetBudgetById(id string) (*Budget, error)
	GetBudgetByDate(month int, year int) (*Budget, error)
}

type BudgetUpdate interface {
	UpdateBudget(value Budget) error
}

type BudgetRepository interface {
	GetByDate(ctx context.Context, month int, year int) (*Budget, error)
	GetById(ctx context.Context, id string) (*Budget, error)
	GetAll(ctx context.Context) ([]*Budget, error)
	Save(ctx context.Context, value *Budget) error
	Update(ctx context.Context, value *Budget) error
}
