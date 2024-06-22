package api

import "banana-back/domain"

type CreateAccountCommandView struct {
	Name string `json:"name"`
}

type AccountIdPathParam struct {
	ID int64 `param:"id" query:"id"`
}

type ArrayResponse[T any] struct {
	Data []domain.Account `json:"data"`
}
