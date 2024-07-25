package domain

import (
	"fmt"
)

type Account struct {
	ID    int64
	Name  string
	Lines []AccountLine
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
