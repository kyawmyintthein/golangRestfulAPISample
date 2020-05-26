package logging

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
)

type LoggerInterface interface {
	GetStructuredLogger(ctx context.Context) *StructuredLoggerEntry
	SetLogLevel(level string) error
	SetReportCaller(bool)
	GetLogrus() *logrus.Logger
	NewRequestStructuredLogger() func(next http.Handler) http.Handler
}
