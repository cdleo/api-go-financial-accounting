package entity

import "context"

type BudgetSummary struct {
	Type        BudgetType             `json:"type"`        /* Monthly or Travel */
	Month       int                    `json:"month"`       /* Month of the Budget */
	Year        int                    `json:"year"`        /* Year of the Budget */
	Description string                 `json:"description"` /* Optional: Descriptive text*/
	Accounts    []BudgetSummaryAccount `json:"accounts"`
}

type BudgetSummaryAccount struct {
	AccountHeader
	TotalExpected  float32          `json:"totalExpected"`  /* Expected expenses total */
	CurrentBalance float32          `json:"currentBalance"` /* Account balance (calculated from movements) */
	Status         []SummaryDetails `json:"status"`
}

type SummaryDetails struct {
	Category RecordCategory `json:"category"` /* Bill, Fuel, Groceries, Entertainment */
	Expected float32        `json:"expected"` /* The expected expenses for this type */
	Current  float32        `json:"current"`  /* Current expenses for this type  */
}

type BudgetSummaryRetrieve interface {
	GetBudgetSummaryByDate(ctx context.Context, month int, year int) (*BudgetSummary, error)
}
