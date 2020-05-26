package cllogging

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"time"
)

type CustomLogEntry interface {
	Write(int, int, http.Header, time.Duration, interface {})
	WriteError(status, bytes int, elapsed time.Duration)
	WriteWarn(status, bytes int, elapsed time.Duration)
	Panic(v interface{}, stack []byte)
}

// LogFormatter initiates the beginning of a new LogEntry per request.
// See DefaultLogFormatter for an example implementation.
type LogFormatter interface {
	NewLogEntry(r *http.Request) middleware.LogEntry
	NewStructuredEntry(r *http.Request) CustomLogEntry
}

func GinRequestLogger(logger LogFormatter) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		logEntry := logger.NewStructuredEntry(c.Request)
		ctx := context.WithValue(c.Request.Context(), middleware.LogEntryCtxKey.String(), logEntry)
		c.Request = c.Request.WithContext(ctx)
		c.Set(middleware.LogEntryCtxKey.String(), logEntry)
		c.Next()
		stop := time.Since(start)
		statusCode := c.Writer.Status()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		if statusCode > 499 {
			logEntry.WriteError(statusCode, dataLength, stop)
		} else if statusCode > 399 {
			logEntry.WriteWarn(statusCode, dataLength, stop)
		} else {
			logEntry.Write(statusCode, dataLength, nil, stop, nil)
		}
	}
}
