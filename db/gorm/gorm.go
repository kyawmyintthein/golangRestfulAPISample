package gorm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/SQLite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"golangRestfulAPISample/bootstrap"
)

var (
	db *gorm.DB
)

// initialize database
func Init() {
	var adapter string
	adapter = bootstrap.App.DBConfig.String("adapter")
	switch adapter {
	case "mysql":
		mysqlConn()
		break
	case "postgres":
		postgresConn()
		break
	}
}

// setupPostgresConn: setup postgres database connection using the configuration from database.yaml
func postgresConn() {
	var (
		connectionString string
		err              error
	)
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", bootstrap.App.DBConfig.String("user"), postgres.Password, postgres.Host, postgres.Port, postgres.Name, postgres.SSLMode)
	if db, err = gorm.Open("postgres", connectionString); err != nil {
		panic(err)
	}

	if err = db.DB().Ping(); err != nil {
		panic(err)
	}

	db.LogMode(true)
	db.Exec("CREATE EXTENSION postgis")

	db.DB().SetMaxIdleConns(beego.AppConfig.DefaultInt(runmode+"::db_max_idle_conns", 10))
	db.DB().SetMaxOpenConns(beego.AppConfig.DefaultInt(runmode+"::db_max_open_conns", 100))
}

// mysqlConn: setup mysql database connection using the configuration from database.yaml
func mysqlConn() {
	var (
		connectionString string
		err              error
	)
	if mysql.Password == "" {
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.Name)
	} else {
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.Name)
	}
	if db, err = gorm.Open("mysql", connectionString); err != nil {
		panic(err)
	}
	if err = db.DB().Ping(); err != nil {
		panic(err)
	}

	db.LogMode(true)
}

/*
 * Gorm: return GORM's postgres database connection instance.
 */
func DBManager() *gorm.DB {
	return db
}
