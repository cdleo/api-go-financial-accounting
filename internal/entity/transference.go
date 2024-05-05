package entity

import "context"

type TransferDetails struct {
	Date        Date    `json:"date"`
	Description string  `json:"description"`
	Amount      float32 `json:"amount"`
	Rate        float32 `json:"conversionRate"`
	FromFee     float32 `json:"fromFee"`
	ToFee       float32 `json:"toFee"`
}

type Transfer struct {
	FromID  string          `json:"fromAccount"`
	ToID    string          `json:"toAccount"`
	Details TransferDetails `json:"details"`
}

type MakeTransference interface {
	MakeTransference(ctx context.Context, transference Transfer) error
}
