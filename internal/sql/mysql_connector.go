package sql

import (
	"bitbucket.org/libertywireless/circles-framework/cllogging"
	"bitbucket.org/libertywireless/circles-framework/cloption"
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "gopkg.in/goracle.v2"
	"time"
)

type DBDriver string

const (
	Mysql    DBDriver = `mysql`
	Oralce   DBDriver = `goracle`
	Postgres DBDriver = `postgres`
	Sqlite3  DBDriver = `sqlite3`
)

type SqlDBConfig struct {
	Driver          DBDriver      `mapstructure:"driver" json:"driver"`
	DBName          string        `mapstructure:"db_name" json:"db_name"`
	DBHost          string        `mapstructure:"db_host" json:"db_host"`
	DialTimeOut     time.Duration `mapstructure:"dial_timeout" json:"dial_timeout"` // second
	MaxIdleConns    int           `mapstructure:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" json:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_life_time" json:"conn_max_life_time"` // second
	Username        string        `mapstructure:"username" json:"username"`
	Password        string        `mapstructure:"password" json:"password"`
	SSLMode         string        `mapstructure:"ssl_mode" json:"ssl_mode"`
	URI             string        `mapstructure:"uri" json:"uri"`

	// Oracle
	ServiceName string `mapstructure:"service_name" json:"service_name,omitempty"`
	SID         string `mapstructure:"sid" json:"sid,omitempty"`
}

type SqlDBConnector interface {
	DB(context.Context) *sqlx.DB
	Config() SqlDBConfig
}

type sqlConnector struct {
	cfg    *SqlDBConfig
	db     *sqlx.DB
	logger cllogging.KVLogger
}

func NewSQLConnector(cfg *SqlDBConfig, opts ...cloption.Option) (SqlDBConnector, error) {
	options := cloption.NewOptions(opts...)
	connectionString := cfg.URI
	if connectionString == "" {
		switch cfg.Driver {
		case Mysql:
			connectionString = getMySQLConnectionString(cfg)
			break
		case Oralce:
			connectionString = getOracleConnectionString(cfg)
		case Postgres:
			connectionString = getPostgresConnectionString(cfg)
			break
		case Sqlite3:
			return nil, fmt.Errorf("connection string is empty! Use URI to set sqlite DB")
		}
	}

	if cfg.SSLMode == ""{
		cfg.SSLMode = "disable"
	}
	sqlConnector := &sqlConnector{
		cfg: cfg,
	}

	// set logger
	lgr, ok := options.Context.Value(loggerKey{}).(cllogging.KVLogger)
	if lgr == nil && !ok {
		sqlConnector.logger = cllogging.DefaultKVLogger()
	} else {
		sqlConnector.logger = lgr
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeOut*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, string(cfg.Driver), connectionString)
	if err != nil {
		sqlConnector.logger.Errorf(err, "Failed to connect to mysql database: %s", cfg.DBName)
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime * time.Second)
	err = db.Ping()
	if err != nil {
		sqlConnector.logger.Errorf(err, "Failed to connect to database: %s", cfg.DBName)
		return sqlConnector, err
	}
	sqlConnector.db = db
	sqlConnector.logger.Infof("Successfully connected to database: %v", cfg.DBName)
	return sqlConnector, nil
}

func (mysql *sqlConnector) DB(ctx context.Context) *sqlx.DB {
	log := cllogging.GetStructuredLogger(ctx)
	err := mysql.db.Ping() // reconnect is happen when Ping is called.
	if err != nil {
		log.WithError(err).Errorln("Failed to reconnect db")
	}
	return mysql.db
}

func (mysql *sqlConnector) Config() SqlDBConfig {
	return *mysql.cfg
}

func getMySQLConnectionString(cfg *SqlDBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true", cfg.Username, cfg.Password, cfg.DBHost, cfg.DBName)
}

func getOracleConnectionString(cfg *SqlDBConfig) string {
	var service string
	if cfg.ServiceName != "" {
		service = cfg.ServiceName
	} else if cfg.SID != "" {
		service = cfg.SID
	}
	return fmt.Sprintf("%s/%s@%s/%s", cfg.Username, cfg.Password, cfg.DBHost, service)
}

func getPostgresConnectionString(cfg *SqlDBConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DBHost, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
}

// TODO:: complete this for other drivers Oracle, Postgres, Sqlite3
// Note:: only support for mysql
func IsDuplicateError(ctx context.Context, err error) bool {
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}

	if mysqlErr.Number == 1062 {
		return true
	}
	return false
}

// TODO:: complete this for other drivers Oracle, Postgres, Sqlite3
// Note:: only support for mysql
func IsNotFoundError(ctx context.Context, err error) bool {
	if sql.ErrNoRows.Error() == err.Error() {
		return true
	}
	return false
}
