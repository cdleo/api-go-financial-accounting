package entity

type CurrencyCode string

const (
	ARS CurrencyCode = "ARS"
	USD CurrencyCode = "USD"
	EUR CurrencyCode = "EUR"
	BRL CurrencyCode = "BRL"
)

var CurrencyCodes = []CurrencyCode{
	ARS,
	USD,
	EUR,
	BRL,
}
