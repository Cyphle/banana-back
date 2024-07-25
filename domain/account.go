package domain

import (
	"fmt"
)

type Account struct {
	ID   int64
	Name string
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
