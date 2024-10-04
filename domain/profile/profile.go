package profile

import (
	"fmt"
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

func ValidateProfileUsername(command *CreateProfileCommand, existingProfile *Profile) (*CreateProfileCommand, error) {
	if existingProfile != nil {
		return nil, fmt.Errorf("username already taken")
	}
	return command, nil
}
