package repositories

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"testing"
)

func (s *RepositorySuite) TestAccountRepository_GetByID() {
	type args struct {
		id int64
	}
	tests := []*struct {
		name    string
		args    args
		seed    func(t *testing.T, client bun.IDB)
		want    *AccountEntity
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "account exists",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&AccountEntity{
					Name: "I am an account",
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				id: 1,
			},
			want: &AccountEntity{
				ID:   1,
				Name: "I am an account",
			},
			wantErr: assert.NoError,
		},
		{
			name: "account does not exist",
			seed: func(_ *testing.T, _ bun.IDB) {},
			args: args{
				id: 1,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrAccountNotFound)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			trx, err := s.client.Begin()
			s.Require().NoError(err)
			defer func() {
				s.Require().NoError(trx.Rollback())
			}()
			tt.seed(t, trx)
			r := NewAccountRepository(trx)
			got, err := r.GetByID(context.Background(), tt.args.id)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func (s *RepositorySuite) TestAccountRepository_List() {
	type args struct {
	}
	tests := []*struct {
		name    string
		seed    func(t *testing.T, client bun.IDB)
		args    args
		want    []AccountEntity
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "list accounts",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&AccountEntity{
					ID:   1,
					Name: "My account",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&AccountEntity{
					ID:   2,
					Name: "My other account",
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args:    args{},
			wantErr: assert.NoError,
			want: []AccountEntity{
				{
					ID:   1,
					Name: "My account",
				},
				{
					ID:   2,
					Name: "My other account",
				},
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			trx, err := s.client.Begin()
			s.Require().NoError(err)
			defer func() {
				s.Require().NoError(trx.Rollback())
			}()
			tt.seed(t, trx)
			r := NewAccountRepository(trx)
			got, err := r.List(context.Background())

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func (s *RepositorySuite) TestAccountRepository_Create() {
	type args struct {
		input *AccountEntityCreateParams
	}
	tests := []*struct {
		name    string
		seed    func(t *testing.T, client bun.IDB)
		args    args
		wantErr assert.ErrorAssertionFunc
		want    func(t *testing.T, client bun.IDB)
	}{
		{
			name: "create account",
			seed: func(_ *testing.T, _ bun.IDB) {},
			args: args{
				input: &AccountEntityCreateParams{
					Name: "Je suis un nouveau compte",
				},
			},
			wantErr: assert.NoError,
			want: func(t *testing.T, client bun.IDB) {
				t.Helper()
				var account AccountEntity
				err := client.
					NewSelect().
					Model(&account).
					Column("id", "name").
					Where("id = ?", 1).
					Scan(context.Background())
				require.NoError(t, err)
				assert.Equal(t, AccountEntity{
					ID:   1,
					Name: "Je suis un nouveau compte",
				}, account)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			trx, err := s.client.Begin()
			s.Require().NoError(err)
			defer func() {
				s.Require().NoError(trx.Rollback())
			}()
			tt.seed(t, trx)
			r := NewAccountRepository(trx)

			tt.wantErr(t, r.Create(context.Background(), tt.args.input))
			tt.want(t, trx)
		})
	}
}

func (s *RepositorySuite) TestAccountRepository_Update() {
	type args struct {
		id    int
		input *AccountEntityUpdateParams
	}
	tests := []*struct {
		name    string
		seed    func(t *testing.T, client bun.IDB)
		args    args
		wantErr assert.ErrorAssertionFunc
		want    func(t *testing.T, client bun.IDB)
	}{
		{
			name: "update account",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&AccountEntity{
					ID:   10009,
					Name: "Account 1",
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				id: 10009,
				input: &AccountEntityUpdateParams{
					Name: "New name",
				},
			},
			wantErr: assert.NoError,
			want: func(t *testing.T, client bun.IDB) {
				t.Helper()
				var account AccountEntity
				err := client.NewSelect().Model(&account).
					Column("id", "name").
					Where("id = ?", 10009).
					Scan(context.Background())
				require.NoError(t, err)
				assert.Equal(t, AccountEntity{
					ID:   10009,
					Name: "New name",
				}, account)
			},
		},
		{
			name: "account does not exist",
			seed: func(_ *testing.T, _ bun.IDB) {},
			args: args{
				id: 10009,
				input: &AccountEntityUpdateParams{
					Name: "New name",
				},
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrAccountNotFound)
			},
			want: func(t *testing.T, client bun.IDB) {
				t.Helper()
				var account AccountEntity
				err := client.NewSelect().Model(&account).
					Column("id", "name").
					Where("id = ?", 10009).
					Scan(context.Background())
				require.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			trx, err := s.client.Begin()
			s.Require().NoError(err)
			defer func() {
				s.Require().NoError(trx.Rollback())
			}()
			tt.seed(t, trx)
			r := NewAccountRepository(trx)

			tt.wantErr(t, r.Update(context.Background(), tt.args.id, tt.args.input))
			tt.want(t, trx)
		})
	}
}

func (s *RepositorySuite) TestStakeholderRepository_Delete() {
	type args struct {
		id int
	}
	tests := []*struct {
		name    string
		seed    func(t *testing.T, client bun.IDB)
		args    args
		wantErr assert.ErrorAssertionFunc
		want    func(t *testing.T, client bun.IDB)
	}{
		{
			name: "delete account",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&AccountEntity{
					ID:   1,
					Name: "My account",
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				id: 1,
			},
			wantErr: assert.NoError,
			want: func(t *testing.T, client bun.IDB) {
				t.Helper()
				var account AccountEntity
				err := client.NewSelect().Model(&account).
					Column("id", "name").
					Where("id = ?", 1).
					Scan(context.Background())
				require.Error(t, err)
			},
		},
		{
			name: "account does not exist",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&AccountEntity{
					ID:   1,
					Name: "My account",
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				id: 2,
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrAccountNotFound)
			},
			want: func(t *testing.T, client bun.IDB) {
				t.Helper()
				var account AccountEntity
				err := client.NewSelect().
					Model(&account).
					Column("id", "name").
					Where("id = ?", 1).
					Scan(context.Background())
				require.NoError(t, err)
				assert.Equal(t, AccountEntity{
					ID:   1,
					Name: "My account",
				}, account)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			trx, err := s.client.Begin()
			s.Require().NoError(err)
			defer func() {
				s.Require().NoError(trx.Rollback())
			}()
			tt.seed(t, trx)
			r := NewAccountRepository(trx)

			tt.wantErr(t, r.Delete(context.Background(), tt.args.id))
			tt.want(t, trx)
		})
	}
}
