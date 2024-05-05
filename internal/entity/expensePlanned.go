package entity

type ExpenseItem struct {
	Date        *Date   `bson:"dueDate,omitempty" json:"dueDate,omitempty"`         /* Optional: Due date of the bill */
	Description string  `bson:"description,omitempty" json:"description,omitempty"` /* Optional: Descriptive text*/
	Amount      float32 `bson:"amount" json:"amount"`                               /* Transaction amount (the currency is inherited from the parent account)*/
}

type ExpensePlanned struct {
	Category RecordCategory `bson:"category" json:"category"` /* Bill, Fuel, Groceries, Entertainment */
	Expenses []*ExpenseItem `bson:"expenses" json:"expenses"` /* Expected transactions in the account */
}
