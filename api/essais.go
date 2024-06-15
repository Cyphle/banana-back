package api

import (
	"banana-back/repositories"
	"context"
)

type MockAccountRepository struct {
}

func (m MockAccountRepository) List(ctx context.Context) ([]repositories.AccountEntity, error) {
	res := make([]repositories.AccountEntity, 0)
	return res, nil
}

func Other() {
	// Le truc qui reçoit l interface dans sa structure peut être instantié avec un pointer d'une implémentation de l'interface
	//mock := &MockAccountRepository{}
	//pouet := NewHttpHandler(mock)
}
