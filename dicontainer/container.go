package dicontainer

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/controller"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/datastore"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/service"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/infrastructure"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
)

/*
 * ServiceContainer resolve all dependencies between controller, service, infrastructure except
 * application level dependencies such us logging, config and etc ...
 */
type ServiceContainer struct{
	config *config.GeneralConfig
	logger logging.Logger

	// controllers
	HealthController *controller.HealthController
	HttpErrorController *controller.HttpErrorController

	UserController  *controller.UserController
}

func NewServiceContainer(config* config.GeneralConfig, logger logging.Logger) *ServiceContainer{
	return &ServiceContainer{
		config: config,
		logger: logger,
	}
}

func (container *ServiceContainer) InitDependenciesInjection() {
	log := container.logger.Logrus()

	// init mongodb
	mongoStore, err := infrastructure.NewMongoStore(container.config, container.logger)
	if err != nil{
		log.Fatalf("Failed to connect to database: %v", container.config.Mongodb.Database)
	}

	// datastores
	userDatastore := datastore.NewUserDatastore(context.Background(), container.config, container.logger, mongoStore)

	// services
	healthService := &service.HealthService{container.config, mongoStore}
	userService := &service.UserService{container.config, container.logger, userDatastore}

	// controllers
	baseController := controller.BaseController{Config: container.config, Logging:  container.logger}
	container.HttpErrorController = &controller.HttpErrorController{baseController}
	container.HealthController = &controller.HealthController{baseController, healthService}
	container.UserController = &controller.UserController{baseController, userService}
}