package repositories

import (
	"banana-back/domain/account"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"strings"
)

type BunAccountRepository struct {
	dbClient bun.IDB
}

func NewAccountRepository(dbClient bun.IDB) *BunAccountRepository {
	return &BunAccountRepository{
		dbClient: dbClient,
	}
}

var (
	ErrAccountNotFound = errors.New("account not found")
)

func (r *BunAccountRepository) FindById(ctx context.Context, id int64) (*account.Account, error) {
	var accountEntity AccountEntity
	err := r.dbClient.
		NewSelect().
		Column("id", "name").
		Model(&accountEntity).
		Where("id = ?", id).
		Scan(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrAccountNotFound
	case err != nil:
		return nil, fmt.Errorf("failed to query account: %w", err)
	default:
		return &account.Account{
			ID:   accountEntity.ID,
			Name: accountEntity.Name,
		}, nil
	}
}

func (r *BunAccountRepository) FindOneByField(ctx context.Context, field string, value string) (*account.Account, error) {
	var accountEntity AccountEntity
	err := r.dbClient.
		NewSelect().
		Column("id", "name").
		Model(&accountEntity).
		Where("UPPER("+field+") = ?", strings.ToUpper(value)).
		Scan(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrAccountNotFound
	case err != nil:
		return nil, fmt.Errorf("failed to query account: %w", err)
	default:
		return &account.Account{
			ID:   accountEntity.ID,
			Name: accountEntity.Name,
		}, nil
	}
}

func (r *BunAccountRepository) FindAll(
	ctx context.Context,
) ([]account.Account, error) {
	var accountEntities []AccountEntity
	query := r.
		dbClient.
		NewSelect().
		Column("id", "name").
		Model(&accountEntities)

	if err := query.Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}

	accounts := make([]account.Account, 0, len(accountEntities))
	for _, accountEntity := range accountEntities {
		accounts = append(accounts, account.Account{
			ID:   accountEntity.ID,
			Name: accountEntity.Name,
		})
	}

	return accounts, nil
}

func (r *BunAccountRepository) Create(ctx context.Context, input *account.CreateAccountCommand) error {
	err := r.dbClient.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.
			NewInsert().
			Model(input).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}
		return nil
	})
	return err
}

func (r *BunAccountRepository) Update(
	ctx context.Context,
	input *account.Account,
) error {
	res, err := r.
		dbClient.
		NewUpdate().
		Model(new(AccountEntity)).
		Set("name = ?", input.Name).
		Where("id = ?", input.ID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}
	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get updated rows: %w", err)
	}
	if rowsUpdated == 0 {
		return ErrAccountNotFound
	}
	return nil
}

func (r *BunAccountRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.
		dbClient.
		NewDelete().
		Model(new(AccountEntity)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}
	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get updated rows: %w", err)
	}
	if rowsUpdated == 0 {
		return ErrAccountNotFound
	}
	return nil
}
