package entity

type ExpenseExpectedItem struct {
	Date        *Date   `bson:"dueDate,omitempty" json:"dueDate,omitempty"` /* Optional: Due date of the bill */
	Description string  `bson:"description" json:"description"`             /* Optional: Descriptive text*/
	Amount      float32 `bson:"amount" json:"amount"`                       /* Transaction amount (the currency is inherited from the parent account)*/
}

type ExpenseExpected struct {
	Type     ExpenseType            `bson:"type" json:"type"`         /* Bill, Fuel, Groceries, Entertainment */
	Expenses []*ExpenseExpectedItem `bson:"expenses" json:"expenses"` /* Expected transactions in the account */
}
