package app

import (
	"github.com/go-chi/chi"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
)

func (app *restApiApplication) routes() {
	router := app.router
	router.Group(func(r chi.Router) {
		r.Use(logging.NewRequestStructuredLogger())
		router.Get("/health", app.HealthHandler.HealthCheck)
		router.Get("/stop", app.ShutdownHandler.Stop)
	})
}
