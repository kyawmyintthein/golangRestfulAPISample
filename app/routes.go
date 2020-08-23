package app

import (
	"github.com/go-chi/chi"
)

func (app *restApiApplication) routes() {
	router := app.router
	router.Group(func(r chi.Router) {
		r.Use(app.logger.NewRequestLogger())
		router.Get("/health", app.HealthHandler.HealthCheck)
		router.Get("/stop", app.ShutdownHandler.Stop)

		router.Post("/articles", app.ArticleHandler.CreateArticle)
		router.Get("/articles/{url}", app.ArticleHandler.GetArticle)
	})
}
