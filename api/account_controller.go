package api

import (
	"github.com/labstack/echo/v4"
)

func ActivateAccountRoutes(e *echo.Echo, handler HttpHandler) {
	handler.Logger.Info("Activating account routes")
	e.GET("/accounts", handler.getAccounts)
}
