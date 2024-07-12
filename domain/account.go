package domain

type Account struct {
	ID   int64
	Name string
}

type CreateAccountCommand struct {
	Name string
}

func CreateAccount(repository Repository[Account], request *CreateAccountCommand) {
	// TODO ici ça devrait être une couche métier qui reçoit une request, qui transforme en commande si ok et qui save dans le repo
	// TODO some business logic
	// Dans repository il faudrait, findOneByField(fieldName, criteria) => select * from documentstorage.documents where UPPER(name) LIKE UPPER('%dpe%');
}
