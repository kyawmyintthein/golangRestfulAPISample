package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/delivery/api"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/repository"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/service"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"

	//"fmt"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	//"net/http"
)

// AppInterface is object which wrap the necessary modules for this system.
// It will export high level interface function.
type App interface {
	// start all dependencies services
	Init() error
	// Start http server
	StartHttpServer() error
}

type restApiApplication struct {
	config     *config.GeneralConfig
	router     *chi.Mux
	httpServer *http.Server

	userRepository repository.UserRepository
	userService    service.UserService

	// API handlers
	HealthHandler   *api.HealthHandler
	ShutdownHandler *api.ShutdownHandler
	UserHandler     *api.UserHandler
	ArticleHandler  *api.ArticleHandler
}

// Start other services if necessary
func (app *restApiApplication) Init() error {
	logging.Init(&app.config.Log)
	app.routes()
	return nil
}

// Start serve http server
func (app *restApiApplication) StartHttpServer() error {
	log := logging.GetStructuredLogger(context.Background())
	if app.config.GracefulShutdown.Enabled {
		go func() {
			err := app.httpServer.ListenAndServe()
			if err != nil {
				log.WithError(err).Errorf("Failed to start http server at port :%d ", app.config.App.HttpPort)
			}
		}()
		log.Infof("Server started at port with graceful shutdown :%d ", app.config.App.HttpPort)
		app.startGracefulShutdownChan()
		return nil
	}
	log.Infof("Server started at port :%d ", app.config.App.HttpPort)
	return app.httpServer.ListenAndServe()
}

func (app *restApiApplication) startGracefulShutdownChan() {
	log := logging.GetStructuredLogger(context.Background())

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGINT)
	// interrupt signal sent from terminal
	signal.Notify(gracefulStop, os.Interrupt)
	// sigterm signal sent from kubernetes
	signal.Notify(gracefulStop, syscall.SIGTERM)
	log.Info("Enabled Graceful Shutdown")

	select {
	case sig := <-gracefulStop:
		log.Warnf("caught sig: %+v", sig)
	case <-app.ShutdownHandler.ShutdownSignal:
		log.Warn("stop api called")
	}

	app.shutDownProcesses()
}

func (app *restApiApplication) shutDownProcesses() {

	log := logging.GetStructuredLogger(context.Background())

	// We received an interrupt signal, server shut down.
	if err := app.httpServer.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		log.WithError(err).Error("HTTP server Shutdown: %v", err)
	}
	log.Infof("Http server shut down finished.")
	// this sleep is to give buffer so that any pending process can completes gracefully,.
	log.Infof("Wait for %v to finish processing", app.config.GracefulShutdown.Timeout*time.Second)
	time.Sleep(app.config.GracefulShutdown.Timeout * time.Second)
	log.Infof("Shutting down.")
	os.Exit(0)
}
