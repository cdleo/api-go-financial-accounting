package service

import (
	"github.com/cdleo/api-go-financial-accounting/internal/entity"
)

type BudgetSummaryRetrieve interface {
	GetBudgetSummaryByDate(month int, year int) (*entity.BudgetSummary, error)
}

type budgetSummaryRetrieve struct {
	service entity.BudgetRetrieve
}

func NewBudgetSummaryRetrieve(service entity.BudgetRetrieve) BudgetSummaryRetrieve {
	return &budgetSummaryRetrieve{
		service: service,
	}
}

func (s *budgetSummaryRetrieve) GetBudgetSummaryByDate(month int, year int) (*entity.BudgetSummary, error) {

	budget, err := s.service.GetBudgetByDate(month, year)
	if budget == nil || err != nil {
		return nil, err
	}

	var summary entity.BudgetSummary

	summary.Type = budget.Type
	summary.Month = budget.Month
	summary.Year = budget.Year
	summary.Description = budget.Description

	//var accountSummary entity.BudgetSummaryAccount

	for _, account := range budget.Accounts {
		accountSummary := entity.BudgetSummaryAccount{
			AccountHeader:  account.AccountHeader,
			CurrentBalance: 0,
			TotalExpected:  0,
		}

		summaryMap := make(map[entity.ExpenseType]entity.ExpenseSummary)
		// Sumamos lo esperado
		var totalExpenses float32 = 0
		for _, expected := range account.Expected {

			var totalExpensesForType float32 = 0
			for _, expenses := range expected.Expenses {
				totalExpensesForType += expenses.Amount
			}
			totalExpenses += totalExpensesForType

			typeSummary := entity.ExpenseSummary{
				Type:     expected.Type,
				Expected: totalExpensesForType,
				Current:  0,
			}

			summaryMap[typeSummary.Type] = typeSummary
		}
		accountSummary.TotalExpected = totalExpenses

		// Sumamos lo acontecido por tipo de gasto y calculamos el balance total
		for _, movement := range account.Movements {

			if movement.Kind == entity.Credit {
				accountSummary.CurrentBalance += movement.Amount
				continue
			}

			typeSummary, ok := summaryMap[movement.SubType]
			if !ok {
				typeSummary = entity.ExpenseSummary{
					Type:     movement.SubType,
					Expected: 0,
					Current:  movement.Amount,
				}
			} else {
				typeSummary.Current += movement.Amount
			}
			accountSummary.CurrentBalance -= movement.Amount
			summaryMap[movement.SubType] = typeSummary
		}

		//Agregamos los status de todos los tipo al summary de la cuenta
		for _, v := range summaryMap {
			accountSummary.Status = append(accountSummary.Status, v)
		}
		//Agregamos la cuenta sumarizada a la respuesta
		summary.Accounts = append(summary.Accounts, accountSummary)
	}

	return &summary, nil
}
