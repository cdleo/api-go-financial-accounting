package entity

type AccountType string

const (
	Bank      AccountType = "Bank"
	Cash      AccountType = "Cash"
	Asset     AccountType = "Asset"
	Liability AccountType = "Liability"
	//Investment AccountType = "Investment"
	//Income     AccountType = "Income"
	Trading    AccountType = "Trading"
	Equity     AccountType = "Equity"
	Receivable AccountType = "Receivable"
	Payable    AccountType = "Payable"
)

var AccountTypes = []AccountType{
	Bank,
	Cash,
	Asset,
	Liability,
	//Investment,
	//Income,
	Trading,
	Equity,
	Receivable,
	Payable,
}
