package main

import (
	"golangRestfulAPISample/app"
	"golangRestfulAPISample/app/models"
	"golangRestfulAPISample/config"
	"golangRestfulAPISample/db/gorm"
)

func main() {
	// init server
	app.Init()

	// init database
	gorm.Init()
	autoDropTables()
	autoCreateTables()
	autoMigrateTables()

	// run server
	app.Server.Logger.Fatal(app.Server.Start(":1313"))
}

// autoCreateTables: create database tables using GORM
func autoCreateTables() {
	if !gorm.DBManager().HasTable(&models.User{}) {
		gorm.DBManager().CreateTable(&models.User{})
	}
}

// autoMigrateTables: migrate table columns using GORM
func autoMigrateTables() {
	gorm.DBManager().AutoMigrate(&models.User{})
}

// auto drop tables on dev mode
func autoDropTables() {
	if config.AppConfig.ENV == "dev" {
		gorm.DBManager().DropTableIfExists(&models.User{}, &models.User{})
	}
}
