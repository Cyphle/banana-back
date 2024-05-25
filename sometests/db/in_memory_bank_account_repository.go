package db

type BankAccount struct {
	Name string
}

var bankAccounts []BankAccount

func List() []BankAccount {
	return bankAccounts
}

func Add(account BankAccount) {
	bankAccounts = append(bankAccounts, account)
}

func ToString() string {
	names := ""
	for _, i := range bankAccounts {
		names = names + i.Name + ", "
	}

	return names
}
