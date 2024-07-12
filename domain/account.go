package domain

import (
	"context"
	"fmt"
)

type Account struct {
	ID   int64
	Name string
}

type CreateAccountCommand struct {
	Name string
}

// TODO to be tested and injected
func CreateAccount(repository Repository[Account], ctx context.Context, request *CreateAccountCommand) (*Account, error) {
	if existingAccount, err := repository.FindOneByField(ctx, "name", request.Name); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	} else {
		if existingAccount != nil {
			return nil, fmt.Errorf("name already exists")
		}
	}

	if err := repository.Create(ctx, &Account{
		Name: request.Name,
	}); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return repository.FindOneByField(ctx, "name", request.Name)
}
