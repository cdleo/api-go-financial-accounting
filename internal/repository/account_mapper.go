package repository

import (
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

	Account := &entity.Account{
		ID:            b.ID.Hex(),
		ParentID:      b.ParentID.Hex(),
		AccountHeader: b.AccountHeader,
		Balance:       b.Balance,
		Movements:     b.Movements,
	}
	return Account
}

func mapToAccountSummary(b *AccountDTO) *entity.AccountSummary {

	return &entity.AccountSummary{
		ID:            b.ID.Hex(),
		AccountHeader: b.AccountHeader,
		Balance:       b.Balance,
	}
}
