package entity

import (
	"fmt"
	"strings"
	"time"
)

const (
	Undefined int = -1
	DDMMYYYY      = "02/01/2006"
)

type Date struct {
	time.Time
}

func (ct *Date) UnmarshalJSON(data []byte) (err error) {
	s := strings.Trim(string(data), "\"")
	if s == "null" || s == "" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(DDMMYYYY, s)
	return err
}

func (ct *Date) MarshalJSON() ([]byte, error) {

	if ct.IsZero() {
		return []byte("\"\""), nil
	}
	date := []byte(fmt.Sprintf("%q", ct.Format(DDMMYYYY)))
	return date, nil
}

type TransactionType string

const (
	Credit TransactionType = "Credit"
	Debit  TransactionType = "Debit"
)

func (t TransactionType) Modifier() float32 {
	if t == Credit {
		return 1
	} else {
		return -1
	}
}

var TransactionTypes = []TransactionType{
	Credit,
	Debit,
}

type ExpenseType string

const (
	Bill          ExpenseType = "Bill"
	Fuel          ExpenseType = "Fuel"
	Goods         ExpenseType = "Goods"
	Groceries     ExpenseType = "Groceries"
	Entertainment ExpenseType = "Entertainment"
	Investment    ExpenseType = "Investment"
	Income        ExpenseType = "Income"
)

var ExpenseTypes = []ExpenseType{
	Bill,
	Fuel,
	Groceries,
	Entertainment,
	Investment,
	Income,
}
