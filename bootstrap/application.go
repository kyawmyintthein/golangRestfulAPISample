package bootstrap

import (
	"github.com/fsnotify/fsnotify"
	"gl.tzv.io/spf13/viper"
	"fmt"
)

var App *Application

type Config viper.Viper

type Application struct {
	Name      string  `json:"name"`
	Version   string  `json:"version"`
	ENV       string  `json:"env"`
	AppConfig Config `json:"application_config"`
	DBConfig  Config `json:"database_config"`
}

func init() {
	App = &Application{}
	App.Name = APP_NAME
	App.Version = APP_VERSION
	App.loadENV()
	App.loadAppConfig()
	App.loadDBConfig()
}

// loadAppConfig: read application config and build viper object
func (app *Application) loadAppConfig() {
	var (
		appConfig *viper.Viper
		err       error
	)
	appConfig = viper.New()
	appConfig.SetEnvKeyReplacer(REPLACER)
	appConfig.SetEnvPrefix(APP_CONFIG_PREFIX)
	appConfig.AutomaticEnv()
	appConfig.SetConfigName(APP_CONFIG_NAME)
	appConfig.AddConfigPath(CONFIG_PATH)
	appConfig.SetConfigType(CONFIG_FILE_TYPE)
	if err = appConfig.ReadInConfig(); err != nil {
		panic(err)
	}
	appConfig.WatchConfig()
	appConfig.OnConfigChange(func(e fsnotify.Event) {
		//	glog.Info("App Config file changed %s:", e.Name)
	})
	app.AppConfig = Config(*appConfig)
}

// loadDBConfig: read application config and build viper object
func (app *Application) loadDBConfig() {
	var (
		dbConfig *viper.Viper
		err      error
	)
	dbConfig = viper.New()
	dbConfig.SetEnvKeyReplacer(REPLACER)
	dbConfig.SetEnvPrefix(DB_CONFIG_PREFIX)
	dbConfig.AutomaticEnv()
	dbConfig.SetConfigName(DB_CONFIG_NAME)
	dbConfig.AddConfigPath(CONFIG_PATH)
	dbConfig.SetConfigType(CONFIG_FILE_TYPE)
	if err = dbConfig.ReadInConfig(); err != nil {
		panic(err)
	}
	dbConfig.WatchConfig()
	dbConfig.OnConfigChange(func(e fsnotify.Event) {
		//	glog.Info("App Config file changed %s:", e.Name)
	})
	app.DBConfig = Config(*dbConfig)
}

// loadENV
func (app *Application) loadENV() {
	var APPENV string
	var appConfig viper.Viper
	appConfig = viper.Viper(app.AppConfig)
	APPENV = appConfig.GetString("env")
	switch APPENV {
	case DEV_ENV:
		app.ENV = DEV_ENV
		break
	case STAGING_ENV:
		app.ENV = STAGING_ENV
		break
	case PROD_ENV:
		app.ENV = PROD_ENV
		break
	default:
		app.ENV = DEV_ENV
		break
	}
}

// String: read string value from viper.Viper
func (config *Config) String(key string) string {
	var viperConfig viper.Viper
	viperConfig = viper.Viper(*config)
	return viperConfig.GetString(fmt.Sprintf("%s.%s", App.ENV, key))
}

// Int: read int value from viper.Viper
func (config *Config) Int(key string) int {
	var viperConfig viper.Viper
	viperConfig = viper.Viper(*config)
	return viperConfig.GetInt(fmt.Sprintf("%s.%s", App.ENV, key))
}

// Boolean: read boolean value from viper.Viper
func (config *Config) Boolean(key string) bool {
	var viperConfig viper.Viper
	viperConfig = viper.Viper(*config)
	return viperConfig.GetBool(fmt.Sprintf("%s.%s", App.ENV, key))
}
