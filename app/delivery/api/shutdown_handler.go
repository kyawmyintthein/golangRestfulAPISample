package api

import (
	"net/http"
	"strings"

	"github.com/kyawmyintthein/orange-contrib/logx"
)

type ShutdownHandler struct {
	*BaseHandler
	ShutdownSignal chan int
}

func ProvideShutdownHandler(baseHandler *BaseHandler) *ShutdownHandler {
	return &ShutdownHandler{
		baseHandler,
		make(chan int),
	}
}

//curl localhost:3000/stop
func (h *ShutdownHandler) Stop(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.Host, "localhost") {
		h.ShutdownSignal <- 1
	} else {
		logx.Warnf(r.Context(), "Stop API only works for localhost, calling host is %s", r.Host)
	}
	return
}
