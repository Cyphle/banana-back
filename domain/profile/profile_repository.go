package profile

import "context"

type ProfileRepository interface {
	FindBy(ctx context.Context, username string) (*Profile, error)
	Create(ctx context.Context, command *CreateProfileCommand) error
}
