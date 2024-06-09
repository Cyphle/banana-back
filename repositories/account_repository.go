package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
)

type AccountRepository struct {
	dbClient bun.IDB
}

func NewAccountRepository(dbClient bun.IDB) *AccountRepository {
	return &AccountRepository{
		dbClient: dbClient,
	}
}

var (
	ErrAccountNotFound = errors.New("account not found")
)

func (r *AccountRepository) GetByID(ctx context.Context, id int64) (*AccountEntity, error) {
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
		return &accountEntity, nil
	}
}

func (r *AccountRepository) List(
	ctx context.Context,
) ([]AccountRepository, error) {
	var accounts []AccountRepository
	query := r.
		dbClient.
		NewSelect().
		Column("id", "name").
		Model(&accounts)

	if err := query.Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}

	return accounts, nil
}

func (r *AccountRepository) Create(ctx context.Context, input *AccountEntityCreateParams) error {
	if _, err := r.dbClient.
		NewInsert().
		Model(input).
		Exec(ctx); err != nil {
		return fmt.Errorf("failed to create stakeholder: %w", err)
	}
	return nil
}

func (r *AccountRepository) Update(
	ctx context.Context,
	id int,
	input *AccountEntityUpdateParams,
) error {
	res, err := r.
		dbClient.
		NewUpdate().
		Model(new(AccountEntity)).
		Set("name = ?", input.Name).
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

func (r *AccountRepository) Delete(ctx context.Context, id int) error {
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