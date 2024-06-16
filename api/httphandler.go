package api

import (
	"banana-back/domain/account"
	"banana-back/repositories"
	"log/slog"
)

type HttpHandler[T any] struct {
	Logger     *slog.Logger
	Repository repositories.Repository[T]
}

func NewHttpHandler(repository repositories.Repository[account.Account]) HttpHandler[account.Account] {
	log := slog.Default()
	return HttpHandler[account.Account]{
		Logger:     log,
		Repository: repository,
	}
}
