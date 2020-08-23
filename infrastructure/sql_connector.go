package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kyawmyintthein/orange-contrib/logx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DBDriver string

const (
	Mysql    DBDriver = `mysql`
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
}

type SqlDBConnector interface {
	DB(context.Context) *sqlx.DB
	Config() SqlDBConfig
}

type sqlConnector struct {
	cfg *SqlDBConfig
	db  *sqlx.DB
}

func NewSQLConnector(cfg *SqlDBConfig) (SqlDBConnector, error) {
	connectionString := cfg.URI
	if connectionString == "" {
		switch cfg.Driver {
		case Mysql:
			connectionString = getMySQLConnectionString(cfg)
			break
		case Postgres:
			connectionString = getPostgresConnectionString(cfg)
			break
		case Sqlite3:
			return nil, fmt.Errorf("connection string is empty! Use URI to set sqlite DB")
		}
	}

	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}
	sqlConnector := &sqlConnector{
		cfg: cfg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeOut*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, string(cfg.Driver), connectionString)
	if err != nil {
		logx.Errorf(context.Background(), err, "Failed to connect to mysql database: %s", cfg.DBName)
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime * time.Second)
	err = db.Ping()
	if err != nil {
		logx.Errorf(context.Background(), err, "Failed to connect to database: %s", cfg.DBName)
		return sqlConnector, err
	}
	sqlConnector.db = db
	logx.Infof(context.Background(), "Successfully connected to database: %v", cfg.DBName)
	return sqlConnector, nil
}

func (mysql *sqlConnector) DB(ctx context.Context) *sqlx.DB {
	err := mysql.db.Ping() // reconnect is happen when Ping is called.
	if err != nil {
		logx.Errorf(context.Background(), err, "Failed to reconnect db")
	}
	return mysql.db
}

func (mysql *sqlConnector) Config() SqlDBConfig {
	return *mysql.cfg
}

func getMySQLConnectionString(cfg *SqlDBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true", cfg.Username, cfg.Password, cfg.DBHost, cfg.DBName)
}

func getPostgresConnectionString(cfg *SqlDBConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DBHost, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
}
