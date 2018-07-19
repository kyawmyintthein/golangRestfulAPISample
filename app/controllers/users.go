package controllers

import (
	"net/http"
	"golangRestfulAPISample/app"
	"github.com/labstack/echo"
)

app.Server.POST("/users", func(c echo.Context) error{
	return  c.JSON(http.StatusOk, nil)
})
