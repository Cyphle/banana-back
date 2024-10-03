package repositories

import (
	"banana-back/domain/profile"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"testing"
)

func (s *RepositorySuite) TestProfileRepository_FindByUsername() {
	type args struct {
		username string
	}
	tests := []*struct {
		name    string
		args    args
		seed    func(t *testing.T, client bun.IDB)
		want    *profile.Profile
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "profile exists",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&ProfileEntity{
					ID:       1,
					Username: "johndoe",
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				username: "johndoe",
			},
			want: &profile.Profile{
				ID:       1,
				Username: "johndoe",
			},
			wantErr: assert.NoError,
		},
		{
			name: "profile does not exist",
			seed: func(_ *testing.T, _ bun.IDB) {},
			args: args{
				username: "johndoe",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrProfileNotFound)
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
			r := NewProfileRepository(trx)

			got, err := r.FindBy(context.Background(), tt.args.username)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func (s *RepositorySuite) TestProfileRepository_Create() {
	type args struct {
		input *profile.CreateProfileCommand
	}
	tests := []*struct {
		name    string
		seed    func(t *testing.T, client bun.IDB)
		args    args
		wantErr assert.ErrorAssertionFunc
		want    func(t *testing.T, client bun.IDB)
	}{
		{
			name: "create profile",
			seed: func(_ *testing.T, _ bun.IDB) {},
			args: args{
				input: &profile.CreateProfileCommand{
					Username: "johnsmith",
				},
			},
			wantErr: assert.NoError,
			want: func(t *testing.T, client bun.IDB) {
				t.Helper()
				var profiles []ProfileEntity
				err := client.
					NewSelect().
					Model(&profiles).
					Column("id", "username").
					Scan(context.Background())
				require.NoError(t, err)
				assert.Equal(t, ProfileEntity{
					ID:       1,
					Username: "johnsmith",
				}, profiles[0])
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
			r := NewProfileRepository(trx)

			secondErr := r.Create(context.Background(), tt.args.input)

			tt.wantErr(t, secondErr, tt.args.input)
			tt.want(t, trx)
		})
	}
}
