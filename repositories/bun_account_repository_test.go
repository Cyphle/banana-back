package repositories

import (
	"banana-back/domain/account"
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
		want    *account.Account
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "account exists",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&ProfileEntity{
					ID:        10004,
					Username:  "johndoe",
					Email:     "johndoe@banana.fr",
					FirstName: "John",
					LastName:  "Doe",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&AccountEntity{
					ID:        1001,
					Name:      "I am an account",
					ProfileId: 10004,
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				id: 1001,
			},
			want: &account.Account{
				ID:   1001,
				Name: "I am an account",
			},
			wantErr: assert.NoError,
		},
		{
			name: "account does not exist",
			seed: func(_ *testing.T, _ bun.IDB) {},
			args: args{
				id: 1001,
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

			got, err := r.FindById(context.Background(), tt.args.id)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func (s *RepositorySuite) TestAccountRepository_FindOneByField() {
	type args struct {
		field string
		value string
	}
	tests := []*struct {
		name    string
		args    args
		seed    func(t *testing.T, client bun.IDB)
		want    *account.Account
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "account exists",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&ProfileEntity{
					ID:        10004,
					Username:  "johndoe",
					Email:     "johndoe@banana.fr",
					FirstName: "John",
					LastName:  "Doe",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&AccountEntity{
					ID:        1001,
					Name:      "I am an account",
					ProfileId: 10004,
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				field: "name",
				value: "I am an account",
			},
			want: &account.Account{
				ID:   1001,
				Name: "I am an account",
			},
			wantErr: assert.NoError,
		},
		{
			name: "account exists with case insensitive",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&ProfileEntity{
					ID:        10004,
					Username:  "johndoe",
					Email:     "johndoe@banana.fr",
					FirstName: "John",
					LastName:  "Doe",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&AccountEntity{
					ID:        1001,
					Name:      "I am an account",
					ProfileId: 10004,
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				field: "name",
				value: "I am an Account",
			},
			want: &account.Account{
				ID:   1001,
				Name: "I am an account",
			},
			wantErr: assert.NoError,
		},
		{
			name: "account does not exist",
			seed: func(_ *testing.T, _ bun.IDB) {},
			args: args{
				field: "name",
				value: "I am an account",
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
			got, err := r.FindOneByField(context.Background(), tt.args.field, tt.args.value)

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
		want    []account.Account
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "list accounts",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&ProfileEntity{
					ID:        10004,
					Username:  "johndoe",
					Email:     "johndoe@banana.fr",
					FirstName: "John",
					LastName:  "Doe",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&AccountEntity{
					ID:        1001,
					Name:      "My account",
					ProfileId: 10004,
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&AccountEntity{
					ID:        1002,
					Name:      "My other account",
					ProfileId: 10004,
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args:    args{},
			wantErr: assert.NoError,
			want: []account.Account{
				{
					ID:   1001,
					Name: "My account",
				},
				{
					ID:   1002,
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
			got, err := r.FindAll(context.Background())

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func (s *RepositorySuite) TestAccountRepository_Create() {
	type args struct {
		input *account.CreateAccountCommand
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
			seed: func(t *testing.T, client bun.IDB) {
				_, err := client.NewInsert().Model(&ProfileEntity{
					ID:        10004,
					Username:  "johndoe",
					Email:     "johndoe@banana.fr",
					FirstName: "John",
					LastName:  "Doe",
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				input: &account.CreateAccountCommand{
					Name:      "Je suis un nouveau compte",
					ProfileId: 10004,
				},
			},
			wantErr: assert.NoError,
			want: func(t *testing.T, client bun.IDB) {
				t.Helper()
				var accounts []AccountEntity
				err := client.
					NewSelect().
					Model(&accounts).
					Column("id", "name").
					Scan(context.Background())
				require.NoError(t, err)
				assert.Equal(t, AccountEntity{
					ID:   1,
					Name: "Je suis un nouveau compte",
				}, accounts[0])
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

			secondErr := r.Create(context.Background(), tt.args.input)

			tt.wantErr(t, secondErr, tt.args.input)
			tt.want(t, trx)
		})
	}
}

func (s *RepositorySuite) TestAccountRepository_Update() {
	type args struct {
		input *account.Account
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
				_, err := client.NewInsert().Model(&ProfileEntity{
					ID:        10004,
					Username:  "johndoe",
					Email:     "johndoe@banana.fr",
					FirstName: "John",
					LastName:  "Doe",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&AccountEntity{
					ID:        10009,
					Name:      "Account 1",
					ProfileId: 10004,
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				input: &account.Account{
					ID:   10009,
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
				input: &account.Account{
					ID:   10009,
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

			tt.wantErr(t, r.Update(context.Background(), tt.args.input))

			tt.want(t, trx)
		})
	}
}

func (s *RepositorySuite) TestAccountRepository_Delete() {
	type args struct {
		id int64
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
				_, err := client.NewInsert().Model(&ProfileEntity{
					ID:        10004,
					Username:  "johndoe",
					Email:     "johndoe@banana.fr",
					FirstName: "John",
					LastName:  "Doe",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&AccountEntity{
					ID:        1001,
					Name:      "My account",
					ProfileId: 10004,
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				id: 1001,
			},
			wantErr: assert.NoError,
			want: func(t *testing.T, client bun.IDB) {
				t.Helper()
				var account AccountEntity
				err := client.NewSelect().Model(&account).
					Column("id", "name").
					Where("id = ?", 1001).
					Scan(context.Background())
				require.Error(t, err)
			},
		},
		{
			name: "account does not exist",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&ProfileEntity{
					ID:        10004,
					Username:  "johndoe",
					Email:     "johndoe@banana.fr",
					FirstName: "John",
					LastName:  "Doe",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&AccountEntity{
					ID:        1001,
					Name:      "My account",
					ProfileId: 10004,
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				id: 1002,
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
					Where("id = ?", 1001).
					Scan(context.Background())
				require.NoError(t, err)
				assert.Equal(t, AccountEntity{
					ID:   1001,
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
