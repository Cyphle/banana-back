package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"time"
)

type AccountRepository struct {
	dbClient bun.IDB
}

type AccountEntity struct {
	bun.BaseModel `bun:"table:accounts"`
	ID            int64      `bun:"id,pk,autoincrement"       json:"-"`
	Name          string     `json:"name,omitempty"`
	CreatedAt     time.Time  `json:"-"`
	UpdatedAt     *time.Time `json:"-"`
	DeletedAt     *time.Time `bun:",soft_delete"              json:"-"`
}

func NewAccountRepository(dbClient bun.IDB) *AccountRepository {
	return &AccountRepository{
		dbClient: dbClient,
	}
}

var (
	AccountNotFound = errors.New("account not found")
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
		return nil, AccountNotFound
	case err != nil:
		return nil, fmt.Errorf("failed to query account: %w", err)
	default:
		return &accountEntity, nil
	}
}

/*

// Create creates a new stakeholder.
func (r *StakeholderRepository) Create(ctx context.Context, input *domain.StakeholderCreateParams) error {
	if _, err := r.client.NewInsert().Model(input).Exec(ctx); err != nil {
		return fmt.Errorf("failed to create stakeholder: %w", err)
	}
	return nil
}

// Update updates a stakeholder.
func (r *StakeholderRepository) Update(
	ctx context.Context,
	uid uuid.UUID,
	orgID int,
	input *domain.StakeholderUpdateParams,
) error {
	res, err := r.client.NewUpdate().
		Model(new(domain.Stakeholder)).
		Set("roles = ?", pgdialect.Array(input.Roles)).
		Set("properties = ?", input.Properties).
		Where("uid = ? and organization_id = ?", uid, orgID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update stakeholder: %w", err)
	}
	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get updated rows: %w", err)
	}
	if rowsUpdated == 0 {
		return domain.ErrStakeholderNotFound
	}
	return nil
}

// Delete deletes a stakeholder.
func (r *StakeholderRepository) Delete(ctx context.Context, uid uuid.UUID, orgID int) error {
	res, err := r.client.NewDelete().
		Model(new(domain.Stakeholder)).
		Where("uid = ? and organization_id = ?", uid, orgID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update stakeholder: %w", err)
	}
	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get updated rows: %w", err)
	}
	if rowsUpdated == 0 {
		return domain.ErrStakeholderNotFound
	}
	return nil
}

// List returns all stakeholders.
func (r *StakeholderRepository) List(
	ctx context.Context,
	params *domain.StakeholderListParams,
) ([]domain.Stakeholder, error) {
	var stakeholders []domain.Stakeholder
	query := r.client.NewSelect().
		Column("uid", "properties", "roles", "organization_id").
		Model(&stakeholders).
		Where("organization_id = ?", params.OrganizationID)

	if params.Freesearch != "" {
		query.Where("freesearch ~* ?", strings.ReplaceAll(params.Freesearch, " ", "|"))
	}

	if err := query.Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to query stakeholders: %w", err)
	}

	return stakeholders, nil
}
*/
