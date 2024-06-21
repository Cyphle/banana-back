package api

import "banana-back/domain"

type CreateAccountCommandView struct {
	Name string `json:"name"`
}

type ArrayResponse[T any] struct {
	Data []domain.Account `json:"data"`
}
