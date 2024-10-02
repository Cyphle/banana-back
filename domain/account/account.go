package account

import (
	"fmt"
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
	Name string
	Type AccountType
}

func CreateAccount(command *CreateAccountCommand, existingAccount *Account) (*Account, error) {
	if existingAccount != nil {
		return nil, fmt.Errorf("name already exists")
	}

	return &Account{Name: command.Name}, nil
}
