package app

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var Server *echo.Echo

// init echo web server
func Init() {
	Server = echo.New()
	// Middleware
	Server.Use(middleware.Logger())
	Server.Use(middleware.Recover())

	//CORS
	Server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	//static file serviing
	Server.Static("/static", "assets")
}
