package domain

import (
	"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/type/date"
)

type AccountLineType string

const (
	Expense AccountLineType = "Expense"
	Budget  AccountLineType = "Budget"
	Charge  AccountLineType = "Charge"
	Credit  AccountLineType = "Credit"
)

type AccountLine struct {
	ID              int64
	Type            AccountLineType
	EventDate       *date.Date
	ApplicationDate *date.Date
	Description     string
	Amount          decimal.Decimal
}

type AddAccountLineCommand struct {
	Type            AccountLineType
	EventDate       *date.Date
	ApplicationDate *date.Date
	Description     string
	Amount          decimal.Decimal
}
