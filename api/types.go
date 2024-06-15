package api

import (
	"banana-back/repositories"
	"log/slog"
)

type HttpHandler struct {
	Logger            *slog.Logger
	AccountRepository *repositories.AccountRepository
}
