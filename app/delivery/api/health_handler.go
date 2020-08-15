package api

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/app/viewmodel"
	"net/http"
)

type HealthHandler struct{
	*BaseHandler
}

func ProvideHealthHandler(baseHandler *BaseHandler) *HealthHandler {
	return &HealthHandler{
		baseHandler,
	}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.RenderJSON(r, w, http.StatusOK,
		viewmodel.HealthCheckVM{
		Environment: h.Config.App.Environment,
	})
}