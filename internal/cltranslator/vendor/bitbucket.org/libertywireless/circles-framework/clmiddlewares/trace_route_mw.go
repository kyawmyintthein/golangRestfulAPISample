package clmiddlewares

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

func TrackRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		rctx := chi.RouteContext(r.Context())
		routePattern := fmt.Sprintf("[%s] %s", rctx.RouteMethod, rctx.RoutePattern())
		if span := opentracing.SpanFromContext(r.Context()); span != nil {
			span.SetOperationName(routePattern)
			if reqID := middleware.GetReqID(r.Context()); reqID != "" {
				span.SetTag("request_id", reqID)
			}
		}
	})
}
