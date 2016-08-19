package main

import (
	"echo_rest_api/db/gorm"
	"echo_rest_api/config"
	"echo_rest_api/models"
	"echo_rest_api/controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	// init database
	gorm.Init()
	autoDropTables()
	autoCreateTables()
	autoMigrateTables()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	//static file serviing
	e.Static("/static", "assets")

	// Routers
	e.POST("/users", controllers.CreateUser)
	e.GET("/users/:id", controllers.ShowUser)
	e.GET("/users", controllers.AllUsers)
	e.PUT("/users/:id", controllers.UpdateUser)
	e.DELETE("/users/:id",controllers.DeleteUser)

	// Server
	e.Run(standard.New(":1323"))
}

// autoCreateTables: create database tables using GORM
func autoCreateTables() {
	if !gorm.MysqlConn().HasTable(&models.User{}) {
		gorm.MysqlConn().CreateTable(&models.User{})
	}
}

// autoMigrateTables: migrate table columns using GORM
func autoMigrateTables() {
	gorm.MysqlConn().AutoMigrate(&models.User{})
}


// auto drop tables on dev mode
func autoDropTables() {
    if config.AppConfig.ENV == "dev" {
        gorm.MysqlConn().DropTableIfExists(&models.User{}, &models.User{})
    }
}

