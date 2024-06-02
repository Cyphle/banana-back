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
		//{
		//	name: "stakeholder does not exist",
		//	seed: func(_ *testing.T, _ bun.IDB) {},
		//	args: args{
		//		uid:   uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
		//		orgID: 1,
		//	},
		//	want: nil,
		//	wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
		//		return assert.ErrorIs(t, err, domain.ErrStakeholderNotFound)
		//	},
		//},
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

/*
func (s *RepositorySuite) TestStakeholderRepository_List() {
	type args struct {
		params *domain.StakeholderListParams
	}
	tests := []*struct {
		name    string
		seed    func(t *testing.T, client bun.IDB)
		args    args
		want    []domain.Stakeholder
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "list stakeholders by organization",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&domain.Organization{
					ID:   10005,
					Code: "org1",
					Name: "Organization 1",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&domain.Organization{
					ID:   10006,
					Code: "org2",
					Name: "Organization 2",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&domain.Stakeholder{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
					OrganizationID: 10005,
					Properties:     []byte(`{"name": "John Doe"}`),
					Roles:          []string{"admin"},
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&domain.Stakeholder{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b9"),
					OrganizationID: 10006,
					Properties:     []byte(`{"name": "Jane Doe"}`),
					Roles:          []string{"user"},
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&domain.Stakeholder{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b0"),
					OrganizationID: 10005,
					Properties:     []byte(`{"name": "Jack Doe"}`),
					Roles:          []string{"user"},
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				params: &domain.StakeholderListParams{
					OrganizationID: 10005,
				},
			},
			wantErr: assert.NoError,
			want: []domain.Stakeholder{
				{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
					OrganizationID: 10005,
					Properties:     []byte(`{"name": "John Doe"}`),
					Roles:          []string{"admin"},
				},
				{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b0"),
					OrganizationID: 10005,
					Properties:     []byte(`{"name": "Jack Doe"}`),
					Roles:          []string{"user"},
				},
			},
		},
		{
			name: "list stakeholders by organization and freesearch",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&domain.Organization{
					ID:   10007,
					Code: "org3",
					Name: "Organization 3",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&domain.Organization{
					ID:   10008,
					Code: "org4",
					Name: "Organization 4",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&domain.Stakeholder{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
					OrganizationID: 10007,
					Properties:     []byte(`{"name": "John Doe"}`),
					Roles:          []string{"admin"},
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&domain.Stakeholder{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b9"),
					OrganizationID: 10008,
					Properties:     []byte(`{"name": "Jane Doe"}`),
					Roles:          []string{"user"},
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&domain.Stakeholder{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b0"),
					OrganizationID: 10007,
					Properties:     []byte(`{"name": "Jack Doe"}`),
					Roles:          []string{"user"},
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				params: &domain.StakeholderListParams{
					OrganizationID: 10007,
					Freesearch:     "Jack",
				},
			},
			wantErr: assert.NoError,
			want: []domain.Stakeholder{
				{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b0"),
					OrganizationID: 10007,
					Properties:     []byte(`{"name": "Jack Doe"}`),
					Roles:          []string{"user"},
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
			r := NewStakeholderRepository(trx)
			got, err := r.List(context.Background(), tt.args.params)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func (s *RepositorySuite) TestStakeholderRepository_Create() {
	type args struct {
		input *domain.StakeholderCreateParams
	}
	tests := []*struct {
		name    string
		seed    func(t *testing.T, client bun.IDB)
		args    args
		wantErr assert.ErrorAssertionFunc
		want    func(t *testing.T, client bun.IDB)
	}{
		{
			name: "create stakeholder",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&domain.Organization{
					ID:   10008,
					Code: "org1",
					Name: "Organization 1",
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				input: &domain.StakeholderCreateParams{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
					OrganizationID: 10008,
					Properties:     []byte(`{"name": "John Doe"}`),
					Roles:          []string{"admin"},
				},
			},
			wantErr: assert.NoError,
			want: func(t *testing.T, client bun.IDB) {
				t.Helper()
				var stakeholder domain.Stakeholder
				err := client.NewSelect().Model(&stakeholder).
					Column("uid", "organization_id", "properties", "roles").
					Where("uid = ?", uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8")).
					Scan(context.Background())
				require.NoError(t, err)
				assert.Equal(t, domain.Stakeholder{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
					OrganizationID: 10008,
					Properties:     []byte(`{"name": "John Doe"}`),
					Roles:          []string{"admin"},
				}, stakeholder)
			},
		},
		{
			name: "organization does not exist",
			seed: func(_ *testing.T, _ bun.IDB) {},
			args: args{
				input: &domain.StakeholderCreateParams{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
					OrganizationID: 1,
					Properties:     []byte(`{"name": "John Doe"}`),
					Roles:          []string{"admin"},
				},
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				var pgErr *pgconn.PgError
				return assert.ErrorAs(t, err, &pgErr) && pgErr.Code == "23503"
			},
			want: func(t *testing.T, client bun.IDB) {
				t.Helper()
				var stakeholder domain.Stakeholder
				err := client.NewSelect().Model(&stakeholder).
					Column("id", "uid", "organization_id", "properties", "roles").
					Where("uid = ?", uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8")).
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
			r := NewStakeholderRepository(trx)

			tt.wantErr(t, r.Create(context.Background(), tt.args.input))
			tt.want(t, trx)
		})
	}
}

func (s *RepositorySuite) TestStakeholderRepository_Update() {
	type args struct {
		uid   uuid.UUID
		orgID int
		input *domain.StakeholderUpdateParams
	}
	tests := []*struct {
		name    string
		seed    func(t *testing.T, client bun.IDB)
		args    args
		wantErr assert.ErrorAssertionFunc
		want    func(t *testing.T, client bun.IDB)
	}{
		{
			name: "update stakeholder",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&domain.Organization{
					ID:   10009,
					Code: "org1",
					Name: "Organization 1",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&domain.Stakeholder{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
					OrganizationID: 10009,
					Properties:     []byte(`{"name": "John Doe"}`),
					Roles:          []string{"admin"},
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				uid:   uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
				orgID: 10009,
				input: &domain.StakeholderUpdateParams{
					Properties: []byte(`{"name": "Jane Doe"}`),
					Roles:      []string{"user"},
				},
			},
			wantErr: assert.NoError,
			want: func(t *testing.T, client bun.IDB) {
				t.Helper()
				var stakeholder domain.Stakeholder
				err := client.NewSelect().Model(&stakeholder).
					Column("uid", "organization_id", "properties", "roles").
					Where("uid = ?", uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8")).
					Scan(context.Background())
				require.NoError(t, err)
				assert.Equal(t, domain.Stakeholder{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
					OrganizationID: 10009,
					Properties:     []byte(`{"name": "Jane Doe"}`),
					Roles:          []string{"user"},
				}, stakeholder)
			},
		},
		{
			name: "stakeholder does not exist",
			seed: func(_ *testing.T, _ bun.IDB) {},
			args: args{
				uid:   uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
				orgID: 10009,
				input: &domain.StakeholderUpdateParams{
					Properties: []byte(`{"name": "Jane Doe"}`),
					Roles:      []string{"user"},
				},
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, domain.ErrStakeholderNotFound)
			},
			want: func(t *testing.T, client bun.IDB) {
				t.Helper()
				var stakeholder domain.Stakeholder
				err := client.NewSelect().Model(&stakeholder).
					Column("id", "uid", "organization_id", "properties", "roles").
					Where("uid = ?", uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8")).
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
			r := NewStakeholderRepository(trx)

			tt.wantErr(t, r.Update(context.Background(), tt.args.uid, tt.args.orgID, tt.args.input))
			tt.want(t, trx)
		})
	}
}

func (s *RepositorySuite) TestStakeholderRepository_Delete() {
	type args struct {
		uid   uuid.UUID
		orgID int
	}
	tests := []*struct {
		name    string
		seed    func(t *testing.T, client bun.IDB)
		args    args
		wantErr assert.ErrorAssertionFunc
		want    func(t *testing.T, client bun.IDB)
	}{
		{
			name: "delete stakeholder",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&domain.Organization{
					ID:   10010,
					Code: "org1",
					Name: "Organization 1",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&domain.Stakeholder{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
					OrganizationID: 10010,
					Properties:     []byte(`{"name": "John Doe"}`),
					Roles:          []string{"admin"},
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				uid:   uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
				orgID: 10010,
			},
			wantErr: assert.NoError,
			want: func(t *testing.T, client bun.IDB) {
				t.Helper()
				var stakeholder domain.Stakeholder
				err := client.NewSelect().Model(&stakeholder).
					Column("uid", "organization_id", "properties", "roles").
					Where("uid = ?", uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8")).
					Scan(context.Background())
				require.Error(t, err)
			},
		},
		{
			name: "stakeholder does not exist",
			seed: func(t *testing.T, client bun.IDB) {
				t.Helper()
				_, err := client.NewInsert().Model(&domain.Organization{
					ID:   10011,
					Code: "org1",
					Name: "Organization 1",
				}).Exec(context.Background())
				require.NoError(t, err)
				_, err = client.NewInsert().Model(&domain.Stakeholder{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
					OrganizationID: 10011,
					Properties:     []byte(`{"name": "John Doe"}`),
					Roles:          []string{"admin"},
				}).Exec(context.Background())
				require.NoError(t, err)
			},
			args: args{
				uid:   uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b9"),
				orgID: 10011,
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, domain.ErrStakeholderNotFound)
			},
			want: func(t *testing.T, client bun.IDB) {
				t.Helper()
				var stakeholder domain.Stakeholder
				err := client.NewSelect().Model(&stakeholder).
					Column("uid", "organization_id", "properties", "roles").
					Where("uid = ?", uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8")).
					Scan(context.Background())
				require.NoError(t, err)
				assert.Equal(t, domain.Stakeholder{
					UID:            uuid.MustParse("1ee25aa2-6165-4f16-ad35-4d5810e4a7b8"),
					OrganizationID: 10011,
					Properties:     []byte(`{"name": "John Doe"}`),
					Roles:          []string{"admin"},
				}, stakeholder)
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
			r := NewStakeholderRepository(trx)

			tt.wantErr(t, r.Delete(context.Background(), tt.args.uid, tt.args.orgID))
			tt.want(t, trx)
		})
	}
}
*/
