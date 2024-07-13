package domain

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	t.Run("should validate that account does not already exists and return a command", func(t *testing.T) {
		command := CreateAccountCommand{
			Name: "My new account",
		}

		result, _ := CreateAccount(&command, nil)

		assert.Equal(
			t,
			&Account{Name: "My new account"},
			result,
		)
	})

	t.Run("should return an error when account name already exist", func(t *testing.T) {
		command := CreateAccountCommand{
			Name: "My new account",
		}
		existingAccount := Account{
			Name: "My new account",
		}

		_, err := CreateAccount(&command, &existingAccount)

		assert.Equal(
			t,
			fmt.Errorf("name already exists"),
			err,
		)
	})
}
