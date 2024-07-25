package domain

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/type/date"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	t.Run("should validate that account does not already exists and returns an account to create", func(t *testing.T) {
		command := CreateAccountCommand{
			Name: "My new account",
		}

		result, _ := CreateAccount(&command, nil)

		assert.Equal(
			t,
			&Account{
				Name: "My new account",
			},
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

func TestManagingLinesInAccount(t *testing.T) {
	t.Run("should add a new line to an account", func(t *testing.T) {
		account := Account{
			Name: "My new account",
		}
		amount, _ := decimal.NewFromString("10.3")
		command := AddAccountLineCommand{
			Type: Expense,
			EventDate: &date.Date{
				Year:  2024,
				Month: 7,
				Day:   25,
			},
			Description: "Some expense",
			Amount:      amount,
		}

		AddLineToAccount(&account, &command)

		assert.Equal(
			t,
			&Account{
				Name: "My new account",
				Lines: []AccountLine{
					{
						Type: Expense,
						EventDate: &date.Date{
							Year:  2024,
							Month: 7,
							Day:   25,
						},
						Description: "Some expense",
						Amount:      amount,
					},
				},
			},
			&account,
		)
	})
}
