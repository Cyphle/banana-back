package domain

import (
	"banana-back/api"
)

type Account struct {
	ID   int64
	Name string
}

type CreateAccountCommand struct {
	Name string
}

func CreateAccount(repository Repository[Account], request *api.CreateAccountRequest) {
	// TODO ici ça devrait être une couche métier qui reçoit une request, qui transforme en commande si ok et qui save dans le repo
	// TODO some business logic
}
