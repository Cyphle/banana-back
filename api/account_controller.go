package api

import (
	"banana-back/repositories"
	"github.com/labstack/echo/v4"
)

func ActivateAccountRoutes(e *echo.Echo, handler HttpHandler[repositories.AccountEntity]) {
	//handler.Logger.Info("Activating account routes")
	e.GET("/accounts", handler.getAccounts)
}
