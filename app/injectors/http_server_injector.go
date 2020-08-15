package injectors

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"net/http"
)

func ProvideHttpServer(config *config.GeneralConfig, router *chi.Mux) (*http.Server) {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.App.HttpPort),
		Handler: router,
	}
	return httpServer
}

