package profile

import "errors"

type Profile struct {
	ID        int64
	Username  string
	Email     string
	firstName string
	lastName  string
}

type CreateProfileCommand struct {
	Username  string
	Email     string
	firstName string
	lastName  string
}

func ValidateProfileUsername(command *CreateProfileCommand, existingProfiles []Profile) (*CreateProfileCommand, error) {
	for _, existingProfile := range existingProfiles {
		if existingProfile.Username == command.Username {
			return command, errors.New("Username is already taken")
		}
	}
	return command, nil
}
