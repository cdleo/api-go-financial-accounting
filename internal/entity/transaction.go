package entity

type TrxDetails struct {
	Date        Date            `bson:"date" json:"date"`                           /* Date of the record */
	Kind        TransactionType `bson:"kind" json:"kind"`                           /* Credit or Debit */
	SubType     ExpenseType     `bson:"subtype,omitempty" json:"subtype,omitempty"` /* Optional (If Debit): Bill, Fuel, Groceries, Entertainment, Goods */
	Description string          `bson:"description" json:"description"`             /* Optional: Descriptive text*/
	Amount      float32         `bson:"amount" json:"amount"`                       /* Transaction amount (the currency is inherited from the parent account)*/
}

type Transaction struct {
	AccountID string     `json:"account_id"` /* Account identifier */
	Details   TrxDetails `json:"details"`    /* Transaction details */
}
