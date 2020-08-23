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
	"github.com/kyawmyintthein/golangRestfulAPISample/app/usecase"
	"github.com/kyawmyintthein/orange-contrib/logx"

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
	logger     logx.Logger

	userRepository repository.UserRepository
	userUsecase    usecase.UserUsecase

	// API handlers
	HealthHandler   *api.HealthHandler
	ShutdownHandler *api.ShutdownHandler
	UserHandler     *api.UserHandler
	ArticleHandler  *api.ArticleHandler
}

// Start other services if necessary
func (app *restApiApplication) Init() error {
	app.routes()
	return nil
}

// Start serve http server
func (app *restApiApplication) StartHttpServer() error {
	if app.config.GracefulShutdown.Enabled {
		go func() {
			err := app.httpServer.ListenAndServe()
			if err != nil {
				logx.Errorf(context.Background(), err, "Failed to start http server at port :%d ", app.config.App.HttpPort)
			}
		}()
		logx.Infof(context.Background(), "Server started at port with graceful shutdown :%d ", app.config.App.HttpPort)
		app.startGracefulShutdownChan()
		return nil
	}
	logx.Infof(context.Background(), "Server started at port :%d ", app.config.App.HttpPort)
	return app.httpServer.ListenAndServe()
}

func (app *restApiApplication) startGracefulShutdownChan() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGINT)
	// interrupt signal sent from terminal
	signal.Notify(gracefulStop, os.Interrupt)
	// sigterm signal sent from kubernetes
	signal.Notify(gracefulStop, syscall.SIGTERM)
	logx.Info(context.Background(), "Enabled Graceful Shutdown")

	select {
	case sig := <-gracefulStop:
		logx.Warnf(context.Background(), "caught sig: %+v", sig)
	case <-app.ShutdownHandler.ShutdownSignal:
		logx.Warn(context.Background(), "stop api called")
	}

	app.shutDownProcesses()
}

func (app *restApiApplication) shutDownProcesses() {

	// We received an interrupt signal, server shut down.
	if err := app.httpServer.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		logx.Error(context.Background(), err, "HTTP server Shutdown")
	}
	logx.Infof(context.Background(), "Http server shut down finished.")
	// this sleep is to give buffer so that any pending process can completes gracefully,.
	logx.Infof(context.Background(), "Wait for %v to finish processing", app.config.GracefulShutdown.Timeout*time.Second)
	time.Sleep(app.config.GracefulShutdown.Timeout * time.Second)
	logx.Infof(context.Background(), "Shutting down.")
	os.Exit(0)
}
