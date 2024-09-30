package account

import (
	"banana-back/domain"
	"time"
)

type TransactionType string

const (
	CHARGE  TransactionType = "CHARGE"
	EXPENSE TransactionType = "EXPENSE"
	CREDIT  TransactionType = "CREDIT"
	BUDGET  TransactionType = "BUDGET"
)

type Transaction struct {
	ID          int64
	ExecutedAt  time.Time
	AppliedAt   time.Time
	Description string
	Amount      float64
	Frequency   domain.Frequency
	Type        TransactionType
	BudgetId    int64
}
