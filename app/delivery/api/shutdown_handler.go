package api

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"net/http"
	"strings"
)

type ShutdownHandler struct{
	*BaseHandler
	ShutdownSignal chan int
}

func ProvideShutdownHandler(baseHandler *BaseHandler) *ShutdownHandler{
	return &ShutdownHandler{
		baseHandler,
		make(chan int),
	}
}

//curl localhost:3000/stop
func (h *ShutdownHandler) Stop(w http.ResponseWriter, r *http.Request) {
	log := logging.GetStructuredLogger(r.Context())
	if strings.HasPrefix(r.Host, "localhost") {
		h.ShutdownSignal <- 1
	} else {
		log.Warnf("Stop API only works for localhost, calling host is %s", r.Host)
	}
	return
}

