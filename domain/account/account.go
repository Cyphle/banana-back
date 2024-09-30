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
}

func CreateAccount(command *CreateAccountCommand, existingAccount *Account) (*Account, error) {
	if existingAccount != nil {
		return nil, fmt.Errorf("name already exists")
	}

	return &Account{Name: command.Name}, nil
}

func AddLineToAccount(account *Account, command *AddAccountLineCommand) error {
	account.Lines = append(account.Lines, AccountLine{
		Type:            command.Type,
		EventDate:       command.EventDate,
		ApplicationDate: command.ApplicationDate,
		Description:     command.Description,
		Amount:          command.Amount,
	})

	return nil
}

func UpdateApplicationDateOfLine(account *Account, lineId int64, date *date.Date) error {
	for i, line := range account.Lines {
		if line.ID == lineId {
			account.Lines[i].ApplicationDate = date
			break
		}
	}

	return nil
}
