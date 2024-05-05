package repository

import (
	"context"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func mapToBudgetDTO(b *entity.Budget) *BudgetDTO {

	id, _ := primitive.ObjectIDFromHex(b.ID)
	BudgetDTO := &BudgetDTO{
		ID:          id,
		Type:        b.Type,
		Month:       b.Month,
		Year:        b.Year,
		Description: b.Description,
	}
	for _, account := range b.Accounts {
		BudgetDTO.Accounts = append(BudgetDTO.Accounts, mapToBudgetAccountDTO(account))
	}
	return BudgetDTO
}

func mapToBudgetAccountDTO(b *entity.BudgetAccount) *BudgetAccountDTO {

	id, _ := primitive.ObjectIDFromHex(b.ID)
	return &BudgetAccountDTO{
		AccountID: id,
		Planned:   b.Planned,
	}
}

func mapToBudget(b *BudgetDTO, accountRepo entity.AccountRepository) (*entity.Budget, error) {

	budget := &entity.Budget{
		ID:          b.ID.Hex(),
		Type:        b.Type,
		Month:       b.Month,
		Year:        b.Year,
		Description: b.Description,
	}
	for _, account := range b.Accounts {
		budgetAccount, err := mapToBudgetAccount(account, accountRepo)
		if err != nil {
			return nil, err
		}

		budget.Accounts = append(budget.Accounts, budgetAccount)
	}
	return budget, nil
}

func mapToBudgetAccount(b *BudgetAccountDTO, accountRepo entity.AccountRepository) (*entity.BudgetAccount, error) {

	account, err := accountRepo.GetByID(context.TODO(), b.AccountID.Hex())
	if err != nil {
		return nil, err
	}
	return &entity.BudgetAccount{
		Account: *account,
		Planned: b.Planned,
	}, nil
}

func mapToBudgetInfo(b *BudgetDTO) *entity.BudgetInfo {

	return &entity.BudgetInfo{
		ID:          b.ID.Hex(),
		Description: b.Description,
	}
}
