package profile

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateProfile(t *testing.T) {
	t.Run("should return the profile when validating that the username is not already taken when creating a profile", func(t *testing.T) {
		command := CreateProfileCommand{
			Username: "johndoe",
		}

		result, _ := ValidateProfileUsername(&command, []Profile{})

		assert.Equal(
			t,
			&CreateProfileCommand{
				Username: "johndoe",
			},
			result,
		)
	})
}
