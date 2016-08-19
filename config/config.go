package config

import (
	"echo_rest_api/library/configor"
)

var ENVs = [...]string{
	"prod",
	"dev",
	"test",
}

type (
	Postgres struct {
		Name     string 	`json:"name"`
		User     string 	`default:"postgres"`
		Password string 	`required:"true" env:"db_password"`
		Host     string 	`json:"host"`
		Port     string 		`default:"3306"`
		SSLMode  bool 		`default:"false"`
		MaxIdleConns int 	`default:"10" json:"max_idle_conns"`
	}

	Mysql struct {
		Name     string `json:"name"`
		User     string `json:"user" default:"root"`
		Password string `json:"password" required:"true" env:"db_password"`
		Host     string `json:"host"`
		Port     string `default:"3306" json:"port"`
	}

	SQLite struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"db_password"`
		Host     string
		Port     string `default:"3306"`
		SSLMode  bool `default:"false"`
	}

	Mongo struct {
		Name    string
		MaxPool uint `default:"10"`
		Path    string
	}

	Redis struct {
		Address  string
		DB       int `default:"0"`
		Password string
	}

	Production struct {
		Mysql    Mysql
		Postgres Postgres
		SQLite   SQLite
		Mongo    Mongo
		Redis    Redis
	}

	Development struct {
		Mysql    Mysql
		Postgres Postgres
		SQLite   SQLite
		Mongo    Mongo
		Redis    Redis
	}

	Test struct {
		Mysql    Mysql
		Postgres Postgres
		SQLite   SQLite
		Mongo    Mongo
		Redis    Redis
	}
)


var AppConfig = struct {
		AppName 	string `json:"APP_NAME" default:"app name"`
		ENV     	string `json:"ENV" default:"dev"`
}{}


var DBConfig = struct {
		Production 	Production	 `json:"production"`
		Development Development  `json:"development"`
		Test 		Test		 `json:"test"`
}{}

func init() {
	configor.Load(&AppConfig, "config/application.json")
	configor.Load(&DBConfig, "config/database.json")
}
