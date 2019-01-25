package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/dicontainer"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"github.com/swaggo/http-swagger"
)

type RouterInterface interface {
	Routes(serviceContainer *dicontainer.ServiceContainer)
	RouteMultiplexer() *chi.Mux
}

type router struct {
	config *config.GeneralConfig
	logger logging.Logger
	mux    *chi.Mux
}

func NewRouter(generalConfig *config.GeneralConfig, logger logging.Logger) RouterInterface {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RealIP)
	mux.Use(SetJSON)
	mux.Use(logger.NewStructuredLogger())
	mux.Use(cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler)
	return &router{
		mux:    mux,
		config: generalConfig,
		logger: logger,
	}
}

func (h *router) RouteMultiplexer() *chi.Mux {
	return h.mux
}

func (h *router) Routes(container *dicontainer.ServiceContainer) {
	h.mux.Group(func(r chi.Router) {
		r.Get("/health", container.HealthController.HealthCheck)
		r.Get("/health/db", container.HealthController.DBHealthCheck)

		//users
		r.Post("/users", container.UserController.Register)
		r.Put("/users", container.UserController.Update)
		r.Delete("/users/{user_id}", container.UserController.Remove)
		r.Get("/users/{user_id}", container.UserController.GetProfile)
		r.Get("/users", container.UserController.GetAllUsers)
	})

	h.mux.NotFound(container.HttpErrorController.ResourceNotFound)
	h.mux.With(RemoveContextTypeJSON).Get("/swagger/*", httpSwagger.WrapHandler)
}
