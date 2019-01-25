package app

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/dicontainer"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"github.com/kyawmyintthein/golangRestfulAPISample/router"
	"github.com/kyawmyintthein/golangRestfulAPISample/docs"
	"net/http"
)

// AppInterface is object which wrap the necessary modules for this system.
// It will export high level interface function.
type AppInterface interface {
	// start all dependencies services
	Init()

	// Start http server
	Start(serverPort string) error
}


type RestApiApplication struct{
	config *config.GeneralConfig
	logger logging.Logger

	serviceContainer *dicontainer.ServiceContainer
	router 		     router.RouterInterface
}

// New : build new RestApiApplication object with config and logging
func New(configFilePaths ...string) AppInterface {

	// init config
	generalConfig := config.Loadconfig(configFilePaths...)
	logger := logging.InitializeLogger(
		generalConfig.Log.LogLevel,
		generalConfig.Log.LogFilePath,
		generalConfig.Log.JsonLogFormat,
		generalConfig.Log.LogRotation)

	return &RestApiApplication{
		config: generalConfig,
		logger: logger,
	}
}

func (app *RestApiApplication) Init() {

	// start dependencies injection
	app.serviceContainer = dicontainer.NewServiceContainer(app.config, app.logger)
	app.serviceContainer.InitDependenciesInjection()

	// init new handler
	app.router = router.NewRouter(app.config, app.logger)
	app.router.Routes(app.serviceContainer)
	app.overrideSwaggerInfo()
}

// init swagger after registered all endpoints
func (app *RestApiApplication) overrideSwaggerInfo() {
	docs.SwaggerInfo.Host = app.config.Swagger.Host
	docs.SwaggerInfo.Version = app.config.Swagger.Version
	docs.SwaggerInfo.BasePath = app.config.Swagger.BasePath
}

// Start serve http server
func (app *RestApiApplication) Start(serverPort string) error{
	app.logger.Logrus().Infoln("############################## Server Started ##############################")
	return http.ListenAndServe(":"+serverPort, app.router.RouteMultiplexer())
}


