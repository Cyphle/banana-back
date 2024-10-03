package profile

import (
	"errors"
	"github.com/uptrace/bun"
)

type Profile struct {
	ID        int64
	Username  string
	Email     string
	FirstName string
	LastName  string
}

type CreateProfileCommand struct {
	bun.BaseModel `bun:"table:profiles"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
}

func ValidateProfileUsername(command *CreateProfileCommand, existingProfiles []Profile) (*CreateProfileCommand, error) {
	for _, existingProfile := range existingProfiles {
		if existingProfile.Username == command.Username {
			return command, errors.New("Username is already taken")
		}
	}
	return command, nil
}
