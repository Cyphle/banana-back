package api

import (
	"banana-back/domain/account"
	"github.com/labstack/echo/v4"
)

func ActivateAccountRoutes(e *echo.Echo, handler HttpHandler[account.Account]) {
	handler.Logger.Info("Activating account routes")
	g := e.Group("/accounts")
	g.GET("", handler.getAccounts)
	g.POST("", handler.createAccount)
}
