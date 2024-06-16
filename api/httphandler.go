package api

import (
	"banana-back/repositories"
	"log/slog"
)

type HttpHandler[T any] struct {
	Logger     *slog.Logger
	Repository repositories.Repository[T]
}

func NewHttpHandler(repository repositories.Repository[repositories.AccountEntity]) HttpHandler[repositories.AccountEntity] {
	return HttpHandler[repositories.AccountEntity]{
		Repository: repository,
	}
}
