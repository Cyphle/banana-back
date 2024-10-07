package account

import "context"

type AccountRepository interface {
	FindAll(ctx context.Context) ([]Account, error)
	FindById(ctx context.Context, id int64) (*Account, error)
	FindOneByField(ctx context.Context, field string, value string) (*Account, error)
	Create(ctx context.Context, input *CreateAccountCommand) error
	Update(ctx context.Context, input *Account) error
	Delete(ctx context.Context, id int64) error
}
