package repository

import (
	"context"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func mapToAccountDTO(b *entity.Account) *AccountDTO {

	id, _ := primitive.ObjectIDFromHex(b.ID)
	parentId, _ := primitive.ObjectIDFromHex(b.ParentID)
	AccountDTO := &AccountDTO{
		ID:            id,
		ParentID:      parentId,
		AccountHeader: b.AccountHeader,
		Balance:       b.Balance,
		Movements:     b.Movements,
	}
	return AccountDTO
}

func mapToAccount(b *AccountDTO) *entity.Account {

	var parentID string = ""
	if b.ParentID != primitive.NilObjectID {
		parentID = b.ParentID.Hex()
	}

	Account := &entity.Account{
		ID:            b.ID.Hex(),
		ParentID:      parentID,
		AccountHeader: b.AccountHeader,
		Balance:       b.Balance,
		Movements:     b.Movements,
	}
	return Account
}

func mapToAccountInfo(b *AccountDTO, budgetRepo entity.BudgetRepository) *entity.AccountInfo {

	var budgetInfo *entity.BudgetInfo = nil
	if b.ParentID != primitive.NilObjectID {
		budgetInfo, _ = budgetRepo.GetInfoById(context.TODO(), b.ParentID.Hex())
	}

	return &entity.AccountInfo{
		ID:            b.ID.Hex(),
		AccountHeader: b.AccountHeader,
		Balance:       b.Balance,
		Budget:        budgetInfo,
	}
}
