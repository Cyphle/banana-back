package repositories

import (
	"banana-back/domain/profile"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"strings"
)

type BunProfileRepository struct {
	dbClient bun.IDB
}

var (
	ErrProfileNotFound = errors.New("profile not found")
)

func NewProfileRepository(dbClient bun.IDB) *BunProfileRepository {
	return &BunProfileRepository{
		dbClient: dbClient,
	}
}

func (r *BunProfileRepository) FindBy(ctx context.Context, username string) (*profile.Profile, error) {
	var profileEntity ProfileEntity
	err := r.dbClient.
		NewSelect().
		Column("id", "username", "email", "first_name", "last_name").
		Model(&profileEntity).
		Where("UPPER(username) = ?", strings.ToUpper(username)).
		Scan(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrProfileNotFound
	case err != nil:
		return nil, fmt.Errorf("failed to query account: %w", err)
	default:
		return &profile.Profile{
			ID:        profileEntity.ID,
			Username:  profileEntity.Username,
			Email:     profileEntity.Email,
			FirstName: profileEntity.FirstName,
			LastName:  profileEntity.LastName,
		}, nil
	}
}

func (r *BunProfileRepository) Create(ctx context.Context, command *profile.CreateProfileCommand) error {
	err := r.dbClient.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.
			NewInsert().
			Model(command).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to create profile: %w", err)
		}
		return nil
	})
	return err
}
