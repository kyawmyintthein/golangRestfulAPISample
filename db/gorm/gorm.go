package gorm
import (
	"echo_rest_api/config"
	"github.com/jinzhu/gorm"	
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/SQLite"
	_ "github.com/lib/pq"
)

var (
	postgresConn *gorm.DB
	mysqlConn    *gorm.DB
	sqliteConn   *gorm.DB
	err          error
)

// initialize database 
func Init() {
	switch config.AppConfig.ENV {
	case "dev":
		bootstrap(config.DBConfig.Development)
	case "prod":
		bootstrap(config.DBConfig.Production)
	case "test":
		bootstrap(config.DBConfig.Test)
	}
}

// bootstrap: setup database based on  using the enviroment configuration from application.yml
func bootstrap(ienv interface{}) {
	switch ienv.(type) {
	case config.Production:
		prod := config.DBConfig.Production
		if prod.Postgres != (config.Postgres{}) {
			setupPostgresConn(&prod.Postgres)
		}
		if prod.Mysql != (config.Mysql{}) {
			setupMysqlConn(&prod.Mysql)
		}
		if prod.SQLite != (config.SQLite{}) {
			setupSQLiteConn(&prod.SQLite)
		}
	case config.Development:
		dev := config.DBConfig.Development
		if dev.Postgres != (config.Postgres{}) {
			setupPostgresConn(&dev.Postgres)
		}

		if dev.Mysql != (config.Mysql{}) {
			setupMysqlConn(&dev.Mysql)
		}

		if dev.SQLite != (config.SQLite{}) {
			setupSQLiteConn(&dev.SQLite)
		}
	case config.Test:
		test := config.DBConfig.Test
		if test.Postgres != (config.Postgres{}) {
			setupPostgresConn(&test.Postgres)
		}

		if test.Mysql != (config.Mysql{}) {
			setupMysqlConn(&test.Mysql)
		}

		if test.SQLite != (config.SQLite{}) {
			setupSQLiteConn(&test.SQLite)
		}
	}
}


// setupPostgresConn: setup postgres database connection using the configuration from appliation.yml
func setupPostgresConn(postgres *config.Postgres) {
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", postgres.User, postgres.Password, postgres.Host, postgres.Port, postgres.Name, postgres.SSLMode)
	postgresConn, err := gorm.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	err = postgresConn.DB().Ping()
	if err != nil {
		panic(err)
	}
	postgresConn.LogMode(true)
	postgresConn.DB().SetMaxIdleConns(postgres.MaxIdleConns)
}

// setupMysqlConn: setup mysql database connection using the configuration from appliation.yml
func setupMysqlConn(mysql *config.Mysql) {
	var connectionString string
	if mysql.Password == ""{
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.Name)
	}else{
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.Name)
	}
	mysqlConn, err = gorm.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	err = mysqlConn.DB().Ping()
	if err != nil {
		panic(err)
	}
	mysqlConn.LogMode(true)
	// mysqlConn.DB().SetMaxIdleConns(mysql.MaxIdleConns)
}

// setupSQLiteConn: setup SQLite database connection using the configuration from application.yml
func setupSQLiteConn(SQLite *config.SQLite) {
}

// PostgresConn: return postgres connection from gorm ORM
func PostgresConn() *gorm.DB {
	return postgresConn
}

// MysqlConn: return mysql connection from gorm ORM
func MysqlConn() *gorm.DB {
	return mysqlConn
}

// SQLiteConn: return SQLite connection from gorm ORM
func SQLiteConn() *gorm.DB {
	return sqliteConn
}
