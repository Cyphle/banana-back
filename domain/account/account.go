package account

import (
	"fmt"
	"github.com/uptrace/bun"
	"google.golang.org/genproto/googleapis/type/date"
)

type AccountType string

const (
	PERSONAL AccountType = "PERSONAL"
	SHARED   AccountType = "SHARED"
	SAVINGS  AccountType = "SAVINGS"
)

type Account struct {
	ID              int64
	Name            string
	Type            AccountType
	CreationDate    date.Date
	StartingBalance float64
	Budgets         []Budget
	Transactions    []Transaction
}

type CreateAccountCommand struct {
	bun.BaseModel   `bun:"table:profiles"`
	Name            string      `json:"name"`
	Type            AccountType `json:"type"`
	StartingBalance float64     `json:"starting_balance"`
}

func CreateAccount(command *CreateAccountCommand, existingAccount *Account) (*CreateAccountCommand, error) {
	if existingAccount != nil {
		return nil, fmt.Errorf("name already exists")
	}

	return command, nil
}
